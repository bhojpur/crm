package uiApiCr

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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bhojpur/crm/pkg/controllers/utilsCr"
	u "github.com/bhojpur/crm/pkg/utils"
)

// Уведомление от Я.Кассы
func PaymentYandexWebHook(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	// Проверяем, что такой способ платеж вообще есть, иначе реджект
	paymentYandexHashId, ok := utilsCr.GetSTRVarFromRequest(r, "paymentYandexHashId")
	if !ok {
		fmt.Println("WebHook: не найден paymentYandexHashId", err)
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка - не найден paymentYandexHashId"}))
		return
	}
	_, err = account.GetPaymentYandexByHashId(paymentYandexHashId)
	if err != nil {
		fmt.Println("WebHook: Не удалось найти способ оплаты", err)
		u.Respond(w, u.MessageError(err, "Не удалось найти способ оплаты"))
		return
	}

	//////////////////////

	type ObjectPayment struct {
		Id     string `json:"id"`
		Status string `json:"status"`
		Paid   bool   `json:"paid"`
	}
	// Читаем вход
	var input struct {
		Event  string        `json:"event"`
		Object ObjectPayment `json:"object"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println("WebHook:  проверьте обязательные поля и типы переменных", err)
		u.Respond(w, u.MessageError(err, "Ошибка в запросе: проверьте обязательные поля и типы переменных"))
		return
	}

	// fmt.Println("Input from Yandex Payment: ", input)
	// fmt.Println("ObjectPayment Event: ", input.Event)
	// fmt.Println("ObjectPayment ID: ", input.Object.Id)
	// fmt.Println("ObjectPayment Status: ", input.Object.Status)
	// fmt.Println("ObjectPayment Paid: ", input.Object.Paid)

	payment, err := account.GetPaymentByExternalId(input.Object.Id)
	if err != nil {
		fmt.Println("WebHook: Не удалось найти платежку", err)
		u.Respond(w, u.MessageError(err, "Не удалось найти платежку"))
		return
	}

	var m map[string]interface{}

	if input.Object.Paid {
		m = map[string]interface{}{
			"status": input.Object.Status,
			"paid":   input.Object.Paid,
			"paidAt": time.Now().UTC(), // обновляем время платежа
		}
	} else {
		m = map[string]interface{}{
			"status": input.Object.Status,
			"paid":   input.Object.Paid,
		}
	}
	err = account.UpdateEntity(payment, m, nil)
	if err != nil {
		fmt.Println("Эх! Не удалось обновить платежку: ", err)
		u.Respond(w, u.MessageError(err, "Не удалось обновить платежку"))
		return
	}

	resp := u.Message(true, "Payment Option Yandex")
	u.Respond(w, resp)
}
