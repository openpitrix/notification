package notification

import (
	"context"
	"testing"
	"time"

	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

func TestCreateNotification(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("testing env disabled")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"openpitrix@163.com\"]}"
	contentStr := "{\"threshold\":80,\"time_series_metrics\":[{\"T\":1243465,\"V\":\"435.4354\"},{\"T\":1243465,\"V\":\"435.4354\"}]}"

	var req = &pb.CreateNotificationRequest{
		ContentType: pbutil.ToProtoString("ContentType"),
		Title:       pbutil.ToProtoString("handler_test.go sends an email."),
		//Content:      pbutil.ToProtoString("Content:handler_test.go sends an email."),
		Content:      pbutil.ToProtoString(contentStr),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoUInt32(0),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		AddressInfo:  pbutil.ToProtoString(testAddrsStr),
	}
	s.CreateNotification(ctx, req)

}

func TestSetServiceConfig(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("testing env disabled")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	emailcfg := &pb.EmailServiceConfig{
		Protocol:     pbutil.ToProtoString("xx"),
		EmailHost:    pbutil.ToProtoString("testhost"),
		Port:         pbutil.ToProtoString("111"),
		DisplayEmail: pbutil.ToProtoString("test@op.notification.com"),
		Email:        pbutil.ToProtoString("test@op.notification.com"),
		Password:     pbutil.ToProtoString("Email"),
		SslEnable:    pbutil.ToProtoBool(false),
	}

	var req = &pb.ServiceConfig{
		EmailServiceConfig: emailcfg,
	}
	s.SetServiceConfig(ctx, req)

}

func TestGetServiceConfig(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("testing env disabled")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var scTypes []string
	scTypes = append(scTypes, "email")

	var req = &pb.GetServiceConfigRequest{
		ServiceType: scTypes,
	}
	s.GetServiceConfig(ctx, req)
}
