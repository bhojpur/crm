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

	"github.com/bhojpur/crm/pkg/controllers/uiApiCr"
	"github.com/gorilla/mux"
)

/**
* [UI-API] - группа роутов для работы публичного UI/API через ui.api.crm.bhojpur.net.
* Авторизации по-умолчанию не требуется!!!
*
* В контексте issuerAccountId = accountId всегда, т.к. доступ к нескольким аккаунтам не предусматриваются.

* Оба роутера монтируются в точку /accounts/{accountId} имеют в контексте account & accountId
* rUiApi - маршруты без проверки JWT
* rUiApiAuth - маршрут с проверкой JWT, а также на совпадение {accountId} с accountId указаном в токене
 */

// ... /accountHashId}/...
var UiApiRoutesV1 = func(rFree *mux.Router) {

	// rFree.HandleFunc("/users", appCr.UserRegistration).Methods(http.MethodPost, http.MethodOptions)
	// rFree.HandleFunc("/users/auth/username", appCr.UserAuthByUsername).Methods(http.MethodPost, http.MethodOptions)
	// rFree.HandleFunc("/users/auth/email", appCr.UserAuthByEmail).Methods(http.MethodPost, http.MethodOptions)

	rFree.HandleFunc("/web-sites/{webSiteId:[0-9]+}/deliveries", uiApiCr.DeliveryGetListByShop).Methods(http.MethodGet, http.MethodOptions)

	rFree.HandleFunc("/web-sites/{webSiteId:[0-9]+}/deliveries-calculate", uiApiCr.DeliveryCalculateDeliveryCost).Methods(http.MethodPost, http.MethodOptions)

	rFree.HandleFunc("/web-sites/{webSiteId:[0-9]+}/deliveries-code-list", uiApiCr.DeliveryCodeList).Methods(http.MethodGet, http.MethodOptions)

	// YandexPayment
	// Адрес для вебхуков от Яндекс.Кассы. Код ответа 200 в случае обработки.
	// rFree.HandleFunc("/payments/yandex-payment/{yandexPayment:[0-9]+}/notifications/", uiApiCr.DeliveryListOptions).Methods(http.MethodGet, http.MethodOptions)

	// URL для яндекса: вставляется hashId магазина, а не id - чтобы защититься от атак.
	rFree.HandleFunc("/yandex-payment/{paymentYandexHashId}/web-hooks", uiApiCr.PaymentYandexWebHook).Methods(http.MethodPost, http.MethodOptions)

	rFree.HandleFunc("/orders", uiApiCr.UiApiOrderCreate).Methods(http.MethodPost, http.MethodOptions)

	// rFree.HandleFunc("/subscribe", uiApiCr.UiApiSubscribe).Methods(http.MethodPost, http.MethodOptions)
	// rFree.HandleFunc("/subscribe", uiApiCr.UiApiSubscribe).Methods(http.MethodPost, http.MethodOptions)

	// subscribe & question
	rFree.HandleFunc("/form", uiApiCr.UiApiForm).Methods(http.MethodPost, http.MethodOptions)

	// rFree.HandleFunc("/test", uiApiCr.Test).Methods(http.MethodGet, http.MethodOptions)
	// rFree.NotFoundHandler = NotFoundHandler()
	// rFree.MethodNotAllowedHandler = NotFoundHandler()

}
