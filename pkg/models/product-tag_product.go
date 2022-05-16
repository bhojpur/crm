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
type ProductTagProduct struct {
	ProductId    uint `json:"product_id" gorm:"type:int;index;not null;"`
	ProductTagId uint `json:"product_tag_id" gorm:"type:int;index;not null;"`
}

func (ProductTagProduct) PgSqlCreate() {
	if err := db.Migrator().CreateTable(&ProductTagProduct{}); err != nil {
		log.Fatal(err)
	}
	err := db.Exec("ALTER TABLE product_tag_products \n    ADD CONSTRAINT product_tag_products_id_fkey FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE ON UPDATE CASCADE,\n    ADD CONSTRAINT product_tag_products_product_tag_id_fkey FOREIGN KEY (product_tag_id) REFERENCES product_tags(id) ON DELETE CASCADE ON UPDATE CASCADE,\n    DROP CONSTRAINT IF EXISTS fk_product_tag_products_product_tag,\n    DROP CONSTRAINT IF EXISTS fk_product_tag_products_product,\n    DROP CONSTRAINT IF EXISTS fk_product_tag_products_tag;\n\ncreate unique index uix_product_tag_products_product_id_product_tag_id ON product_tag_products (product_id,product_tag_id);").Error
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
