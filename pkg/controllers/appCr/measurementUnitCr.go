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

func MeasurementUnitStatusCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.MeasurementUnit
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	measurementUnit, err := account.CreateEntity(&input.MeasurementUnit)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания статуса"}))
		return
	}

	resp := u.Message(true, "POST MeasurementUnit Created")
	resp["measurement_unit"] = measurementUnit
	u.Respond(w, resp)
}

func MeasurementUnitStatusGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	measurementUnitId, err := utilsCr.GetUINTVarFromRequest(r, "measurementUnitId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке web site Id"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var measurementUnit models.MeasurementUnit
	err = account.LoadEntity(&measurementUnit, measurementUnitId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "GET MeasurementUnit")
	resp["measurement_unit"] = measurementUnit
	u.Respond(w, resp)
}

func MeasurementUnitStatusGetList(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		return
	}

	var total int64 = 0
	measurementUnits := make([]models.Entity, 0)

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	measurementUnits, total, err = account.GetListEntity(&models.MeasurementUnit{}, "id", preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить список"))
		return
	}

	resp := u.Message(true, "GET MeasurementUnit List")
	resp["total"] = total
	resp["measurement_units"] = measurementUnits
	u.Respond(w, resp)
}

func MeasurementUnitStatusUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	measurementUnitId, err := utilsCr.GetUINTVarFromRequest(r, "measurementUnitId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var measurementUnit models.MeasurementUnit
	err = account.LoadEntity(&measurementUnit, measurementUnitId, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
		return
	}

	// Проверка прав
	if measurementUnit.AccountId != account.Id {
		u.Respond(w, u.MessageError(u.Error{Message: "У вас нет прав на создание/изменение объектов этого типа"}))
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	err = account.UpdateEntity(&measurementUnit, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH MeasurementUnit Update")
	resp["measurement_unit"] = measurementUnit
	u.Respond(w, resp)
}

func MeasurementUnitStatusDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	measurementUnitId, err := utilsCr.GetUINTVarFromRequest(r, "measurementUnitId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var measurementUnit models.MeasurementUnit
	err = account.LoadEntity(&measurementUnit, measurementUnitId, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Failed to load данные"))
		return
	}

	// Проверка на права изменения
	if measurementUnit.AccountId != account.Id {
		u.Respond(w, u.MessageError(u.Error{Message: "У вас нет прав на создание/изменение объектов этого типа"}))
		return
	}

	if err = account.DeleteEntity(&measurementUnit); err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении"))
		return
	}

	resp := u.Message(true, "DELETE MeasurementUnit Successful")
	u.Respond(w, resp)
}
