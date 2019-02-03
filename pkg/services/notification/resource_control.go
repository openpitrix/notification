// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
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

func RegisterTask(ctx context.Context, task *models.Task) error {
	db := globalcfg.GetInstance().GetDB()
	tx := db.Begin()
	err := tx.Create(&task).Error
	if err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Insert task failed, [%+v]", err)
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

func UpdateTaskStatus(taskId, status string) error {
	db := globalcfg.GetInstance().GetDB()
	task := &models.Task{
		TaskId: taskId,
	}
	tx := db.Begin()
	err := db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Error
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func GetStatusTasks(notificationId string, status []string) []*models.Task {
	db := globalcfg.GetInstance().GetDB()
	var tasks []*models.Task
	tx := db.Begin()
	db.Where("notification_id = ? AND status in (?)", notificationId, status).Find(&tasks)
	tx.Commit()
	return tasks
}

func GetTask(taskId string) (*models.Task, error) {
	db := globalcfg.GetInstance().GetDB()
	task := new(models.Task)
	err := db.Where("task_id = ?", taskId).First(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
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

func GetEmailServiceConfig(ctx context.Context) *pb.EmailServiceConfig {
	mycfg := config.GetInstance()

	protocol := mycfg.Email.Protocol
	emailHost := mycfg.Email.EmailHost
	port := mycfg.Email.Port
	displayEmail := mycfg.Email.DisplayEmail
	email := mycfg.Email.Email
	password := mycfg.Email.Password
	sslEnable := mycfg.Email.SSLEnable

	emailCfg := &pb.EmailServiceConfig{
		Protocol:     pbutil.ToProtoString(protocol),
		EmailHost:    pbutil.ToProtoString(emailHost),
		Port:         pbutil.ToProtoString(string(port)),
		DisplayEmail: pbutil.ToProtoString(displayEmail),
		Email:        pbutil.ToProtoString(email),
		Password:     pbutil.ToProtoString(password),
		SslEnable:    pbutil.ToProtoBool(sslEnable),
	}

	return emailCfg
}
