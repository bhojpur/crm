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

func InventoryCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.Inventory
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	inventory, err := account.CreateEntity(&input.Inventory)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания"}))
		return
	}

	resp := u.Message(true, "POST Inventory Created")
	resp["inventory"] = inventory
	u.Respond(w, resp)
}

func InventoryGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	inventoryId, err := utilsCr.GetUINTVarFromRequest(r, "inventoryId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке inventory Id"))
		return
	}

	var inventory models.Inventory
	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	// 2. Узнаем, какой id учитывается нужен
	publicOk := utilsCr.GetQueryBoolVarFromGET(r, "public_id")

	if publicOk {
		err = account.LoadEntityByPublicId(&inventory, inventoryId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
			return
		}
	} else {
		err = account.LoadEntity(&inventory, inventoryId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Failed to load магазин"))
			return
		}
	}

	resp := u.Message(true, "GET Inventory")
	resp["inventory"] = inventory
	u.Respond(w, resp)
}

func InventoryListPaginationGet(w http.ResponseWriter, r *http.Request) {

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
	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	// 2. Узнаем, какой список нужен
	all := utilsCr.GetQueryBoolVarFromGET(r, "all")

	// Узнаем нужен ли фильтр
	filter := map[string]interface{}{}
	/*
		webSiteId, _filterWebSite := utilsCr.GetQueryUINTVarFromGET(r, "webSiteId")
		if _filterWebSite {
			filter["web_site_id"] = webSiteId
		}*/

	var total int64 = 0
	inventorys := make([]models.Entity, 0)

	if all && len(filter) < 1 {
		inventorys, total, err = account.GetListEntity(&models.Inventory{}, sortBy, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список страниц"))
			return
		}
	} else {
		// webHooks, total, err = account.GetWebHooksPaginationList(offset, limit, search)
		inventorys, total, err = account.GetPaginationListEntity(&models.Inventory{}, offset, limit, sortBy, search, filter, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список страниц"))
			return
		}
	}

	resp := u.Message(true, "GET product Cards PaginationList")
	resp["inventorys"] = inventorys
	resp["total"] = total
	u.Respond(w, resp)
}

func InventoryUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	inventoryId, err := utilsCr.GetUINTVarFromRequest(r, "inventoryId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}
	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")
	var inventory models.Inventory

	err = account.LoadEntity(&inventory, inventoryId, nil)

	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}

	var input map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = account.UpdateEntity(&inventory, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH Inventory Update")
	resp["inventory"] = inventory
	u.Respond(w, resp)
}

func InventoryDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	inventoryId, err := utilsCr.GetUINTVarFromRequest(r, "inventoryId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var inventory models.Inventory
	err = account.LoadEntity(&inventory, inventoryId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить магазин"))
		return
	}

	if err = account.DeleteEntity(&inventory); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении инвентаризации"))
		return
	}

	resp := u.Message(true, "DELETE Inventory Successful")
	u.Respond(w, resp)
}
func InventoryAppendProduct(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	inventoryId, err := utilsCr.GetUINTVarFromRequest(r, "inventoryId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id emailQueueId"))
		return
	}

	productId, err := utilsCr.GetUINTVarFromRequest(r, "productId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке productId"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var inventory models.Inventory
	if err = account.LoadEntity(&inventory, inventoryId, preloads); err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка в загрузке карточки товара"}))
		return
	}

	product := models.Product{}
	if err = account.LoadEntity(&product, productId, nil); err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка в загрузке товара"}))
		return
	}

	// Декодируем параметры
	var input struct {
		volumeFact float64 `json:"volume_fact"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	if err = inventory.AppendProduct(product, input.volumeFact); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка добавления продукта"))
		return
	}

	if err = account.LoadEntity(&inventory, inventoryId, preloads); err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка загрузки инвентаризации склада"}))
		return
	}

	resp := u.Message(true, "PATCH Inventory Append Product")
	resp["inventory"] = inventory
	u.Respond(w, resp)
}

func InventoryRemoveProduct(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	inventoryId, err := utilsCr.GetUINTVarFromRequest(r, "inventoryId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id emailQueueId"))
		return
	}

	productId, err := utilsCr.GetUINTVarFromRequest(r, "productId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке productId"))
		return
	}
	product := models.Product{}
	if err = account.LoadEntity(&product, productId, nil); err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка в загрузке товара"}))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var inventory models.Inventory
	if err = account.LoadEntity(&inventory, inventoryId, preloads); err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка в загрузке инвентаризации"}))
		return
	}

	if err = inventory.RemoveProduct(product); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка удаления продукта из карточки товара"))
		return
	}

	// Обновляем данные карточки
	if err = account.LoadEntity(&inventory, inventoryId, preloads); err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка в загрузке карточки товара"}))
		return
	}

	resp := u.Message(true, "PATCH Inventory Remove Products")
	resp["inventory"] = inventory
	u.Respond(w, resp)
}
