// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import "time"

type AddressList struct {
	AddressListId string    `gorm:"column:address_list_id"`
	Name          string    `gorm:"column:name"`
	Extra         string    `gorm:"column:extra"`
	Status        string    `gorm:"column:status"`
	CreateTime    time.Time `gorm:"column:create_time"`
	StatusTime    time.Time `gorm:"column:status_time"`
}
