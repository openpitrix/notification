package services

import (
	"golang.org/x/net/context"
	"log"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

// SayHello implements nf.RegisterNotificationServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	logger.Debugf(nil,"step5:call s.nfhandler.SayHello")
	s.nfhandler.SayHello(ctx, in)
	return &pb.HelloReply{Message: "Hello,use function SayHello at server end. " + in.Name}, nil
}

func (s *Server) CreateNfWaddrs(ctx context.Context, in *pb.CreateNfWaddrsRequest) (*pb.CreateNfResponse, error) {
	res,err:= s.nfhandler.CreateNfWaddrs(ctx, in)
	return res, err
}

func (s *Server) CreateNfWUserFilter(ctx context.Context, in *pb.CreateNfWUserFilterRequest) (*pb.CreateNfResponse, error) {
	logger.Debugf(nil,"Hello,use function CreateNfWUserFilter at server end.")
	return &pb.CreateNfResponse{NfPostId: pbutil.ToProtoString("testID4CreateNfWUserFilter")}, nil
}

func (s *Server) CreateNfWAppFilter(ctx context.Context, in *pb.CreateNfWAppFilterRequest) (*pb.CreateNfResponse, error) {
	logger.Debugf(nil,"Hello,use function CreateNfWAppFilter at server end.")
	return &pb.CreateNfResponse{NfPostId: pbutil.ToProtoString("testID4CreateNfWAppFilter")}, nil
}

func (s *Server) DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	log.Println("Hello,use function DescribeNfs at server end.")
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeNfs at server end. "}, nil
}

func (s *Server) DescribeUserNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	log.Println("Hello,use function DescribeUserNfs at server end.")
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeUserNfs at server end. "}, nil
}
