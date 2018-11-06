package nf

import (
	"golang.org/x/net/context"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/pb"
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
	logger.Debugf(nil,"Step6:call h.nfservice.SayHello")
	h.nfsc.SayHello("222")
	return nil
}

func (h *handler) CreateNfWaddrs(ctx context.Context, in *pb.CreateNfWaddrsRequest) error {
	logger.Debugf(nil,"Call handlerImpl.CreateNfWaddrs")
	var (
		parser = &NfHandlerModelParser{}
	)
	nf, err := parser.CreateNfWaddrs(in)

	err = h.nfsc.CreateNfWaddrs(nf)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return  err
	}
	return nil
}

func (h *handler) DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	nfId := ""
	nf, err := h.nfsc.DescribeNfs(nfId)
	logger.Debugf(nil, "%+v",nf)
	if err != nil {
		logger.Warnf(nil, "%+v", err)
		return  nil, nil
	}
	return nil, nil
}
