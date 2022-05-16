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
	"fmt"
	"net/http"

	"github.com/bhojpur/crm/pkg/controllers/utilsCr"
	"github.com/bhojpur/crm/pkg/models"
	u "github.com/bhojpur/crm/pkg/utils"
)

func WebSiteCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.WebSite
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	webSite, err := account.CreateEntity(&input.WebSite)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания"}))
		return
	}

	resp := u.Message(true, "POST WebSite Created")
	resp["web_site"] = webSite
	u.Respond(w, resp)
}

func WebSiteGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	webSiteId, err := utilsCr.GetUINTVarFromRequest(r, "webSiteId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке webSite Id"))
		return
	}

	var webSite models.WebSite

	// 2. Узнаем, какой id учитывается нужен
	publicOk := utilsCr.GetQueryBoolVarFromGET(r, "public_id")

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	if publicOk {
		err = account.LoadEntityByPublicId(&webSite, webSiteId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
			return
		}
	} else {
		err = account.LoadEntity(&webSite, webSiteId, preloads)
		if err != nil {
			fmt.Println(err)
			u.Respond(w, u.MessageError(err, "Failed to load магазин"))
			return
		}
	}

	resp := u.Message(true, "GET WebSite ")
	resp["web_site"] = webSite
	u.Respond(w, resp)
}

func WebSiteListPaginationGet(w http.ResponseWriter, r *http.Request) {

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

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var total int64 = 0
	webSites := make([]models.Entity, 0)

	if all {
		webSites, total, err = account.GetListEntity(&models.WebSite{}, sortBy, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список сайтов"))
			return
		}
	} else {
		// webHooks, total, err = account.GetWebHooksPaginationList(offset, limit, search)
		webSites, total, err = account.GetPaginationListEntity(&models.WebSite{}, offset, limit, sortBy, search, nil, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список сайтов"))
			return
		}
	}

	resp := u.Message(true, "GET Web Sites PaginationList")
	resp["web_sites"] = webSites
	resp["total"] = total
	u.Respond(w, resp)
}

func WebSiteUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	webSiteId, err := utilsCr.GetUINTVarFromRequest(r, "webSiteId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var webSite models.WebSite
	err = account.LoadEntity(&webSite, webSiteId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}

	var input map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	// webSite, err := account.UpdateWebSite(webSiteId, &input.WebSite)
	err = account.UpdateEntity(&webSite, input, preloads)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH WebSite Update")
	resp["web_site"] = webSite
	u.Respond(w, resp)
}

func WebSiteDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	webSiteId, err := utilsCr.GetUINTVarFromRequest(r, "webSiteId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var webSite models.WebSite
	err = account.LoadEntity(&webSite, webSiteId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить магазин"))
		return
	}

	if err = account.DeleteEntity(&webSite); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении WebSite"))
		return
	}

	resp := u.Message(true, "DELETE WebSite Successful")
	u.Respond(w, resp)
}
