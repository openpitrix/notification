// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"
	"testing"
	"time"

	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

func TestSetServiceConfig(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	controller, _ := NewController()
	s := &Server{controller: controller}
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
	controller, _ := NewController()
	s := &Server{controller: controller}
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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	//testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"openpitrix@163.com\"]}"
	//testAddrsStr := "{\"web\": [\"test_huojiao1\", \"test_huojiao2\"]}"

	testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"openpitrix@163.com\"],\"websocket\": [\"system\", \"huojiao\"]}"
	//testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"openpitrix@163.com\"] }"
	//testAddrsStr := "[\"adl-LQ2WQlJRzBo8\"]"
	//testAddrsStr := "[\"adl-VDP0l9x1z6k4\"]"
	//TimeFormat := "15:04:05"
	//availableStartTime, _ := time.Parse(TimeFormat, "00:00:00")
	//availableEndTime, _ := time.Parse(TimeFormat, "24:00:00")

	testExtra := "{\"ws_service\": \"op\",\"ws_message_type\": \"event\"}"

	var req = &pb.CreateNotificationRequest{
		ContentType:        pbutil.ToProtoString("other"),
		Title:              pbutil.ToProtoString("handler_test.go Title_test."),
		Content:            pbutil.ToProtoString("Content:handler_test.go Content_test."),
		ShortContent:       pbutil.ToProtoString("ShortContent"),
		ExpiredDays:        pbutil.ToProtoUInt32(0),
		Owner:              pbutil.ToProtoString("HuoJiao"),
		AddressInfo:        pbutil.ToProtoString(testAddrsStr),
		AvailableStartTime: pbutil.ToProtoString(""),
		AvailableEndTime:   pbutil.ToProtoString(""),
		Extra:              pbutil.ToProtoString(testExtra),
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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var nfIds []string
	nfIds = append(nfIds, "nf-R3E9xWV7yXnj")
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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var taskIds []string
	taskIds = append(taskIds, "t-Plj28mV7yXnj")

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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var nfIds []string
	nfIds = append(nfIds, "nf-YA4v2ANNMX0D")

	//var contentTypes []string
	//contentTypes = append(contentTypes, "other")

	//var owners []string
	//owners = append(owners, "HuoJiao")

	var statuses []string
	statuses = append(statuses, "successful")
	statuses = append(statuses, "pending")
	statuses = append(statuses, "failed")

	var req = &pb.DescribeNotificationsRequest{
		NotificationId: nfIds,
		//ContentType:    contentTypes,
		//Owner:          owners,
		Status: statuses,
		Limit:  20,
		Offset: 0,
		//SearchWord: pbutil.ToProtoString("successful"),
		SortKey: pbutil.ToProtoString("status"),
		Reverse: pbutil.ToProtoBool(true),
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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var nfIds []string
	nfIds = append(nfIds, "nf-EwMK3Kn55AWZ")
	nfIds = append(nfIds, "nf-PnWMnjJGyXnj")

	var taskIds []string
	taskIds = append(taskIds, "t-B10JwjJGyXnj")

	var statuses []string
	statuses = append(statuses, "successful")

	var req = &pb.DescribeTasksRequest{
		TaskId:         taskIds,
		NotificationId: nfIds,
		ErrorCode:      nil,
		Status:         statuses,
		Limit:          20,
		Offset:         0,
		SearchWord:     pbutil.ToProtoString("successful"),
		SortKey:        pbutil.ToProtoString("status"),
		Reverse:        pbutil.ToProtoBool(true),
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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &pb.CreateAddressRequest{
		Address:          pbutil.ToProtoString("openpitrix@163.com"),
		Remarks:          pbutil.ToProtoString("sss2"),
		VerificationCode: pbutil.ToProtoString("sss3"),
		NotifyType:       pbutil.ToProtoString("email"),
	}

	//var req = &pb.CreateAddressRequest{
	//	Address:          pbutil.ToProtoString("huojiao"),
	//	Remarks:          pbutil.ToProtoString("test_Remarks"),
	//	VerificationCode: pbutil.ToProtoString("test_VerificationCode"),
	//	NotifyType:       pbutil.ToProtoString("websocket"),
	//}

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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var AddressIds []string
	AddressIds = append(AddressIds, "addr-xPgQPnOJM36K")

	var AddressListIds []string
	AddressListIds = append(AddressListIds, "adl-lW4DmnoJWp98")

	var statuses []string
	statuses = append(statuses, "active")

	var nfTypes []string
	nfTypes = append(nfTypes, "email")

	var addrs []string
	addrs = append(addrs, "openpitrix@foxmail.com")

	var req = &pb.DescribeAddressesRequest{
		AddressId:     AddressIds,
		AddressListId: AddressListIds,
		Address:       addrs,
		NotifyType:    nfTypes,
		Status:        statuses,
		Limit:         20,
		Offset:        0,
		SearchWord:    pbutil.ToProtoString("successful"),
		SortKey:       pbutil.ToProtoString("status"),
		Reverse:       pbutil.ToProtoBool(true),
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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var req = &pb.ModifyAddressRequest{
		AddressId:        pbutil.ToProtoString("addr-wRKQzOy7jAWZ"),
		Address:          pbutil.ToProtoString("hello@openpitrix.com"),
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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var addressIds []string
	addressIds = append(addressIds, "addr-wRKQzOy7jAWZ")

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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var addressIds []string
	//addressIds = append(addressIds, "addr-5mwjWpjrZ1YD")
	//addressIds = append(addressIds, "addr-79ME9JRwyM94")
	//addressIds = append(addressIds, "addr-vApolRp19pnj")

	addressIds = append(addressIds, "addr-4MMNPyVzqMEr")
	//addressIds = append(addressIds, "addr-zP8zA0mZqYvj")

	var req = &pb.CreateAddressListRequest{
		AddressListName: pbutil.ToProtoString("通知列表1"),
		//Extra:           pbutil.ToProtoString("{\"email\": [\"openpitrix@163.com\", \"513590612@qq.com\"]}"),
		//Extra:     pbutil.ToProtoString("{}"),
		AddressId: addressIds,
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
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var addressListIds []string
	addressListIds = append(addressListIds, "adl-EA6x97WKl7WZ")

	var statuses []string
	statuses = append(statuses, "active")

	var addressListNames []string
	addressListNames = append(addressListNames, "test")

	var extras []string
	extras = append(extras, "test")

	var req = &pb.DescribeAddressListRequest{
		AddressListId:   addressListIds,
		AddressListName: addressListNames,
		Extra:           extras,
		Status:          statuses,
		Limit:           20,
		Offset:          0,
		SearchWord:      pbutil.ToProtoString("successful"),
		SortKey:         pbutil.ToProtoString("status"),
		Reverse:         pbutil.ToProtoBool(true),
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
	controller, _ := NewController()
	s := &Server{controller: controller}
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
	controller, _ := NewController()
	s := &Server{controller: controller}
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

func TestValidateEmailService(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	emailcfg := &pb.EmailServiceConfig{
		Protocol:      pbutil.ToProtoString("POP3"),
		EmailHost:     pbutil.ToProtoString("smtp.qq.com"),
		Port:          pbutil.ToProtoUInt32(25),
		DisplaySender: pbutil.ToProtoString("OpenPitrix"),
		Email:         pbutil.ToProtoString("openpitrix@foxmail.com"),
		Password:      pbutil.ToProtoString("*********"),
		SslEnable:     pbutil.ToProtoBool(false),
	}

	var req = &pb.ServiceConfig{
		EmailServiceConfig: emailcfg,
	}
	resp, err := s.ValidateEmailService(ctx, req)
	if err != nil {
		t.Fatalf("Test ValidateEmailService failed")
	}

	t.Log(nil, "Test Passed, Test ValidateEmailService", resp.IsSucc)

}
