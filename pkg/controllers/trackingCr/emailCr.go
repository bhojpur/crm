package trackingCr

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
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bhojpur/crm/pkg/controllers/utilsCr"
)

//  Query:
//  ?u={userHashId}   <<< Минимальный пакет данных
// 	?u={userHashId}&?i={mtaHistoryId}&hi={hashId}
// 	?u={userHashId}&hi={hashId}
func UnsubscribeUser(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	userHashId, ok := utilsCr.GetQuerySTRVarFromGET(r, "u")
	if !ok {
		retHTML(w, "Необходимо указать пользователя")
		return
	}

	// Получаем пользователя
	user, err := account.GetUserByHashId(userHashId)
	if err != nil {
		retHTML(w, "Пользователь не найден!")
		return
	}

	// Если пользователь уже отписан
	if !user.Subscribed {
		retHTML(w, "Пользователь уже отписан от всех рассылок")
		return
	}

	if err := user.Unsubscribing(); err != nil {
		retHTML(w, "Ошибка во время отписки пользователя")
		return
	}

	// ####### - Обновляем контент отписки - #######

	// 1. Находим событие по по hashId, в следствии которого, клиент получал письмо
	mtaHistoryHashId, ok := utilsCr.GetQuerySTRVarFromGET(r, "hi")
	if ok {
		mtaHistory, err := account.GetMTAHistoryByHashId(mtaHistoryHashId)
		if err == nil {
			// если hashId Отправки совпадает, то отписываем. Это как проверочный код.
			if mtaHistory.HashId == mtaHistoryHashId {
				if err := mtaHistory.UpdateSetUnsubscribeUser(utilsCr.GetIP(r)); err != nil {
					log.Printf("Ошибка обновления истории отписки пользователя: %v\n", err)
				}
			}
		} else {
			log.Printf("Ошибка обновления истории отписки пользователя (2): %v\n", err)
		}
	}

	retHTML(w, "Пользователь успешно отписан от всех рассылок")
	return
}

func OpenEmailByPixelUser(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err == nil && account != nil {
		mtaHistoryHashId, ok := utilsCr.GetQuerySTRVarFromGET(r, "hi")
		if ok {

			mtaHistory, err := account.GetMTAHistoryByHashId(mtaHistoryHashId)
			if err == nil {

				if err := mtaHistory.UpdateOpenUser(utilsCr.GetIP(r)); err != nil {
					log.Printf("Ошибка обновления счетчика открытий: %v\n", err)
				}
			}
		}
	}

	// ####### - Обновляем контент открытий - #######

	// 1. Находим событие по по hashId, в следствии которого, клиент получал письмо

	const transPixel = "\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B"

	// Pixel
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	w.Header().Set("Expires", "Wed, 11 Nov 1998 11:11:11 GMT")
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	fmt.Fprintf(w, transPixel)
}

func retHTML(w http.ResponseWriter, message string) {
	body := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Bhojpur CRM - Пользователь отписан</title>
</head>
<body style="background-color: #F4F4F4;" leftmargin="0" marginwidth="0" topmargin="0" marginheight="0" offset="0">
	<div style="padding: 5px 15px;"><h5 style="font-size: 18px;color: #4a4949;">%v</h5></div>
</body>
</html>`, message)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(body))

	// fmt.Fprint(w, body)
	return
}
