package appCr

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
	"net/http"

	"github.com/bhojpur/crm/pkg/controllers/utilsCr"
	"github.com/bhojpur/crm/pkg/models"
	u "github.com/bhojpur/crm/pkg/utils"
)

func PaymentMethodCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.PaymentMethod
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	// хз хз
	paymentMethod, err := account.CreateEntity(input.PaymentMethod)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания метода"}))
		return
	}

	resp := u.Message(true, "POST PaymentMethod Created")
	resp["payment_method"] = paymentMethod
	u.Respond(w, resp)
}

func PaymentMethodGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	// узнаем ID элемента
	paymentMethodId, err := utilsCr.GetUINTVarFromRequest(r, "paymentMethodId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	// 2. Узнаем, какой список нужен
	code, ok := utilsCr.GetQuerySTRVarFromGET(r, "code")
	if !ok {
		u.Respond(w, u.MessageError(err, "Необходимо указать тип метода"))
		return
	}

	paymentMethod, err := account.GetPaymentMethodByCode(code, paymentMethodId)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "GET PaymentMethod")
	resp["payment_method"] = paymentMethod
	u.Respond(w, resp)
}

func PaymentMethodGetList(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	// todo: get webSiteid: ?webSite=1,2,3

	paymentMethods, err := account.GetPaymentMethods()
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при получении списка методов оплаты"))
		return
	}

	resp := u.Message(true, "GET PaymentMethod List")
	resp["payment_methods"] = paymentMethods
	u.Respond(w, resp)
}

func PaymentMethodUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// узнаем ID элемента
	paymentMethodId, err := utilsCr.GetUINTVarFromRequest(r, "paymentMethodId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	// 2. Узнаем, какой список нужен
	code, ok := utilsCr.GetQuerySTRVarFromGET(r, "code")
	if !ok {
		u.Respond(w, u.MessageError(err, "Необходимо указать тип метода"))
		return
	}

	paymentMethod, err := account.GetPaymentMethodByCode(code, paymentMethodId)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = account.UpdateEntity(paymentMethod, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH PaymentMethod Update")
	resp["payment_method"] = paymentMethod
	u.Respond(w, resp)
}

func PaymentMethodDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// узнаем ID элемента
	paymentMethodId, err := utilsCr.GetUINTVarFromRequest(r, "paymentMethodId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	// 2. Узнаем, какой список нужен
	code, ok := utilsCr.GetQuerySTRVarFromGET(r, "code")
	if !ok {
		u.Respond(w, u.MessageError(err, "Необходимо указать тип метода"))
		return
	}

	paymentMethod, err := account.GetPaymentMethodByCode(code, paymentMethodId)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}

	if err = account.DeleteEntity(paymentMethod); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении"))
		return
	}

	resp := u.Message(true, "DELETE PaymentMethod Successful")
	u.Respond(w, resp)
}
