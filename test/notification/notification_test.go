// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build integration

package notification

import (
	"context"
	"testing"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testAddrsStr := "{\"email\": [\"huojiao2006@163.com\", \"513590612@qq.com\"]}"
	contentStr := "{\"threshold\":80,\"time_series_metrics\":[{\"T\":1243465,\"V\":\"435.4354\"},{\"T\":1243465,\"V\":\"435.4354\"}]}"

	var req = &pb.CreateNotificationRequest{
		ContentType: pbutil.ToProtoString("ContentType"),
		Title:       pbutil.ToProtoString("notification_test.go sends an email."),
		//Content:      pbutil.ToProtoString("Content:handler_test.go sends an email."),
		Content:      pbutil.ToProtoString(contentStr),
		ShortContent: pbutil.ToProtoString("ShortContent"),
		ExpiredDays:  pbutil.ToProtoUInt32(0),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		AddressInfo:  pbutil.ToProtoString(testAddrsStr),
	}

	nfID, err := client.CreateNotification(ctx, req)
	if err != nil {
		t.Log(err)
		t.Fatalf("failed to CreateNotification.")
	}
	t.Log(nfID)
	//t.Log("CreateNotification successfully,nfID = [%s]", nfID)

	t.Log("test notification finish, all tests is ok")

}
