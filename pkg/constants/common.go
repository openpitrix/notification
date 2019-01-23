// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

const (
	EtcdPrefix              = "notification/"
	NotificationTaskTopic   = "nft"
	NotificationTopic       = "nf"
	MaxWorkingTasks         = 5
	MaxWorkingNotifications = 5
	MaxTaskRetryTimes       = 5
)

const (
	NotificationIdPrefix = "nf-"
	TaskIdPrefix         = "t-"
)

const (
	NotifyTypeWeb    = "web"
	NotifyTypeMobile = "mobile"
	NotifyTypeEmail  = "email"
	NotifyTypeSms    = "sms"
	NotifyTypeWeChat = "wechat"
)

const (
	StatusPending    = "pending"
	StatusSending    = "sending"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

const (
	ServiceName = "Notification"
	prefix      = "notification-"

	NotificationManagerHost = prefix + "manager"
	NotificationManagerPort = 9201
)
