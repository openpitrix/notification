// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type Notification struct {
	NotificationId string    `gorm:"column:notification_id"`
	ContentType    string    `gorm:"column:content_type"`
	Title          string    `gorm:"column:title"`
	Content        string    `gorm:"column:content"`
	ShortContent   string    `gorm:"column:short_content"`
	ExpiredDays    uint32    `gorm:"column:expired_days"`
	AddressInfo    string    `gorm:"column:address_info"`
	Owner          string    `gorm:"column:owner"`
	Status         string    `gorm:"column:status"`
	CreateTime     time.Time `gorm:"column:create_time"`
	StatusTime     time.Time `gorm:"column:status_time"`
}

//table name
const (
	TableNotification = "notification"
)
const (
	NotificationIdPrefix = "nf-"
)

//field name
//Nf is short for notification.
const (
	NfColId          = "notification_id"
	NfColStatus      = "status"
	NfColContentType = "content_type"
	NfColOwner       = "owner"
	NfColTitle       = "title"
	NfColAddressInfo = "address_info"
	NfColCreateTime  = "create_time"
)

func NewNotificationId() string {
	return idutil.GetUuid(NotificationIdPrefix)
}

func NewNotification(contentType, title, content, shortContent, addressInfo, owner string, expiredDays uint32) *Notification {
	notification := &Notification{
		NotificationId: NewNotificationId(),
		ContentType:    contentType,
		Title:          title,
		Content:        content,
		ShortContent:   shortContent,
		ExpiredDays:    expiredDays,
		AddressInfo:    addressInfo,
		Owner:          owner,
		Status:         constants.StatusPending,
		CreateTime:     time.Now(),
		StatusTime:     time.Now(),
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
	return &pbNotification
}

func ParseNfSet2PbSet(inNfs []*Notification) []*pb.Notification {
	var pbNfs []*pb.Notification
	for _, inNf := range inNfs {
		pbNf := NotificationToPb(inNf)
		pbNfs = append(pbNfs, pbNf)
	}
	return pbNfs
}
