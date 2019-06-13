// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource_control

import (
	"context"
	"time"

	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	nfdb "openpitrix.io/notification/pkg/db"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/services/websocket"
	wstypes "openpitrix.io/notification/pkg/services/websocket/types"
	"openpitrix.io/notification/pkg/util/jsonutil"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func RegisterNotification(ctx context.Context, notification *models.Notification) error {
	tx := global.GetInstance().GetDB().Begin()
	addressInfo := notification.AddressInfo
	//Step1: check addressInfo format is like address_info = {"email": ["xxx@abc.com", "xxx@xxx.com"],"websocket": ["system", "huojiao"]}
	_, err := models.DecodeAddressInfo(addressInfo)

	//Step2: check addressInfo format is like address_info = ["adl-xxxx1", "adl-xxxx2"]
	if err != nil {
		addressListIds, err := models.DecodeAddressListIds(addressInfo)
		if err != nil {
			return err
		}

		//Step3: insert data into nf_address_list
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

	//Step4: insert data into notification
	err = tx.Create(&notification).Error
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

	//Step1: check addressInfo format is like address_info = {"email": ["xxx@abc.com", "xxx@xxx.com"],"websocket": ["system", "huojiao"]}
	_, decodeMapErr := models.DecodeAddressInfo(notification.AddressInfo)

	//Step2: check addressInfo format is like address_info = ["adl-xxxx1", "adl-xxxx2"]
	if decodeMapErr == nil {
		tasks, err := processsAddressInfo4AddressMap(ctx, notification)
		if err != nil {
			return nil, err
		} else {
			return tasks, nil
		}
	} else {
		tasks, err := processsAddressInfo4AddressListIds(ctx, notification)
		if err != nil {
			return nil, err
		} else {
			return tasks, nil
		}

	}
}

//address_info = ["adl-xxxx1", "adl-xxxx2"]
func processsAddressInfo4AddressListIds(ctx context.Context, notification *models.Notification) ([]*models.Task, error) {
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
			NotificationId:     notification.NotificationId,
			Address:            address.Address,
			NotifyType:         address.NotifyType,
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
			address.NotifyType,
		)
		//if websocket message,just push to ws pubsub queue,no need to create task record in DB.
		if address.NotifyType == constants.NotifyTypeWebsocket {
			err = pushTask2WsPubSub(ctx, task, notification)
			if err != nil {
				return nil, err
			}
		} else {
			logger.Debugf(ctx, "Split notification into task[%s] successfully. ", task.TaskId)
			tasks = append(tasks, task)
		}

	}
	return tasks, nil
}

func processsAddressInfo4AddressMap(ctx context.Context, notification *models.Notification) ([]*models.Task, error) {
	addressInfo, err := models.DecodeAddressInfo(notification.AddressInfo)
	var tasks []*models.Task
	//address_info = {"email": ["xxx@abc.com", "xxx@xxx.com"],"websocket": ["system", "huojiao"]}
	for notifyType, addresses := range *addressInfo {
		for _, address := range addresses {
			directive := &models.TaskDirective{
				NotificationId:     notification.NotificationId,
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
				notifyType,
			)

			if notifyType == constants.NotifyTypeWebsocket {
				err = pushTask2WsPubSub(ctx, task, notification)
				if err != nil {
					return nil, err
				}
			} else {
				logger.Debugf(ctx, "Split notification into task[%s] successfully. ", task.TaskId)
				tasks = append(tasks, task)
			}
		}
	}
	return tasks, nil
}

func pushTask2WsPubSub(ctx context.Context, task *models.Task, nf *models.Notification) error {
	service, messageType := models.DecodeNotificationExtra4ws(nf.Extra)

	//if notify type is websocket,call websocket PushWsMessage to pubsub.
	if task.NotifyType == constants.NotifyTypeWebsocket {
		taskDirective, err := models.DecodeTaskDirective(task.Directive)
		if err != nil {
			return err
		}
		userId := taskDirective.Address

		msgDetail := wstypes.MessageDetail{
			MessageId:      wstypes.NewWsMessageId(),
			UserId:         userId,
			Service:        service,
			MessageType:    messageType,
			MessageContent: task.Directive,
		}

		userMsg := wstypes.UserMessage{
			UserId:        userId,
			Service:       service,
			MessageType:   messageType,
			MessageDetail: msgDetail,
		}
		err = websocket.PushWsMessage(context.Background(), config.GetInstance().PubSub.Type, global.GetInstance().GetPubSubClient(), &userMsg)
		if err != nil {
			logger.Errorf(ctx, "Push user [%s] websocket message id [%s] failed: %+v", userId, task.TaskId, err)
			return err
		}
		logger.Infof(ctx, "Push user [%s] websocket message id [%s] ,Directive[%s] to etcd .", userId, task.TaskId, task.Directive)

	}
	return nil

}
