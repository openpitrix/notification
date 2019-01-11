package notification

import (
	"golang.org/x/net/context"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/globalcfg"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/services/notification/service/notification"
	"openpitrix.io/notification/pkg/services/notification/service/task"
	"openpitrix.io/notification/pkg/util/pbutil"
	"testing"
	"time"
)

func TestCreateNfWithAddrs4handler(t *testing.T) {
	nfservice := notification.NewService()
	taskservice := task.NewService()
	queue := globalcfg.GetInstance().GetEtcd().NewQueue(constants.EmailQueue)
	handler := NewHandler(nfservice, taskservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testAddrsStr := "huojiao2006@163.com;513590162@qq.com"

	var req = &pb.CreateNfWithAddrsRequest{
		ContentType:  pbutil.ToProtoString("Information"),
		SentType:     pbutil.ToProtoString("email"),
		AddrsStr:     pbutil.ToProtoString(testAddrsStr),
		Title:        pbutil.ToProtoString("Title Test"),
		Content:      pbutil.ToProtoString("Content"),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoString("7"),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		Status:       pbutil.ToProtoString("New"),
	}

	handler.CreateNfWithAddrs(ctx, req, queue)
}
