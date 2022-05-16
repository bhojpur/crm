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

// Ui / API there!
func EmailQueueCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.EmailQueue
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	emailQueue, err := account.CreateEntity(&input.EmailQueue)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания серии"}))
		return
	}

	resp := u.Message(true, "POST EmailQueue Created")
	resp["email_queue"] = emailQueue
	u.Respond(w, resp)
}

func EmailQueueGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	// ThisIs PublicID or inside
	emailQueueId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	var emailQueue models.EmailQueue

	// 2. Узнаем, какой id учитывается нужен
	publicOk := utilsCr.GetQueryBoolVarFromGET(r, "public_id")

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	if publicOk {
		err = account.LoadEntityByPublicId(&emailQueue, emailQueueId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
			return
		}
	} else {
		err = account.LoadEntity(&emailQueue, emailQueueId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
			return
		}
	}

	resp := u.Message(true, "GET EmailQueue")
	resp["email_queue"] = emailQueue
	u.Respond(w, resp)
}

func EmailQueueGetListPagination(w http.ResponseWriter, r *http.Request) {

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

	// возвращаемые переменные
	var total int64 = 0
	emailQueues := make([]models.Entity, 0)

	// 2. Узнаем, какой список нужен
	all := utilsCr.GetQueryBoolVarFromGET(r, "all")

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	if all {
		emailQueues, total, err = account.GetListEntity(&models.EmailQueue{}, sortBy, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список"))
			return
		}
	} else {
		emailQueues, total, err = account.GetPaginationListEntity(&models.EmailQueue{}, offset, limit, sortBy, search, nil, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список"))
			return
		}
	}

	resp := u.Message(true, "GET EmailQueue Pagination List")
	resp["total"] = total
	resp["email_queues"] = emailQueues
	u.Respond(w, resp)
}

func EmailQueueUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailQueueId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailQueue models.EmailQueue
	err = account.LoadEntity(&emailQueue, emailQueueId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = account.UpdateEntity(&emailQueue, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH EmailQueue Update")
	resp["email_queue"] = emailQueue
	u.Respond(w, resp)
}

func EmailQueueDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailQueueId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var emailQueue models.EmailQueue
	err = account.LoadEntity(&emailQueue, emailQueueId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}
	if err = account.DeleteEntity(&emailQueue); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении серии писем"))
		return
	}

	resp := u.Message(true, "DELETE EmailQueue Successful")
	u.Respond(w, resp)
}

func EmailQueueChangeStatus(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailQueueId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailQueue models.EmailQueue
	err = account.LoadEntity(&emailQueue, emailQueueId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить серию писем"))
		return
	}

	var input struct {
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = emailQueue.ChangeWorkStatus(input.Status, input.Reason)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "GET Email Campaign Execute")
	resp["email_queue"] = emailQueue
	u.Respond(w, resp)
}

func EmailQueueAppendAll(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailQueueId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailQueue models.EmailQueue
	err = account.LoadEntity(&emailQueue, emailQueueId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить серию писем"))
		return
	}

	/*var input struct{
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}*/

	count, err := emailQueue.AppendAllUsers()
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	err = account.LoadEntity(&emailQueue, emailQueueId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить серию писем"))
		return
	}

	resp := u.Message(true, "PATH Email Queue Append All Users")
	resp["email_queue"] = emailQueue
	resp["count"] = count
	u.Respond(w, resp)
}
