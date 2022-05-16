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

import "testing"

// ### Testing CRUD functions

func TestUser_create(t *testing.T) {
	// see: TestAccount_UserCreate in account_test.go
}

func TestUser_getUserById(t *testing.T) {
	account, err := Account{Name: "TestAccount_getUserById"}.create()
	if err != nil {
		t.Fatalf("Не удалось создать тестовый аккаунт: %v", err)
	}
	defer account.HardDelete()

	user := User{
		IssuerAccountId: account.Id,
		Username:        "TestUser_getUserById",
	}
	u, err := user.create()
	if err != nil {
		t.Fatalf("Не удалось создать пользователя, %v", err)
	}
	defer u.hardDelete()
}

func TestUser_hardDelete(t *testing.T) {
	account, err := Account{Name: "TestAccount_hardDelete"}.create()
	if err != nil {
		t.Fatalf("Не удалось создать тестовый аккаунт: %v", err)
	}
	defer account.HardDelete()

	testUser := &User{
		IssuerAccountId: account.Id,
		Username:        "TestUser_getUserById",
	}
	user, err := testUser.create()
	if err != nil {
		t.Fatalf("Не удалось создать пользователя, %v", err)
	}
	if err := user.hardDelete(); err != nil {
		t.Fatalf("Не удалось удалить пользователя, %v", err)
	}
	// проверяем, что пользователя нет
	_, err = getUserById(user.Id)
	if err == nil {
		t.Fatalf("Пользователь на самом деле не удалился")
	}

}

func TestUser_softDelete(t *testing.T) {
	account, err := Account{Name: "TestUser_softDelete"}.create()
	if err != nil {
		t.Fatalf("Не удалось создать тестовый аккаунт: %v", err)
	}
	defer account.HardDelete()

	testUser := &User{
		IssuerAccountId: account.Id,
		Username:        "TestUser_getUserById",
	}
	user, err := testUser.create()
	if err != nil {
		t.Fatalf("Не удалось создать пользователя, %v", err)
	}
	defer user.hardDelete()

	if err := user.softDelete(); err != nil {
		t.Fatalf("Не удалось удалить пользователя, %v", err)
	}

	// проверяем, что пользователя нет
	_, err = getUserById(user.Id)
	if err == nil {
		t.Fatalf("Пользователь на самом деле не удалился")
	}

	// а вот тут пользователь должен был найтись
	fUser, err := getUnscopedUserById(user.Id)
	if err != nil {
		t.Fatalf("Пользователь должен был найтись")
	}

	if fUser.Id != user.Id {
		t.Fatalf("Id пользователей не совпадают")
	}

}
