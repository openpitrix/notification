// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

const (
	DefaultSelectLimit = 200
)

const (
	DefaultOffset = uint32(0)
	DefaultLimit  = uint32(20)
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
	NotifyTypeEmail  = "email"
	NotifyTypeWeb    = "web"
	NotifyTypeMobile = "mobile"
	NotifyTypeSms    = "sms"
	NotifyTypeWeChat = "wechat"
)

const (
	ServiceTypeEmail  = "email"
	ServiceTypeSms    = "sms"
	ServiceTypeWeChat = "wechat"
)

const (
	StatusPending    = "pending"
	StatusSending    = "sending"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

const (
	ServiceName = "Notification"
)

const (
	ServiceCfgProtocol     = "protocol"
	ServiceCfgEmailHost    = "email_host"
	ServiceCfgPort         = "port"
	ServiceCfgDisplayEmail = "display_email"
	ServiceCfgEmail        = "email"
	ServiceCfgPassword     = "password"
)

const (
	DESC = "desc"
	ASC  = "asc"
)

const (
	TagName = "json"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusDeleted  = "deleted"
)

const (
	ContentTypeInvite   = "invite"
	ContentTypeverify   = "verify"
	ContentTypeFee      = "fee"
	ContentTypeBusiness = "business"
	ContentTypeOther    = "other"
)
