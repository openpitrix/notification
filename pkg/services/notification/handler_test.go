// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"testing"
	"time"

	"openpitrix.io/notification/pkg/constants"

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
		Protocol:      pbutil.ToProtoString("POP3"),
		EmailHost:     pbutil.ToProtoString("testhost"),
		Port:          pbutil.ToProtoUInt32(888),
		DisplaySender: pbutil.ToProtoString("tester"),
		Email:         pbutil.ToProtoString("test@op.notification.com"),
		Password:      pbutil.ToProtoString("passwordtest"),
		SslEnable:     pbutil.ToProtoBool(false),
	}

	var req = &pb.ServiceConfig{
		EmailServiceConfig: emailcfg,
	}
	resp, err := s.SetServiceConfig(ctx, req)
	if err != nil {
		t.Fatalf("Test SetServiceConfig failed")
	}

	t.Log(nil, "Test Passed, Test SetServiceConfig", resp.IsSucc)

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
	scTypes = append(scTypes, constants.NotifyTypeEmail)

	var req = &pb.GetServiceConfigRequest{
		ServiceType: scTypes,
	}
	resp, err := s.GetServiceConfig(ctx, req)
	if err != nil {
		t.Fatalf("Test SetServiceConfig failed")
	}

	t.Log(nil, "Test Passed, Test SetServiceConfig", resp.EmailServiceConfig)
}

func TestCreateNotification(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"513590612@qq.com\"]}"
	//testAddrsStr := "{\"email\": [ \"513590612@qq.com\"]}"
	//testAddrListIds := "[\"adl-RWAZ8kZ39wzn\"]"
	var req = &pb.CreateNotificationRequest{
		ContentType:  pbutil.ToProtoString("other"),
		Title:        pbutil.ToProtoString("handler_test.go sends an email."),
		Content:      pbutil.ToProtoString("Content:handler_test.go sends an email."),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoUInt32(0),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		//AddressInfo:  pbutil.ToProtoString(testAddrListIds),
		AddressInfo: pbutil.ToProtoString(testAddrsStr),
	}
	resp, err := s.CreateNotification(ctx, req)
	if err != nil {
		t.Fatalf("Test CreateNotification failed")
	}

	t.Log(nil, "Test Passed, Test CreateNotification", resp)
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
	nfIds = append(nfIds, "nf-ByERxoV2lAZO")
	var req = &pb.RetryNotificationsRequest{
		NotificationId: nfIds,
	}
	resp, err := s.RetryNotifications(ctx, req)
	if err != nil {
		t.Fatalf("Test Retry Notifications failed[%s]", nfIds)
	}
	t.Log(nil, "Test Passed, Test RetryNotifications", resp.NotificationSet)

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
	taskIds = append(taskIds, "t-o4WE3lPx8R98")

	var req = &pb.RetryTasksRequest{
		TaskId: taskIds,
	}
	resp, err := s.RetryTasks(ctx, req)
	if err != nil {
		t.Fatalf("TestRetryTasks failed[%s]", taskIds)
	}

	t.Log(nil, "Test Passed, TestRetryTasks", resp.TaskSet)

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
	nfIds = append(nfIds, "nf-0YK516ArgM94")

	var contentTypes []string
	contentTypes = append(contentTypes, "other")

	var owners []string
	owners = append(owners, "HuoJiao")

	var statuses []string
	statuses = append(statuses, "successful")
	statuses = append(statuses, "pending")

	var req = &pb.DescribeNotificationsRequest{
		NotificationId: nfIds,
		ContentType:    contentTypes,
		Owner:          owners,
		Status:         statuses,
		//Limit:      20,
		//Offset:     0,
		SearchWord: pbutil.ToProtoString("successful"),
		SortKey:    pbutil.ToProtoString("status"),
		Reverse:    pbutil.ToProtoBool(true),
	}
	resp, err := s.DescribeNotifications(ctx, req)
	if err != nil {
		t.Fatalf("Test Describe Notifications failed[%s]", nfIds)
	}

	t.Log(nil, "Test Passed, TestDescribeNotifications", resp)
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
	resp, err := s.DescribeTasks(ctx, req)
	if err != nil {
		t.Fatalf("Test DescribeTasks failed[%s]", nfIds)
	}

	t.Log(nil, "Test Passed, TestDescribeTasks", resp)

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
	resp, err := s.CreateAddress(ctx, req)
	if err != nil {
		t.Fatalf("TestCreateAddress failed")
	}
	t.Log(nil, "Test Passed, TestCreateAddress", resp)

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
	AddressIds = append(AddressIds, "addr-xPgQPnOJM36K")

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
	resp, err := s.DescribeAddresses(ctx, req)
	if err != nil {
		t.Fatalf("TestDescribeAddresses failed[%s]", AddressIds)
	}
	t.Log(nil, "Test Passed, TestDescribeAddresses", resp)

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
		AddressId:        pbutil.ToProtoString("addr-W2LE6NoKGpnj"),
		Address:          pbutil.ToProtoString("hello1@openpitrix.com"),
		Remarks:          pbutil.ToProtoString("测试Remarks2211"),
		VerificationCode: pbutil.ToProtoString("VerificationCode test"),
		NotifyType:       pbutil.ToProtoString("email"),
	}

	resp, err := s.ModifyAddress(ctx, req)
	if err != nil {
		t.Fatalf("TestModifyAddress failed.")
	}
	t.Log(nil, "Test Passed, TestModifyAddress", resp)

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
	addressIds = append(addressIds, "addr-xPgQPnOJM36K11")

	var req = &pb.DeleteAddressesRequest{
		AddressId: addressIds,
	}

	resp, err := s.DeleteAddresses(ctx, req)
	if err != nil {
		t.Fatalf("TestDeleteAddresses failed[%s]", addressIds)
	}
	t.Log(nil, "Test Passed, TestDeleteAddresses", resp)

}

