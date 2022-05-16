package models

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
	"testing"
	"time"
)

func TestEmailAccessToken_Expired(t *testing.T) {

	// todo дописать список тестов
	var testList = []struct {
		eat         EmailAccessToken
		expected    bool
		description string
	}{
		{EmailAccessToken{CreatedAt: time.Now()}, false, "Т.к. время текущее, то токен точно не должен быть просрочен"},
		{EmailAccessToken{CreatedAt: time.Now().Truncate(time.Hour * 24 * 100)}, true, "Токен должен быть просрочен т.к. ему уже 100 дней"},
	}

	for _, v := range testList {

		if v.eat.Expired() != v.expected {
			t.Fatalf("Не верная работа функции для: %v \n Описание: %v", v, v.description)
		}
	}
}
