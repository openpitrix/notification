// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"context"

	"google.golang.org/grpc"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/services/notification/service/notification"
	"openpitrix.io/notification/pkg/services/notification/service/task"
	"openpitrix.io/openpitrix/pkg/manager"
)

type Server struct {
	handler    *Handler
	controller *Controller
}

func NewServer() (*Server, error) {
	s := new(Server)

	nfService := notification.NewService()
	taskService := task.NewService()

	nfHandler := NewHandler(nfService, taskService)
	s.handler = &nfHandler

	taskController := NewController(nfService, taskService)
	s.controller = &taskController

	return s, nil
}

func Serve() {
	cfg := config.GetInstance().LoadConf()
	s, _ := NewServer()

	go s.controller.Serve()

	manager.NewGrpcServer("notification-manager", constants.NotificationManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		Serve(func(server *grpc.Server) {
			pb.RegisterNotificationServer(server, s)
		})
}

func (s *Server) DescribeNfs(ctx context.Context, req *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeNfs at server end. "}, nil
}

func (s *Server) CreateNfWithAddrs(ctx context.Context, req *pb.CreateNfWithAddrsRequest) (*pb.CreateNfWithAddrsResponse, error) {
	q := s.controller.queue
	res, err := s.handler.CreateNfWithAddrs(ctx, req, q)
	return res, err
}
