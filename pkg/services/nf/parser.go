package nf

import (
	"log"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"time"
)

type NfHandlerModelParser struct{
}

const NfPostIDPrifix = "nf-"
const JobPostIDPrifix = "job-"

func CreatenfPostID() string {
	return idutil.GetUuid(NfPostIDPrifix)
}

func CreatenfJobID() string {
	return idutil.GetUuid(JobPostIDPrifix)
}

func  (parser *NfHandlerModelParser)CreateNfWaddrs(in *pb.CreateNfWaddrsRequest) (*models.NotificationCenterPost, error) {
	nf := &models.NotificationCenterPost{
		NfPostID:        CreatenfPostID(),
		NfPostType:  in.GetNfPostType().GetValue(),
		AddrsStr:in.GetAddrsStr().GetValue(),
		Title:  in.GetTitle().GetValue(),
		Content:  in.GetContent().GetValue(),
		ShortContent :  in.GetShortContent().GetValue(),
		ExporedDays :2,
		Owner :  in.GetOwner().GetValue(),
		Status:"New",
		CreatedAt:time.Now(),
		UpdatedAt:time.Now(),
		DeletedAt:time.Now(),
	}
	return nf, nil
}

func (parser *NfHandlerModelParser)GenJobfromReq(nf *models.NotificationCenterPost) (*models.Job, error){

	job:=&models.Job{
		JobID:CreatenfJobID(),
		NfPostID:nf.NfPostID,
		JobType:"Email",
		AddrsStr: nf.AddrsStr,
		JobAction: "Job Action Test",
		ExeCondition :"Job Action Test",
		TotalTaskCount: 2,
		TaskSuccCount :2,
		Result :"succ",
		ErrorCode:  0,
		Status : "done",
		CreatedAt:time.Now(),
		UpdatedAt:time.Now(),
		DeletedAt:time.Now(),
	}

	log.Println(job.JobID)
	log.Println(job.NfPostID)
	log.Println(job.AddrsStr)
	log.Println(job.JobType)

	return job,nil
}