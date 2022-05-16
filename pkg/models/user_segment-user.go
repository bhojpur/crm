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
type UsersSegmentUser struct {
	UserId         uint `json:"user_id" gorm:"type:int;index;not null;"`
	UsersSegmentId uint `json:"users_segment_id" gorm:"type:int;index;not null;"`
}

func (UsersSegmentUser) PgSqlCreate() {
	if err := db.Migrator().CreateTable(&UsersSegmentUser{}); err != nil {
		log.Fatal(err)
	}
	err := db.Exec("ALTER TABLE users_segment_users \n    ADD CONSTRAINT users_segment_users_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,\n    ADD CONSTRAINT users_segment_users_users_segment_id_fkey FOREIGN KEY (users_segment_id) REFERENCES users_segments(id) ON DELETE CASCADE ON UPDATE CASCADE,\n    DROP CONSTRAINT IF EXISTS fk_users_segment_users_users,\n    DROP CONSTRAINT IF EXISTS fk_users_segment_users_user_segments;").Error
	if err != nil {
		log.Fatal("Error: ", err)
	}

}

func (UsersSegmentUser) TableName() string {
	return "users_segment_users"
}
