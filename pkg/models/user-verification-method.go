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
	"log"

	"github.com/bhojpur/crm/pkg/utils"
	"gorm.io/gorm"
)

//UserVerificationMethod - верификация пользователя
type UserVerificationMethod struct {
	Id          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" gorm:"type:varchar(255)"`                // Регистрация по email, ...
	Tag         string  `json:"tag" gorm:"type:varchar(50);unique;not null;"` // email, phone, email-phone - одна из заранее определенных констант!
	Description *string `json:"description" gorm:"type:varchar(255);"`        // краткое описание
}

const (
	VerificationMethodEmail         = "email"
	VerificationMethodPhone         = "phone"
	VerificationMethodEmailAndPhone = "email+phone"
)

func (UserVerificationMethod) PgSqlCreate() {
	if db.Migrator().HasTable(&UserVerificationMethod{}) {
		return
	}
	if err := db.Migrator().CreateTable(&UserVerificationMethod{}); err != nil {
		log.Fatal(err)
	}

	var verificationMethods = []UserVerificationMethod{
		{Name: "Email-верификация", Tag: VerificationMethodEmail, Description: utils.STRp("Пользователю будет необходимо перейти по ссылке в email.")},
		{Name: "SMS-верификация", Tag: VerificationMethodPhone, Description: utils.STRp("Пользователю необходимо будет ввести код из SMS.")},
		{Name: "Двойная Email+SMS верификация", Tag: VerificationMethodEmailAndPhone, Description: utils.STRp("Пользователю необходимо будет ввести код из SMS в специальной форме по ссылке в email.")},
	}
	for _, v := range verificationMethods {
		_, err := v.Create()
		if err != nil {
			log.Fatalf("Не удалось создать тип верификации: %v", v)
		}
	}
}

// Пользователь проходит верификацию, когда поля, указанные в методе верификации пользователя в аккаунте, - надежно подтверждены самим пользователем.
func (uvt UserVerificationMethod) Create() (*UserVerificationMethod, error) {

	if len([]rune(uvt.Name)) < 1 {
		return nil, utils.Error{Message: "Не верно указаны данные", Errors: map[string]interface{}{"name": "Введите описание типа верификации"}}
	}

	if len([]rune(uvt.Tag)) < 2 {
		return nil, utils.Error{Message: "Не верно указаны данные", Errors: map[string]interface{}{"code": "код должен быть не менее 2х символов"}}
	}

	if err := db.Create(&uvt).Error; err != nil {
		return nil, err
	}

	return &uvt, nil
}

func GetUserVerificationTypeById(id uint) (*UserVerificationMethod, error) {
	uvt := UserVerificationMethod{}
	err := db.First(&uvt, id).Error
	if err == gorm.ErrRecordNotFound {
		err = utils.Error{Message: "Метод верификации не найден"}
	}
	return &uvt, err
}

func GetUserVerificationTypeByCode(tag string) (*UserVerificationMethod, error) {
	uvt := UserVerificationMethod{}
	err := db.First(&uvt, "tag = ?", tag).Error
	return &uvt, err
}

func (uvt UserVerificationMethod) Delete() error {
	return db.Model(&UserVerificationMethod{}).Where("id = ?", uvt.Id).Delete(uvt).Error
}
