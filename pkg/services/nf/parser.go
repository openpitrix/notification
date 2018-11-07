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



func (parser *NfHandlerModelParser) CreateNfWaddrs(in *pb.CreateNfWaddrsRequest) (*models.NotificationCenterPost, error) {
	nf := &models.NotificationCenterPost{
		NfPostID:     idutil.GetUuid(constants.NfPostIDPrifix),
		NfPostType:   in.GetNfPostType().GetValue(),
		NotifyType:   in.GetNotifyType().GetValue(),
		AddrsStr:     in.GetAddrsStr().GetValue(),
		Title:        in.GetTitle().GetValue(),
		Content:      in.GetContent().GetValue(),
		ShortContent: in.GetShortContent().GetValue(),
		ExporedDays:  2,
		Owner:        in.GetOwner().GetValue(),
		Status:       "New",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return nf, nil
}

func (parser *NfHandlerModelParser) GenJobfromNf(nf *models.NotificationCenterPost) (*models.Job, error) {
	//todo check eamil string
	emailsArray := strings.Split(nf.AddrsStr, ";")
	taskcnt := int64(len(emailsArray))
	job := &models.Job{
		JobID:         idutil.GetUuid(constants.JobPostIDPrifix),
		NfPostID:       nf.NfPostID,
		JobType:        nf.NotifyType,
		AddrsStr:       nf.AddrsStr,
		JobAction:      "Job Action Test",
		ExeCondition:   "Job Action Test",
		TotalTaskCount: taskcnt,
		TaskSuccCount:  0,
		Result:         "N",
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
			AddrsStr:   email,
			TaskAction: "",
			Result:     "Ready",
			ErrorCode:  0,
			Status:     "New",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}
	return tasks, nil
}
