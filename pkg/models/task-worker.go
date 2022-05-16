package models

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"log"
	"time"
)

// обходит email_queue_workflows каждые N секунд
// Считывает по 100 актуальных задач
// Получает данные для отправки писем
// Добавляет письма в MTA-server

func init() {
	go taskWorker()
}

func taskWorker() {

	for {
		if db == nil {
			time.Sleep(time.Millisecond * 2000)
			continue
		}

		tasks := make([]TaskScheduler, 0)

		// Получаем задачи у которых статус planned
		err := db.Model(&TaskScheduler{}).
			Where("expected_time_to_start <= ? AND status = ?", time.Now().UTC(), WorkStatusPlanned).Limit(10).
			Find(&tasks).Error
		if err != nil {
			log.Printf("TaskScheduler:  %v", err)
			time.Sleep(time.Second * 10)
			continue
		}

		// Подготавливаем отправку
		for i := range tasks {

			// 1. Перед началом меняем статус у задачи с pending на planned
			if err := tasks[i].SetStatus(WorkStatusPlanned); err != nil {
				log.Printf("taskWorker: Ошибка установки статуса Pending задачи [%v]: %v", tasks[i].Id, err)
				continue
			}

			// 2. Запускаем выполнение задачи и ожидаем результата...
			err = tasks[i].Execute()

			if err != nil {

				// Если задача провалена - ставим статус failed
				if err := tasks[i].SetStatus(WorkStatusFailed, err.Error()); err != nil {
					log.Printf("taskWorker: Ошибка установки статуса Failed у задачи [%v]: %v", tasks[i].Id, err)
				}

			} else {
				// Если задача выполнена - ставим статус completed
				if err := tasks[i].SetStatus(WorkStatusCompleted); err != nil {
					log.Printf("taskWorker: Ошибка установки статуса Completed задачи [%v]: %v", tasks[i].Id, err)
				}
			}

			continue
		}

		if len(tasks) > 100 {
			time.Sleep(time.Minute * 2)
			continue
		}

		time.Sleep(time.Second * 20)
		continue
	}
}
