// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"context"
	"time"

	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/gerr"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"openpitrix.io/notification/pkg/util/jsonutil"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type Notification struct {
	NotificationId     string    `gorm:"column:notification_id"`
	ContentType        string    `gorm:"column:content_type"`
	Title              string    `gorm:"column:title"`
	Content            string    `gorm:"column:content"`
	ShortContent       string    `gorm:"column:short_content"`
	ExpiredDays        uint32    `gorm:"column:expired_days"`
	AddressInfo        string    `gorm:"column:address_info"`
	Owner              string    `gorm:"column:owner"`
	Status             string    `gorm:"column:status"`
	CreateTime         time.Time `gorm:"column:create_time"`
	StatusTime         time.Time `gorm:"column:status_time"`
	AvailableStartTime string    `gorm:"column:available_start_time"`
	AvailableEndTime   string    `gorm:"column:available_end_time"`
	Extra              string    `gorm:"column:extra"`
}

//table name
const (
	TableNotification = "notification"
)

//ID Prefix
const (
	NotificationIdPrefix = "nf-"
)

//field name
//Nf is short for notification.
const (
	NfColId           = "notification_id"
	NfColContentType  = "content_type"
	NfColTitle        = "title"
	NfColContent      = "content"
	NfColShortContent = "short_content"
	NfColExpiredDays  = "expired_days"
	NfColAddressInfo  = "address_info"
	NfColOwner        = "owner"
	NfColStatus       = "status"
	NfColCreateTime   = "create_time"
	NfColStatusTime   = "status_time"
)

//constants
const (
	ContentTypeInvite   = "invite"
	ContentTypeVerify   = "verify"
	ContentTypeFee      = "fee"
	ContentTypeBusiness = "business"
	ContentTypeAlerting = "alert"
	ContentTypeOther    = "other"
	ContentTypeEvent    = "event"
)

var ContentTypes = []string{
	ContentTypeInvite,
	ContentTypeVerify,
	ContentTypeFee,
	ContentTypeBusiness,
	ContentTypeAlerting,
	ContentTypeOther,
	ContentTypeEvent,
}

func NewNotificationId() string {
	return idutil.GetUuid(NotificationIdPrefix)
}

func NewNotification(contentType, title, content, shortContent, addressInfo, owner string, expiredDays uint32, availableStartTimeStr string, availableEndTimeStr string, extra string) *Notification {
	if extra == "" {
		extra = "{}"
	}

	notification := &Notification{
		NotificationId:     NewNotificationId(),
		ContentType:        contentType,
		Title:              title,
		Content:            content,
		ShortContent:       shortContent,
		ExpiredDays:        expiredDays,
		AddressInfo:        addressInfo,
		Owner:              owner,
		Status:             constants.StatusPending,
		CreateTime:         time.Now(),
		StatusTime:         time.Now(),
		AvailableStartTime: availableStartTimeStr,
		AvailableEndTime:   availableEndTimeStr,
		Extra:              extra,
	}
	return notification
}

func NotificationToPb(notification *Notification) *pb.Notification {
	pbNotification := pb.Notification{}
	pbNotification.NotificationId = pbutil.ToProtoString(notification.NotificationId)
	pbNotification.ContentType = pbutil.ToProtoString(notification.ContentType)
	pbNotification.Title = pbutil.ToProtoString(notification.Title)
	pbNotification.Content = pbutil.ToProtoString(notification.Content)
	pbNotification.ShortContent = pbutil.ToProtoString(notification.ShortContent)
	pbNotification.ExpiredDays = pbutil.ToProtoUInt32(notification.ExpiredDays)
	pbNotification.AddressInfo = pbutil.ToProtoString(notification.AddressInfo)
	pbNotification.Owner = pbutil.ToProtoString(notification.Owner)
	pbNotification.Status = pbutil.ToProtoString(notification.Status)
	pbNotification.CreateTime = pbutil.ToProtoTimestamp(notification.CreateTime)
	pbNotification.StatusTime = pbutil.ToProtoTimestamp(notification.StatusTime)
	pbNotification.AvailableStartTime = pbutil.ToProtoString(notification.AvailableStartTime)
	pbNotification.AvailableEndTime = pbutil.ToProtoString(notification.AvailableEndTime)
	pbNotification.Extra = pbutil.ToProtoString(notification.Extra)
	return &pbNotification
}

func NotificationSet2PbSet(inNfs []*Notification) []*pb.Notification {
	var pbNfs []*pb.Notification
	for _, inNf := range inNfs {
		pbNf := NotificationToPb(inNf)
		pbNfs = append(pbNfs, pbNf)
	}
	return pbNfs
}

type NotificationExtra map[string]string

func DecodeNotificationExtra(data string) (*map[string]string, error) {
	extra := new(map[string]string)
	err := jsonutil.Decode([]byte(data), extra)
	if err != nil {
		logger.Errorf(nil, "Decode [%s] into notification extra failed: %+v", data, err)
	}
	return extra, err
}

func CheckExtra(ctx context.Context, notification *Notification) error {
	extraStr := notification.Extra
	if extraStr == "" {
		logger.Errorf(ctx, "Failed to validate addressInfo, extra is blank: [%s].", extraStr)
		return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationExtra, extraStr)
	} else {
		//check Extra:  "extra": "{\"ws_message_type\": \"ws_op_nf\"}"
		nfExtraMap, err := DecodeNotificationExtra(extraStr)
		if err != nil {
			logger.Errorf(ctx, "Failed to validate notification extra [%s], should be: {\"ws_message_type\": \"xxx\"}", extraStr)
			return err
		}
		_, ok := (*nfExtraMap)[constants.WsMessageType]
		if !ok {
			logger.Errorf(ctx, "Failed to validate notification extra [%s], should be: {\"ws_message_type\": \"xxx\"}", extraStr)
			return gerr.New(ctx, gerr.InvalidArgument, gerr.ErrorIllegalNotificationExtra, extraStr)
		}
	}
	return nil
}
