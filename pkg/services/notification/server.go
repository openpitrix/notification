// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"golang.org/x/net/context"
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
	s := &Server{}

	nfservice := notification.NewService()
	taskservice := task.NewService()

	nfhandler := NewHandler(nfservice, taskservice)
	s.handler = &nfhandler

	taskController := NewController(nfservice, taskservice)
	s.controller = &taskController

	return s, nil
}

func Serve() {
	cfg := config.GetInstance().LoadConf()
	s, _ := NewServer()

	go s.controller.Serve()

	manager.NewGrpcServer("notification-manager", constants.NotificationManagerPort).
		ShowErrorCause(cfg.Grpc.ShowErrorCause).
		//WithBuilder(nil).
		Serve(func(server *grpc.Server) {
			pb.RegisterNotificationServer(server, s)
		})

}

func (s *Server) DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeNfs at server end. "}, nil
}

func (s *Server) CreateNfWithAddrs(ctx context.Context, in *pb.CreateNfWithAddrsRequest) (*pb.CreateNfResponse, error) {
	q := s.controller.queue
	res, err := s.handler.CreateNfWithAddrs(ctx, in, q)
	return res, err
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello,use function SayHello at server end. "}, nil
}
