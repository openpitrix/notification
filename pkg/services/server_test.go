package services

import (
	"golang.org/x/net/context"
	"openpitrix.io/logger"
	notification "openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	logger.SetLevelByString("debug")
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.SayHello(ctx, &notification.HelloRequest{Name: "unit_test2"})
}

func TestSayHello(t *testing.T) {
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.SayHello(ctx, &notification.HelloRequest{Name: "unit_test2"})
}

func TestCreateNfWaddrs(t *testing.T) {
	InitGlobelSetting()
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &notification.CreateNfWaddrsRequest{
		NfPostType:   pbutil.ToProtoString("Information"),
		NotifyType:   pbutil.ToProtoString("Email"),
		AddrsStr:     pbutil.ToProtoString("johuo@yunify.com;johuo@yunify.com;johuo@yunify.com;huojiao2006@163.com;huojiao2006@163.com;huojiao2006@163.com"),
		Title:        pbutil.ToProtoString("Title Test"),
		Content:      pbutil.ToProtoString("Content"),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoString("7"),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		Status:       pbutil.ToProtoString("New"),
	}
	s.CreateNfWaddrs(ctx, req)
}
