// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package notification

import (
	"context"
	"strconv"
	"testing"
	"time"

	"openpitrix.io/logger"

	nfclient "openpitrix.io/notification/pkg/client/notification"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

const Service = "notification-service"

func TestNotification(t *testing.T) {
	client, err := nfclient.NewClient()
	if err != nil {
		t.Fatalf("failed to create nfclient.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	testAddrsStr := "{\"email\": [\"openpitrix@163.com\"]}"
	contentStr := "test content"

	var req = &pb.CreateNotificationRequest{
		ContentType:  pbutil.ToProtoString("other"),
		Title:        pbutil.ToProtoString("handler_test.go Title_test."),
		Content:      pbutil.ToProtoString(contentStr),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoUInt32(0),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		AddressInfo:  pbutil.ToProtoString(testAddrsStr),
	}

	_, err = client.CreateNotification(ctx, req)
	if err != nil {
		t.Log(err)
		t.Fatalf("failed to CreateNotification.")
	}

	t.Log("create notification successfully.")

}

func createNF2(i string) {
	client, err := nfclient.NewClient()
	if err != nil {
		logger.Errorf(nil, "failed to create nfclient.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	testAddrsStr := "{\"email\": [\"admin@app-center.com.cn\"]}"
	contentStr := "content for pressure Testing"
	var req = &pb.CreateNotificationRequest{
		ContentType:  pbutil.ToProtoString("other"),
		Title:        pbutil.ToProtoString("Pressure Testing"),
		Content:      pbutil.ToProtoString(contentStr),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoUInt32(0),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		AddressInfo:  pbutil.ToProtoString(testAddrsStr),
	}

	s := req.Content.GetValue() + ",第" + i + "封邮件"
	req.Content = pbutil.ToProtoString(s)
	_, err = client.CreateNotification(ctx, req)
	if err != nil {
		logger.Errorf(nil, "failed to CreateNotification,err= %+v", err)
	}
	logger.Infof(nil, "create notification successfully,i= %+v", i)

}

const (
	Maxtasks2 = 1000
)

func TestCreateNotificationByPressure(t *testing.T) {
	for i := 0; i < Maxtasks2; i++ {
		go createNF2(strconv.Itoa(i))
	}

	for {
		time.Sleep(time.Second * 1)
	}
}
