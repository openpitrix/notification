package nf

import (
	"golang.org/x/net/context"
	"log"
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

func (h *handler) SayHello(ctx context.Context, in *pb.HelloRequest) (error) {
	log.Println("Step6:call h.nfservice.SayHello")
	h.nfsc.SayHello("222")
	return nil
}

func (h *handler) CreateNfWaddrs(ctx context.Context, in *pb.CreateNfWaddrsRequest) (error) {
	log.Println("Call handlerImpl.CreateNfWaddrs")
	var (
		parser = &NfHandlerModelParser{}
	)
	nf, _ := parser.CreateNfWaddrs(in)

	log.Print(nf.AddrsStr)
	err := h.nfsc.CreateNfWaddrs(nf)
	if err != nil {
		log.Println("something is wrong")
	}
	return nil
}

