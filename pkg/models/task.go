// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/util/idutil"
	"openpitrix.io/notification/pkg/util/jsonutil"
)

func NewTaskId() string {
	return idutil.GetUuid(constants.TaskIdPrefix)
}

type Task struct {
	TaskId         string    `gorm:"column:task_id"`
	NotificationId string    `gorm:"column:notification_id"`
	ErrorCode      int64     `gorm:"column:error_code"`
	Status         string    `gorm:"column:status"`
	CreateTime     time.Time `gorm:"column:create_time"`
	StatusTime     time.Time `gorm:"column:status_time"`
	Directive      string    `gorm:"column:directive"`
}

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
