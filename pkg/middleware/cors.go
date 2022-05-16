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
	"net/http"
)

func CorsAccessControl(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Проверка источника
		// allowAuth := []string{"http://crm.bhojpur.net","http://localhost:8090","http://crm.bhojpur.net","https://crm.bhojpur.net"} //List of endpoints that doesn't require auth

		// 2. Frontend request host
		// requestHost := r.Header.Get("Origin") //current request path

		// w.Header().Add("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Max-Age", "86400") // max 600
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT,PATCH,DELETE")
		// w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type,User-Agent")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Accept-Language,Cache-Control,Content-Type,Content-Length,Accept,Origin,X-Requested-With,Access-Control-Request-Headers,Access-Control-Request-Method,Access-Control-Allow-Credentials,Host, Origin, User-Agent, Referer")

		// 3. Выставляем cors-политику, если источник в разрешенных хостах
		/*for _, value := range allowAuth {

			if value == requestHost {
				w.Header().Add("Content-Type", "application/json;charset=UTF-8")
				//w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT,PATCH")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization,Accept-Language,Cache-Control,Content-Type,Content-Length,Accept,Origin,X-Requested-With,Access-Control-Request-Headers,Access-Control-Request-Method,Access-Control-Allow-Credentials,Host, Origin, User-Agent, Referer")
				w.Header().Set("Access-Control-Max-Age", "600") // max for FireFox, Chrome max 600
				break
			}
		}*/

		// 3. Не передаем запрос дальше, если мето OPTIONS
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func CorsAPIAccessControl(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Max-Age", "86400") // max 600
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT,PATCH,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Accept-Language,Cache-Control,Content-Type,Content-Length,Accept,Origin,X-Requested-With,Access-Control-Request-Headers,Access-Control-Request-Method,Access-Control-Allow-Credentials,Host, Origin, User-Agent, Referer")

		// 3. Не передаем запрос дальше, если метод OPTIONS
		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
