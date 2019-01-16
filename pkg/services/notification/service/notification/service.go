// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/openpitrix/pkg/etcd"
)

type Service interface {
	CreateNfWithAddrs(nf *models.Notification, q *etcd.Queue) (string, error)
	DescribeNfs(nfID string) (*models.Notification, error)
	UpdateStatus2FinishedById(nfId string) (bool, error)
	UpdateNfJob2SendingByIds(nfId string, jobId string) (bool, error)
	UpdateTask2SendingById(taskId string) (bool, error)
}

type nfService struct {
}

func NewService() Service {
	return &nfService{}
}

func (sc *nfService) CreateNfWithAddrs(nf *models.Notification, q *etcd.Queue) (string, error) {
	db := globalcfg.GetInstance().GetDB()
	var err error
	var job *models.Job

	tx := db.Begin()

	if err = tx.Create(&nf).Error; err != nil {
		tx.Rollback()
		logger.Errorf(nil, "Cannot insert notification data to db, [%+v]", err)
		return "", err
	}

	parser := &models.ModelParser{}
	job, err = parser.GenJobfromNf(nf)
	if err := tx.Create(&job).Error; err != nil {
		tx.Rollback()
		logger.Errorf(nil, "Cannot insert job data to db, [%+v]", err)
		return "", err
	}

	tasks, err := parser.GenTasksfromJob(job)
	for _, task := range tasks {
		if err := tx.Create(&task).Error; err != nil {
			tx.Rollback()
			logger.Errorf(nil, "Cannot insert task data to db, [%+v]", err)
			return "", err
		}
	}

	if err != nil {
		logger.Errorf(nil, "CreateNfWithAddrs failed, [%+v]", err)
		return "", err
	}

	tx.Commit()

	//After write DB,then write to Etcd.
	//The format to write to Etcd is nf.NotificationId + "," + task.TaskID.
	for _, task := range tasks {
		nfTaskIdStr := nf.NotificationId + "," + task.TaskID
		err = q.Enqueue(nfTaskIdStr)
		if err != nil {
			logger.Errorf(nil, "failed to push task ID into ETCD, error: [%+v]", err)
		}
		logger.Debugf(nil, "%+s", "success to push task ID into ETCD")
		//After send one task id to etcd then need to this one task status to sending.
		_, err = sc.UpdateTask2SendingById(task.TaskID)
		if err != nil {
			logger.Errorf(nil, "failed to UpdateTask2SendingById, error: [%+v]", err)
		}
		logger.Debugf(nil, "%+s", "success to UpdateTask2SendingById")
	}

	//After send all task ids to etcd then need to update nf and job status to sending.
	_, err = sc.UpdateNfJob2SendingByIds(nf.NotificationId, job.JobID)
	if err != nil {
		logger.Errorf(nil, "failed to UpdateNfJob2SendingByIds, error: [%+v]", err)
	}
	logger.Debugf(nil, "%+s", "success to UpdateNfJob2SendingByIds")

	return nf.NotificationId, nil
}

func (sc *nfService) DescribeNfs(nfID string) (*models.Notification, error) {
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

func (sc *nfService) UpdateStatus2FinishedById(nfId string) (bool, error) {
	db := globalcfg.GetInstance().GetDB()
	nf := &models.Notification{
		NotificationId: nfId,
	}
	status := constants.StatusFinished
	err := db.Model(&nf).Where("notification_id = ?", nfId).Update("status", status).Error
	if err != nil {
		logger.Errorf(nil, "%+v", err)
		return false, err
	}

	return true, nil
}

func (sc *nfService) UpdateNfJob2SendingByIds(nfId string, jobId string) (bool, error) {
	db := globalcfg.GetInstance().GetDB()

	job := &models.Job{
		JobID: jobId,
	}
	//task := &models.Task{
	//	TaskID: taskId,
	//}
	nf := &models.Notification{
		NotificationId: nfId,
	}

	tx := db.Begin()
	status := constants.StatusSending
	//err := db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Error
	err := db.Model(&job).Where("job_id = ?", jobId).Update("status", status).Error
	err = db.Model(&nf).Where("notification_id = ?", nfId).Update("status", status).Error
	if err != nil {
		logger.Errorf(nil, "%+v", err)
		return false, err
	}
	tx.Commit()

	return true, nil
}

func (sc *nfService) UpdateTask2SendingById(taskId string) (bool, error) {
	db := globalcfg.GetInstance().GetDB()
	task := &models.Task{
		TaskID: taskId,
	}
	status := constants.StatusSending
	err := db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Error
	if err != nil {
		logger.Errorf(nil, "%+v", err)
		return false, err
	}
	return true, nil
}
