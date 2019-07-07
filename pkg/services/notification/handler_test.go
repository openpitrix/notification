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

	config.GetInstance()
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	emailcfg := &pb.EmailServiceConfig{
		Protocol:      pbutil.ToProtoString("SMTP"),
		EmailHost:     pbutil.ToProtoString("smtp.qq.com"),
		Port:          pbutil.ToProtoUInt32(25),
		DisplaySender: pbutil.ToProtoString("tester0"),
		Email:         pbutil.ToProtoString("openpitrix@foxmail.com"),
		Password:      pbutil.ToProtoString("iusjafvwmjhddeaf"),
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

	config.GetInstance()
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

	testContentStr := "{\"html\":\"test_content_html\",  \"normal\":\"test_content_normal\"}"

	//testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"openpitrix@163.com\"]}"
	//testAddrsStr := "{\"websocket\": [\"system\", \"huojiao\"]}"
	testAddrsStr := "{\"email\": [\"openpitrix@163.com\", \"openpitrix@163.com\"],\"websocket\": [\"system\", \"huojiao\"]}"
	//testAddrsStr := "[\"adl-LQ2WQlJRzBo8\"]"
	//testAddrsStr := "[\"adl-VDP0l9x1z6k4\"]"

	//TimeFormat := "15:04:05"
	//availableStartTime, _ := time.Parse(TimeFormat, "00:00:00")
	//availableEndTime, _ := time.Parse(TimeFormat, "24:00:00")

	testExtra := "{\"ws_service\": \"ks\",\"ws_message_type\": \"event\"}"

	var req = &pb.CreateNotificationRequest{
		ContentType:        pbutil.ToProtoString("other"),
		Title:              pbutil.ToProtoString("Title_test."),
		Content:            pbutil.ToProtoString(testContentStr),
		ShortContent:       pbutil.ToProtoString("ShortContent_test"),
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
		Address:          pbutil.ToProtoString("openpitrix2@163.com"),
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
	AddressIds = append(AddressIds, "addr-4MMNPyVzqMEr")

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
		Address:          "addr-Q0rqVLlzXGwZ",
		AddressDetail:    pbutil.ToProtoString("test@email.com"),
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
	//addressIds = append(addressIds, "addr-YZrP836RR3ZO")
	//addressIds = append(addressIds, "addr-xNNWLDLvR36K")

	var req = &pb.ModifyAddressListRequest{
		Addresslist:          "adl-BRGg8yq0EZYD",
		AddressListName:      pbutil.ToProtoString("updateTes"),
		Extra:                nil,
		Status:               pbutil.ToProtoString(constants.StatusActive),
		AddressId:            addressIds,
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
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

	config.GetInstance()
	controller, _ := NewController()
	s := &Server{controller: controller}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	iconStr := `data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDI4IiBoZWlnaHQ9IjkwIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxnIGZpbGw9Im5vbmUiIGZpbGwtcnVsZT0iZXZlbm9kZCI+PGcgZmlsbD0iIzQ0M0Q0RSI+PHBhdGggZD0iTTMwNy41MzQgMjNoLTMuNjY2Yy0uOCAwLTEuNDY4LjY4NC0xLjQ2OCAxLjUwNFY0Mi4wMUgyOTAuNlYyNC41MDRjMC0uODItLjYtMS41MDQtMS40NjYtMS41MDRoLTMuNjY4Yy0uOCAwLTEuNDY3LjY4NC0xLjQ2NyAxLjUwNHY0MS45OWMwIC44MjMuNjY3IDEuNTA2IDEuNDY3IDEuNTA2aDMuNjY4Yy44NjYgMCAxLjQ2Ni0uNjgzIDEuNDY2LTEuNTA1VjQ4Ljc4SDMwMi40djE3LjcxNWMwIC44MjIuNjY3IDEuNTA1IDEuNDY4IDEuNTA1aDMuNjY2Yy44NjcgMCAxLjQ2Ni0uNjgzIDEuNDY2LTEuNTA1di00MS45OWMwLS44Mi0uNi0xLjUwNS0xLjQ2Ni0xLjUwNU0yMTggNjIuNzM0YzAtLjgyMS0uNjc3LTEuNTA1LTEuNDg5LTEuNTA1aC0xNi44MVY0OC44ODZoMTYuODFjLjgxMiAwIDEuNDg4LS42MTUgMS40ODgtMS41MDV2LTMuNzZjMC0uODIyLS42NzctMS41MDYtMS40ODktMS41MDZoLTE2LjgwOVYyOS43N2gxNi44MWMuODEgMCAxLjQ4OC0uNjE1IDEuNDg4LTEuNTA0di0zLjc2YzAtLjgyMy0uNjc3LTEuNTA2LTEuNDg5LTEuNTA2aC0xOC4yOTUtMy43MjVjLS44MTQgMC0xLjQ5LjY4My0xLjQ5IDEuNTA0djQxLjk5YzAgLjgyMi42NzYgMS41MDYgMS40OSAxLjUwNkgyMTYuNTFjLjgxMiAwIDEuNDg5LS42MTYgMS40ODktMS41MDV2LTMuNzYxek0zNDAgNjIuNzM0YzAtLjgyMS0uNjc3LTEuNTA1LTEuNDg5LTEuNTA1aC0xNi44MVY0OC44ODZoMTYuODFjLjgxIDAgMS40ODgtLjYxNSAxLjQ4OC0xLjUwNXYtMy43NmMwLS44MjItLjY3Ny0xLjUwNi0xLjQ4OS0xLjUwNmgtMTYuODA5VjI5Ljc3aDE2LjgxYy44MSAwIDEuNDg4LS42MTUgMS40ODgtMS41MDR2LTMuNzZjMC0uODIzLS42NzctMS41MDYtMS40ODktMS41MDZoLTE4LjI5NS0zLjcyNWMtLjgxNCAwLTEuNDkuNjgzLTEuNDkgMS41MDR2NDEuOTljMCAuODIyLjY3NiAxLjUwNiAxLjQ5IDEuNTA2SDMzOC41MWMuODEyIDAgMS40ODktLjYxNiAxLjQ4OS0xLjUwNXYtMy43NjF6TTQwNCA2Mi43MzRjMC0uODIxLS42NzctMS41MDUtMS40ODktMS41MDVoLTE2LjgxVjQ4Ljg4NmgxNi44MWMuODEyIDAgMS40ODgtLjYxNSAxLjQ4OC0xLjUwNXYtMy43NmMwLS44MjItLjY3Ny0xLjUwNi0xLjQ4OS0xLjUwNmgtMTYuODA5VjI5Ljc3aDE2LjgxYy44MSAwIDEuNDg4LS42MTUgMS40ODgtMS41MDR2LTMuNzZjMC0uODIzLS42NzctMS41MDYtMS40ODktMS41MDZoLTE4LjI5NS0zLjcyNWMtLjgxMyAwLTEuNDkuNjgzLTEuNDkgMS41MDR2NDEuOTljMCAuODIyLjY3NyAxLjUwNiAxLjQ5IDEuNTA2SDQwMi41MWMuODEyIDAgMS40ODktLjYxNiAxLjQ4OS0xLjUwNXYtMy43NjF6TTIyOS40NjcgMzQuNzcydjIuNTg2YzAgMS45NzQgMS44NjcgMy4xMyA0IDQuMDE1bDguODY2IDMuODFjMy4zMzMgMS40OTggNS42NjcgNC4xNSA1LjY2NyA4LjIzNHY0LjgzQzI0OCA2Mi4yNjQgMjM5LjczMiA2OSAyMzUuMzMyIDY5Yy0yLjk5OSAwLTcuNTMyLTIuMTc4LTExLjA2NC00LjE1MS0uODAxLS40MDgtMS4zMzQtMS4zNjItLjg2Ny0yLjQ1bDEuMTM0LTIuMzE0Yy40NjYtLjk1MiAxLjQ2NS0xLjE1NyAyLjMzMi0uNjggMi44NjYgMS40OTcgNi42NjYgMy4zMzQgOC40NjUgMy4zMzQgMiAwIDYuMjAxLTMuNDcgNi4yMDEtNS40NDR2LTIuOTkzYzAtMi4zMTQtMS43MzMtMy42MDctNC4xMzMtNC41NTlsLTguMjY3LTMuNTRjLTMuMi0xLjM2LTYuMTMzLTQuMTUtNi4xMzMtNy45NlYzMy44MmMwLTQuNDI0IDguNzMzLTEwLjgyIDEyLjkzNC0xMC44MiAyLjggMCA3LjMzMiAyLjExIDEwLjMzMyAzLjc0My45MzMuNDc3IDEuMiAxLjU2NS43OTggMi4zODFMMjQ2IDMxLjQzOGMtLjMzMy44MTYtMS40IDEuMDktMi4zMzIuNjgtMi4yNjctMS4wODgtNi4wNjctMi44NTctNy43MzMtMi44NTctMS45MzQgMC02LjQ2NyAzLjI2Ni02LjQ2NyA1LjUxMU0zNTkuNDUyIDI4Ljg4MmgtNi42NTd2MTUuMDQ2aDYuNzI1YzIuODE1IDAgNi4xMDktNC45OTIgNi4xMDktNy42NiAwLTIuNTk4LTMuNDMyLTcuMzg2LTYuMTc3LTcuMzg2bS0uMjc1IDIwLjkyN2gtNi4zODJ2MTYuNjg4YzAgLjg4OC0uNjE4IDEuNTAzLTEuNTggMS41MDNoLTMuNjM2Yy0uODI0IDAtMS41NzktLjYxNS0xLjU3OS0xLjUwM3YtNDAuMzVjMC0xLjU3MyAxLjg1NC0zLjE0NyAzLjM2My0zLjE0N2gxMC45ODJDMzY1LjM1NCAyMyAzNzIgMzAuNzg1IDM3MiAzNS45MTRjMCAzLjQyLTMuMDA4IDguNTYtNi4yMzMgMTAuNTQ0IDIuNjA3IDEuMzY3IDYuMTc2IDQuOTkyIDYuMTc2IDguMzQ0djExLjY5NWExLjUyIDEuNTIgMCAwIDEtMS41MSAxLjUwM2gtMy43NzRhMS41MiAxLjUyIDAgMCAxLTEuNTEtMS41MDN2LTEwLjY3YzAtMi4zOTQtMy4zNjItNi4wMTgtNS45NzItNi4wMThNMjY3LjA4OCAyOC44ODJoLTYuNDc3djE1LjA0Nmg2LjU0M2MyLjczOSAwIDUuNTYtNC45OTIgNS41Ni03LjY2IDAtMi41OTgtMi45NTYtNy4zODYtNS42MjYtNy4zODZtLS4yNjcgMjAuOTI3aC02LjIxdjE2LjY4OGMwIC44ODgtLjYwMiAxLjUwMy0xLjUzNiAxLjUwM2gtMy41NGMtLjggMC0xLjUzNS0uNjE1LTEuNTM1LTEuNTAzdi00MC4zNWMwLTEuNTczIDEuODAzLTMuMTQ3IDMuMjcyLTMuMTQ3aDEwLjY4NEMyNzIuODMgMjMgMjc5IDMxLjAyOCAyNzkgMzYuMTU2YzAgMy40Mi0yLjQ3MyA4LjAyLTUuMzMgMTAuNjgtMS40MDUgMS4zMDgtMy42MzYgMi45NzMtNi44NDkgMi45NzNNMTQyLjQ4NCA2OUMxMzguMTM2IDY5IDEzMCA2My4yMzcgMTMwIDU4LjI2M1YyNC44ODRjMC0uODI4LjcyMi0xLjUxNCAxLjUxNy0xLjUxNGgzLjQ4NWMuODUyIDAgMS41MS42ODYgMS41MSAxLjUxNFY1Ny4zN2MwIDIuNDE4IDMuOTMgNS4zNDkgNS45NzIgNS4zNDkgMi4wNDIgMCA1Ljk5OC0zLjMwMiA1Ljk5OC01LjcyVjI0LjUxNGMwLS44My41OTItMS41MTQgMS4zODgtMS41MTRoMy42MTNjLjc5NSAwIDEuNTE3LjY4NSAxLjUxNyAxLjUxNHYzMy4zOEMxNTUgNjIuNzE3IDE0Ni44MjcgNjkgMTQyLjQ4NCA2OU0xNzQuMTU3IDYyLjEyaC02LjM5NFY0Ny4wNzJoNi40NjJjMi44MDEgMCA2LjA4MSA0LjczIDYuMDgxIDcuMzk2IDAgMi42LTMuNDE2IDcuNjUtNi4xNDkgNy42NXptLTYuMzk0LTMzLjI0aDYuMzk0YzIuNzMzIDAgNi4xNDkgMy45MTggNi4xNDkgNi41MTggMCAyLjY2Ny0zLjI4IDcuMDIyLTYuMDgxIDcuMDIyaC02LjQ2MlYyOC44OHptMTIuOTg4IDE1Ljg2N2MzLjExMy0yLjA1MiA2LjI0OS02LjAzOSA2LjI0OS05LjM1QzE4NyAzMC4yNyAxODAuMDMyIDIzIDE3NS4wNDYgMjNoLTEwLjdjLTEuNTAxIDAtMy4zNDYgMS41NzMtMy4zNDYgMy4xNDV2NDAuMzUyYzAgLjg5Ljc1MSAxLjUwMyAxLjU3MSAxLjUwM2gxMi43MDZDMTgwLjI2MyA2OCAxODcgNTkuNiAxODcgNTQuNDdjMC0zLjMxLTMuMTM2LTcuNjQ3LTYuMjQ5LTkuNzIzek0xMDguOTA5IDQ1Ljc1bDE0LjcxNy0xNy4yNDRjLjUzLS42MjIuNTEyLTEuNTQtLjE1My0yLjEyMmwtMi44MTQtMi40NThhMS40ODggMS40ODggMCAwIDAtMi4wOTguMTU1bC0xMi44NyAxNS4wNzhWMjQuNTA1YzAtLjgyLS42MDgtMS41MDUtMS40ODYtMS41MDVoLTMuNzE4Yy0uODExIDAtMS40ODcuNjg0LTEuNDg3IDEuNTA1djQxLjk5MmMwIC44Mi42NzYgMS41MDMgMS40ODcgMS41MDNoMy43MThjLjg3OCAwIDEuNDg3LS42ODMgMS40ODctMS41MDNWNTIuMzM5bDEyLjg2OSAxNS4wOGExLjQ5IDEuNDkgMCAwIDAgMi4wOTguMTU0bDIuODE0LTIuNDU4Yy42NjUtLjU4Mi42ODMtMS41MDEuMTUzLTIuMTIxbC0xNC43MTctMTcuMjQ1ek00MjAuNTEgMjYuNDI1aC0xLjg0djIuNDM0aDEuODRjLjU5NCAwIDEuMjQ1LS40ODEgMS4yNDUtMS4xOSAwLS43NjMtLjY1LTEuMjQ0LTEuMjQ1LTEuMjQ0em0xLjEzMSA2LjAyOGwtMS43ODItMi43MTdoLTEuMTl2Mi43MTdoLS45NjF2LTYuODc3aDIuODAyYzEuMTYgMCAyLjIzNS44MiAyLjIzNSAyLjA5NCAwIDEuNTI5LTEuMzU4IDIuMDM4LTEuNzU0IDIuMDM4bDEuODQgMi43NDVoLTEuMTl6TTQyMCAyMy45MDZBNS4wNTYgNS4wNTYgMCAwIDAgNDE0LjkwNiAyOSA1LjA5IDUuMDkgMCAwIDAgNDIwIDM0LjA5NGE1LjA5IDUuMDkgMCAwIDAgNS4wOTQtNS4wOTNBNS4wNTYgNS4wNTYgMCAwIDAgNDIwIDIzLjkwNnpNNDIwIDM1Yy0zLjMxMSAwLTYtMi42ODktNi02IDAtMy4zMzkgMi42ODktNiA2LTYgMy4zNCAwIDYgMi42NjEgNiA2IDAgMy4zMTEtMi42NiA2LTYgNnoiLz48L2c+PHBhdGggZmlsbD0iIzAwQTk3MSIgZD0iTTY1IDcxLjY0N2wtMTktMTF2MjJ6TTY1IDE5LjY0N2wtMTktMTF2MjJ6TTE5LjY3OCA0NS42NDdMMzcgMzUuNTU2VjMuNjQ3TDEgMjQuNjJ2NDIuMDU1bDM2IDIwLjk3MlY1NS43Mzl6Ii8+PHBhdGggZmlsbD0iIzAwQTk3MSIgZD0iTTM3IDQ1LjY0N2wzNyAyMXYtNDJ6Ii8+PC9nPjwvc3ZnPg==`

	emailcfg := &pb.EmailServiceConfig{
		Protocol:        pbutil.ToProtoString("POP3"),
		EmailHost:       pbutil.ToProtoString("smtp.qq.com"),
		Port:            pbutil.ToProtoUInt32(25),
		DisplaySender:   pbutil.ToProtoString("OpenPitrix"),
		Email:           pbutil.ToProtoString("openpitrix@foxmail.com"),
		Password:        pbutil.ToProtoString("iusjafvwmjhddeaf"),
		SslEnable:       pbutil.ToProtoBool(false),
		ValidationIcon:  pbutil.ToProtoString(iconStr),
		ValidationTitle: pbutil.ToProtoString("[KubeSphere] 测试邮件"),
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
