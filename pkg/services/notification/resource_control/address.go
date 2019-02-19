// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.
package resource_control

import (
	"context"
	"time"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/dbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func CreateAddress(ctx context.Context, addr *models.Address) (string, error) {
	db := globalcfg.GetInstance().GetDB()
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
	req.Status = stringutil.SimplifyStringList(req.NotifyType)

	limit := dbutil.GetLimit(req.Limit)
	offset := dbutil.GetOffset(req.Offset)

	var nfs []*models.Address
	var count uint64

	if err := dbutil.GetChain(globalcfg.GetInstance().GetDB().Table(models.TableAddress)).
		AddQueryOrderDir(req, models.AddrColCreateTime).
		BuildFilterConditions(req, models.TableAddress).
		Offset(offset).
		Limit(limit).
		Find(&nfs).Error; err != nil {
		logger.Errorf(ctx, "Describe Addresses failed: %+v", err)
		return nil, 0, err
	}

	if err := dbutil.GetChain(globalcfg.GetInstance().GetDB().Table(models.TableAddress)).
		BuildFilterConditions(req, models.TableAddress).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe Addresses count failed: %+v", err)
		return nil, 0, err
	}

	return nfs, count, nil
}

func ModifyAddress(ctx context.Context, req *pb.ModifyAddressRequest) (string, error) {
	addressId := req.AddressId.GetValue()

	attributes := make(map[string]interface{})

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

	if err := dbutil.GetChain(globalcfg.GetInstance().GetDB().Table(models.TableAddress)).
		Where(models.AddrColId+" = ?", addressId).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update Address [%s] failed: %+v", addressId, err)
		return "", err
	}

	return addressId, nil
}
