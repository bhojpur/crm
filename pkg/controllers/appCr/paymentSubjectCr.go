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

func PaymentSubjectCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	if !account.IsMainAccount() {
		u.Respond(w, u.MessageError(u.Error{Message: "У вас нет прав на создание/изменение объектов этого типа"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.PaymentSubject
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	paymentSubject, err := account.CreateEntity(&input.PaymentSubject)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания"}))
		return
	}

	resp := u.Message(true, "POST PaymentSubject Created")
	resp["payment_subject"] = paymentSubject
	u.Respond(w, resp)
}

func PaymentSubjectGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	paymentSubjectId, err := utilsCr.GetUINTVarFromRequest(r, "paymentSubjectId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var paymentSubject models.PaymentSubject
	err = account.LoadEntity(&paymentSubject, paymentSubjectId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}

	resp := u.Message(true, "GET PaymentSubject")
	resp["payment_subject"] = paymentSubject
	u.Respond(w, resp)
}

func PaymentSubjectGetListPagination(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	limit, ok := utilsCr.GetQueryINTVarFromGET(r, "limit")
	if !ok {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}
	offset, ok := utilsCr.GetQueryINTVarFromGET(r, "offset")
	if !ok || offset < 0 {
		offset = 0
	}
	sortDesc := utilsCr.GetQueryBoolVarFromGET(r, "sortDesc") // обратный или нет порядок
	sortBy, ok := utilsCr.GetQuerySTRVarFromGET(r, "sortBy")
	if !ok {
		sortBy = "id"
	}
	if sortDesc {
		sortBy += " desc"
	}

	search, ok := utilsCr.GetQuerySTRVarFromGET(r, "search")
	if !ok {
		search = ""
	}

	var total int64 = 0
	paymentSubjects := make([]models.Entity, 0)

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	all := utilsCr.GetQueryBoolVarFromGET(r, "all")

	if all {
		paymentSubjects, total, err = account.GetListEntity(&models.PaymentSubject{}, sortBy, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список"))
			return
		}
	} else {
		paymentSubjects, total, err = account.GetPaginationListEntity(&models.PaymentSubject{}, offset, limit, sortBy, search, nil, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список"))
			return
		}
	}

	resp := u.Message(true, "GET PaymentSubject pagination list")
	resp["total"] = total
	resp["payment_subjects"] = paymentSubjects
	u.Respond(w, resp)
}

func PaymentSubjectUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	if !account.IsMainAccount() {
		u.Respond(w, u.MessageError(u.Error{Message: "У вас нет прав на создание/изменение объектов этого типа"}))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	paymentSubjectId, err := utilsCr.GetUINTVarFromRequest(r, "paymentSubjectId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var paymentSubject models.PaymentSubject
	err = account.LoadEntity(&paymentSubject, paymentSubjectId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = account.UpdateEntity(&paymentSubject, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH PaymentSubject Update")
	resp["payment_subject"] = paymentSubject
	u.Respond(w, resp)
}

func PaymentSubjectDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	if !account.IsMainAccount() {
		u.Respond(w, u.MessageError(u.Error{Message: "У вас нет прав на создание/изменение объектов этого типа"}))
		return
	}

	paymentSubjectId, err := utilsCr.GetUINTVarFromRequest(r, "paymentSubjectId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var paymentSubject models.PaymentSubject
	err = account.LoadEntity(&paymentSubject, paymentSubjectId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}
	if err = account.DeleteEntity(&paymentSubject); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении типа товара"))
		return
	}

	resp := u.Message(true, "DELETE PaymentSubject Successful")
	u.Respond(w, resp)
}
