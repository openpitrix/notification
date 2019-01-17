// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package notification

import (
	"golang.org/x/net/context"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/services/notification/service/notification"
	"openpitrix.io/notification/pkg/services/notification/service/task"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/etcd"
)

type Handler struct {
	nfService   notification.Service
	taskService task.Service
}

func NewHandler(nfService notification.Service, taskService task.Service) Handler {
	return Handler{
		nfService:   nfService,
		taskService: taskService,
	}
}

func (h *Handler) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello,use function SayHello at server end. "}, nil
}

func (h *Handler) CreateNfWithAddrs(ctx context.Context, in *pb.CreateNfWithAddrsRequest, q *etcd.Queue) (*pb.CreateNfResponse, error) {
	parser := &models.ModelParser{}
	nf, err := parser.CreateNfWithAddrs(in)
	if err != nil {
		logger.Warnf(ctx, "Failed to parser.CreateNfWithAddrs, error:[%+v]", err)
		return nil, err
	}
	logger.Debugf(ctx, "Success to  parser.CreateNfWithAddrs, NotificationId:[%+s]", nf.NotificationId)

	nfId, err := h.nfService.CreateNfWithAddrs(nf, q)
	if err != nil {
		logger.Warnf(ctx, "Failed to service.CreateNfWithAddrs, error:[%+v]", err)
		return nil, err
	}
	logger.Debugf(ctx, "Success to  service.CreateNfWithAddrs, NotificationId:[%+s]", nf.NotificationId)

	res := &pb.CreateNfResponse{
		NotificationId: pbutil.ToProtoString(nfId),
	}
	return res, nil
}

func (h *Handler) DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	nfId := ""
	nf, err := h.nfService.DescribeNfs(nfId)
	logger.Debugf(ctx, "%+v", nf)
	if err != nil {
		logger.Warnf(ctx, "%+v", err)
		return nil, err
	}
	return nil, nil
}
