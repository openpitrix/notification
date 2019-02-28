// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.
package resource_control

import (
	"context"
	"time"

	"openpitrix.io/notification/pkg/util/pbutil"

	"openpitrix.io/logger"
	nfdb "openpitrix.io/notification/pkg/db"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func CreateAddressListWithAddrIDs(ctx context.Context, addrList *models.AddressList, addrIds []string) (string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&addrList).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Insert address list failed, [%+v]", err)
		return "", err
	}

	_, err = UpdateAddressAddrLsIdByIds(ctx, addrList.AddressListId, addrIds)
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Update Address AddrLsId ByIds failed, [%+v]", err)
		return "", err
	}
	tx.Commit()
	return addrList.AddressListId, nil
}

func CreateAddressList(ctx context.Context, addrList *models.AddressList) (string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&addrList).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Insert address list failed, [%+v]", err)
		return "", err
	}
	tx.Commit()
	return addrList.AddressListId, nil
}

func DescribeAddressLists(ctx context.Context, req *pb.DescribeAddressListRequest) ([]*models.AddressList, uint64, error) {
	req.AddressListId = stringutil.SimplifyStringList(req.AddressListId)
	req.AddressListName = stringutil.SimplifyStringList(req.AddressListName)
	req.Extra = stringutil.SimplifyStringList(req.Extra)
	req.Status = stringutil.SimplifyStringList(req.Status)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var addressLists []*models.AddressList
	var count uint64

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddressList)).
		AddQueryOrderDir(req, models.AddrColCreateTime).
		BuildFilterConditions(req, models.TableAddressList).
		Offset(offset).
		Limit(limit).
		Find(&addressLists).Error; err != nil {
		logger.Errorf(ctx, "Describe Addresses Lists failed: %+v", err)
		return nil, 0, err
	}

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddressList)).
		BuildFilterConditions(req, models.TableAddressList).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe Addresses Lists count failed: %+v", err)
		return nil, 0, err
	}

	return addressLists, count, nil
}

func ModifyAddressList(ctx context.Context, req *pb.ModifyAddressListRequest) (string, error) {
	addressListId := req.AddressListId.GetValue()
	attributes := make(map[string]interface{})

	if req.AddressListName.GetValue() != "" {
		attributes[models.AddrLsName] = req.AddressListName.GetValue()
	}
	if req.Extra.GetValue() != "" {
		attributes[models.AddrLsExtra] = req.Extra.GetValue()
	}

	attributes[models.AddrColStatusTime] = time.Now()

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddressList)).
		Where(models.AddrLsColId+" = ?", addressListId).
		Updates(attributes).Error; err != nil {
		logger.Errorf(ctx, "Update Address List [%s] failed: %+v", addressListId, err)
		return "", err
	}

	return addressListId, nil
}

func DeleteAddressLists(ctx context.Context, addressListIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	tx.Delete(models.AddressList{}, models.AddrLsColId+" in (?)", addressListIds)
	if err := tx.Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Delete Address Lists failed: %+v", err)
		return nil, err
	}
	tx.Commit()
	return addressListIds, nil
}
