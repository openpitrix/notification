// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"google.golang.org/grpc"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/manager"
	"openpitrix.io/notification/pkg/pb"
	rs "openpitrix.io/notification/pkg/services/notification/resource_control"
)

type Server struct {
	controller *Controller
}

func Serve() {
	cfg := config.GetInstance()

	//read email config data from db.
	// check the data in data is default data or not
	// 1.if data in DB is default data, use cfg from ENV to update data in DB.
	// 2.if data in DB is not default data, use the data in DB.
	err := rs.ResetEmailCfg(cfg)
	if err != nil {
		logger.Errorf(nil, "Failed to reset email config: %+v.", err)
	}

	controller, err := NewController()
	if err != nil {
		logger.Criticalf(nil, "Failed to start serve: %+v.", err)
	}
	s := &Server{controller: controller}

	/**********************************************************
	** start controller **
	**********************************************************/
	logger.Infof(nil, "[%s]", "/**********************************************************")
	logger.Infof(nil, "[%s]", "** start controller **")
	logger.Infof(nil, "[%s]", "**********************************************************/")
	go s.controller.Serve()

	/**********************************************************
	** start ServeApiGateway **
	**********************************************************/
	logger.Infof(nil, "[%s]", "/**********************************************************")
	logger.Infof(nil, "[%s]", "** start ServeApiGateway **")
	logger.Infof(nil, "[%s]", "**********************************************************/")
	go ServeApiGateway()

	manager.NewGrpcServer(cfg.App.Host, cfg.App.Port).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		WithChecker(s.Checker).
		Serve(func(server *grpc.Server) {
			pb.RegisterNotificationServer(server, s)
		})
}