func TestCreateAddressList(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//var addressIds []string
	//addressIds = append(addressIds, "addr-2pZGl631lAYr")
	//addressIds = append(addressIds, "addr-79ME9JRwyM94")
	//addressIds = append(addressIds, "addr-vApolRp19pnj")

	var req = &pb.CreateAddressListRequest{
		AddressListName: pbutil.ToProtoString("邮件通知列表1"),
		Extra:           pbutil.ToProtoString("{\"email\": [\"openpitrix@163.com\", \"513590612@qq.com\"]}"),
		//AddressId:       addressIds,
	}
	resp, err := s.CreateAddressList(ctx, req)
	if err != nil {
		t.Fatalf("TestCreateAddressList failed")
	}
	t.Log(nil, "Test Passed, TestCreateAddressList", resp)
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
	addressListIds = append(addressListIds, "adl-EA6x97WKl7WZ")

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
	resp, err := s.DescribeAddressList(ctx, req)
	if err != nil {
		t.Fatalf("TestDescribeAddressList failed[%s]", addressListIds)
	}
	t.Log(nil, "Test Passed, TestDescribeAddressList", resp)

}

func TestModifyAddressList(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	s := &Server{controller: NewController()}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var addressIds []string
	addressIds = append(addressIds, "addr-vVlnxXyRR7WZ")
	addressIds = append(addressIds, "addr-YZrP836RR3ZO")
	addressIds = append(addressIds, "addr-xNNWLDLvR36K")

	var req = &pb.ModifyAddressListRequest{
		AddressListId:   pbutil.ToProtoString("adl-2Enmg1466AYr"),
		AddressListName: pbutil.ToProtoString("updateTes"),
		Status:          pbutil.ToProtoString(constants.StatusActive),
		AddressId:       addressIds,
	}

	resp, err := s.ModifyAddressList(ctx, req)
	if err != nil {
		t.Fatalf("TestModifyAddressList failed[%s]", addressIds)
	}
	t.Log(nil, "Test Passed, TestModifyAddressList", resp)

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
	adls = append(adls, "adl-2Enmg1466AYr")

	var req = &pb.DeleteAddressListRequest{
		AddressListId: adls,
	}

	resp, err := s.DeleteAddressList(ctx, req)
	if err != nil {
		t.Fatalf("TestDeleteAddressList failed[%s]", adls)
	}
	t.Log(nil, "Test Passed, TestDeleteAddressList", resp)

}
