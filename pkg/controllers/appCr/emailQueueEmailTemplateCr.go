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
func EmailQueueEmailTemplateCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.EmailQueueEmailTemplate
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	emailQueueEmailTemplate, err := account.CreateEntity(&input.EmailQueueEmailTemplate)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания серии"}))
		return
	}

	resp := u.Message(true, "POST EmailQueueEmailTemplate Created")
	resp["email_queue_email_template"] = emailQueueEmailTemplate
	u.Respond(w, resp)
}

func EmailQueueEmailTemplateGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	// ThisIs PublicID or inside
	emailQueueEmailTemplateId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueEmailTemplateId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailQueueEmailTemplate models.EmailQueueEmailTemplate

	err = account.LoadEntity(&emailQueueEmailTemplate, emailQueueEmailTemplateId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список заказов"))
		return
	}

	resp := u.Message(true, "GET EmailQueueEmailTemplate")
	resp["email_queue_email_template"] = emailQueueEmailTemplate
	u.Respond(w, resp)
}

func EmailQueueEmailTemplateGetListPagination(w http.ResponseWriter, r *http.Request) {

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

	emailQueueId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id emailQueueId"))
		return
	}

	filter := make(map[string]interface{}, 0)
	filter["email_queue_id"] = emailQueueId

	var total int64 = 0
	emailQueueEmailTemplates := make([]models.Entity, 0)

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	emailQueueEmailTemplates, total, err = account.GetPaginationListEntity(&models.EmailQueueEmailTemplate{}, offset, limit, sortBy, search, filter, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "GET EmailQueueEmailTemplate Pagination List")
	resp["total"] = total
	resp["email_queue_email_templates"] = emailQueueEmailTemplates
	u.Respond(w, resp)
}

func EmailQueueEmailTemplateUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailQueueEmailTemplateId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueEmailTemplateId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailQueueEmailTemplate models.EmailQueueEmailTemplate
	err = account.LoadEntity(&emailQueueEmailTemplate, emailQueueEmailTemplateId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = account.UpdateEntity(&emailQueueEmailTemplate, input, preloads)
	if err != nil {
		// fmt.Println(err)
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH EmailQueueEmailTemplate Update")
	resp["email_queue_email_template"] = emailQueueEmailTemplate
	u.Respond(w, resp)
}

func EmailQueueEmailTemplateMassUpdates(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailQueueId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id emailQueueId"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	emailQueue := models.EmailQueue{}
	if err = account.LoadEntity(&emailQueue, emailQueueId, preloads); err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Не найдена серия писем"}))
		return
	}

	var input struct {
		EmailQueueEmailTemplates []models.MassUpdateEmailQueueTemplate `json:"email_queue_email_templates"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	if err = emailQueue.UpdateOrderEmailTemplates(input.EmailQueueEmailTemplates); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе 2"))
		return
	}

	filter := make(map[string]interface{}, 0)
	filter["email_queue_id"] = emailQueueId

	var total int64 = 0
	emailQueueEmailTemplates := make([]models.Entity, 0)

	emailQueueEmailTemplates, total, err = account.GetPaginationListEntity(&models.EmailQueueEmailTemplate{}, 0, 100, "step", "", filter, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "PATCH EmailQueueEmailTemplate MassUpdates")
	resp["email_queue_email_templates"] = emailQueueEmailTemplates
	resp["total"] = total
	u.Respond(w, resp)
}

func EmailQueueEmailTemplateDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailQueueEmailTemplateId, err := utilsCr.GetUINTVarFromRequest(r, "emailQueueEmailTemplateId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var emailQueueEmailTemplate models.EmailQueueEmailTemplate
	err = account.LoadEntity(&emailQueueEmailTemplate, emailQueueEmailTemplateId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}
	if err = account.DeleteEntity(&emailQueueEmailTemplate); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении шаблона"))
		return
	}

	resp := u.Message(true, "DELETE EmailQueueEmailTemplate Successful")
	u.Respond(w, resp)
}
