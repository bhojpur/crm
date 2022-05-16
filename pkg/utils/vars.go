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
	"reflect"
	"time"
)

func UINTp(x uint) *uint {
	return &x
}
func INTp(x int) *int {
	return &x
}
func ParseUINTp(x *uint) uint {
	if x == nil {
		return 0
	}
	return *x
}

func FL64p(x float64) *float64 {
	return &x
}
func STRp(s string) *string {
	return &s
}
func TimeP(t time.Time) *time.Time {
	return &t
}
func ConvertMapVarsToUINT(input *map[string]interface{}, keys []string) error {

	for _, key := range keys {
		if _vI, ok := (*input)[key]; ok && _vI != nil {
			switch reflect.TypeOf(_vI).String() {

			// case "float64":
			case "uint":
				(*input)[key] = _vI
			case "int":
				_vFInt, ok := _vI.(int)
				if !ok {
					return Error{Message: "Техническая ошибка в запросе", Errors: map[string]interface{}{key: "не удается разобрать значение"}}
				}
				(*input)[key] = uint(_vFInt)
			case "string":
				(*input)[key] = nil
				// fmt.Println("Char")
			default:
				_vF64, ok := _vI.(float64)
				if !ok {
					return Error{Message: "Техническая ошибка в запросе", Errors: map[string]interface{}{key: "не удается разобрать значение"}}
				}
				(*input)[key] = uint(_vF64)
			}
		}
	}
	return nil
}

func ConvertMapVarsToFloat64(input *map[string]interface{}, keys []string) error {

	for _, key := range keys {
		if _vI, ok := (*input)[key]; ok && _vI != nil {
			switch reflect.TypeOf(_vI).String() {

			// case "float64":
			case "uint":
				(*input)[key] = _vI
			case "int":
				_vFInt, ok := _vI.(int)
				if !ok {
					return Error{Message: "Техническая ошибка в запросе", Errors: map[string]interface{}{key: "не удается разобрать значение"}}
				}
				(*input)[key] = uint(_vFInt)
			default:
				_vF64, ok := _vI.(float64)
				if !ok {
					return Error{Message: "Техническая ошибка в запросе", Errors: map[string]interface{}{key: "не удается разобрать значение"}}
				}
				(*input)[key] = uint(_vF64)
			}
		}
	}
	return nil
}

// возвращает только разрешенные ключи
func FilterAllowedKeySTRArray(input []string, keys []string) []string {

	retArr := make([]string, 0)
	for inK := range input {
		for key := range keys {
			if ToCamel(input[inK]) == ToCamel(keys[key]) {
				retArr = append(retArr, keys[key])
			}
		}
	}
	return retArr
}
