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

func EmailBoxCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.EmailBox
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	emailBox, err := account.CreateEntity(&input.EmailBox)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания"}))
		return
	}

	resp := u.Message(true, "POST EmailBox Created")
	resp["email_box"] = emailBox
	u.Respond(w, resp)
}

func EmailBoxGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	emailBoxId, err := utilsCr.GetUINTVarFromRequest(r, "emailBoxId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке emailBox Id"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailBox models.EmailBox

	// 2. Узнаем, какой id учитывается нужен
	publicOk := utilsCr.GetQueryBoolVarFromGET(r, "public_id")

	if publicOk {
		err = account.LoadEntityByPublicId(&emailBox, emailBoxId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
			return
		}
	} else {
		err = account.LoadEntity(&emailBox, emailBoxId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Failed to load магазин"))
			return
		}
	}

	resp := u.Message(true, "GET EmailBox ")
	resp["email_box"] = emailBox
	u.Respond(w, resp)
}

func EmailBoxListPaginationGet(w http.ResponseWriter, r *http.Request) {

	var account *models.Account
	var err error
	// 1. Получаем рабочий аккаунт в зависимости от источника (автома. сверка с {hashId}.)

	account, err = utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	limit, ok := utilsCr.GetQueryINTVarFromGET(r, "limit")
	if !ok {
		limit = 25
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
	// 2. Узнаем, какой список нужен
	all := utilsCr.GetQueryBoolVarFromGET(r, "all")

	// Узнаем нужен ли фильтр
	filter := map[string]interface{}{}
	webSiteId, _filterWebSite := utilsCr.GetQueryUINTVarFromGET(r, "webSiteId")
	if _filterWebSite {
		filter["web_site_id"] = webSiteId
	}

	var total int64 = 0
	emailBoxes := make([]models.Entity, 0)

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	if all && len(filter) < 1 {
		emailBoxes, total, err = account.GetListEntity(&models.EmailBox{}, sortBy, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список почтовых ящиков"))
			return
		}
	} else {
		emailBoxes, total, err = account.GetPaginationListEntity(&models.EmailBox{}, offset, limit, sortBy, search, filter, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список почтовых ящиков"))
			return
		}
	}

	resp := u.Message(true, "GET Web Sites PaginationList")
	resp["email_boxes"] = emailBoxes
	resp["total"] = total
	u.Respond(w, resp)
}

func EmailBoxUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailBoxId, err := utilsCr.GetUINTVarFromRequest(r, "emailBoxId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailBox models.EmailBox
	err = account.LoadEntity(&emailBox, emailBoxId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}

	var input map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	// emailBox, err := account.UpdateWebSite(emailBoxId, &input.EmailBox)
	err = account.UpdateEntity(&emailBox, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH EmailBox Update")
	resp["email_box"] = emailBox
	u.Respond(w, resp)
}

func EmailBoxDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailBoxId, err := utilsCr.GetUINTVarFromRequest(r, "emailBoxId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var emailBox models.EmailBox
	err = account.LoadEntity(&emailBox, emailBoxId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить магазин"))
		return
	}

	if err = account.DeleteEntity(&emailBox); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении email box"))
		return
	}

	resp := u.Message(true, "DELETE EmailBox Successful")
	u.Respond(w, resp)
}
