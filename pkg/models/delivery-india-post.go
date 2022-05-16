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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bhojpur/crm/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeliveryIndiaPost struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	AccountId uint   `json:"-" gorm:"index;not null"`                                      // аккаунт-владелец ключа
	WebSiteId uint   `json:"web_site_id" gorm:"type:int;index;"`                           // магазин, к которому относится
	Code      string `json:"code" gorm:"type:varchar(16);default:'indiaPost';"`            // Для идентификации во фронтенде
	Type      string `json:"type" gorm:"type:varchar(32);default:'delivery_india_posts';"` // Для идентификации

	Enabled bool   `json:"enabled" gorm:"type:bool;default:true"` // активен ли способ доставки
	Name    string `json:"name" gorm:"type:varchar(255);"`        // "Курьерская доставка", "Почта России", "Самовывоз"

	AccessToken        string  `json:"access_token" gorm:"type:varchar(255);"`         // accessToken
	XUserAuthorization string  `json:"x_user_authorization" gorm:"type:varchar(255);"` // XUserAuthorization в base64
	MaxWeight          float64 `json:"max_weight" gorm:"type:numeric;default:20"`      // максимальная масса в кг

	PostalCodeFrom       string `json:"postal_code_from" gorm:"type:varchar(255);"`            // индекс отправки с почты России
	MailCategory         string `json:"mail_category" gorm:"type:varchar(50);"`                // https://otpravka.pochta.ru/specification#/enums-base-mail-category
	MailType             string `json:"mail_type" gorm:"type:varchar(50);"`                    // https://otpravka.pochta.ru/specification#/enums-base-mail-type
	Fragile              bool   `json:"fragile" gorm:"type:bool;default:false"`                // отметка "Осторожно хрупкое"
	WithElectronicNotice bool   `json:"with_electronic_notice"" gorm:"type:bool;default:true"` // отметка "Осторожно хрупкое"
	WithOrderOfNotice    bool   `json:"with_order_of_notice" gorm:"type:bool;default:true"`    // отметка "Осторожно хрупкое"
	WithSimpleNotice     bool   `json:"with_simple_notice" gorm:"type:bool;default:false"`     // отметка "Осторожно хрупкое"

	AddressRequired    bool `json:"address_required" gorm:"type:bool;default:true"`     // Требуется ли адрес доставки
	PostalCodeRequired bool `json:"postal_code_required" gorm:"type:bool;default:true"` // Требуется ли индекс в адресе доставки

	// Признак предмета расчета
	PaymentSubjectId *uint          `json:"payment_subject_id" gorm:"type:int;default:1;"` //
	PaymentSubject   PaymentSubject `json:"payment_subject"`

	VatCodeId uint    `json:"vat_code_id" gorm:"type:int;not null;default:1;"` // товар или услуга ? [вид номенклатуры]
	VatCode   VatCode `json:"vat_code"`

	// загружаемый интерфейс
	PaymentMethods []PaymentMethod `json:"payment_methods" gorm:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DeliveryIndiaPost) PgSqlCreate() {
	db.Migrator().CreateTable(&DeliveryIndiaPost{})

	// db.Model(&DeliveryIndiaPost{}).AddForeignKey("account_id", "accounts(id)", "CASCADE", "CASCADE")
	// db.Model(&DeliveryIndiaPost{}).AddForeignKey("payment_subject_id", "payment_subjects(id)", "RESTRICT", "CASCADE")
	// db.Model(&DeliveryIndiaPost{}).AddForeignKey("vat_code_id", "vat_codes(id)", "RESTRICT", "CASCADE")
	err := db.Exec("ALTER TABLE delivery_india_posts " +
		"ADD CONSTRAINT delivery_india_posts_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE ON UPDATE CASCADE," +
		"ADD CONSTRAINT delivery_india_posts_payment_subject_id_fkey FOREIGN KEY (payment_subject_id) REFERENCES payment_subjects(id) ON DELETE RESTRICT ON UPDATE CASCADE," +
		"ADD CONSTRAINT delivery_india_posts_vat_code_id_fkey FOREIGN KEY (vat_code_id) REFERENCES vat_codes(id) ON DELETE RESTRICT ON UPDATE CASCADE;").Error
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

// ############# Entity interface #############
func (deliveryIndiaPost DeliveryIndiaPost) GetId() uint          { return deliveryIndiaPost.Id }
func (deliveryIndiaPost *DeliveryIndiaPost) setId(id uint)       { deliveryIndiaPost.Id = id }
func (deliveryIndiaPost *DeliveryIndiaPost) setPublicId(id uint) {}
func (deliveryIndiaPost DeliveryIndiaPost) GetAccountId() uint {
	return deliveryIndiaPost.AccountId
}
func (deliveryIndiaPost DeliveryIndiaPost) GetWebSiteId() uint {
	return deliveryIndiaPost.WebSiteId
}
func (deliveryIndiaPost *DeliveryIndiaPost) setAccountId(id uint) {
	deliveryIndiaPost.AccountId = id
}
func (deliveryIndiaPost *DeliveryIndiaPost) setWebSiteId(webSiteId uint) {
	deliveryIndiaPost.WebSiteId = webSiteId
}
func (deliveryIndiaPost DeliveryIndiaPost) SystemEntity() bool { return false }

func (deliveryIndiaPost DeliveryIndiaPost) GetCode() string {
	return deliveryIndiaPost.Code
}
func (deliveryIndiaPost DeliveryIndiaPost) GetType() string {
	return deliveryIndiaPost.Type
}
func (deliveryIndiaPost DeliveryIndiaPost) GetPaymentSubject() PaymentSubject {
	return deliveryIndiaPost.PaymentSubject
}

// ############# Entity interface #############

// ###### GORM Functional #######
func (deliveryIndiaPost *DeliveryIndiaPost) BeforeCreate(tx *gorm.DB) error {
	deliveryIndiaPost.Id = 0
	return nil
}
func (deliveryIndiaPost *DeliveryIndiaPost) AfterFind(tx *gorm.DB) (err error) {

	methods, err := GetPaymentMethodsByDelivery(deliveryIndiaPost)
	if err != nil {
		return err
	}
	deliveryIndiaPost.PaymentMethods = methods

	return nil
}

// ###### End of GORM Functional #######

// ############# CRUD Entity interface #############
func (deliveryIndiaPost DeliveryIndiaPost) create() (Entity, error) {

	_item := deliveryIndiaPost
	if err := db.Create(&_item).Error; err != nil {
		return nil, err
	}

	if err := _item.GetPreloadDb(false, false, nil).First(&_item, _item.Id).Error; err != nil {
		return nil, err
	}

	var entity Entity = &_item

	return entity, nil
}

func (DeliveryIndiaPost) get(id uint, preloads []string) (Entity, error) {

	var item DeliveryIndiaPost

	err := item.GetPreloadDb(false, false, preloads).First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (deliveryIndiaPost *DeliveryIndiaPost) load(preloads []string) error {
	if deliveryIndiaPost.Id < 1 {
		return utils.Error{Message: "Невозможно загрузить DeliveryIndiaPost - не указан  Id"}
	}

	err := deliveryIndiaPost.GetPreloadDb(false, false, preloads).First(deliveryIndiaPost, deliveryIndiaPost.Id).Error
	if err != nil {
		return err
	}
	return nil
}
func (deliveryIndiaPost *DeliveryIndiaPost) loadByPublicId(preloads []string) error {
	return errors.New("Нет возможности загрузить объект по Public Id")
}

func (DeliveryIndiaPost) getList(accountId uint, sortBy string, preload []string) ([]Entity, int64, error) {
	return DeliveryIndiaPost{}.getPaginationList(accountId, 0, 100, sortBy, "", nil, preload)
}
func (DeliveryIndiaPost) getListByShop(accountId, websiteId uint) ([]DeliveryIndiaPost, error) {

	deliveryIndiaPosts := make([]DeliveryIndiaPost, 0)

	err := (&DeliveryIndiaPost{}).GetPreloadDb(false, false, []string{"PaymentMethods"}).
		Limit(100).Where("account_id = ? AND web_site_id = ?", accountId, websiteId).
		Find(&deliveryIndiaPosts).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return deliveryIndiaPosts, nil
}

func (DeliveryIndiaPost) getPaginationList(accountId uint, offset, limit int, sortBy, search string, filter map[string]interface{}, preloads []string) ([]Entity, int64, error) {

	deliveryIndiaPosts := make([]DeliveryIndiaPost, 0)
	var total int64

	// if need to search
	if len(search) > 0 {

		// string pattern
		search = "%" + search + "%"

		err := (&DeliveryIndiaPost{}).GetPreloadDb(false, false, preloads).Limit(limit).Offset(offset).Order(sortBy).Where("account_id = ?", accountId).
			Find(&deliveryIndiaPosts, "name ILIKE ? OR code ILIKE ? OR postal_code_from ILIKE ?", search, search, search).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}

		// Определяем total
		err = (&DeliveryIndiaPost{}).GetPreloadDb(false, false, nil).
			Where("account_id = ? AND name ILIKE ? OR code ILIKE ? OR postal_code_from ILIKE ?", accountId, search, search, search).
			Count(&total).Error
		if err != nil {
			return nil, 0, utils.Error{Message: "Ошибка определения объема базы"}
		}

	} else {

		err := (&DeliveryIndiaPost{}).GetPreloadDb(false, false, preloads).Limit(limit).Offset(offset).Order(sortBy).Where("account_id = ?", accountId).
			Find(&deliveryIndiaPosts).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, 0, err
		}

		// Определяем total
		err = (&DeliveryIndiaPost{}).GetPreloadDb(false, false, nil).Where("account_id = ?", accountId).Count(&total).Error
		if err != nil {
			return nil, 0, utils.Error{Message: "Ошибка определения объема базы"}
		}
	}

	// Преобразуем полученные данные
	entities := make([]Entity, len(deliveryIndiaPosts))
	for i := range deliveryIndiaPosts {
		entities[i] = &deliveryIndiaPosts[i]
	}

	return entities, total, nil
}

func (deliveryIndiaPost *DeliveryIndiaPost) update(input map[string]interface{}, preloads []string) error {
	delete(input, "payment_subject")
	delete(input, "vat_code")
	delete(input, "payment_methods")
	utils.FixInputHiddenVars(&input)
	if err := utils.ConvertMapVarsToUINT(&input, []string{"web_site_id", "payment_subject_id", "vat_code_id"}); err != nil {
		return err
	}

	if err := deliveryIndiaPost.GetPreloadDb(false, false, nil).Where("id = ?", deliveryIndiaPost.Id).Omit("id", "account_id").Updates(input).
		Error; err != nil {
		return err
	}

	err := deliveryIndiaPost.GetPreloadDb(false, false, preloads).First(deliveryIndiaPost, deliveryIndiaPost.Id).Error
	if err != nil {
		return err
	}

	return nil

}

func (deliveryIndiaPost *DeliveryIndiaPost) delete() error {
	return deliveryIndiaPost.GetPreloadDb(true, false, nil).Where("id = ?", deliveryIndiaPost.Id).Delete(deliveryIndiaPost).Error
}

// ########## End of CRUD Entity interface ###########
func (deliveryIndiaPost DeliveryIndiaPost) CalculateDelivery(deliveryData DeliveryData, weight float64) (float64, error) {

	if weight == 0 {
		// Грязный хак, в кг
		weight = 1
		// return 0, utils.Error{Message: "Ошибка расчета стоимости доставки: отсутствует вес товара"}
	}
	// базовые данные для запроса в api почта россиии
	url := "https://india-api.post.in/1.0/tariff"
	Authorization := "AccessToken " + deliveryIndiaPost.AccessToken
	XUserAuthorization := "Basic " + deliveryIndiaPost.XUserAuthorization

	// Формируем json для запроса
	rawJson := utils.MapToRawJson(map[string]interface{}{
		"index-from":    deliveryIndiaPost.PostalCodeFrom,
		"index-to":      deliveryData.PostalCode,
		"mail-category": deliveryIndiaPost.MailCategory,
		"mail-type":     deliveryIndiaPost.MailType,
		"mass":          weight * float64(1000), // масса в граммах (*1000)
		/*"dimension": map[string]interface{}{
			"height": 90, // в см.
			"length": 30, // в см.
			"width": 30, // в см.
		},*/
		"fragile":                deliveryIndiaPost.Fragile,              // отметка "Осторожно хрупкое"
		"with-electronic-notice": deliveryIndiaPost.WithElectronicNotice, // уведомление на емейл
		"with-order-of-notice":   deliveryIndiaPost.WithOrderOfNotice,    // уведомление заказное
		"with-simple-notice":     deliveryIndiaPost.WithSimpleNotice,     // уведомление заказное
	})

	// response, err := http.Post(url, "application/json", strings.NewReader(""))
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(rawJson))
	if err != nil {
		return 0, utils.Error{Message: "Ошибка связи с сервисом Почты России"}
	}

	request.Header.Set("Authorization", Authorization)
	request.Header.Set("X-User-Authorization", XUserAuthorization)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return 0, utils.Error{Message: "Ошибка связи с сервисом Почты России"}
	}
	defer response.Body.Close()

	// 1. Сначала узнаем статус запроса
	if response.StatusCode == 200 {
		var input struct {
			TotalRate float64 `json:"total-rate"`
			TotalVat  float64 `json:"total-vat"`
		}

		if err := json.NewDecoder(response.Body).Decode(&input); err != nil {
			return 0, utils.Error{Message: "Ошибка данных со стороны Почты России"}
		}

		return float64(input.TotalRate / 100), nil // т.е. в копейках
		// deliveryData.TotalCost = input.TotalRate / 100 // т.е. в копейках
		// return &deliveryData, nil

	} else {
		var input struct {
			Desc string `json:"desc"`
		}

		if err := json.NewDecoder(response.Body).Decode(&input); err != nil {
			return 0, err
		}

		return 0, utils.Error{Message: input.Desc}
	}

	return 0, utils.Error{Message: "Ошибка расчета стоимости"}
}

func (deliveryIndiaPost DeliveryIndiaPost) checkMaxWeight(weight float64) error {
	// проверяем максимальную массу:
	if weight > deliveryIndiaPost.MaxWeight {
		return utils.Error{Message: fmt.Sprintf("Превышен максимальный вес посылки в %vкг.", deliveryIndiaPost.MaxWeight),
			Errors: map[string]interface{}{"delivery": fmt.Sprintf("Превышен максимальный вес посылки в %vкг.", deliveryIndiaPost.MaxWeight)}}
	}

	return nil
}

func (deliveryIndiaPost DeliveryIndiaPost) GetName() string {
	return deliveryIndiaPost.Name
}
func (deliveryIndiaPost DeliveryIndiaPost) GetVatCode() (*VatCode, error) {
	if err := deliveryIndiaPost.load([]string{"VatCode"}); err != nil {
		return nil, err
	}
	return &deliveryIndiaPost.VatCode, nil
}

func (deliveryIndiaPost DeliveryIndiaPost) AppendPaymentMethods(paymentMethods []PaymentMethod) error {
	return nil
}
func (deliveryIndiaPost DeliveryIndiaPost) RemovePaymentMethods(paymentMethods []PaymentMethod) error {
	return nil
}
func (deliveryIndiaPost DeliveryIndiaPost) ExistPaymentMethod(paymentMethod PaymentMethod) bool {
	return true
}

func (deliveryIndiaPost DeliveryIndiaPost) CreateDeliveryOrder(deliveryData DeliveryData, cost float64, order Order) (Entity, error) {

	status, err := DeliveryStatus{}.GetStatusNew()
	if err != nil {
		return nil, err
	}
	customerId := uint(0)
	if order.CustomerId != nil {
		customerId = *order.CustomerId
	}
	webSiteId := uint(0)
	if order.WebSiteId != nil {
		webSiteId = *order.WebSiteId
	}
	deliveryOrder := DeliveryOrder{
		AccountId:  deliveryIndiaPost.AccountId,
		OrderId:    &order.Id,
		CustomerId: customerId,
		WebSiteId:  webSiteId,
		Code:       deliveryIndiaPost.Code,
		MethodId:   deliveryIndiaPost.Id,
		Address:    utils.STRp(deliveryData.Address),
		PostalCode: utils.STRp(deliveryData.PostalCode),
		Cost:       cost,
		StatusId:   status.Id,
	}

	return deliveryOrder.create()

}

func (deliveryIndiaPost *DeliveryIndiaPost) GetPreloadDb(getModel bool, autoPreload bool, preloads []string) *gorm.DB {

	_db := db

	if getModel {
		_db = _db.Model(deliveryIndiaPost)
	} else {
		_db = _db.Model(&DeliveryIndiaPost{})
	}

	if autoPreload {
		return _db.Preload(clause.Associations)
	} else {

		allowed := utils.FilterAllowedKeySTRArray(preloads, []string{"PaymentSubject", "VatCode"})

		for _, v := range allowed {
			_db.Preload(v)
		}
		return _db
	}

}
