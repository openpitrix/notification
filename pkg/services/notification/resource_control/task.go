// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource_control

import (
	"context"
	"time"

	"openpitrix.io/logger"

	nfdb "openpitrix.io/notification/pkg/db"
	"openpitrix.io/notification/pkg/global"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/util/stringutil"
)

func RegisterTask(ctx context.Context, task *models.Task) error {
	db := global.GetInstance().GetDB()
	err := db.Create(&task).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to insert task, %+v.", err)
		return err
	}
	return nil
}

func UpdateTasksStatus(ctx context.Context, taskIds []string, status string) error {
	db := global.GetInstance().GetDB()
	err := db.Table(models.TableTask).Where(models.TaskColTaskId+" in (?)", taskIds).Updates(map[string]interface{}{models.TaskColStatus: status, models.NfColStatusTime: time.Now()}).Error

	if err != nil {
		logger.Errorf(ctx, "Failed to update task status to [%s], %+v.", status, err)
		return err
	}
	return nil
}

func GetTasksByStatus(ctx context.Context, notificationId string, status []string) []*models.Task {
	db := global.GetInstance().GetDB()
	var tasks []*models.Task
	err := db.Where("notification_id = ? AND status in (?)", notificationId, status).Find(&tasks).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to get tasks by status[%+v], %+v.", status, err)
		return nil
	}
	return tasks
}

func GetTasksByTaskIds(ctx context.Context, taskIds []string) ([]*models.Task, error) {
	db := global.GetInstance().GetDB()
	var tasks []*models.Task
	err := db.Where("task_id in( ? )", taskIds).Find(&tasks).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to get tasks by taskIds[%+v], %+v.", taskIds, err)
		return nil, err
	}
	return tasks, nil
}

func GetTaskIdsByNfIds(ctx context.Context, nfIds []string) ([]string, error) {
	var tasks []*models.Task
	err := global.GetInstance().GetDB().Table(models.TableTask).Where(models.TaskColNfId+" in (?)", nfIds).Find(&tasks).Error
	if err != nil {
		logger.Errorf(ctx, "Failed to get task Ids by NfIds[%+v], %+v.", nfIds, err)
		return nil, err
	}

	var taskIds []string
	for _, task := range tasks {
		taskId := task.TaskId
		taskIds = append(taskIds, taskId)
	}

	return taskIds, nil

}

func DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) ([]*models.Task, uint64, error) {
	req.NotificationId = stringutil.SimplifyStringList(req.NotificationId)
	req.TaskId = stringutil.SimplifyStringList(req.TaskId)
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
		logger.Errorf(ctx, "Failed to describe tasks, %+v.", err)
		return nil, 0, err
	}

	if err := nfdb.GetChain(global.GetInstance().GetDB().Table(models.TableTask)).
		BuildFilterConditions(req, models.TableTask).
		Count(&count).Error; err != nil {
		logger.Errorf(ctx, "Failed to describe task count, %+v", err)
		return nil, 0, err
	}

	return tasks, count, nil

}
