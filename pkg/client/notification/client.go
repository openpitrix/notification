// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file

package notification

import (
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/manager"
	"openpitrix.io/notification/pkg/pb"
)

type Client struct {
	pb.NotificationClient
}

func NewClient() (*Client, error) {
	cfg := config.GetInstance().LoadConf()
	conn, err := manager.NewClient(cfg.App.Host, cfg.App.Port)
	if err != nil {
		return nil, err
	}
	return &Client{
		NotificationClient: pb.NewNotificationClient(conn),
	}, nil
}
