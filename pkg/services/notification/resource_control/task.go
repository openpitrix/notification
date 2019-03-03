// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource_control

import (
	"context"
	"time"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/dbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

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

func UpdateTaskStatus(taskId, status string) error {
	db := globalcfg.GetInstance().GetDB()
	task := &models.Task{
		TaskId: taskId,
	}
	tx := db.Begin()
	err := db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Update("status_time", time.Now()).Error
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func GetTasksByStatus(notificationId string, status []string) []*models.Task {
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

func DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) ([]*models.Task, uint64, error) {
	req.NotificationId = stringutil.SimplifyStringList(req.NotificationId)
	req.TaskId = stringutil.SimplifyStringList(req.TaskId)
	req.TaskAction = stringutil.SimplifyStringList(req.TaskAction)
	req.ErrorCode = stringutil.SimplifyStringList(req.ErrorCode)
	req.Status = stringutil.SimplifyStringList(req.Status)

	limit := dbutil.GetLimit(req.Limit)
	offset := dbutil.GetOffset(req.Offset)

	var tasks []*models.Task
	var count uint64

	if err := dbutil.GetChain(globalcfg.GetInstance().GetDB().Table(models.TableTask)).
		AddQueryOrderDir(req, models.TaskColCreateTime).
		BuildFilterConditions(req, models.TableTask).
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error; err != nil {
		logger.Errorf(ctx, "Describe Tasks failed: %+v", err)
		return nil, 0, err
	}

	if err := dbutil.GetChain(globalcfg.GetInstance().GetDB().Table(models.TableTask)).
		BuildFilterConditions(req, models.TableTask).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe Tasks count failed: %+v", err)
		return nil, 0, err
	}

	return tasks, count, nil

}
