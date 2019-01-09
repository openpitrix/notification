package notification

import (
	"golang.org/x/net/context"
	"openpitrix.io/logger"
	notification "openpitrix.io/notification/pkg/pb_bak"
	"openpitrix.io/notification/pkg/services/test"
	"openpitrix.io/notification/pkg/util/pbutil"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	logger.Debugf(nil, "Test func NewHandler")

	db, q := test.GetTestDBAndEtcd4Test()
	nfservice := NewService(db, q)
	handler := NewHandler(nfservice)
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	logger.Infof(nil, "[%+v]", handler)
}

func TestCreateNfWithAddrs4handler(t *testing.T) {
	logger.Debugf(nil, "Test func NewHandler")
	db, q := test.GetTestDBAndEtcd4Test()
	nfservice := NewService(db, q)

	handler := NewHandler(nfservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testAddrsStr := "johuo@yunify.com;huojiao2006@163.com;513590162@qq.com"

	var req = &notification.CreateNfWithAddrsRequest{
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

	handler.CreateNfWithAddrs(ctx, req)
}
