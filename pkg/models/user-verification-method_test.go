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

	"github.com/bhojpur/crm/pkg/utils"
)

func TestUserVerificationType_Create(t *testing.T) {
	listTest := []struct {
		uvt         UserVerificationMethod
		expected    bool
		description string
	}{
		{UserVerificationMethod{Name: "TestCreate Type", Tag: utils.RandStringBytesMaskImprSrcUnsafe(10, false)}, true, "Все есть"},
		{UserVerificationMethod{Name: "", Tag: utils.RandStringBytesMaskImprSrcUnsafe(9, false)}, false, "Нет имени"},
		{UserVerificationMethod{Name: "TestCreate Type", Tag: "Ф"}, false, "Слишком короткий код"},
	}

	for i, v := range listTest {
		uvt, err := v.uvt.Create()

		if v.expected == true && err != nil {
			t.Fatalf("Не удалось создать настройки верификации для варианта %v, ошибка: %v", i, err)
		}

		if v.expected == false && err == nil {
			t.Fatalf("Удалось создать настройки верификации, хотя они не должны были создаться для варианта %v", i)

		}

		if uvt != nil && err == nil {
			defer uvt.Delete()
		}

	}
}

func TestGetUserVerificationTypeById(t *testing.T) {
	uvt, err := UserVerificationMethod{Name: "TestDelete", Tag: utils.RandStringBytesMaskImprSrcUnsafe(5, false)}.Create()
	if err != nil {
		t.Fatalf("Cant create ver %v", err)
	}
	defer uvt.Delete()

	uvtF, err := GetUserVerificationTypeById(uvt.Id)
	if err != nil {
		t.Fatalf("Cant find ver type by id %v", err)
	}

	if uvt.Tag != uvtF.Tag {
		t.Fatalf("А коды то разные мужик!")
	}
}

func TestGetUserVerificationTypeByCode(t *testing.T) {
	uvt, err := UserVerificationMethod{Name: "TestDelete", Tag: utils.RandStringBytesMaskImprSrcUnsafe(5, false)}.Create()
	if err != nil {
		t.Fatalf("Cant create ver %v", err)
	}
	defer uvt.Delete()

	uvtF, err := GetUserVerificationTypeByCode(uvt.Tag)
	if err != nil {
		t.Fatalf("Cant find ver type by Code %v", err)
	}

	if uvt.Id != uvtF.Id {
		t.Fatalf("А Id то разные мужик!")
	}
}

func TestUserVerificationType_Delete(t *testing.T) {
	uvt, err := UserVerificationMethod{Name: "TestDelete", Tag: "testCode"}.Create()
	if err != nil {
		t.Fatalf("Cant create ver %v", err)
	}
	defer uvt.Delete()

	// убеждаемся, что код верфикации есть
	_, err = GetUserVerificationTypeById(uvt.Id)
	if err != nil {
		t.Fatalf("Cant find ver type by id %v", err)
	}

	if err := uvt.Delete(); err != nil {
		t.Fatalf("Ай ай, удаление не прошло: %v", err)
	}

	// убеждаемся, что кода верфикации нет
	_, err = GetUserVerificationTypeById(uvt.Id)
	if err == nil {
		t.Fatalf("Нашли код, хотя он должен был быть удален")
	}

	_, err = GetUserVerificationTypeByCode(uvt.Tag)
	if err == nil {
		t.Fatalf("Нашли код, хотя он должен был быть удален")
	}
}
