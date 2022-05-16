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

func EmailTemplateCreate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	// Get JSON-request
	var input struct {
		models.EmailTemplate
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	// template, err := account.CreateEmailTemplate(models.EmailTemplate{Name: v.Name, Code: string(v.Code)})
	emailTemplate, err := account.CreateEntity(&input.EmailTemplate)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при создании шаблона"))
		return
	}

	resp := u.Message(true, "POST Email Templates Created")
	resp["email_template"] = emailTemplate
	u.Respond(w, resp)
}

/* Возвращает список шаблонов для текущего аккаунта */
func EmailTemplateGet(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailTemplateId, err := utilsCr.GetUINTVarFromRequest(r, "emailTemplateId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailTemplate models.EmailTemplate

	// 2. Узнаем, какой id учитывается нужен
	publicOk := utilsCr.GetQueryBoolVarFromGET(r, "public_id")

	if publicOk {
		err = account.LoadEntityByPublicId(&emailTemplate, emailTemplateId, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить объект"))
			return
		}
	} else {
		if err = account.LoadEntity(&emailTemplate, emailTemplateId, preloads); err != nil {
			u.Respond(w, u.MessageError(err, "Шаблон не найден"))
			return
		}

	}

	resp := u.Message(true, "GET Email template")
	resp["email_template"] = emailTemplate
	u.Respond(w, resp)
}

func EmailTemplateGetListPagination(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil {
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
	emailTemplates := make([]models.Entity, 0)

	if all {
		emailTemplates, total, err = account.GetListEntity(&models.EmailTemplate{}, sortBy, preloads)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список email-шаблонов"))
			return
		}
	} else {
		// webHooks, total, err = account.GetWebHooksPaginationList(offset, limit, search)
		emailTemplates, total, err = account.GetPaginationListEntity(&models.EmailTemplate{}, offset, limit, sortBy, search, nil, nil)
		if err != nil {
			u.Respond(w, u.MessageError(err, "Не удалось получить список email-шаблонов"))
			return
		}
	}

	resp := u.Message(true, "GET Email Template Pagination List")
	resp["total"] = total
	resp["email_templates"] = emailTemplates
	u.Respond(w, resp)
}

func EmailTemplateUpdate(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailTemplateId, err := utilsCr.GetUINTVarFromRequest(r, "emailTemplateId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	preloads := utilsCr.GetQueryStringArrayFromGET(r, "preloads")

	var emailTemplate models.EmailTemplate
	if err := account.LoadEntity(&emailTemplate, emailTemplateId, preloads); err != nil {
		u.Respond(w, u.MessageError(err, "Шаблон не найден"))
		return
	}

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	// err = account.EmailTemplateUpdate(tpl, input)
	err = account.UpdateEntity(&emailTemplate, input, preloads)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка при обновлении шаблона"))
		return
	}

	resp := u.Message(true, "Email Template Updated")
	resp["email_template"] = emailTemplate
	u.Respond(w, resp)
}

func EmailTemplateDelete(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w, r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message: "Ошибка авторизации"}))
		return
	}

	emailTemplateId, err := utilsCr.GetUINTVarFromRequest(r, "emailTemplateId")
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id шаблона"))
		return
	}

	var emailTemplate models.EmailTemplate
	if err := account.LoadEntity(&emailTemplate, emailTemplateId, nil); err != nil {
		u.Respond(w, u.MessageError(err, "Шаблон не найден"))
		return
	}

	err = account.DeleteEntity(&emailTemplate)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка удаления шаблона"))
		return
	}

	resp := u.Message(true, "Email templates delete")
	u.Respond(w, resp)
}

