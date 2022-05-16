//go:build !server
// +build !server

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

package main

import (
	"fmt"
	"log"

	"github.com/bhojpur/crm/pkg/database/base"
	"github.com/bhojpur/crm/pkg/models"
	"github.com/bhojpur/crm/pkg/routes"
	"github.com/joho/godotenv"
	"github.com/ttacon/libphonenumber"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error load .env file", err)
	}
}

func main() {
	models.ConnectDb()
	// defer db.Close()
	// defer pool.Close()

	// base.Test()

	if false {
		base.RefreshTablesPart_I()
		base.UploadTestDataPart_I()
		base.LoadImagesAiroClimate(13)
		base.LoadArticlesAiroClimate()
		base.LoadProductDescriptionAiroClimate()
		base.LoadProductCategoryDescriptionAiroClimate()
		base.UploadTestDataPart_II()
		base.UploadTestDataPart_III()
		base.UploadTestDataPart_IV()
		base.UploadBroUserData()

		base.UploadTestDataPart_V()
		base.Upload357grData()
	}

	// base.MigrateQuestions()

	if err := models.SettingsDb(); err != nil {
		log.Fatal(err)
	}

	// base.Migrate_I()

	if err := (models.EventListener{}).ReloadEventHandlers(); err != nil {
		log.Fatal(fmt.Sprintf("Failed to register EventHandler: %v", err))
	}

	models.RunHttpServer(routes.Handlers())
}

func examplePhone(numToParse string) {

	//num, err := libphonenumber.Get
	num, err := libphonenumber.Parse(numToParse, "IN")
	if err != nil {
		// Handle error appropriately.
		log.Fatal("Err: ", err)
	}
	formattedNum := libphonenumber.Format(num, libphonenumber.NATIONAL)

	fmt.Println("CountryCode: ", *num.CountryCode)
	fmt.Println("National Number: ", *num.NationalNumber)
	fmt.Println("National Formatted: ", formattedNum)
	fmt.Println("RFC3966: ", libphonenumber.Format(num, libphonenumber.RFC3966))
	fmt.Println("INTERNATIONAL: ", libphonenumber.Format(num, libphonenumber.INTERNATIONAL)) // наиболее популярный
	fmt.Println("E164: ", libphonenumber.Format(num, libphonenumber.E164))

	// num is a *libphonenumber.PhoneNumber
}
