// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
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

	manager.NewGrpcServer(cfg.App.Host, cfg.App.Port).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		Serve(func(server *grpc.Server) {
			pb.RegisterNotificationServer(server, s)
		})
}
