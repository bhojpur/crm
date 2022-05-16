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
	"errors"
	"fmt"

	"github.com/bhojpur/crm/pkg/utils"
)

// "payment_cash", "payment_yandex"
type PaymentMethod interface {
	Entity

	GetType() string
	GetCode() string

	// Функция запускающая процесс создания платежа под Order (Заказ)
	CreatePaymentByOrder(order Order, mode PaymentMode) (*Payment, error)
	PrepaymentCheck(payment *Payment, order Order) (*Payment, error)
	GetWebSiteId() uint

	// true - моментальная передача товара, false- разнесено во времени; нужно ли разносить оплату и доставку
	IsInstantDelivery() bool
}

func (account Account) GetPaymentMethods() ([]PaymentMethod, error) {

	// Находим все необходимые методы
	var paymentCashes []PaymentCash
	paymentCashesEntity, _, err := PaymentCash{}.getList(account.Id, "id", nil)
	if err != nil {
		return nil, err
	}
	for i := range paymentCashesEntity {
		_p, ok := paymentCashesEntity[i].(*PaymentCash)
		if !ok {
			return nil, utils.Error{Message: "Ошибка при получении списка PaymentCash"}
		}
		paymentCashes = append(paymentCashes, *_p)
	}

	var paymentYandexes []PaymentYandex
	paymentYandexEntity, _, err := PaymentYandex{}.getList(account.Id, "id", nil)
	if err != nil {
		return nil, err
	}
	for i := range paymentYandexEntity {
		_p, ok := paymentYandexEntity[i].(*PaymentYandex)
		if !ok {
			return nil, utils.Error{Message: "Ошибка при получении списка PaymentCash"}
		}

		paymentYandexes = append(paymentYandexes, *_p)
	}

	methods := make([]PaymentMethod, len(paymentCashes)+len(paymentYandexes))
	for i := range paymentCashes {
		methods[i] = &paymentCashes[i]
	}
	for i := range paymentYandexes {
		methods[i+len(paymentCashes)] = &paymentYandexes[i]
	}

	return methods, nil
}

func (account Account) GetPaymentMethodByCode(code string, methodId uint) (PaymentMethod, error) {

	if code == "" || methodId < 1 {
		return nil, utils.Error{Message: "Не верно указан код типа оплаты"}
	}
	// 1. Получаем все варианты доставки (обычно их мало). Можно через switch, но лень потом исправлять баг с новыми типом доставки
	methods, err := account.GetPaymentMethods()
	if err != nil {
		return nil, err
	}

	// Ищем наш вариант доставки
	var method PaymentMethod
	for _, v := range methods {
		// fmt.Println("Code: ", v.GetCode(), code)
		// fmt.Println("v.GetId(): ", v.GetId(), methodId)
		if v.GetCode() == code && v.GetId() == methodId {
			method = v
			break
		}
	}

	// Проверяем, удалось ли найти выбранный вариант оплаты
	if method == nil {
		return nil, utils.Error{Message: "Не верно указан тип оплаты"}
	}

	return method, nil
}

func (account Account) GetPaymentMethodByType(codeType string, methodId uint) (PaymentMethod, error) {

	if codeType == "" || methodId < 1 {
		fmt.Println("codeType: ", codeType)
		fmt.Println("methodId: ", methodId)
		return nil, errors.New("Не верно указаны данные типа оплаты")
		// return nil, utils.Error{Message: "Не верно указаны данные типа оплаты"}
	}
	// 1. Получаем все варианты доставки (обычно их мало). Можно через switch, но лень потом исправлять баг с новыми типом доставки
	methods, err := account.GetPaymentMethods()
	if err != nil {
		return nil, err
	}

	// Ищем наш вариант доставки
	var method PaymentMethod
	for _, v := range methods {
		if v.GetType() == codeType && v.GetId() == methodId {
			method = v
			break
		}
	}

	// Проверяем, удалось ли найти выбранный вариант оплаты
	if method == nil {
		return nil, utils.Error{Message: "Не верно указан тип оплаты"}
	}

	return method, nil
}
