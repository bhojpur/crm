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
	"net/http"

	u "github.com/bhojpur/crm/pkg/utils"
)

func CheckApi(w http.ResponseWriter, r *http.Request) {
	resp := u.Message(true, "Welcome to REST JSON API of Bhojpur CRM")
	resp["help"] = "You can read more about that: https://dev.crm.bhojpur.net/#api"
	u.Respond(w, resp)
	return
}

func CheckAppUiApi(w http.ResponseWriter, r *http.Request) {
	resp := u.Message(true, "This is UI-API of APP Bhojpur CRM")
	resp["help"] = "Most likely, you were looking for Public UI-API. Read more: https://dev.crm.bhojpur.net/#ui-api"
	u.Respond(w, resp)
	return
}
func CheckUiApi(w http.ResponseWriter, r *http.Request) {
	resp := u.Message(true, "Welcome to Public UI-API of Bhojpur CRM")
	resp["help"] = "You can read more about that: https://dev.crm.bhojpur.net/#ui-api"
	u.Respond(w, resp)
	return
}
