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
	"database/sql"
	"log"

	"github.com/bhojpur/crm/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Расчетный счет - не публичны, привязан к аккаунту
type PaymentAccount struct {
	Id        uint `json:"id" gorm:"primaryKey"`
	PublicId  uint `json:"public_id" gorm:"type:int;index;not null;"`
	AccountId uint `json:"-" gorm:"type:int;index;not null;"`

	BankId uint `json:"bank_id" gorm:"type:int;"`
	Bank   Bank `json:"bank"`

	// Расчетный счет
	AccAt string `json:"acc_at" gorm:"type:varchar(20);"`

	// Описание / Назначение счета (если надо)
	Description *string `json:"description" gorm:"type:varchar(255);"`
}

func (PaymentAccount) PgSqlCreate() {
	if err := db.Migrator().CreateTable(&PaymentAccount{}); err != nil {
		log.Fatal(err)
	}
	err := db.Exec("ALTER TABLE payment_accounts " +
		"ADD CONSTRAINT payment_accounts_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE," +
		"ADD CONSTRAINT payment_accounts_bank_id_fkey FOREIGN KEY (bank_id) REFERENCES banks(id) ON DELETE CASCADE ON UPDATE CASCADE;").Error
	if err != nil {
		log.Fatal("Error: ", err)
	}

}
func (paymentAccount *PaymentAccount) BeforeCreate(tx *gorm.DB) error {
	paymentAccount.Id = 0

	// PublicId
	var lastIdx sql.NullInt64
	err := db.Model(&PaymentAccount{}).Where("account_id = ?", paymentAccount.AccountId).
		Select("max(public_id)").Row().Scan(&lastIdx)
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	paymentAccount.PublicId = 1 + uint(lastIdx.Int64)

	return nil
}
func (paymentAccount *PaymentAccount) AfterFind(tx *gorm.DB) (err error) {

	return nil
}
func (paymentAccount *PaymentAccount) GetPreloadDb(getModel bool, autoPreload bool, preloads []string) *gorm.DB {

	_db := db

	if getModel {
		_db = _db.Model(paymentAccount)
	} else {
		_db = _db.Model(&PaymentAccount{})
	}

	if autoPreload {
		return _db.Preload(clause.Associations)
	} else {

		allowed := utils.FilterAllowedKeySTRArray(preloads, []string{"Bank"})

		for _, v := range allowed {
			_db.Preload(v)
		}
		return _db
	}
}

// ############# Entity interface #############
func (paymentAccount PaymentAccount) GetId() uint                { return paymentAccount.Id }
func (paymentAccount *PaymentAccount) setId(id uint)             { paymentAccount.Id = id }
func (paymentAccount *PaymentAccount) setPublicId(publicId uint) { paymentAccount.PublicId = publicId }
func (paymentAccount PaymentAccount) GetAccountId() uint         { return paymentAccount.AccountId }
func (paymentAccount *PaymentAccount) setAccountId(id uint)      { paymentAccount.AccountId = id }
func (PaymentAccount) SystemEntity() bool                        { return false }

// ############# End Entity interface #############

// ######### CRUD Functions ############
func (paymentAccount PaymentAccount) create() (Entity, error) {

	en := paymentAccount

	if err := db.Create(&en).Error; err != nil {
		return nil, err
	}

	err := en.GetPreloadDb(false, false, nil).First(&en, en.Id).Error
	if err != nil {
		return nil, err
	}

	var newItem Entity = &en

	return newItem, nil
}
func (PaymentAccount) get(id uint, preloads []string) (Entity, error) {

	var paymentAccount PaymentAccount

	err := paymentAccount.GetPreloadDb(false, false, preloads).First(&paymentAccount, id).Error
	if err != nil {
		return nil, err
	}
	return &paymentAccount, nil
}
func (paymentAccount *PaymentAccount) load(preloads []string) error {

	err := paymentAccount.GetPreloadDb(false, false, preloads).First(paymentAccount, paymentAccount.Id).Error
	if err != nil {
		return err
	}
	return nil
}
func (paymentAccount *PaymentAccount) loadByPublicId(preloads []string) error {

	if paymentAccount.PublicId < 1 {
		return utils.Error{Message: "Невозможно загрузить PaymentAccount - не указан  Id"}
	}
	if err := paymentAccount.GetPreloadDb(false, false, preloads).First(paymentAccount, "account_id = ? AND public_id = ?", paymentAccount.AccountId, paymentAccount.PublicId).Error; err != nil {
		return err
	}

	return nil
}
func (PaymentAccount) getList(accountId uint, sortBy string, preload []string) ([]Entity, int64, error) {
	return PaymentAccount{}.getPaginationList(accountId, 0, 100, sortBy, "", nil, preload)
}
func (PaymentAccount) getPaginationList(accountId uint, offset, limit int, sortBy, search string, filter map[string]interface{}, preloads []string) ([]Entity, int64, error) {

	paymentAccounts := make([]PaymentAccount, 0)
	var total int64

	// if need to search
	if len(search) > 0 {

		search = "%" + search + "%"

		err := (&PaymentAccount{}).GetPreloadDb(false, false, preloads).Limit(limit).Limit(limit).Offset(offset).Order(sortBy).Where("account_id = ?", accountId).
			Where(filter).Find(&paymentAccounts, "acc_at ILIKE ? ", search).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}

		// Определяем total
		err = db.Model(&PaymentAccount{}).
			Where("account_id = ? AND acc_at ILIKE ?", search, search).
			Count(&total).Error
		if err != nil {
			return nil, 0, utils.Error{Message: "Ошибка определения объема базы"}
		}

	} else {

		err := (&PaymentAccount{}).GetPreloadDb(false, false, preloads).Limit(limit).Offset(offset).Order(sortBy).Where("account_id = ?", accountId).
			Where(filter).Find(&paymentAccounts).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}

		// Определяем total
		err = db.Model(&PaymentAccount{}).Where("account_id = ?", accountId).Count(&total).Error
		if err != nil {
			return nil, 0, utils.Error{Message: "Ошибка определения объема базы"}
		}
	}

	// Преобразуем полученные данные
	entities := make([]Entity, len(paymentAccounts))
	for i := range paymentAccounts {
		entities[i] = &paymentAccounts[i]
	}

	return entities, total, nil
}
func (paymentAccount *PaymentAccount) update(input map[string]interface{}, preloads []string) error {

	delete(input, "bank")
	utils.FixInputHiddenVars(&input)
	if err := utils.ConvertMapVarsToUINT(&input, []string{"public_id", "bank_id"}); err != nil {
		return err
	}

	if err := paymentAccount.GetPreloadDb(false, false, nil).Where(" id = ?", paymentAccount.Id).
		Omit("id", "account_id", "public_id").Updates(input).Error; err != nil {
		return err
	}

	err := paymentAccount.GetPreloadDb(false, false, preloads).First(paymentAccount, paymentAccount.Id).Error
	if err != nil {
		return err
	}

	return nil
}
func (paymentAccount *PaymentAccount) delete() error {
	return paymentAccount.GetPreloadDb(true, false, nil).Where("id = ?", paymentAccount.Id).Delete(paymentAccount).Error
}

// ######### END CRUD Functions ############
