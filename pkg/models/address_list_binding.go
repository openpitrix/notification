// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/notification/pkg/util/idutil"
)

type AddressListBinding struct {
	BindingId     string    `gorm:"column:binding_id"`
	AddressListId string    `gorm:"column:address_list_id"`
	AddressId     string    `gorm:"column:address_id"`
	CreateTime    time.Time `gorm:"column:create_time"`
}

const (
	AddressListBindingIdPrefix = "bid-"
)

//table name
const (
	TableAddressListBinding = "address_list_binding"
)

//field name
//Bind is short for AddressListBinding.
const (
	BindColId         = "binding_id"
	BindColAddrListId = "address_list_id"
	BindColAddrId     = "address_id"
	BindColCreateTime = "create_time"
)

func NewAddressListBindingId() string {
	return idutil.GetUuid(AddressListBindingIdPrefix)
}
