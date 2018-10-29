package nf

import (
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"time"
)

type NfHandlerModelParser struct{
}

const NfPostIDPrifix = "nf-"


func CreatenfPostID() string {
	return idutil.GetUuid(NfPostIDPrifix)
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

