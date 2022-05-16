package middleware

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
	"context"
	"net/http"

	"github.com/bhojpur/crm/pkg/models"
	u "github.com/bhojpur/crm/pkg/utils"
	"github.com/gorilla/mux"
)

//Тут собраны посредники определяющие issuerAccountId и добавляющие его в контекст

// Add to Context(r) issuerAccountId (= 1, Bhojpur CRM)
func AddContextMainAccount(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		issuerAccount, err := models.GetMainAccount() // BhojpurCRM
		if err != nil {
			u.Respond(w, u.MessageError(nil, "An account with the specified hash Id was not found"))
			return
		}

		// For future
		r = r.WithContext(context.WithValue(r.Context(), "issuer", "app"))
		r = r.WithContext(context.WithValue(r.Context(), "issuerAccountId", issuerAccount.Id))
		r = r.WithContext(context.WithValue(r.Context(), "issuerAccount", issuerAccount))

		next.ServeHTTP(w, r)
	})

}

// Вставляет в контекст issuerAccountId из hashId (раскрытие issuer account)
func ContextMuxVarAccountHashId(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		accountHashId := mux.Vars(r)["accountHashId"] // защищаемся от парсинга / спама

		if len(accountHashId) != 12 {
			u.Respond(w, u.MessageError(nil, "The hashId length must be 12 symbols"))
			return
		}

		issuerAccount, err := models.GetAccountByHash(accountHashId)
		if err != nil {
			u.Respond(w, u.MessageError(nil, "An account with the specified hash Id was not found"))
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "auth", "ui/api"))
		r = r.WithContext(context.WithValue(r.Context(), "accountId", issuerAccount.Id))
		r = r.WithContext(context.WithValue(r.Context(), "account", issuerAccount))

		next.ServeHTTP(w, r)
	})

}
