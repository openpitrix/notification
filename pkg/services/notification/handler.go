package notification

import (
	"golang.org/x/net/context"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type Handler interface {
	CreateNfWithAddrs(ctx context.Context, in *pb.CreateNfWithAddrsRequest) (*pb.CreateNfResponse, error)
	DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error)
}

type handler struct {
	nfsc Service
}

func NewHandler(nfService Service) Handler {
	return &handler{
		nfsc: nfService,
	}
}

func (h *handler) CreateNfWithAddrs(ctx context.Context, in *pb.CreateNfWithAddrsRequest) (*pb.CreateNfResponse, error) {
	parser := &models.ModelParser{}
	nf, err := parser.CreateNfWithAddrs(in)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}
	logger.Debugf(nil, "%+v", nf.NotificationId)

	nfId, err := h.nfsc.CreateNfWithAddrs(nf)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}

	res := &pb.CreateNfResponse{
		NotificationId: pbutil.ToProtoString(nfId),
	}
	return res, nil
}

func (h *handler) DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	nfId := ""
	nf, err := h.nfsc.DescribeNfs(nfId)
	logger.Debugf(nil, "%+v", nf)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}
	return nil, nil
}
