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

func CreateAddressList(ctx context.Context, addrList *models.AddressList) (string, error) {
	err := global.GetInstance().GetDB().Create(&addrList).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to insert address_list, %+v.", err)
		return "", err
	}
	return addrList.AddressListId, nil
}

func CreateAddressListWithAddrIDs(ctx context.Context, addrList *models.AddressList, addrIds []string) (string, error) {
	tx := global.GetInstance().GetDB().Begin()
	err := tx.Create(&addrList).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to insert address_list, %+v.", err)
		return "", err
	}

	for _, addrId := range addrIds {
		addressListBinding := &models.AddressListBinding{
			BindingId:     models.NewAddressListBindingId(),
			AddressListId: addrList.AddressListId,
			AddressId:     addrId,
			CreateTime:    time.Now(),
		}
		err := tx.Create(&addressListBinding).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "Failed to insert address_list_binding, %+v.", err)
			return "", err
		}
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
		BuildFilterConditions(req, models.TableAddressList, "and").
		Offset(offset).
		Limit(limit).
		Find(&addressLists).Error; err != nil {
		logger.Errorf(ctx, "Describe addresses list failed,[%+v]", err)
		return nil, 0, err
	}

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableAddressList)).
		BuildFilterConditions(req, models.TableAddressList, "and").
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe addresses list count failed, %+v.", err)
		return nil, 0, err
	}

	return addressLists, count, nil
}

func ModifyAddressList(ctx context.Context, addressListId string, attributes map[string]interface{}) error {
	tx := global.GetInstance().GetDB().Begin()
	addrs, err := GetAddressesListByIds(ctx, tx, []string{addressListId})
	if len(addrs) == 0 {
		tx.Rollback()
		err := gerr.New(ctx, gerr.NotFound, gerr.ErrorResourceNotExist, addressListId)
		logger.Errorf(ctx, "Failed to update address list[%s],address list does not exits, %+v.", addressListId, err)
		return err
	}

	err = tx.Table(models.TableAddressList).Where(models.AddrLsColId+" = ?", addressListId).Updates(attributes).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to update address [%s], %+v.", addressListId, err)
		return err
	}
	tx.Commit()
	return nil
}

func GetAddressesListByIds(ctx context.Context, tx *gorm.DB, addressListIds []string) ([]*models.AddressList, error) {
	var addrLists []*models.AddressList
	err := tx.Table(models.TableAddressList).
		Select("*").
		Where(models.AddrLsColId+" in ( ? )", addressListIds).
		Set("gorm:query_option", "FOR UPDATE").
		Scan(&addrLists).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to get address list by ids [%+v], %+v.", addressListIds, err)
		return nil, err
	}
	return addrLists, nil
}

func ModifyAddressListWithAddrIDs(ctx context.Context, addressListId string, attributes map[string]interface{}, addrIds []string) error {
	var err error
	tx := global.GetInstance().GetDB().Begin()
	addressListIds := []string{addressListId}
	addrLists, err := GetAddressListsByIds(ctx, tx, addressListIds)
	addrIds = stringutil.Unique(addrIds)
	if len(addrLists) == 0 {
		tx.Rollback()
		err := gerr.New(ctx, gerr.NotFound, gerr.ErrorResourceNotExist, strings.Trim(fmt.Sprint(addressListIds), "[]"))
		logger.Errorf(ctx, "Failed to update address_list[%s],address_list does not exits, %+v.", addressListIds, err)
		return err
	}

	err = tx.Table(models.TableAddressList).Where(models.AddrLsColId+" = ?", addressListId).Updates(attributes).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to update address_list[%s], %+v.", addressListId, err)
		return err
	}

	err = tx.Delete(models.AddressListBinding{}, models.BindColAddrListId+" in (?)", addressListId).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to delete address_list_binding by address_list_id[%s], %+v.", addressListId, err)
		return err
	}

	for _, addrId := range addrIds {
		addressListBinding := &models.AddressListBinding{
			BindingId:     models.NewAddressListBindingId(),
			AddressListId: addressListId,
			AddressId:     addrId,
			CreateTime:    time.Now(),
		}
		err := tx.Create(&addressListBinding).Error
		if err != nil {
			tx.Rollback()
			logger.Errorf(ctx, "Failed to insert address_list_binding, %+v.", err)
			return err
		}
	}
	tx.Commit()
	return nil
}

func DeleteAddressLists(ctx context.Context, addressListIds []string) error {
	tx := global.GetInstance().GetDB().Begin()

	err := tx.Table(models.TableAddressList).Where(models.AddrLsColId+" in (?)", addressListIds).Updates(map[string]interface{}{models.AddrLsColStatus: constants.StatusDeleted, models.AddrLsColStatusTime: time.Now()}).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to update address_list status to deleted, %+v.", err)
		return err
	}

	//at the same,need to delete all this address_list_id rows in address_list_binding.
	err = tx.Delete(models.AddressListBinding{}, models.BindColAddrListId+" in (?)", addressListIds).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to delete address_list_binding by address_list_ids [%+v], %+v.", addressListIds, err)
		return err
	}
	tx.Commit()

	return nil
}

func GetAddressListsByIds(ctx context.Context, tx *gorm.DB, addressListIds []string) ([]*models.AddressList, error) {
	var addrLists []*models.AddressList
	err := tx.Table(models.TableAddressList).
		Select("*").
		Where(models.AddrLsColId+" in ( ? )", addressListIds).
		Set("gorm:query_option", "FOR UPDATE").
		Scan(&addrLists).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to get address_list by ids [%+v], %+v.", addrLists, err)
		return nil, err
	}
	return addrLists, nil
}

func GetActiveAddressesListsByIds(ctx context.Context, addressListIds []string) ([]*models.AddressList, error) {
	var addrLists []*models.AddressList
	db := global.GetInstance().GetDB()

	err := db.Table(models.TableAddressList).
		Select("*").
		Where(models.AddrLsColStatus+" in ( '"+constants.StatusActive+"' )").
		Where(models.AddrLsColId+" in ( ? )", addressListIds).
		Scan(&addrLists).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to get active address list by ids[%+v], %+v.", addrLists, err)
		return nil, err
	}
	return addrLists, nil
}

func GetDeletedAddressListsByIds(ctx context.Context, tx *gorm.DB, addressListIds []string) ([]*models.AddressList, error) {
	var addrLists []*models.AddressList
	err := tx.Table(models.TableAddressList).
		Select("*").
		Where(models.AddrLsColId+" in ( ? )", addressListIds).
		Where(models.AddrColStatus+" in ( '"+constants.StatusDeleted+"' )").
		Set("gorm:query_option", "FOR UPDATE").
		Scan(&addrLists).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to get address_list by ids [%+v], %+v.", addrLists, err)
		return nil, err
	}
	return addrLists, nil
}
