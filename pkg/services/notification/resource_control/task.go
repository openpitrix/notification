// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource_control

import (
	"context"
	"time"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	nfdb "openpitrix.io/notification/pkg/db"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func RegisterTask(ctx context.Context, task *models.Task) error {
	db := global.GetInstance().GetDB()
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
	db := global.GetInstance().GetDB()
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
	db := global.GetInstance().GetDB()
	var tasks []*models.Task
	tx := db.Begin()
	db.Where("notification_id = ? AND status in (?)", notificationId, status).Find(&tasks)
	tx.Commit()
	return tasks
}

func GetTask(taskId string) (*models.Task, error) {
	db := global.GetInstance().GetDB()
	task := new(models.Task)
	err := db.Where("task_id = ?", taskId).First(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func GetTasksByTaskIds(taskIds []string) ([]*models.Task, error) {
	db := global.GetInstance().GetDB()
	var tasks []*models.Task
	err := db.Where("task_id in( ? )", taskIds).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTasksByNfId(nfIds []string) ([]*models.Task, error) {
	var tasks []*models.Task
	err := nfdb.GetChain(global.GetInstance().GetDB().Where(models.TaskColNfId+" in (?)", nfIds).Find(&tasks)).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil

}

func DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) ([]*models.Task, uint64, error) {
	req.NotificationId = stringutil.SimplifyStringList(req.NotificationId)
	req.TaskId = stringutil.SimplifyStringList(req.TaskId)
	req.TaskAction = stringutil.SimplifyStringList(req.TaskAction)
	req.ErrorCode = stringutil.SimplifyStringList(req.ErrorCode)
	req.Status = stringutil.SimplifyStringList(req.Status)
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	var tasks []*models.Task
	var count uint64

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableTask)).
		AddQueryOrderDir(req, models.TaskColCreateTime).
		BuildFilterConditions(req, models.TableTask).
		Offset(offset).
		Limit(limit).
		Find(&tasks).Error; err != nil {
		logger.Errorf(ctx, "Describe Tasks failed: %+v", err)
		return nil, 0, err
	}

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableTask)).
		BuildFilterConditions(req, models.TableTask).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Describe Tasks count failed: %+v", err)
		return nil, 0, err
	}

	return tasks, count, nil

}

func UpdateTasks2Pending(ctx context.Context, taskIds []string) ([]string, error) {
	db := global.GetInstance().GetDB()
	tx := db.Begin()
	db.Table(models.TableTask).Where(models.TaskColTaskId+" in (?)", taskIds).Updates(map[string]interface{}{models.TaskColStatus: constants.StatusPending, models.NfColStatusTime: time.Now()})

	if err := tx.Error; err != nil {
		tx.Rollback()
		logger.Errorf(ctx, "Update Tasks Status to Pending failed: [%+v].", err)
		return nil, err
	}

	tx.Commit()
	return taskIds, nil
}
