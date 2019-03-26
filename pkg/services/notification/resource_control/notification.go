// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource_control

import (
	"context"
	"time"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	nfdb "openpitrix.io/notification/pkg/db"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/jsonutil"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func RegisterNotification(ctx context.Context, notification *models.Notification) error {
	tx := global.GetInstance().GetDB().Begin()
	addressInfo := notification.AddressInfo
	_, err := models.DecodeAddressInfo(addressInfo)
	if err != nil {
		addressListIds, err := models.DecodeAddressListIds(addressInfo)
		if err != nil {
			return err
		}
		notification.AddressInfo = "{}"

		for _, listId := range []string(*addressListIds) {
			nfAddressList := &models.NFAddressList{
				NFAddressListId: models.NewNFAddressListId(),
				NotificationId:  notification.NotificationId,
				AddressListId:   listId,
			}
			err := tx.Create(&nfAddressList).Error
			if err != nil {
				tx.Rollback()
				logger.Errorf(ctx, "Failed to insert nf_address_list, %+v.", err)
				return err
			}
		}
	}

	err = tx.Create(&notification).Error
	notification.AddressInfo = addressInfo
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Failed to insert notification, %+v.", err)
		return err
	}
	tx.Commit()
	return nil
}

func UpdateNotificationsStatus(ctx context.Context, nfIds []string, status string) error {
	db := global.GetInstance().GetDB()
	err := db.Table(models.TableNotification).Where(models.NfColId+" in (?)", nfIds).Updates(map[string]interface{}{models.NfColStatus: status, models.NfColStatusTime: time.Now()}).Error

	if err != nil {
		logger.Errorf(ctx, "Failed to update notification [%+v] status to [%s] failed, %+v.", nfIds, status, err)
		return err
	}
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

	err := nfdb.GetChain(global.GetInstance().GetDB().
		Table(models.TableNotification)).
		AddQueryOrderDir(req, models.NfColCreateTime).
		BuildFilterConditions(req, models.TableNotification).
		Offset(offset).
		Limit(limit).
		Find(&nfs).Error

	if err != nil {
		logger.Errorf(ctx, "Failed to describe notification, %+v.", err)
		return nil, 0, err
	}

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableNotification)).
		BuildFilterConditions(req, models.TableNotification).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Failed to describe notification count, %+v.", err)
		return nil, 0, err
	}

	return nfs, count, nil
}

func GetNfsByNfIds(ctx context.Context, nfIds []string) ([]*models.Notification, error) {
	db := global.GetInstance().GetDB()
	var nfs []*models.Notification
	err := db.Where("notification_id in( ? )", nfIds).Find(&nfs).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to get notifications by ids [%+v], %+v.", nfIds, err)
		return nil, err
	}
	return nfs, nil
}

func SplitNotificationIntoTasks(ctx context.Context, notification *models.Notification) ([]*models.Task, error) {
	addressInfo, err := models.DecodeAddressInfo(notification.AddressInfo)
	if err != nil {
		addressListIds, err := models.DecodeAddressListIds(notification.AddressInfo)
		if err != nil {
			return nil, err
		}
		addresses, err := GetAddressesByListIds(ctx, []string(*addressListIds))
		if err != nil {
			return nil, err
		}
		var tasks []*models.Task
		for _, address := range addresses {
			directive := &models.TaskDirective{
				Address:            address.Address,
				NotifyType:         constants.NotifyTypeEmail,
				ContentType:        notification.ContentType,
				Title:              notification.Title,
				Content:            notification.Content,
				ShortContent:       notification.ShortContent,
				ExpiredDays:        notification.ExpiredDays,
				AvailableStartTime: notification.AvailableStartTime,
				AvailableEndTime:   notification.AvailableEndTime,
			}
			task := models.NewTask(
				notification.NotificationId,
				jsonutil.ToString(directive),
			)
			logger.Debugf(ctx, "Split notification into task[%s] successfully. ", task.TaskId)
			tasks = append(tasks, task)
		}
		return tasks, nil
	}

	var tasks []*models.Task
	for notifyType, addresses := range *addressInfo {
		for _, address := range addresses {
			directive := &models.TaskDirective{
				Address:            address,
				NotifyType:         notifyType,
				ContentType:        notification.ContentType,
				Title:              notification.Title,
				Content:            notification.Content,
				ShortContent:       notification.ShortContent,
				ExpiredDays:        notification.ExpiredDays,
				AvailableStartTime: notification.AvailableStartTime,
				AvailableEndTime:   notification.AvailableEndTime,
			}
			task := models.NewTask(
				notification.NotificationId,
				jsonutil.ToString(directive),
			)
			logger.Debugf(ctx, "Split notification into task[%s] successfully. ", task.TaskId)
			tasks = append(tasks, task)

		}
	}
	return tasks, nil
}
