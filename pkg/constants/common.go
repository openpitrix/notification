// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

const (
	ServiceName = "Notification"
)

const (
	TagName = "json"
)

const (
	EtcdPrefix                  = "notification/"
	NotificationTaskTopicPrefix = "nf-task"
	NotificationTopicPrefix     = "nf-job"
	MaxWorkingTasks             = 5
	MaxWorkingNotifications     = 5
	MaxTaskRetryTimes           = 5
)

const (
	DESC = "desc"
	ASC  = "asc"
)

const (
	DefaultOffset = uint32(0)
	DefaultLimit  = uint32(20)
)

const (
	DefaultSelectLimit = 200
)

const (
	NotifyTypeEmail  = "email"
	NotifyTypeWeb    = "web"
	NotifyTypeMobile = "mobile"
	NotifyTypeSms    = "sms"
	NotifyTypeWeChat = "wechat"
)

var NotifyTypes = []string{
	NotifyTypeEmail,
	NotifyTypeWeb,
	NotifyTypeMobile,
	NotifyTypeSms,
	NotifyTypeWeChat,
}

const (
	StatusPending    = "pending"
	StatusSending    = "sending"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

var SendingStatuses = []string{
	StatusPending,
	StatusSending,
	StatusSuccessful,
	StatusFailed,
}

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusDeleted  = "deleted"
)

var RecordStatuses = []string{
	StatusActive,
	StatusDisabled,
	StatusDeleted,
}
