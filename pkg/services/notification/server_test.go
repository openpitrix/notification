package notification

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
	s, _ := NewServer()
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	logger.Infof(nil, "[%+v]", s)
}

func TestCreateNfWithAddrs(t *testing.T) {
	s, _ := NewServer()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	testAddrsStr := "huojiao2006@163.com;513590612@qq.com"
	var req = &notification.CreateNfWithAddrsRequest{
		ContentType:  pbutil.ToProtoString("Information"),
		SentType:     pbutil.ToProtoString("Email"),
		AddrsStr:     pbutil.ToProtoString(testAddrsStr),
		Title:        pbutil.ToProtoString("Run case"),
		Content:      pbutil.ToProtoString("Run case content"),
		ShortContent: pbutil.ToProtoString("Run case ShortContent"),
		ExpiredDays:  pbutil.ToProtoString("7"),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		Status:       pbutil.ToProtoString("New"),
	}
	s.CreateNfWithAddrs(ctx, req)
}
