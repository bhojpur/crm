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

	"gorm.io/gorm"
)

// Хелп таблица ManyToMany PaymentMethods <> Delivery
type Payment2Delivery struct {
	Id        uint `json:"id" gorm:"primaryKey"`
	AccountId uint `json:"-" gorm:"index;not null"`
	WebSiteId uint `json:"web_site_id" gorm:"type:int;index;"`

	PaymentId   uint   `json:"payment_id" gorm:"type:int;not null;"`
	PaymentType string `json:"payment_type" gorm:"varchar(32);not null;"`

	DeliveryId   uint   `json:"delivery_id" gorm:"type:int;not null;"`
	DeliveryType string `json:"delivery_type" gorm:"varchar(32);not null;"`
}

func (Payment2Delivery) PgSqlCreate() {
	if err := db.Migrator().CreateTable(&Payment2Delivery{}); err != nil {
		log.Fatal(err)
	}
	// db.Model(&Payment2Delivery{}).AddForeignKey("account_id", "accounts(id)", "CASCADE", "CASCADE")
	// db.Model(&Payment2Delivery{}).AddForeignKey("web_site_id", "web_sites(id)", "CASCADE", "CASCADE")
	err := db.Exec("ALTER TABLE payment_to_delivery " +
		"ADD CONSTRAINT payment_to_delivery_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE," +
		"ADD CONSTRAINT payment_to_delivery_web_site_id_fkey FOREIGN KEY (web_site_id) REFERENCES web_sites(id) ON DELETE CASCADE ON UPDATE CASCADE;").Error
	if err != nil {
		log.Fatal("Error: ", err)
	}
}
func (p2d *Payment2Delivery) BeforeCreate(tx *gorm.DB) error {
	p2d.Id = 0
	return nil
}
func (Payment2Delivery) TableName() string {
	return "payment_to_delivery"
}

// Добавляет, если есть - ничего не делает
func (webSite WebSite) AppendPayment2Delivery(paymentId uint, paymentType string, deliveryId uint, deliveryType string) error {
	if err := db.Create(&Payment2Delivery{
		AccountId:    webSite.AccountId,
		WebSiteId:    webSite.Id,
		PaymentId:    paymentId,
		PaymentType:  paymentType,
		DeliveryId:   deliveryId,
		DeliveryType: deliveryType,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (p2d Payment2Delivery) Append(paymentMethod PaymentMethod, delivery Delivery) {

}

func (Payment2Delivery) RemoveByIds(paymentMethodId uint, paymentMethodType string, paymentDeliveryId uint, paymentDeliveryType string) {

}
