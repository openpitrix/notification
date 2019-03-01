// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"strconv"

	"google.golang.org/grpc"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/manager"
	"openpitrix.io/notification/pkg/pb"
)

type Server struct {
	controller *Controller
}

func Serve() {
	cfg := config.GetInstance()
	s := &Server{
		controller: NewController(),
	}

	go s.controller.Serve()
	go ServeApiGateway()

	notificationManagerHost := cfg.App.Host
	notificationManagerPort, _ := strconv.Atoi(cfg.App.Port)

	manager.NewGrpcServer(notificationManagerHost, notificationManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		Serve(func(server *grpc.Server) {
			pb.RegisterNotificationServer(server, s)
		})
}
