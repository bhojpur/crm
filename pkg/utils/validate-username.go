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
	"errors"
	"regexp"
)

// проверяет имя пользователя на соответствие правилам. Не проверяет уникальность
func VerifyUsername(username *string) error {

	if username == nil {
		return nil
	}

	if len(*username) == 0 {
		return errors.New("Поле необходимо заполнить")
	}

	if len(*username) < 3 {
		return errors.New("Имя пользователя слишком короткое")
	}

	if len(*username) > 25 || len([]rune(*username)) > 25 {
		return errors.New("Имя пользователя слишком длинное")
	}

	var rxUsername = regexp.MustCompile("^[a-zA-Z0-9,\\-_]+$")

	if !rxUsername.MatchString(*username) {
		return errors.New("Используйте только a-z,A-Z,0-9 а также символ -")
	}

	return nil
}
