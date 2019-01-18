// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"strings"
	"time"

	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
)

func GenNotificationFromReq(req *pb.CreateNfWithAddrsRequest) (*models.Notification, error) {
	nf := &models.Notification{
		NotificationId: idutil.GetUuid(constants.NfPostIDPrefix),
		ContentType:    req.GetContentType().GetValue(),
		SentType:       req.GetSentType().GetValue(),
		AddrsStr:       req.GetAddrsStr().GetValue(),
		Title:          req.GetTitle().GetValue(),
		Content:        req.GetContent().GetValue(),
		ShortContent:   req.GetShortContent().GetValue(),
		ExporedDays:    2,
		Owner:          req.GetOwner().GetValue(),
		Status:         constants.StatusNew,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return nf, nil
}

func GenJobFromNf(nf *models.Notification) (*models.Job, error) {
	//todo check eamil string
	emailsArray := strings.Split(nf.AddrsStr, ";")
	taskcnt := int64(len(emailsArray))
	job := &models.Job{
		JobID:          idutil.GetUuid(constants.JobPostIDPrefix),
		NotificationId: nf.NotificationId,
		JobType:        nf.SentType,
		AddrsStr:       nf.AddrsStr,
		JobAction:      "Job Action Test",
		ExeCondition:   "Job Action Test",
		TotalTaskCount: taskcnt,
		TaskSuccCount:  0,
		ErrorCode:      0,
		Status:         "Ready",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return job, nil
}

func GenTasksFromJob(job *models.Job) ([]*models.Task, error) {
	emailsArray := strings.Split(job.AddrsStr, ";")
	tasks := make([]*models.Task, 0, len(emailsArray))
	for _, email := range emailsArray {
		tasks = append(tasks, &models.Task{
			TaskID:     idutil.GetUuid(constants.TaskPostIDPrefix),
			JobID:      job.JobID,
			EmailAddr:  email,
			TaskAction: "",
			ErrorCode:  0,
			Status:     "New",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}
	return tasks, nil
}
