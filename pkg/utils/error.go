package utils

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
	"fmt"
)

type Error struct {
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"message"`
}

// Пиздец важная функция
func (e Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Message, e.Errors)
}

func (error Error) HasErrors() (status bool) {

	status = false

	if len(error.Message) > 0 || len(error.Errors) > 0 {
		status = true
	}
	return
}

func (e *Error) AddErrors(key string, value interface{}) {
	if e.Errors == nil {
		//e.Errors = make(map[]interface{})
		e.Errors = make(map[string]interface{})
	}
	e.Errors[key] = value
}

func (e *Error) GetErrors() map[string]interface{} {
	return e.Errors
	/*if e.Errors == nil {
		// e.Errors = make(map[string]interface{})
		e.Errors = nil
	}
	return e.Errors*/
}

func (e *Error) GetError(key string) interface{} {
	if e.Errors == nil {
		e.Errors = make(map[string]interface{})
	}
	return e.Errors[key]
}
