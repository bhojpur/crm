package routes

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
	"net/http"

	"github.com/bhojpur/crm/pkg/controllers/appCr"
	"github.com/bhojpur/crm/pkg/controllers/trackingCr"
	"github.com/gorilla/mux"
)

/**
* [TRACKING] - группа роутов для отслеживания внешних данных
* MiddleWare - Cors & AccountHashId
 */
var TrackingRoutes = func(rFree *mux.Router) {

	// загружаем базовые настройки системы
	rFree.HandleFunc("/", appCr.CheckApi).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	// ?u={userHashId} а аккаунт через accountHashId передан.
	rFree.HandleFunc("/e/unsubscribe", trackingCr.UnsubscribeUser).Methods(http.MethodGet, http.MethodOptions)
	rFree.HandleFunc("/e/open", trackingCr.OpenEmailByPixelUser).Methods(http.MethodGet, http.MethodOptions)

	// rFree.HandleFunc("/e/subscribe", trackingCr.OpenEmailByPixelUser).Methods(http.MethodGet, http.MethodOptions)

	// pixel ?u={userHashId} а аккаунт через accountHashId передан.
	// r.HandleFunc("/e/open", trackingCr.UnsubscribeUser).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
}