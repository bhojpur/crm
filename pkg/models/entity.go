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
	"reflect"
)

type Entity interface {
	GetId() uint
	setId(id uint)
	setPublicId(id uint)
	GetAccountId() uint
	setAccountId(id uint)

	// CRUD model
	create() (Entity, error)
	get(id uint, preloads []string) (Entity, error)
	load(preloads []string) error
	loadByPublicId(preloads []string) error
	getList(accountId uint, order string, preloads []string) ([]Entity, int64, error)
	getPaginationList(accountId uint, offset, limit int, sortBy, search string, filter map[string]interface{}, preloads []string) ([]Entity, int64, error)
	update(input map[string]interface{}, preloads []string) error
	delete() error

	// AppendAssociationMethod(options Entity)
	SystemEntity() bool
}

func Get(v Entity) error {

	// id := v.getId()

	r := reflect.TypeOf(v)

	fmt.Println(r, r.Elem(), r.NumMethod())

	println("We are use v Entity GET function")

	return nil
}

func (account Account) CreateEntity(input Entity) (Entity, error) {
	input.setAccountId(account.Id)
	return input.create()
}

func (account Account) GetEntity(model Entity, id uint, preloads []string) (Entity, error) {

	entity, err := model.get(id, preloads)
	if err != nil {
		return nil, err
	}

	// Тут надо бы показать, что она системная
	if entity.GetAccountId() != account.Id && !entity.SystemEntity() {
		return nil, errors.New("Модель принадлежит другому аккаунту")
	}

	return entity, nil
}

func (account Account) LoadEntity(entity Entity, primaryKey uint, preloads []string) error {

	// На всякий случай
	entity.setId(primaryKey)

	// Загружаем по ссылке
	err := entity.load(preloads)
	if err != nil {
		return err
	}

	// Проверяем принадлежность к аккаунту
	if entity.GetAccountId() != account.Id && !entity.SystemEntity() {
		return errors.New("Модель принадлежит другому аккаунту")
	}

	return nil
}

func (account Account) LoadEntityByPublicId(entity Entity, publicId uint, preloads []string) error {

	// На всякий случай
	entity.setPublicId(publicId)
	entity.setAccountId(account.Id)

	// Загружаем по ссылке
	err := entity.loadByPublicId(preloads)
	if err != nil {
		return err
	}

	// Проверяем принадлежность к аккаунту
	if entity.GetAccountId() != account.Id && !entity.SystemEntity() {
		return errors.New("Модель принадлежит другому аккаунту")
	}

	return nil
}

func (account Account) GetListEntity(model Entity, order string, preload []string) ([]Entity, int64, error) {
	return model.getList(account.Id, order, preload)
}

func (account Account) GetPaginationListEntity(model Entity, offset, limit int, order string, search string, filter map[string]interface{}, preloads []string) ([]Entity, int64, error) {
	return model.getPaginationList(account.Id, offset, limit, order, search, filter, preloads)
}

func (account Account) UpdateEntity(entity Entity, input map[string]interface{}, preloads []string) error {
	if entity.GetAccountId() != account.Id && !entity.SystemEntity() {
		return errors.New("Модель принадлежит другому аккаунту")
	}

	return entity.update(input, preloads)
}

func (account Account) DeleteEntity(entity Entity) error {
	if entity.GetAccountId() != account.Id && !entity.SystemEntity() {
		return errors.New("Модель принадлежит другому аккаунту")
	}

	return entity.delete()
}
