// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file

package notification

import (
	"strconv"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/manager"
	"openpitrix.io/notification/pkg/pb"
)

type Client struct {
	pb.NotificationClient
}

func NewClient() (*Client, error) {
	cfg := config.GetInstance().LoadConf()
	notificationManagerHost := cfg.App.Host
	notificationManagerPort, _ := strconv.Atoi(cfg.App.Port)

	conn, err := manager.NewClient(notificationManagerHost, notificationManagerPort)
	if err != nil {
		return nil, err
	}
	return &Client{
		NotificationClient: pb.NewNotificationClient(conn),
	}, nil
}
