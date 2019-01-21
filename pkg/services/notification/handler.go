// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/etcd"
)

func CreateNotification(nf *models.Notification, taskQueue, jobQueue *etcd.Queue) (string, error) {
	db := globalcfg.GetInstance().GetDB()
	var err error
	var job *models.Job

	tx := db.Begin()

	if err = tx.Create(&nf).Error; err != nil {
		tx.Rollback()
		logger.Errorf(nil, "Cannot insert notification data to db, [%+v]", err)
		return "", err
	}

	job, err = GenJobFromNf(nf)
	if err := tx.Create(&job).Error; err != nil {
		tx.Rollback()
		logger.Errorf(nil, "Cannot insert job data to db, [%+v]", err)
		return "", err
	}

	tasks, err := GenTasksFromJob(job)
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

	// Before send all task ids to etcd then need to update nf and job status to sending.
	err = UpdateNfStatus(nf.NotificationId, constants.StatusSending)
	if err != nil {
		logger.Errorf(nil, "Failed to update nf [%s] status to sending, error: [%+v]", nf.NotificationId, err)
	}
	logger.Debugf(nil, "Succeed to update nf [%s] status to sending", nf.NotificationId)

	err = UpdateJobStatus(job.JobID, constants.StatusSending)
	if err != nil {
		logger.Errorf(nil, "Failed to update job [%s] status to sending, error: [%+v]", job.JobID, err)
	}
	logger.Debugf(nil, "Succeed to update job [%s] status to sending", job.JobID)

	nfIdAndJobIdStr := nf.NotificationId + "," + job.JobID
	err = jobQueue.Enqueue(nfIdAndJobIdStr)

	// After write DB,then write to Etcd.
	// The format to write to Etcd is nf.NotificationId + "," + task.TaskID.
	for _, task := range tasks {
		// After send one task id to etcd then need to change the task status to sending.
		err = UpdateTaskStatus(task.TaskID, constants.StatusSending)
		if err != nil {
			logger.Errorf(nil, "Failed to update task [%s] status to sending, error: [%+v]", task.TaskID, err)
		}
		logger.Debugf(nil, "Succeed to update task [%s] status to sending", task.TaskID)

		err = taskQueue.Enqueue(task.TaskID)
		if err != nil {
			logger.Errorf(nil, "Failed to push task [%s] into etcd, error: [%+v]", task.TaskID, err)
		}
		logger.Debugf(nil, "Succeed to push task [%s] into etcd", task.TaskID)
	}

	return nf.NotificationId, nil
}

func (s *Server) DescribeNfs(ctx context.Context, req *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeNfs at server end. "}, nil
}

func (s *Server) CreateNfWithAddrs(ctx context.Context, req *pb.CreateNfWithAddrsRequest) (*pb.CreateNfWithAddrsResponse, error) {
	nf, err := GenNotificationFromReq(req)
	if err != nil {
		logger.Errorf(ctx, "Failed to parser.CreateNfWithAddrs, error:[%+v]", err)
		return nil, err
	}
	logger.Debugf(ctx, "Succeed to parser.CreateNfWithAddrs, NotificationId:[%s]", nf.NotificationId)

	nfId, err := CreateNotification(nf, s.controller.taskQueue, s.controller.jobQueue)
	if err != nil {
		logger.Errorf(ctx, "Failed to service.CreateNfWithAddrs, error:[%+v]", err)
		return nil, err
	}
	logger.Debugf(ctx, "Succeed to service.CreateNfWithAddrs, NotificationId:[%s]", nf.NotificationId)

	res := &pb.CreateNfWithAddrsResponse{
		NotificationId: pbutil.ToProtoString(nfId),
	}
	return res, nil
}
