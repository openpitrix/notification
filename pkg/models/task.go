// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"openpitrix.io/notification/pkg/util/jsonutil"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type Task struct {
	TaskId         string    `gorm:"column:task_id"`
	NotificationId string    `gorm:"column:notification_id"`
	ErrorCode      int64     `gorm:"column:error_code"`
	Status         string    `gorm:"column:status"`
	CreateTime     time.Time `gorm:"column:create_time"`
	StatusTime     time.Time `gorm:"column:status_time"`
	Directive      string    `gorm:"column:directive"`
}

//table name
const (
	TableTask = "task"
)

const (
	TaskIdPrefix = "t-"
)

//field name
//Nf is short for notification.
const (
	TaskColNfId       = "notification_id"
	TaskColTaskId     = "task_id"
	TaskColStatus     = "status"
	TaskColErrorCode  = "error_code"
	TaskColCreateTime = "create_time"
)

func NewTask(notificationId, directive string) *Task {
	task := &Task{
		TaskId:         NewTaskId(),
		NotificationId: notificationId,
		ErrorCode:      0,
		Status:         constants.StatusPending,
		CreateTime:     time.Now(),
		StatusTime:     time.Now(),
		Directive:      directive,
	}
	return task
}

func NewTaskId() string {
	return idutil.GetUuid(TaskIdPrefix)
}

type TaskDirective struct {
	Address      string
	NotifyType   string
	ContentType  string
	Title        string
	Content      string
	ShortContent string
	ExpiredDays  uint32
}

func DecodeTaskDirective(data string) (*TaskDirective, error) {
	taskDirective := new(TaskDirective)
	err := jsonutil.Decode([]byte(data), taskDirective)
	if err != nil {
		logger.Errorf(nil, "Decode [%s] into task directive failed: %+v", data, err)
	}
	return taskDirective, err
}

type TaskWithNfInfo struct {
	NotificationId string
	JobID          string
	TaskID         string
	Title          string
	ShortContent   string
	Content        string
	EmailAddr      string
}

func TaskToPb(task *Task) *pb.Task {
	pbTask := pb.Task{}
	pbTask.NotificationId = pbutil.ToProtoString(task.NotificationId)
	return &pbTask
}

func TaskSet2PbSet(inTasks []*Task) []*pb.Task {
	var pbTasks []*pb.Task
	for _, inTask := range inTasks {
		pbTask := TaskToPb(inTask)
		pbTasks = append(pbTasks, pbTask)
	}
	return pbTasks
}
