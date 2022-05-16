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

func DeliveryOrderCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	var input struct {
		models.DeliveryOrder
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	deliveryOrder, err := account.CreateEntity(&input.DeliveryOrder)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания заявки"}))
		return
	}

	resp := u.Message(true, "POST DeliveryOrder Created")
	resp["delivery_order"] = deliveryOrder
	u.Respond(w, resp)
}

func DeliveryOrderGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	deliveryOrderId, err := utilsCr.GetUINTVarFromRequest(r, "deliveryOrderId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	publicOk := utilsCr.GetQueryBoolVarFromGET(r, "public_id")

	var deliveryOrder models.DeliveryOrder

	if publicOk {
		// err = account.LoadEntityByPublicId(&emailNotification, emailNotificationId, preloads)
		err = account.LoadEntityByPublicId(&deliveryOrder, deliveryOrderId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список"))
			return
		}
	} else {
		err = account.LoadEntity(&deliveryOrder, deliveryOrderId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список"))
			return
		}
	}

	resp := u.Message(true, "GET DeliveryOrder")
	resp["delivery_order"] = deliveryOrder
	u.Respond(w, resp)
}

func DeliveryOrderGetListPagination(w http.ResponseWriter, r *http.Request) {

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

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var total int64 = 0
	deliveryOrders := make([]models.Entity, 0)

	deliveryOrders, total, err = account.GetPaginationListEntity(&models.DeliveryOrder{}, offset, limit, sortBy, search, nil, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "GET DeliveryOrder Pagination List")
	resp["total"] = total
	resp["delivery_orders"] = deliveryOrders
	u.Respond(w, resp)
}

func DeliveryOrderUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	deliveryOrderId, err := utilsCr.GetUINTVarFromRequest(r, "deliveryOrderId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var deliveryOrder models.DeliveryOrder
	err = account.LoadEntity(&deliveryOrder, deliveryOrderId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = account.UpdateEntity(&deliveryOrder, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH DeliveryOrder Update")
	resp["delivery_order"] = deliveryOrder
	u.Respond(w, resp)
}

func DeliveryOrderDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	deliveryOrderId, err := utilsCr.GetUINTVarFromRequest(r, "deliveryOrderId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var deliveryOrder models.DeliveryOrder
	err = account.LoadEntity(&deliveryOrder, deliveryOrderId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}
	if err = account.DeleteEntity(&deliveryOrder); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении заказа на доставку"))
		return
	}

	resp := u.Message(true, "DELETE DeliveryOrder Successful")
	u.Respond(w, resp)
}
