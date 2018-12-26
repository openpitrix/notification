package nf

import (
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"strings"
	"time"
)

type NfHandlerModelParser struct {
}

func (parser *NfHandlerModelParser) CreateNfWaddrs(in *pb.CreateNfWithAddrsRequest) (*models.Notification, error) {
	nf := &models.Notification{
		NotificationId: idutil.GetUuid(constants.NfPostIDPrifix),
		ContentType:    in.GetContentType().GetValue(),
		SentType:       in.GetSentType().GetValue(),
		AddrsStr:       in.GetAddrsStr().GetValue(),
		Title:          in.GetTitle().GetValue(),
		Content:        in.GetContent().GetValue(),
		ShortContent:   in.GetShortContent().GetValue(),
		ExporedDays:    2,
		Owner:          in.GetOwner().GetValue(),
		Status:         "New",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return nf, nil
}

func (parser *NfHandlerModelParser) GenJobfromNf(nf *models.Notification) (*models.Job, error) {
	//todo check eamil string
	emailsArray := strings.Split(nf.AddrsStr, ";")
	taskcnt := int64(len(emailsArray))
	job := &models.Job{
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
func (parser *NfHandlerModelParser) GenTasksfromJob(job *models.Job) ([]*models.Task, error) {
	emailsArray := strings.Split(job.AddrsStr, ";")
	tasks := make([]*models.Task, 0, len(emailsArray))
	for _, email := range emailsArray {
		tasks = append(tasks, &models.Task{
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
