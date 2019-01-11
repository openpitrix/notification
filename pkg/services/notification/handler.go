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
	nfsc   notification.Service
	tasksc task.Service
}

func NewHandler(nfService notification.Service, tasksc task.Service) Handler {
	return Handler{
		nfsc:   nfService,
		tasksc: tasksc,
	}
}

func (h *Handler) CreateNfWithAddrs(ctx context.Context, in *pb.CreateNfWithAddrsRequest, q *etcd.Queue) (*pb.CreateNfResponse, error) {
	parser := &models.ModelParser{}
	nf, err := parser.CreateNfWithAddrs(in)
	if err != nil {
		logger.Warnf(nil, "Failed to parser.CreateNfWithAddrs, error:[%+v]", err)
		return nil, err
	}
	logger.Debugf(nil, "Success to  parser.CreateNfWithAddrs, NotificationId:[%+s]", nf.NotificationId)

	nfId, err := h.nfsc.CreateNfWithAddrs(nf, q)
	if err != nil {
		logger.Warnf(nil, "Failed to service.CreateNfWithAddrs, error:[%+v]", err)
		return nil, err
	}
	logger.Debugf(nil, "Success to  service.CreateNfWithAddrs, NotificationId:[%+s]", nf.NotificationId)

	res := &pb.CreateNfResponse{
		NotificationId: pbutil.ToProtoString(nfId),
	}
	return res, nil
}

func (h *Handler) DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	nfId := ""
	nf, err := h.nfsc.DescribeNfs(nfId)
	logger.Debugf(nil, "%+v", nf)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}
	return nil, nil
}