// -- TEST --
/*func EmailTemplateSendToUser(w http.ResponseWriter, r *http.Request) {

	account, err := utilsCr.GetWorkAccount(w,r)
	if err != nil || account == nil {
		u.Respond(w, u.MessageError(u.Error{Message:"Ошибка авторизации"}))
		return
	}

	hashId, ok := utilsCr.GetSTRVarFromRequest(r, "emailTemplateHashId")
	if !ok {
		u.Respond(w, u.MessageError(nil, "Ошибка в обработке Id шаблона"))
		return
	}

	template, err := account.EmailTemplateGetByHashId(hashId)
	if err != nil || template == nil {
		u.Respond(w, u.MessageError(err, "Шаблон не найден"))
		return
	}

	// 2. Get JSON-request
	input := &struct {
		UserId uint `json:"user_id"`
		EmailBoxId uint `json:"email_box_id"` // emailBoxId
		Subject string 	`json:"subject"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.Respond(w, u.MessageError(err, "Техническая ошибка в запросе"))
		return
	}

	user, err := account.GetUser(input.UserId)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id пользователя"))
		return
	}

	var ebox models.EmailBox
	err = account.LoadEntity(&ebox, input.EmailBoxId)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в обработке Id Email box"))
		return
	}

	// err = template.Send(*ebox,*user,input.Subject)
	err = template.SendChannel(ebox,user.Email,input.Subject, nil)
	if err != nil {
		u.Respond(w, u.MessageError(err, "Ошибка в отправке письма"))
		return
	}

	resp := u.Message(true, "GET Email Template send to user")
	u.Respond(w, resp)
}*/

// ### --- Public function --- ###
func EmailTemplatePreviewGetHTML(w http.ResponseWriter, r *http.Request) {

	hashId, ok := utilsCr.GetSTRVarFromRequest(r, "emailTemplateHashId")
	if !ok {
		u.Respond(w, u.MessageError(nil, "Ошибка в обработке Id шаблона"))
		return
	}

	template, err := (models.Account{}).EmailTemplateGetSharedByHashId(hashId)
	if err != nil || template == nil {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.Write([]byte(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Bhojpur CRM - Шаблон не может быть отображен</title></head><body><h4 style="color:#5f5f5f;">Данный шаблон не может быть отображен.</h4></body></html>`))
		return
	}

	// Подготавливаем данные для шаблона
	// vData, err := template.PrepareViewData(tempUser())
	account, err := models.GetAccount(template.AccountId)
	if err != nil {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.Write([]byte(`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Bhojpur CRM - Аккаунт не найден</title></head><body><h4 style="color:#5f5f5f;">Данный шаблон не может быть отображен.</h4></body></html>`))
		return
	}
	user := tempUser()
	vData, err := template.PrepareViewData(template.Name, "PreviewText", account.GetTemplateData(&user), "#", u.STRp("#"))
	if err != nil {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.Write(errorHTMLPage("Ошибка подготовки данных для отображения HTML"))
		return
	}

	html, err := template.GetHTML(vData)
	if err != nil {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.Write(errorHTMLPage("Ошибка получения HTML из шаблона"))
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Write([]byte(html))
}

func EmailTemplatePreviewGetRawHTML(w http.ResponseWriter, r *http.Request) {

	hashId, ok := utilsCr.GetSTRVarFromRequest(r, "emailTemplateHashId")
	if !ok {
		u.Respond(w, u.MessageError(nil, "Ошибка в обработке Id шаблона"))
		return
	}

	template, err := (models.Account{}).EmailTemplateGetSharedByHashId(hashId)
	if err != nil || template == nil {
		w.Header().Set("Content-Type", "text/html;charset=UTF-8")
		w.Write(errorHTMLPage("Данный шаблон не может быть отображен"))
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Write([]byte(template.HTMLData))
}

func errorHTMLPage(errorText string) []byte {
	return []byte(
		fmt.Sprintf(
			`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<title>Bhojpur CRM - Шаблон не может быть отображен</title>
</head>
<body>
	<h4 style="color:#5f5f5f;">%s</h4>
</body>
</html>`,
			errorText))
}

func tempUser() models.User {
	return models.User{
		Username:    u.STRp("serName"),
		Name:        u.STRp("Иван"),
		Surname:     u.STRp("Иванов"),
		Email:       u.STRp("info@bhojpur.net"),
		PhoneRegion: u.STRp("IN"),
		Phone:       u.STRp("+918788682735"),
		Password:    u.STRp("kjdfhkdfsr439rrfh39f34"),
	}
}
