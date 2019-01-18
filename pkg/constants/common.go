// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package constants

const (
	EtcdPrefix            = "notification/"
	NotificationTaskTopic = "nft"
	NotificationJobTopic  = "nfj"
	MaxWorkingTasks       = 5
	MaxWorkingJobs        = 5
	MaxTaskRetryTimes     = 5
)

const (
	NfPostIDPrefix   = "nf-"
	JobPostIDPrefix  = "job-"
	TaskPostIDPrefix = "task-"
)

const (
	StatusNew        = "new"
	StatusSending    = "sending"
	StatusSuccessful = "successful"
	StatusFailed     = "failed"
)

const (
	ServiceName    = "Notification"
	prefix         = "notification-"
	ApiGatewayHost = prefix + "api-gateway"
	//ApiGatewayHost = "127.0.0.1"
	ApiGatewayPort = 9200

	NotificationManagerHost = prefix + "manager"
	//NotificationManagerHost = "127.0.0.1"
	//NotificationManagerHost = "192.168.0.3"
	NotificationManagerPort = 9201
)
