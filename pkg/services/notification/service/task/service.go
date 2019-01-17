// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
)

type Service interface {
	GetTaskWithNfContentByID(taskID string) (*models.TaskWithNfInfo, error)
	UpdateStatus2SendingByIds(taskWithNfInfo models.TaskWithNfInfo) (bool, error)
	UpdateJobTaskStatus2FinishedById(taskWithNfInfo models.TaskWithNfInfo) (bool, error)
}

type taskService struct {
}

func NewService() Service {
	return &taskService{}
}

func (sc *taskService) getTaskbyID(taskID string) (*models.Task, error) {
	db := globalcfg.GetInstance().GetDB()
	task := &models.Task{}
	err := db.
		Where("task_id = ?", taskID).
		First(task).Error
	if err != nil {
		//if err != gorm.ErrRecordNotFound {
		//	return nil, err
		//}
		return nil, err
	}
	return task, nil
}

func (sc *taskService) GetTaskWithNfContentByID(taskID string) (*models.TaskWithNfInfo, error) {
	db := globalcfg.GetInstance().GetDB()
	logger.Debugf(nil, "test========taskID=%s", taskID)
	taskWithNfInfo := &models.TaskWithNfInfo{}
	sql := models.GetTaskWithNfContentByIDSQL
	db.Raw(sql, taskID).Scan(&taskWithNfInfo)
	logger.Debugf(nil, "getTaskwithNfContentbyID got a task,TaskID: [%+s]", taskWithNfInfo.TaskID)
	return taskWithNfInfo, nil
}

func (sc *taskService) UpdateStatus2SendingByIds(taskWithNfInfo models.TaskWithNfInfo) (bool, error) {
	db := globalcfg.GetInstance().GetDB()
	jobId := taskWithNfInfo.JobID
	taskId := taskWithNfInfo.TaskID
	nfId := taskWithNfInfo.NotificationId

	job := &models.Job{
		JobID: jobId,
	}
	task := &models.Task{
		TaskID: taskId,
	}
	nf := &models.Notification{
		NotificationId: nfId,
	}

	tx := db.Begin()
	status := constants.StatusSending
	err := db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Error
	err = db.Model(&job).Where("job_id = ?", jobId).Update("status", status).Error
	err = db.Model(&nf).Where("notification_id = ?", nfId).Update("status", status).Error
	if err != nil {
		logger.Errorf(nil, "%+v", err)
		return false, err
	}
	tx.Commit()

	return true, nil
}

func (sc *taskService) UpdateJobTaskStatus2FinishedById(taskWithNfInfo models.TaskWithNfInfo) (bool, error) {
	db := globalcfg.GetInstance().GetDB()

	jobId := taskWithNfInfo.JobID
	taskId := taskWithNfInfo.TaskID

	job := &models.Job{
		JobID: jobId,
	}
	task := &models.Task{
		TaskID: taskId,
	}

	tx := db.Begin()
	status := constants.StatusFinished
	err := db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Error
	err = db.Model(&job).Where("job_id = ?", jobId).Update("status", status).Error
	if err != nil {
		logger.Errorf(nil, "%+v", err)
		return false, err
	}
	tx.Commit()

	return true, nil
}
