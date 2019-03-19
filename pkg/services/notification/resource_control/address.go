// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.
package resource_control

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	nfdb "openpitrix.io/notification/pkg/db"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func CreateAddress(ctx context.Context, addr *models.Address) (string, error) {
	db := global.GetInstance().GetDB()
	err := db.Create(&addr).Error
	if err != nil {
		logger.Errorf(ctx, "Insert into address failed, %+v.", err)
		return "", err
	}
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
		logger.Errorf(ctx, "Describe address failed, %+v.", err)
		return nil, 0, err
	}

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddress)).
		BuildFilterConditions(req, models.TableAddress).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe address count failed, %+v.", err)
		return nil, 0, err
	}

	return addrs, count, nil
}

func ModifyAddress(ctx context.Context, addressId string, attributes map[string]interface{}) error {
	tx := global.GetInstance().GetDB().Begin()
	addrs, err := GetAddressesByIds(ctx, tx, []string{addressId})
	if len(addrs) == 0 {
		tx.Rollback()
		err := gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorAddressNotExist, addressId)
		logger.Errorf(ctx, "Failed to update address [%s],address does not exits, %+v.", addressId, err)
		return err
	}

	err = tx.Table(models.TableAddress).Where(models.AddrColId+" = ?", addressId).Updates(attributes).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to update address [%s], %+v.", addressId, err)
		return err
	}
	tx.Commit()
	return nil
}

func DeleteAddresses(ctx context.Context, addressIds []string) error {
	tx := global.GetInstance().GetDB().Begin()
	addrs, err := GetAddressesByIds(ctx, tx, addressIds)
	if len(addrs) != len(addressIds) {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to delete address [%s], address does not exits, %+v.", addressIds, err)
		err := gerr.NewWithDetail(ctx, gerr.NotFound, err, gerr.ErrorAddressNotExist, strings.Trim(fmt.Sprint(addressIds), "[]"))
		return err
	}

	err = tx.Table(models.TableAddress).Where(models.AddrColId+" in (?)", addressIds).Updates(map[string]interface{}{models.AddrColStatus: constants.StatusDeleted, models.AddrColStatusTime: time.Now()}).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to update address status to deleted, %+v.", err)
		return err
	}

	//at the same,need to delete all this address rows in address_list_binding.
	err = tx.Delete(models.AddressListBinding{}, models.BindColAddrId+" in (?)", addressIds).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to delete address_list_binding by address_list_ids [%+v], %+v.", addressIds, err)
		return err
	}
	tx.Commit()

	return nil
}

func GetAddressesByListIds(ctx context.Context, listIds []string) ([]*models.Address, error) {
	var addrs []*models.Address
	db := global.GetInstance().GetDB()

	err := db.Table("address_list_binding t1").
		Select("t2.address_id,t2.address,t2.remarks,t2.verification_code,t2.create_time,t2.verify_time,t2.status_time,t2.notify_type,t2.status").
		Joins(" join address t2 on t1.address_id=t2.address_id").
		Joins(" join address_list t3 on t1.address_list_id=t3.address_list_id").
		Where(" t3.address_list_id in ( ? )", listIds).
		Scan(&addrs).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to get address by list ids [%+v], %+v.", listIds, err)
		return nil, err
	}

	return addrs, nil
}

func GetAddressesByIds(ctx context.Context, tx *gorm.DB, addressIds []string) ([]*models.Address, error) {
	var addrs []*models.Address
	err := tx.Table(models.TableAddress).
		Select("*").
		Where(models.AddrColId+" in ( ? )", addressIds).
		Set("gorm:query_option", "FOR UPDATE").
		Scan(&addrs).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to get address by ids [%+v], %+v.", addressIds, err)
		return nil, err
	}
	return addrs, nil
}
