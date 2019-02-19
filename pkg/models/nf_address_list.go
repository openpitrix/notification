// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

type NFAddressList struct {
	NFAddressListId string `gorm:"column:nf_address_list_id"`
	NotificationId  string `gorm:"column:notification_id"`
	AddressListId   string `gorm:"column:address_list_id"`
}
