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

	//testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"openpitrix@163.com\"]}"
	testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"513590612@qq.com\"]}"
	var req = &pb.CreateNotificationRequest{
		ContentType:  pbutil.ToProtoString("other"),
		Title:        pbutil.ToProtoString("handler_test.go sends an email."),
		Content:      pbutil.ToProtoString("Content:handler_test.go sends an email."),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoUInt32(0),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		AddressInfo:  pbutil.ToProtoString(testAddrsStr),
	}
	s.CreateNotification(ctx, req)

}

func TestDescribeNotifications(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var nfIds []string
	nfIds = append(nfIds, "nf-VLjBx1nMOB94")

	var contentTypes []string
	contentTypes = append(contentTypes, "other")

	var owners []string
	owners = append(owners, "HuoJiao")

	var statuses []string
	statuses = append(statuses, "successful")

	var req = &pb.DescribeNotificationsRequest{
		NotificationId: nfIds,
		ContentType:    contentTypes,
		Owner:          owners,
		Status:         statuses,
		//Limit:      20,
		//Offset:     0,
		SearchWord: nil,
		SortKey:    pbutil.ToProtoString("status"),
		Reverse:    pbutil.ToProtoBool(false),
	}
	resp, _ := s.DescribeNotifications(ctx, req)
	logger.Infof(nil, "Test Passed,TestDescribeNotifications Notifications = [%+v]", resp)

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
	nfIds = append(nfIds, "nf-6YoKxDk9BLWZ")
	nfIds = append(nfIds, "nf-WpQ8pmVBvJ98")

	var taskIds []string
	taskIds = append(taskIds, "t-7k1BMPnAq3zn")

	var statuses []string
	statuses = append(statuses, "successful")

	var req = &pb.DescribeTasksRequest{
		NotificationId: nfIds,
		TaskId:         taskIds,
		TaskAction:     nil,
		ErrorCode:      nil,
		Status:         statuses,
		SearchWord:     nil,
		SortKey:        nil,
		Reverse:        nil,
	}
	resp, _ := s.DescribeTasks(ctx, req)
	logger.Infof(nil, "Test Passed,Test DescribeTasks Tasks = [%+v]", resp)

}

func TestRetryNotifications(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var nfIds []string
	nfIds = append(nfIds, "nf-41ZYN9zyx3Yr")

	var statuses []string
	statuses = append(statuses, "successful")

	var req = &pb.RetryNotificationsRequest{
		NotificationId: nfIds,
	}
	resp, _ := s.RetryNotifications(ctx, req)

	logger.Infof(nil, "Test Passed,Test RetryNotifications[%+v]", resp.NotificationSet)

}

func TestRetryTasks(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var taskIds []string
	taskIds = append(taskIds, "t-l32A2Yvnq3zn")

	var statuses []string
	statuses = append(statuses, "successful")

	var req = &pb.RetryTasksRequest{
		TaskId: taskIds,
	}
	resp, _ := s.RetryTasks(ctx, req)

	logger.Infof(nil, "Test Passed,Test RetryTasks,Tasks=[%+v]", resp)

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
		Address:          pbutil.ToProtoString("openpitrix@163.com"),
		Remarks:          pbutil.ToProtoString("sss2"),
		VerificationCode: pbutil.ToProtoString("sss3"),
		NotifyType:       pbutil.ToProtoString("email"),
	}
	_, e := s.CreateAddress(ctx, req)
	if e != nil {
		logger.Criticalf(nil, "Test CreateAddress failed...")
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

	var AddressIds []string
	AddressIds = append(AddressIds, "addr-RGN5PwjjDp6K")

	var statuses []string
	statuses = append(statuses, "active")

	var nfTypes []string
	nfTypes = append(nfTypes, "email")

	var req = &pb.DescribeAddressesRequest{
		AddressId:     AddressIds,
		AddressListId: nil,
		Address:       nil,
		Status:        statuses,
		NotifyType:    nfTypes,
		SearchWord:    nil,
		SortKey:       nil,
		Reverse:       nil,
	}
	resp, _ := s.DescribeAddresses(ctx, req)
	logger.Infof(nil, "Test Passed,Test DescribeAddresses Addresses = [%+v]", resp)

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
		AddressId:        pbutil.ToProtoString("addr-RGN5PwjjDp6K"),
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

func TestDeleteAddresses(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var addressIds []string
	addressIds = append(addressIds, "addr-106wppjyPB94")

	var req = &pb.DeleteAddressesRequest{
		AddressId: addressIds,
	}

	_, e := s.DeleteAddresses(ctx, req)
	if e != nil {
		logger.Criticalf(nil, "Test DeleteAddresses failed...")
	}

}

func TestCreateAddressList(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var addressIds []string
	addressIds = append(addressIds, "addr-BPrjMv8Qr4Yr")

	var req = &pb.CreateAddressListRequest{
		AddressListName: pbutil.ToProtoString("邮件通知列表1"),
		Extra:           pbutil.ToProtoString("{\"email\": [\"openpitrix@163.com\", \"513590612@qq.com\"]}"),
		AddressId:       addressIds,
	}
	_, e := s.CreateAddressList(ctx, req)
	if e != nil {
		logger.Criticalf(nil, "Test CreateAddressList failed...")
	}
}

func TestDescribeAddressList(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var addressListIds []string
	addressListIds = append(addressListIds, "adl-6ODJlDYg3LWZ")

	var statuses []string
	statuses = append(statuses, "active")

	var req = &pb.DescribeAddressListRequest{
		AddressListId:   addressListIds,
		AddressListName: nil,
		Extra:           nil,
		Status:          statuses,
		SearchWord:      nil,
		SortKey:         nil,
		Reverse:         nil,
	}
	resp, _ := s.DescribeAddressList(ctx, req)
	logger.Infof(nil, "Test Passed,Test DescribeAddressesList AddressList = [%+v]", resp)

}

func TestModifyAddressList(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &pb.ModifyAddressListRequest{
		AddressListId:   pbutil.ToProtoString("adl-6ODJlDYg3LWZ"),
		AddressListName: pbutil.ToProtoString("测试修改邮件通知列表1"),
		Extra:           nil,
	}

	_, e := s.ModifyAddressList(ctx, req)
	if e != nil {
		logger.Criticalf(nil, "Test ModifyAddressList failed...")
	}

}

func TestDeleteAddressList(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var adls []string
	adls = append(adls, "adl-6ODJlDYg3LWZ")

	var req = &pb.DeleteAddressListRequest{
		AddressListId: adls,
	}

	_, e := s.DeleteAddressList(ctx, req)
	if e != nil {
		logger.Criticalf(nil, "Test DeleteAddressList failed...")
	}

}
