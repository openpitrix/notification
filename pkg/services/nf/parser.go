package nf

import (
	"notification/pkg/models"
	"notification/pkg/pb"
	"notification/pkg/util/idutil"
)

type NfHandlerModelParser struct{
}

const NfPostIDPrifix = "nf-"


func createnfPostID() string {
	return idutil.GetUuid(NfPostIDPrifix)
}

func  (parser *NfHandlerModelParser)CreateNfWaddrs(in *pb.CreateNfWaddrsRequest) (*models.NotificationCenterPost, error) {
	nf := &models.NotificationCenterPost{
		NfPostID:        createnfPostID(),
		NfPostType:  in.GetNfPostType().GetValue(),
		Title:  in.GetTitle().GetValue(),
		Content:  in.GetContent().GetValue(),
		ShortContent :  in.GetShortContent().GetValue(),
		ExporedDays :2,
		Owner :  in.GetOwner().GetValue(),
	}
	return nf, nil
}

