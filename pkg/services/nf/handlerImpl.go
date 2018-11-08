package nf

import (
	"golang.org/x/net/context"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type handler struct {
	nfsc Service
}

func NewHandler(nfService Service) Handler {
	return &handler{
		nfsc: nfService,
	}
}

func (h *handler) SayHello(ctx context.Context, in *pb.HelloRequest) error {
	logger.Debugf(nil, "Step6:call h.nfservice.SayHello")
	h.nfsc.SayHello("222")
	return nil
}

func (h *handler) CreateNfWaddrs(ctx context.Context, in *pb.CreateNfWaddrsRequest) (*pb.CreateNfResponse, error) {
	var (
		parser = &NfHandlerModelParser{}
	)

	nf, err := parser.CreateNfWaddrs(in)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}
	logger.Debugf(nil, "%+v", nf.NfPostID)

	nfPostID, err := h.nfsc.CreateNfWaddrs(nf)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return nil, err
	}

	res := &pb.CreateNfResponse{
		NfPostId: pbutil.ToProtoString(nfPostID),
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
