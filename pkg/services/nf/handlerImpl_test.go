package nf

import (
	"golang.org/x/net/context"
	"log"
	notification "openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/dbutil"
	"openpitrix.io/notification/pkg/util/pbutil"
	"os"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	log.Println("Test func NewHandler")

	issucc := dbutil.GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}
	db := dbutil.GetInstance().GetMysqlDB()
	nfservice := NewService(db)

	handler := NewHandler(nfservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var req = &notification.HelloRequest{Name: "hello world."}

	handler.SayHello(ctx, req)
}

func TestSayHello2(t *testing.T) {

	log.Println("Test func NewHandler")

	issucc := dbutil.GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}
	db := dbutil.GetInstance().GetMysqlDB()
	nfservice := NewService(db)

	handler := NewHandler(nfservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var req = &notification.HelloRequest{Name: "hello world."}

	handler.SayHello(ctx, req)
}

func TestCreateNfWaddrs2(t *testing.T) {
	log.Println("Test func NewHandler")

	issucc := dbutil.GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}

	db := dbutil.GetInstance().GetMysqlDB()
	nfservice := NewService(db)

	handler := NewHandler(nfservice)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &notification.CreateNfWaddrsRequest{
		NfPostType:   pbutil.ToProtoString("Information"),
		NotifyType:   pbutil.ToProtoString("email"),
		AddrsStr:     pbutil.ToProtoString("johuo@yunify.com;danma@yunify.com"),
		Title:        pbutil.ToProtoString("Title Test"),
		Content:      pbutil.ToProtoString("Content"),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoString("7"),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		Status:       pbutil.ToProtoString("New"),
	}

	handler.CreateNfWaddrs(ctx, req)
}
