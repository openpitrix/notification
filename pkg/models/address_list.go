// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type AddressList struct {
	AddressListId   string    `gorm:"column:address_list_id"`
	AddressListName string    `gorm:"column:address_list_name"`
	Extra           string    `gorm:"column:extra"`
	Status          string    `gorm:"column:status"`
	CreateTime      time.Time `gorm:"column:create_time"`
	StatusTime      time.Time `gorm:"column:status_time"`
}

const (
	AddressListIdPrefix = "adl-"
)

//table name
const (
	TableAddressList = "address_list"
)

//field name
//AddrLs is short for AddressList.
const (
	AddrLsColId         = "address_list_id"
	AddrLsColName       = "address_list_name"
	AddrLsColExtra      = "extra"
	AddrLsColStatus     = "status"
	AddrLsColCreateTime = "create_time"
	AddrLsColStatusTime = "status_time"
)

func NewAddressListId() string {
	return idutil.GetUuid(AddressListIdPrefix)
}

func NewAddressList(req *pb.CreateAddressListRequest) *AddressList {
	addressList := &AddressList{
		AddressListId:   NewAddressListId(),
		AddressListName: req.GetAddressListName().GetValue(),
		Extra:           req.GetExtra().GetValue(),
		Status:          constants.StatusActive,
		CreateTime:      time.Now(),
		StatusTime:      time.Now(),
	}
	return addressList
}

func AddressListToPb(addressList *AddressList) *pb.AddressList {
	pbAddressList := pb.AddressList{}
	pbAddressList.AddressListId = pbutil.ToProtoString(addressList.AddressListId)
	pbAddressList.AddressListName = pbutil.ToProtoString(addressList.AddressListName)
	pbAddressList.Extra = pbutil.ToProtoString(addressList.Extra)
	pbAddressList.Status = pbutil.ToProtoString(addressList.Status)
	pbAddressList.StatusTime = pbutil.ToProtoTimestamp(addressList.StatusTime)
	pbAddressList.CreateTime = pbutil.ToProtoTimestamp(addressList.CreateTime)

	return &pbAddressList
}

func AddressListSet2PbSet(inAddrLists []*AddressList) []*pb.AddressList {
	var pbAddrsLists []*pb.AddressList
	for _, inAddrList := range inAddrLists {
		pbAddrList := AddressListToPb(inAddrList)
		pbAddrsLists = append(pbAddrsLists, pbAddrList)
	}
	return pbAddrsLists
}
