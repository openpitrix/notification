package services

import (
	"golang.org/x/net/context"
	"log"
	"notification/pkg/pb"
	"notification/pkg/services/nf"
)


// SayHello implements nf.RegisterNotificationServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println("Hello,use function SayHello at server end.")
	ss:=in.GetName();
	s.nfservice.SayHello("hhhhhhhh:"+ss)
	return &pb.HelloReply{Message: "Hello,use function SayHello at server end. " + in.Name}, nil
}

func (s *Server) CreateNfWUserFilter(ctx context.Context, in *pb.CreateNfWUserFilterRequest) (*pb.CreateNfResponse, error) {
	log.Println("Hello,use function CreateNfWUserFilter at server end.")
	return &pb.CreateNfResponse{Message: "1111 "  }, nil
}

func (s *Server) CreateNfWAppFilter(ctx context.Context, in *pb.CreateNfWAppFilterRequest) (*pb.CreateNfResponse, error) {
	log.Println("Hello,use function CreateNfWAppFilter at server end.")
	return &pb.CreateNfResponse{Message: "Hello,use function CreateNfWAppFilter at server end. " }, nil
}

func (s *Server) DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	log.Println("Hello,use function DescribeNfs at server end.")
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeNfs at server end. " }, nil
}

func (s *Server) DescribeUserNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	log.Println("Hello,use function DescribeUserNfs at server end.")
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeUserNfs at server end. " }, nil
}



func (s *Server) CreateNfWaddrs(ctx context.Context, in *pb.CreateNfWaddrsRequest) (*pb.CreateNfResponse, error) {
	log.Println("Hello,use function CreateNfWaddrs at server end.")

	var (
		parser = &nf.NfHandlerModelParser{}
	)

	nf1, _ :=parser.CreateNfWaddrs(in)

	log.Print(nf1.AddrsStr)

	err :=s.nfservice.CreateNfWaddrs2(nf1)

	if err != nil {
		return nil,err
	}
	return &pb.CreateNfResponse{Message: "Hello,use function CreateNfWaddrs at server end. " }, nil
}

