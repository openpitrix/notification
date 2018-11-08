package nf

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

func TestCreateNfWaddrs4handler(t *testing.T) {
	logger.Debugf(nil, "Test func NewHandler")

	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	nfservice := NewService(db, q)

	handler := NewHandler(nfservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testAddrsStr:="johuo@yunify.com;huojiao2006@163.com"

	var req = &notification.CreateNfWaddrsRequest{
		NfPostType:   pbutil.ToProtoString("Information"),
		NotifyType:   pbutil.ToProtoString("email"),
		AddrsStr:     pbutil.ToProtoString(testAddrsStr),
		Title:        pbutil.ToProtoString("Title Test"),
		Content:      pbutil.ToProtoString("Content"),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoString("7"),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		Status:       pbutil.ToProtoString("New"),
	}

	handler.CreateNfWaddrs(ctx, req)
}

func TestDescribeNfs4handler(t *testing.T) {

}
