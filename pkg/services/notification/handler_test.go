// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"testing"
	"time"

	"openpitrix.io/logger"
	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

func TestSetServiceConfig(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
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
		t.Skip("Local Dev testing env disabled.")
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

func TestCreateNotification(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"openpitrix@163.com\"]}"

	var req = &pb.CreateNotificationRequest{
		ContentType:  pbutil.ToProtoString("ContentType"),
		Title:        pbutil.ToProtoString("handler_test.go sends an email."),
		Content:      pbutil.ToProtoString("Content:handler_test.go sends an email."),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoUInt32(0),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		AddressInfo:  pbutil.ToProtoString(testAddrsStr),
	}
	s.CreateNotification(ctx, req)

}

func TestDescribeNotifications4Handler(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var nfIds []string
	nfIds = append(nfIds, "nf-yM793AqkEmnj")
	nfIds = append(nfIds, "nf-p1J10q82WnZO")

	var contentTypes []string
	contentTypes = append(contentTypes, "email")

	var owners []string
	owners = append(owners, "HuoJiao")

	var statuses []string
	statuses = append(statuses, "successful")

	var displayCols []string
	displayCols = append(displayCols, "")

	var req = &pb.DescribeNotificationsRequest{
		NotificationId: nfIds,
		ContentType:    contentTypes,
		Owner:          owners,
		Status:         statuses,
		Limit:          20,
		Offset:         0,
		SearchWord:     nil,
		SortKey:        pbutil.ToProtoString("status"),
		Reverse:        pbutil.ToProtoBool(false),
		DisplayColumns: displayCols,
	}
	resp, _ := s.DescribeNotifications(ctx, req)
	logger.Infof(nil, "Test Passed,TestDescribeNotifications4Handler TotalCount = %d", resp.GetTotalCount())

}

func TestDescribeTasks(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var nfIds []string
	nfIds = append(nfIds, "nf-yM793AqkEmnj")
	nfIds = append(nfIds, "nf-p1J10q82WnZO")

	var statuses []string
	statuses = append(statuses, "successful")

	var displayCols []string
	displayCols = append(displayCols, "")

	var req = &pb.DescribeTasksRequest{
		NotificationId: nfIds,
		TaskId:         nil,
		TaskAction:     nil,
		ErrorCode:      nil,
		Status:         statuses,
		Limit:          20,
		Offset:         0,
		SearchWord:     nil,
		SortKey:        nil,
		Reverse:        nil,
		DisplayColumns: nil,
	}
	resp, _ := s.DescribeTasks(ctx, req)
	logger.Infof(nil, "Test Passed,Test DescribeTasks TotalCount = %d", resp.GetTotalCount())

}

func TestCreateAddress(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &pb.CreateAddressRequest{
		//AddressListId:    nil,
		Address:          pbutil.ToProtoString("sss1"),
		Remarks:          pbutil.ToProtoString("sss2"),
		VerificationCode: pbutil.ToProtoString("sss3"),
		NotifyType:       pbutil.ToProtoString("sss4"),
	}
	_, e := s.CreateAddress(ctx, req)
	if e != nil {
		logger.Criticalf(nil, "Test CreateAddress failed...")
	}

}

func TestModifyAddress(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &pb.ModifyAddressRequest{
		AddressId:        pbutil.ToProtoString("addr-BPrjMv8Qr4Yr"),
		AddressListId:    pbutil.ToProtoString("xxAddressListId"),
		Address:          pbutil.ToProtoString("hello@openpitrix.com"),
		Remarks:          pbutil.ToProtoString("测试Remarks"),
		VerificationCode: pbutil.ToProtoString("VerificationCode test"),
		NotifyType:       pbutil.ToProtoString("email"),
	}

	_, e := s.ModifyAddress(ctx, req)
	if e != nil {
		logger.Criticalf(nil, "Test ModifyAddress failed...")
	}

}

func TestDescribeAddresses(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var AddressId []string
	AddressId = append(AddressId, "addr-BPrjMv8Qr4Yr")

	var statuses []string
	statuses = append(statuses, "active")

	var displayCols []string
	displayCols = append(displayCols, "")

	var req = &pb.DescribeAddressesRequest{
		AddressId:      AddressId,
		AddressListId:  nil,
		Address:        nil,
		NotifyType:     nil,
		Status:         statuses,
		Limit:          20,
		Offset:         0,
		SearchWord:     nil,
		SortKey:        nil,
		Reverse:        nil,
		DisplayColumns: nil,
	}
	resp, _ := s.DescribeAddresses(ctx, req)
	logger.Infof(nil, "Test Passed,Test DescribeAddresses TotalCount = %d", resp.GetTotalCount())

}
