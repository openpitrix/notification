// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource_control

import (
	"context"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/dbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func RegisterNotification(ctx context.Context, notification *models.Notification) error {
	addressInfo := notification.AddressInfo
	_, err := models.DecodeAddressInfo(addressInfo)
	if err != nil {
		_, err := models.DecodeAddressListIds(addressInfo)
		if err != nil {
			return err
		}
		notification.AddressInfo = ""
		// TODO: register nf_address_list
	}

	db := globalcfg.GetInstance().GetDB()
	tx := db.Begin()
	err = tx.Create(&notification).Error
	notification.AddressInfo = addressInfo
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Insert notification failed, [%+v]", err)
		return err
	}
	tx.Commit()
	return nil
}

func UpdateNotificationStatus(notificationId, status string) error {
	db := globalcfg.GetInstance().GetDB()
	nf := &models.Notification{
		NotificationId: notificationId,
	}
	tx := db.Begin()
	err := db.Model(&nf).Where("notification_id = ?", notificationId).Update("status", status).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func GetNotification(notificationId string) (*models.Notification, error) {
	db := globalcfg.GetInstance().GetDB()
	nf := new(models.Notification)
	err := db.Where("notification_id = ?", notificationId).First(nf).Error
	if err != nil {
		return nil, err
	}
	return nf, nil
}

func DescribeNotifications(ctx context.Context, req *pb.DescribeNotificationsRequest) ([]*models.Notification, uint64, error) {
	req.NotificationId = stringutil.SimplifyStringList(req.NotificationId)
	req.ContentType = stringutil.SimplifyStringList(req.ContentType)
	req.Owner = stringutil.SimplifyStringList(req.Owner)
	req.Status = stringutil.SimplifyStringList(req.Status)

	limit := dbutil.GetLimit(req.Limit)
	offset := dbutil.GetOffset(req.Offset)

	var nfs []*models.Notification
	var count uint64

	if err := dbutil.GetChain(globalcfg.GetInstance().GetDB().Table(models.TableNotification)).
		AddQueryOrderDir(req, models.NfColCreateTime).
		BuildFilterConditions(req, models.TableNotification).
		Offset(offset).
		Limit(limit).
		Find(&nfs).Error; err != nil {
		logger.Errorf(ctx, "Describe Notifications failed: %+v", err)
		return nil, 0, err
	}

	if err := dbutil.GetChain(globalcfg.GetInstance().GetDB().Table(models.TableNotification)).
		BuildFilterConditions(req, models.TableNotification).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe Notifications count failed: %+v", err)
		return nil, 0, err
	}

	return nfs, count, nil
}
