// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"github.com/jinzhu/gorm"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
)

func UpdateNfStatus(nfId, status string) error {
	db := globalcfg.GetInstance().GetDB()

	nf := &models.Notification{
		NotificationId: nfId,
	}

	tx := db.Begin()
	err := db.Model(&nf).Where("notification_id = ?", nfId).Update("status", status).Error
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func UpdateJobStatus(jobId, status string) error {
	db := globalcfg.GetInstance().GetDB()

	job := &models.Job{
		JobID: jobId,
	}

	tx := db.Begin()
	err := db.Model(&job).Where("job_id = ?", jobId).Update("status", status).Error
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func UpdateTaskStatus(taskId, status string) error {
	db := globalcfg.GetInstance().GetDB()

	task := &models.Task{
		TaskID: taskId,
	}

	tx := db.Begin()
	err := db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Error
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func GetStatusTasks(jobId string, status []string) []*models.Task {
	db := globalcfg.GetInstance().GetDB()

	var tasks []*models.Task

	tx := db.Begin()
	db.Where("job_id = ? AND status in (?)", jobId, status).Find(&tasks)

	tx.Commit()

	return tasks
}

func GetTaskWithNfInfo(taskID string) (*models.TaskWithNfInfo, error) {
	db := globalcfg.GetInstance().GetDB()
	taskWithNfInfo := &models.TaskWithNfInfo{}
	sql := models.GetTaskWithNfContentByIDSQL
	db.Raw(sql, taskID).Scan(&taskWithNfInfo)
	logger.Debugf(nil, "Get task: [%s]", taskWithNfInfo.TaskID)
	return taskWithNfInfo, nil
}

func GetNotification(nfID string) (*models.Notification, error) {
	db := globalcfg.GetInstance().GetDB()
	nf := &models.Notification{}
	err := db.
		Where("nf_post_id = ?", nfID).
		First(nf).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, err
	}

	return nf, nil
}
