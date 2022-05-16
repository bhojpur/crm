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
	"github.com/gorilla/mux"
)

// subdomain: cdn.<bhojpur.net>
var CDNRoutes = func(r *mux.Router) {

	// тут надо бы сделать методичку одного окна, чтобы можно было посмотреть и raw и html
	// http://public.crm.local/email/templates/share/4fgjy6lk1kxp
	r.HandleFunc("/emails/preview/raw/{emailTemplateHashId}", appCr.EmailTemplatePreviewGetRawHTML).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/emails/preview/compile/{emailTemplateHashId}", appCr.EmailTemplatePreviewGetHTML).Methods(http.MethodGet, http.MethodOptions)
	// r.HandleFunc("/emails/preview/html/{emailTemplateHashId}", controllers.EmailTemplatePreviewHTMLGet).Methods(http.MethodGet, http.MethodOptions)
	// r.HandleFunc("/emails/preview/raw/{emailTemplateHashId}", controllers.EmailTemplatePreviewRawGet).Methods(http.MethodGet, http.MethodOptions)

	// r.HandleFunc("/storage", controllers.StorageGetList).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/public/{hashId}", appCr.StorageCDNGet).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/products/images/{hashId}", appCr.StorageCDNGet).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/emails/images/{hashId}", appCr.StorageCDNGet).Methods(http.MethodGet, http.MethodOptions)

	r.HandleFunc("/articles/preview/raw/{hashId}", appCr.ArticleRawPreviewCDNGet).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/articles/preview/compile/{hashId}", appCr.ArticleCompilePreviewCDNGet).Methods(http.MethodGet, http.MethodOptions)

}
