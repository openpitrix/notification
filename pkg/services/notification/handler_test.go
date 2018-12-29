package notification

import (
	"golang.org/x/net/context"
	"openpitrix.io/logger"
	notification "openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/services/test"
	"openpitrix.io/notification/pkg/util/pbutil"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	logger.Debugf(nil, "Test func NewHandler")

	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	nfservice := NewService(db, q)

	handler := NewHandler(nfservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var req = &notification.HelloRequest{Name: "hello world."}

	handler.SayHello(ctx, req)
}

func TestSayHello4handler(t *testing.T) {
	logger.Debugf(nil, "Test func NewHandler")

	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	nfservice := NewService(db, q)

	handler := NewHandler(nfservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var req = &notification.HelloRequest{Name: "hello world."}

	handler.SayHello(ctx, req)
}

func TestCreateNfWithAddrs4handler(t *testing.T) {
	logger.Debugf(nil, "Test func NewHandler")

	db := test.GetTestDB()
	q := test.GetEtcdQueue()
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
