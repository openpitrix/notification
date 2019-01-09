package models

import (
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"strings"
	"time"
)

type ModelParser struct {
}

func (parser *ModelParser) CreateNfWithAddrs(in *pb.CreateNfWithAddrsRequest) (*Notification, error) {
	nf := &Notification{
		NotificationId: idutil.GetUuid(constants.NfPostIDPrifix),
		ContentType:    in.GetContentType().GetValue(),
		SentType:       in.GetSentType().GetValue(),
		AddrsStr:       in.GetAddrsStr().GetValue(),
		Title:          in.GetTitle().GetValue(),
		Content:        in.GetContent().GetValue(),
		ShortContent:   in.GetShortContent().GetValue(),
		ExporedDays:    2,
		Owner:          in.GetOwner().GetValue(),
		Status:         constants.StatusNew,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return nf, nil
}

func (parser *ModelParser) GenJobfromNf(nf *Notification) (*Job, error) {
	//todo check eamil string
	emailsArray := strings.Split(nf.AddrsStr, ";")
	taskcnt := int64(len(emailsArray))
	job := &Job{
		JobID:          idutil.GetUuid(constants.JobPostIDPrifix),
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

//GenTaskfromJob
func (parser *ModelParser) GenTasksfromJob(job *Job) ([]*Task, error) {
	emailsArray := strings.Split(job.AddrsStr, ";")
	tasks := make([]*Task, 0, len(emailsArray))
	for _, email := range emailsArray {
		tasks = append(tasks, &Task{
			TaskID:     idutil.GetUuid(constants.TaskPostIDPrifix),
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
