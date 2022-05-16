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
	"testing"
)

type TestStruct struct {
	Username string
	Email    string
	Phone    string
	Name     string
}

func TestAccount_CheckUserInputRequiredFields(t *testing.T) {

	TestList := []struct {
		Enum        []string
		Input       interface{}
		expected    bool
		description string
	}{
		{
			[]string{"username", "email", "phone"},
			&TestStruct{Username: "", Email: "", Phone: ""},
			false,
			"Требуемые поля - пустые",
		},
		{
			[]string{"email", "username"},
			&TestStruct{Username: "TestUser 1", Email: "mail@example.com"},
			true,
			"Формально поля есть",
		},
		{
			[]string{"name", "phone"},
			&TestStruct{Email: "raisb.brahmrishi@gmail.com", Name: "Никита"},
			false,
			"Нет поля с телефоном",
		},
		{
			[]string{"name", "phone", "username"},
			&TestStruct{Username: "", Phone: "+79251952295", Name: "Никита"},
			false,
			"Нет поля с телефоном",
		},
	}

	for i, v := range TestList {
		err := CheckNotNullFields(v.Input, v.Enum)

		// если прошел проверку
		if v.expected == true && err != nil {
			t.Fatalf("Проверка провалена, а должна была пройти:\nПользователь %v : \nОжидалось: %v \n user: %v \nТребуемые поля: %v", i, v.expected, v.Input, v.Enum)
		}

		if v.expected == false && err == nil {
			t.Fatalf("Проверка прошла успешно, но должна быть провалена:\nПользователь %v : \nОжидалось: %v \n user: %v \nТребуемые поля: %v", i, v.expected, v.Input, v.Enum)
		}

	}

}
