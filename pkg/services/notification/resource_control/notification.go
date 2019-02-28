// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource_control

import (
	"context"
	"time"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/util/jsonutil"
	"openpitrix.io/notification/pkg/util/pbutil"

	"openpitrix.io/logger"
	nfdb "openpitrix.io/notification/pkg/db"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
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

	db := global.GetInstance().GetDB()
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
	db := global.GetInstance().GetDB()
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

func DescribeNotifications(ctx context.Context, req *pb.DescribeNotificationsRequest) ([]*models.Notification, uint64, error) {
	req.NotificationId = stringutil.SimplifyStringList(req.NotificationId)
	req.ContentType = stringutil.SimplifyStringList(req.ContentType)
	req.Owner = stringutil.SimplifyStringList(req.Owner)
	req.Status = stringutil.SimplifyStringList(req.Status)

	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var nfs []*models.Notification
	var count uint64

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableNotification)).
		AddQueryOrderDir(req, models.NfColCreateTime).
		BuildFilterConditions(req, models.TableNotification).
		Offset(offset).
		Limit(limit).
		Find(&nfs).Error; err != nil {
		logger.Errorf(ctx, "Describe Notifications failed: %+v", err)
		return nil, 0, err
	}

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableNotification)).
		BuildFilterConditions(req, models.TableNotification).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe Notifications count failed: %+v", err)
		return nil, 0, err
	}

	return nfs, count, nil
}

func UpdateNotifications2Pending(ctx context.Context, nfIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	db.Table(models.TableNotification).Where(models.NfColId+" in (?)", nfIds).Updates(map[string]interface{}{models.NfColStatus: constants.StatusPending, models.NfColStatusTime: time.Now()})

	if err := tx.Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Update Notifications Status to Pending failed: [%+v].", err)
		return nil, err
	}

	tx.Commit()
	return nfIds, nil
}

func SplitNotificationIntoTasks(notification *models.Notification) ([]*models.Task, error) {
	addressInfo, err := models.DecodeAddressInfo(notification.AddressInfo)
	if err != nil {
		return nil, err
	}

	var tasks []*models.Task
	for notifyType, addresses := range *addressInfo {
		for _, address := range addresses {
			directive := &models.TaskDirective{
				Address:      address,
				NotifyType:   notifyType,
				ContentType:  notification.ContentType,
				Title:        notification.Title,
				Content:      notification.Content,
				ShortContent: notification.ShortContent,
				ExpiredDays:  notification.ExpiredDays,
			}
			task := models.NewTask(
				notification.NotificationId,
				jsonutil.ToString(directive),
			)
			logger.Debugf(nil, "Split notification into tasks[%s] successfully. ", task.TaskId)
			tasks = append(tasks, task)

		}
	}
	return tasks, nil
}
