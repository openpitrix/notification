// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import "openpitrix.io/notification/pkg/util/idutil"

type NFAddressList struct {
	NFAddressListId string `gorm:"column:nf_address_list_id"`
	NotificationId  string `gorm:"column:notification_id"`
	AddressListId   string `gorm:"column:address_list_id"`
}

const (
	NFAddressListIdPrefix = "nfa-"
)

//table name
const (
	TableNFAddressList = "nf_address_list"
)

//field name
//NFAddrLs is short for NFAddressList.
const (
	NFAddrLsColId   = "nf_address_list_id"
	NFAddrLsColLsId = "address_list_id"
	NFAddrLsColNfId = "notification_id"
)

func NewNFAddressListId() string {
	return idutil.GetUuid(NFAddressListIdPrefix)
}
