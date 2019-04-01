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

type Address struct {
	AddressId        string    `gorm:"column:address_id"`
	Address          string    `gorm:"column:address"`
	Remarks          string    `gorm:"column:remarks"`
	VerificationCode string    `gorm:"column:verification_code"`
	CreateTime       time.Time `gorm:"column:create_time"`
	VerifyTime       time.Time `gorm:"column:verify_time"`
	StatusTime       time.Time `gorm:"column:status_time"`
	NotifyType       string    `gorm:"column:notify_type"`
	Status           string    `gorm:"column:status"`
}

type AddressWithListId struct {
	AddressId        string    `gorm:"column:address_id"`
	AddressListId    string    `gorm:"column:address_list_id"`
	Address          string    `gorm:"column:address"`
	Remarks          string    `gorm:"column:remarks"`
	VerificationCode string    `gorm:"column:verification_code"`
	CreateTime       time.Time `gorm:"column:create_time"`
	VerifyTime       time.Time `gorm:"column:verify_time"`
	StatusTime       time.Time `gorm:"column:status_time"`
	NotifyType       string    `gorm:"column:notify_type"`
	Status           string    `gorm:"column:status"`
}

const (
	AddressIdPrefix = "addr-"
)

//table name
const (
	TableAddress = "address"
)

//field name
//Addr is short for Address.
const (
	AddrColId               = "address_id"
	AddrColAddress          = "address"
	AddrColRemarks          = "remarks"
	AddrColVerificationCode = "verification_code"
	AddrColCreateTime       = "create_time"
	AddrColVerifyTime       = "verify_time"
	AddrColStatusTime       = "status_time"
	AddrColNotifyType       = "notify_type"
	AddrColStatus           = "status"
)

func NewAddressId() string {
	return idutil.GetUuid(AddressIdPrefix)
}

func NewAddress(req *pb.CreateAddressRequest) *Address {
	address := &Address{
		AddressId:        NewAddressId(),
		Address:          req.GetAddress().GetValue(),
		Remarks:          req.GetRemarks().GetValue(),
		VerificationCode: req.GetVerificationCode().GetValue(),
		CreateTime:       time.Now(),
		VerifyTime:       time.Now(),
		StatusTime:       time.Now(),
		NotifyType:       req.GetNotifyType().GetValue(),
		Status:           constants.StatusActive,
	}
	return address
}

func AddressWithListIdToPb(addressWithListId *AddressWithListId) *pb.Address {
	pbAddress := pb.Address{}
	pbAddress.AddressId = pbutil.ToProtoString(addressWithListId.AddressId)
	pbAddress.AddressListId = pbutil.ToProtoString(addressWithListId.AddressListId)
	pbAddress.Address = pbutil.ToProtoString(addressWithListId.Address)
	pbAddress.Remarks = pbutil.ToProtoString(addressWithListId.Remarks)
	pbAddress.VerificationCode = pbutil.ToProtoString(addressWithListId.VerificationCode)
	pbAddress.Status = pbutil.ToProtoString(addressWithListId.Status)
	pbAddress.CreateTime = pbutil.ToProtoTimestamp(addressWithListId.CreateTime)
	pbAddress.VerifyTime = pbutil.ToProtoTimestamp(addressWithListId.VerifyTime)
	pbAddress.StatusTime = pbutil.ToProtoTimestamp(addressWithListId.StatusTime)
	pbAddress.NotifyType = pbutil.ToProtoString(addressWithListId.NotifyType)
	return &pbAddress
}

func AddressWithListIdSet2PbSet(inAddrs []*AddressWithListId) []*pb.Address {
	var pbAddrs []*pb.Address
	for _, inAddr := range inAddrs {
		pbAddr := AddressWithListIdToPb(inAddr)
		pbAddrs = append(pbAddrs, pbAddr)
	}
	return pbAddrs
}
