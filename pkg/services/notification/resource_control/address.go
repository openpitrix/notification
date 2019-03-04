// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.
package resource_control

import (
	"context"
	"time"

	"openpitrix.io/logger"
	nfdb "openpitrix.io/notification/pkg/db"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func CreateAddress(ctx context.Context, addr *models.Address) (string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&addr).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Insert address failed, [%+v]", err)
		return "", err
	}
	tx.Commit()
	return addr.AddressId, nil
}
func DescribeAddresses(ctx context.Context, req *pb.DescribeAddressesRequest) ([]*models.Address, uint64, error) {
	req.AddressId = stringutil.SimplifyStringList(req.AddressId)
	req.AddressListId = stringutil.SimplifyStringList(req.AddressListId)
	req.Address = stringutil.SimplifyStringList(req.Address)
	req.NotifyType = stringutil.SimplifyStringList(req.NotifyType)
	req.Status = stringutil.SimplifyStringList(req.Status)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var addrs []*models.Address
	var count uint64

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddress)).
		AddQueryOrderDir(req, models.AddrColCreateTime).
		BuildFilterConditions(req, models.TableAddress).
		Offset(offset).
		Limit(limit).
		Find(&addrs).Error; err != nil {
		logger.Errorf(ctx, "Describe Addresses failed: %+v", err)
		return nil, 0, err
	}

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddress)).
		BuildFilterConditions(req, models.TableAddress).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe Addresses count failed: %+v", err)
		return nil, 0, err
	}

	return addrs, count, nil
}

func ModifyAddress(ctx context.Context, req *pb.ModifyAddressRequest) (string, error) {
	addressId := req.AddressId.GetValue()

	attributes := make(map[string]interface{})

	if req.AddressListId.GetValue() != "" {
		attributes[models.AddrColAddrListId] = req.AddressListId.GetValue()
	}

	if req.Address.GetValue() != "" {
		attributes[models.AddrColAddress] = req.Address.GetValue()
	}
	if req.Remarks.GetValue() != "" {
		attributes[models.AddrColRemarks] = req.Remarks.GetValue()
	}
	if req.VerificationCode.GetValue() != "" {
		attributes[models.AddrColVerificationCode] = req.VerificationCode.GetValue()
	}
	if req.NotifyType.GetValue() != "" {
		attributes[models.AddrColNotifyType] = req.NotifyType.GetValue()
	}

	attributes[models.AddrColStatusTime] = time.Now()

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddress)).
		Where(models.AddrColId+" = ?", addressId).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update Address [%s] failed: %+v", addressId, err)
		return "", err
	}

	return addressId, nil
}

func UpdateAddressAddrLsIdByIds(ctx context.Context, addrListId string, addrIds []string) ([]string, error) {
	attributes := make(map[string]interface{})
	attributes[models.AddrColAddrListId] = addrListId
	attributes[models.AddrColStatusTime] = time.Now()

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddress)).
		Where(models.AddrColId+" in (?)", addrIds).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update Address [%s] failed: %+v", addrIds, err)
		return nil, err
	}

	return addrIds, nil
}

func DeleteAddresses(ctx context.Context, addressIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	tx.Delete(models.Address{}, models.AddrColId+" in (?)", addressIds)
	if err := tx.Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Delete address failed: %+v", err)
		return nil, err
	}
	tx.Commit()
	return addressIds, nil
}
