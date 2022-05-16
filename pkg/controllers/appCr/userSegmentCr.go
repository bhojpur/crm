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

func UsersSegmentCreate(w http.ResponseWriter, r *http.Request) {

	var account *models.Account
	var err error
	account, err = utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	var input struct {
		models.UsersSegment
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	entity, err := account.CreateEntity(&input.UsersSegment)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания UsersSegment"}))
		return
	}

	resp := u.Message(true, "POST UsersSegment Create")
	resp["users_segment"] = entity
	u.Respond(w, resp)
}

func UsersSegmentGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	usersSegmentId, err := utilsCr.GetUINTVarFromRequest(r, "usersSegmentId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id usersSegment"))
		return
	}

	var usersSegment models.UsersSegment
	// 2. Узнаем, какой id учитывается нужен
	publicOk := utilsCr.GetQueryBoolVarFromGET(r, "public_id")

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	if publicOk {
		err = account.LoadEntityByPublicId(&usersSegment, usersSegmentId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
			return
		}
	} else {
		err = account.LoadEntity(&usersSegment, usersSegmentId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось найти магазин"))
			return
		}

	}

	resp := u.Message(true, "GET Web Hook")
	resp["users_segment"] = usersSegment
	u.Respond(w, resp)
}

func UsersSegmentPaginationGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
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

	var total int64 = 0
	usersSegments := make([]models.Entity, 0)

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	if all {
		usersSegments, total, err = account.GetListEntity(&models.UsersSegment{}, sortBy, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список пользовательских сегментов"))
			return
		}
	} else {
		// usersSegments, total, err = account.GetUsersSegmentsPaginationList(offset, limit, search)
		usersSegments, total, err = account.GetPaginationListEntity(&models.UsersSegment{}, offset, limit, sortBy, search, nil, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список"))
			return
		}
	}

	resp := u.Message(true, "GET UsersSegments PaginationList")
	resp["users_segments"] = usersSegments
	resp["total"] = total
	u.Respond(w, resp)
}

func UsersSegmentUpdate(w http.ResponseWriter, r *http.Request) {

	var account *models.Account
	var err error
	account, err = utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	usersSegmentId, err := utilsCr.GetUINTVarFromRequest(r, "usersSegmentId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id группы"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	// var input interface{}
	var input map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	var usersSegment models.UsersSegment
	if err = account.LoadEntity(&usersSegment, usersSegmentId, preloads); err != nil {
		u.Respond(w, u.MessageError(err, "WEbHook не найден"))
		return
	}

	// usersSegment, err := account.UpdateUsersSegment(usersSegmentId, input)
	if err = account.UpdateEntity(&usersSegment, input, preloads); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH UsersSegment Update")
	resp["users_segment"] = usersSegment
	u.Respond(w, resp)
}

func UsersSegmentDelete(w http.ResponseWriter, r *http.Request) {

	var account *models.Account
	var err error
	account, err = utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	usersSegmentId, err := utilsCr.GetUINTVarFromRequest(r, "usersSegmentId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id группы"))
		return
	}

	var usersSegment models.UsersSegment
	if err = account.LoadEntity(&usersSegment, usersSegmentId, nil); err != nil {
		u.Respond(w, u.MessageError(err, "WEbHook не найден"))
		return
	}

	// usersSegment, err := account.UpdateUsersSegment(usersSegmentId, input)
	if err = account.DeleteEntity(&usersSegment); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении"))
		return
	}

	resp := u.Message(true, "DELETE UsersSegment Successful")
	u.Respond(w, resp)
}
