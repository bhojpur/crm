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
	"crypto/aes"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
	"unicode/utf8"
	"unsafe"
)

//const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyz"
const letterBytesChar = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const letterBytesNum = "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	return StringWithCharset(n, letterBytesNum)
}

func StringWithCharset(length int, charset string) string {

	var seededRand = rand.New(
		rand.NewSource(time.Now().UTC().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Генерит случайную строку длиной {n} по-взрослому!
func RandStringBytesMaskImprSrcUnsafe(n int, withNum bool) string {

	letterBytes := letterBytesChar

	if withNum {
		letterBytes = letterBytesNum
	}

	var src = rand.NewSource((time.Now().Add(time.Minute * 268)).UnixNano())

	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// Считаем длину ключа 128 битным
func CreateAes128Key() (string, error) {
	str := RandStringBytesMaskImprSrcUnsafe(16, false)
	return str, ValidationAesKey(str)
}
func ValidationAesKey(key string) error {
	if len(key) != 16 {
		return Error{Message: "Некорректная длина ключа AES-128", Errors: map[string]interface{}{"aesKey": "Длина 128-битного ключа должна быть равна 16 символам в UTF-8"}}
	}

	// check utf-8
	r, size := utf8.DecodeRune([]byte(key))
	if r == utf8.RuneError || size < 1 {
		return Error{Message: "Некорректный ключ AES-128", Errors: map[string]interface{}{"aesKey": "Символы должны быть в кодировке UTF-8"}}
	}

	// пробуем создать новый AES Cipher
	_, err := aes.NewCipher([]byte(key))
	if err != nil {
		return Error{Message: "Некорректный ключ AES-128", Errors: map[string]interface{}{"aesKey": "Проверьте символы и кодировку UTF-8"}}
	}

	return nil
}
func CreateHS256Key() string {
	return RandStringBytesMaskImprSrcUnsafe(32, false)
}

func GetMD5Hash(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}
