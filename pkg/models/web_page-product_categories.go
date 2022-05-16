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
)

// HELPER FOR M<>M IN PgSQL
type WebPageProductCategory struct {
	WebPageId         uint `json:"web_page_id"`
	ProductCategoryId uint `json:"product_category_id"`

	WebPage         WebPage         `json:"-"`
	ProductCategory ProductCategory `json:"-"`

	// Порядок отображения категории на странице (нужно ли)
	Priority *int `json:"priority" gorm:"type:int;default:10;"`
}

func (WebPageProductCategory) PgSqlCreate() {
	if err := db.Migrator().CreateTable(&WebPageProductCategory{}); err != nil {
		log.Fatal(err)
	}
	err := db.Exec("ALTER TABLE web_page_product_categories \n    ADD CONSTRAINT web_page_product_categories_web_page_id_fkey FOREIGN KEY (web_page_id) REFERENCES web_pages(id) ON DELETE CASCADE ON UPDATE CASCADE,\n    ADD CONSTRAINT web_page_product_categories_product_category_id_fkey FOREIGN KEY (product_category_id) REFERENCES product_categories(id) ON DELETE CASCADE ON UPDATE CASCADE,\n    DROP CONSTRAINT IF EXISTS fk_web_page_product_categories_web_page,\n    DROP CONSTRAINT IF EXISTS fk_web_page_product_categories_product_category;").Error
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
