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
	"fmt"
	"time"

	"gorm.io/gorm/logger"

	// "github.com/go-gorm/gorm"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// _ "github.com/go-sql-driver/mysql"
	// _ "github.com/lib/pq"
	_ "github.com/jackc/pgx/v4"

	"os"
)

/*func init() {
	connectDb()
}*/

var db *gorm.DB

func GetDB() *gorm.DB {

	if db == nil {
		db = ConnectDb()
	}

	return db
}

func SetDB(p *gorm.DB) {
	db = p
}

func GetPool() *gorm.DB {
	return db
}

func ConnectDb() *gorm.DB {

	if db != nil {
		return nil
	}

	if os.Getenv("ENV_VAR") == "test" {
		e := godotenv.Load("/Volumes/Software/bhojpur/net/crm/.env.test")
		if e != nil {
			log.Fatal("Error loading test .env file")
		}
	} else {
		e := godotenv.Load()
		if e != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	// https://github.com/go-gorm/postgres
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Millisecond * 200, // Slow SQL threshold
			LogLevel:      logger.Error,           // ???????????? ?????????????????????? GORM: Silent, Error, Warn, Info
			Colorful:      true,                   // Disable color
		},
	)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: os.Getenv("DATABASE_URL"),
		// PreferSimpleProtocol: true, // disables implicit prepared statement usage
		PreferSimpleProtocol: false, // disables implicit prepared statement usage
	}), &gorm.Config{
		Logger:                 dbLogger,
		SkipDefaultTransaction: false, // ?????????????? ???????????????????? ????????????
		PrepareStmt:            false, // ?????????? ???????????????????????????? ???????????????????? ?????????????? preparedstmt ?? ???????????????? ???? ?????? ?????????????????? ?????????????? ??????????????
		// DisableForeignKeyConstraintWhenMigrating: true,
		/*NamingStrategy: schema.NamingStrategy{
			TablePrefix: "t_",   // ?????????????? ???????? ????????????, ?????????????? ?????? `User` ?????????? `t_users`
			SingularTable: true, // ???????????????????????? ???????????????????? ?? ???????????????????????? ??????????, ?????????????? ?????? `User` ?????????? `user` ?????? ?????????????????? ???????? ??????????, ?????? `t_user` ?????? TablePrefix = "t_"
		},*/
	})

	if err != nil {
		log.Fatal("Error connect to DB: ", err)
	}

	if db == nil {
		log.Fatal("Error connect to DB == nil")
	}

	// db.DB().SetConnMaxLifetime(0)
	// db.DB().SetMaxIdleConns(10)
	// db.DB().SetMaxOpenConns(100)

	//db = db.Set("gorm:auto_preload", true)

	// db = db.LogMode(true)

	/*if err := SettingsDb(); err != nil {
		log.Println("???????????? ???????????????? ????", err)
	}*/
	SetDB(db)
	fmt.Println("DataBase init full!")
	return db
}

func SettingsDb() error {

	err := db.SetupJoinTable(&ProductCard{}, "Products", &ProductCardProduct{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&Product{}, "ProductCards", &ProductCardProduct{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.SetupJoinTable(&Account{}, "Users", &AccountUser{}); err != nil {
		log.Fatal(err)
	}
	if err := db.SetupJoinTable(&User{}, "Accounts", &AccountUser{}); err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&Warehouse{}, "Products", &WarehouseItem{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&Product{}, "Warehouses", &WarehouseItem{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&Shipment{}, "Products", &ShipmentItem{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&Product{}, "Shipments", &ShipmentItem{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&User{}, "Companies", &CompanyUser{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&Company{}, "Users", &CompanyUser{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&WebPage{}, "ProductCategories", &WebPageProductCategory{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&ProductCategory{}, "WebPages", &WebPageProductCategory{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&Inventory{}, "Products", &InventoryItem{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&Product{}, "Inventories", &InventoryItem{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&ProductCategory{}, "Products", &ProductCategoryProduct{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&Product{}, "ProductCategories", &ProductCategoryProduct{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&Product{}, "ProductTags", &ProductTagProduct{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&ProductTag{}, "Products", &ProductTagProduct{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.SetupJoinTable(&User{}, "UsersSegments", &UsersSegmentUser{})
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupJoinTable(&UsersSegment{}, "Users", &UsersSegmentUser{})
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
