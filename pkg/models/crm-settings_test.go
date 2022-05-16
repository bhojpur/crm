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
)

func TestGetCrmSettings(t *testing.T) {
	settings, err := GetCrmSettings()
	if err != nil || settings == nil {
		t.Fatalf("Failed to load настройки crm-системы: %v", err)
	}

	if settings.Id != 1 {
		t.Fatalf("Id у строки настроек не равен 1: , %v", settings)
	}
}

func TestCrmSetting_Save(t *testing.T) {
	settings, err := GetCrmSettings()
	if err != nil {
		t.Fatal("Не удалось получить настройки CRM")
	}

	settings.ApiEnabled = !settings.ApiEnabled
	if err := settings.Save(); err != nil {
		t.Fatalf("Cant Save CRM settings: %v", err)
	}

	// возвращаем назад настройки
	defer func() {
		settings.ApiEnabled = !settings.ApiEnabled
		if err := settings.Save(); err != nil {
			t.Fatalf("Cannot back CRM settings : %v", err)
		}
	}()

	newSettings, err := GetCrmSettings()
	if err != nil {
		t.Fatal("Не удалось получить настройки CRM (2)")
	}

	if newSettings.ApiEnabled != settings.ApiEnabled {
		t.Fatal("Функция Save CrmSettings реально не сохранила данные")
	}
}

func BenchmarkGetCrmSettings(b *testing.B) {
	settings, err := GetCrmSettings()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		settings, err = GetCrmSettings()
		if err != nil || settings == nil {
			b.Fatalf("Failed to load настройки crm-системы: %v", err)
		}
	}
}
