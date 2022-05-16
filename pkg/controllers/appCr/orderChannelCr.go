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

func OrderChannelCreate(w http.ResponseWriter, r *http.Request) {

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
		models.OrderChannel
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	orderChannel, err := account.CreateEntity(&input.OrderChannel)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания канала"}))
		return
	}

	resp := u.Message(true, "POST OrderChannel Created")
	resp["order_channel"] = orderChannel
	u.Respond(w, resp)
}

func OrderChannelGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	orderChannelId, err := utilsCr.GetUINTVarFromRequest(r, "orderChannelId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var orderChannel models.OrderChannel
	err = account.LoadEntity(&orderChannel, orderChannelId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "GET OrderChannel")
	resp["order_channel"] = orderChannel
	u.Respond(w, resp)
}

func OrderChannelGetListPagination(w http.ResponseWriter, r *http.Request) {

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
	orderChannels := make([]models.Entity, 0)

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	orderChannels, total, err = account.GetPaginationListEntity(&models.OrderChannel{}, offset, limit, sortBy, search, nil, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "GET OrderChannel Pagination List")
	resp["total"] = total
	resp["order_channels"] = orderChannels
	u.Respond(w, resp)
}

func OrderChannelUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	if !account.IsMainAccount() {
		u.Respond(w, u.MessageError(u.Error{Message: "У вас нет прав на создание/изменение объектов этого типа"}))
		return
	}

	orderChannelId, err := utilsCr.GetUINTVarFromRequest(r, "orderChannelId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var orderChannel models.OrderChannel
	err = account.LoadEntity(&orderChannel, orderChannelId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = account.UpdateEntity(&orderChannel, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH OrderChannel Update")
	resp["order_channel"] = orderChannel
	u.Respond(w, resp)
}

func OrderChannelDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	if !account.IsMainAccount() {
		u.Respond(w, u.MessageError(u.Error{Message: "У вас нет прав на создание/изменение объектов этого типа"}))
		return
	}

	orderChannelId, err := utilsCr.GetUINTVarFromRequest(r, "orderChannelId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var orderChannel models.OrderChannel
	err = account.LoadEntity(&orderChannel, orderChannelId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}
	if err = account.DeleteEntity(&orderChannel); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении канала заказов"))
		return
	}

	resp := u.Message(true, "DELETE OrderChannel Successful")
	u.Respond(w, resp)
}
