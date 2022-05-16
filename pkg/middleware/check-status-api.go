package middleware

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
	"net/http"

	"github.com/bhojpur/crm/pkg/models"
	u "github.com/bhojpur/crm/pkg/utils"
)

// Проверяет общий статус API интерфейса для ВСЕХ аккаунтов
func CheckApiStatus(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Подгружаем настройки CRM
		crmSettings, err := models.GetCrmSettings()

		if err != nil {
			u.Respond(w, u.MessageError(nil, "Server is unavailable")) // что это?)
			return
		}

		// Проверяем статус API для всех клиентов
		if !crmSettings.ApiEnabled {
			u.Respond(w, u.MessageError(nil, crmSettings.ApiDisabledMessage))
			return
		}

		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

// Проверяет статус App UI-API (для GUI Bhojpur CRM)
func CheckAppUiApiStatus(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Подгружаем настройки CRM
		crmSettings, err := models.GetCrmSettings()

		if err != nil {
			u.Respond(w, u.MessageError(nil, "Server is unavailable")) // что это?)
			return
		}

		// Проверяем статус UI-API для главного приложения (APP/GUI)
		if !crmSettings.AppUiApiEnabled {
			u.Respond(w, u.MessageError(nil, crmSettings.AppUiApiDisabledMessage))
			return
		}

		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

// Проверяет общий статус UI-API интерфейса для ВСЕХ аккаунтов
func CheckUiApiStatus(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Подгружаем настройки CRM
		crmSettings, err := models.GetCrmSettings()

		if err != nil {
			u.Respond(w, u.MessageError(nil, "Сервер не может обработать запрос")) // что это?)
			return
		}

		// Проверяем статус UI-API для всех клиентов
		if !crmSettings.UiApiEnabled {
			u.Respond(w, u.MessageError(nil, crmSettings.UiApiDisabledMessage))
			return
		}

		next.ServeHTTP(w, r)
	})

}
