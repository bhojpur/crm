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

// CRUD functional <* love

func ApiKeyCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	input := struct {
		Name    string `json:"name"`
		Enabled bool   `json:"enabled"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	apiKey, err := account.ApiKeyCreate(models.ApiKey{Name: input.Name, Enabled: input.Enabled})
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка во время создания ключа"}))
		return
	}

	resp := u.Message(true, "POST Api Key Created")
	resp["api_key"] = *apiKey
	u.Respond(w, resp)
}

func ApiKeyGetList(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	apiKeys, err := account.ApiKeysList()
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка в обработке запроса"}))
		return
	}

	resp := u.Message(true, "GET list Api keys of account")
	resp["api_keys"] = apiKeys
	u.Respond(w, resp)
}

func ApiKeyGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	apiKeyId, err := utilsCr.GetUINTVarFromRequest(r, "apiKeyId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id"))
		return
	}

	apiKey, err := account.ApiKeyGet(apiKeyId)
	if err != nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка в обработке запроса", Errors: map[string]interface{}{"emailTemplates": "Не удалось получить ключ"}}))
		return
	}

	resp := u.Message(true, "GET ApiKey")
	resp["api_key"] = apiKey
	u.Respond(w, resp)
}

func ApiKeyUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	idApiKey, err := utilsCr.GetUINTVarFromRequest(r, "id")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	// Get JSON-request

	var input map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	apiKey, err := account.ApiKeyUpdate(idApiKey, input)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении"))
		return
	}

	resp := u.Message(true, "PATCH Api Key Update")
	resp["api_key"] = apiKey
	u.Respond(w, resp)
}

func ApiKeyDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	idApiKey, err := utilsCr.GetUINTVarFromRequest(r, "id")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	err = account.ApiKeyDelete(idApiKey)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при удалении Api-ключа"))
		return
	}

	resp := u.Message(true, "DELETE Api Key Update")
	u.Respond(w, resp)
}

// END CRUD functional
