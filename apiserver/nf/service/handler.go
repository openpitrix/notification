package service

import (
	"golang.org/x/net/context"
	"log"
	pb "notificationService/apiserver/nf/proto"
)



// SayHello implements notification.RegisterNotificationServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Println("Hello,use function SayHello at server end.")
	return &pb.HelloReply{Message: "Hello,use function SayHello at server end. " + in.Name}, nil
}

func (s *server) CreateNfWUserFilter(ctx context.Context, in *pb.CreateNfWUserFilterRequest) (*pb.CreateNfResponse, error) {
	log.Println("Hello,use function CreateNfWUserFilter at server end.")
	return &pb.CreateNfResponse{Message: "1111 "  }, nil
}

func (s *server) CreateNfWAppFilter(ctx context.Context, in *pb.CreateNfWAppFilterRequest) (*pb.CreateNfResponse, error) {
	log.Println("Hello,use function CreateNfWAppFilter at server end.")
	return &pb.CreateNfResponse{Message: "Hello,use function CreateNfWAppFilter at server end. " }, nil
}

func (s *server) CreateNfWaddrs(ctx context.Context, in *pb.CreateNfWaddrsRequest) (*pb.CreateNfResponse, error) {
	log.Println("Hello,use function CreateNfWaddrs at server end.")
	return &pb.CreateNfResponse{Message: "Hello,use function CreateNfWaddrs at server end. " }, nil

	//s := senderutil.GetSenderFromContext(ctx)
	//// validate req
	//err := validateCreateRuntimeRequest(ctx, req)
	//// TODO: refactor create runtime params
	//if err != nil {
	//	if gerr.IsGRPCError(err) {
	//		return nil, err
	//	} else {
	//		return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorValidateFailed)
	//	}
	//}
	//
	//// create runtime credential
	//runtimeCredentialId, err := createRuntimeCredential(ctx, req.Provider.GetValue(), req.RuntimeCredential.GetValue())
	//if err != nil {
	//	return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	//}
	//
	//// create runtime
	//runtimeId, err := createRuntime(
	//	ctx,
	//	req.GetName().GetValue(),
	//	req.GetDescription().GetValue(),
	//	req.Provider.GetValue(),
	//	req.GetRuntimeUrl().GetValue(),
	//	runtimeCredentialId,
	//	req.Zone.GetValue(),
	//	s.UserId)
	//if err != nil {
	//	return nil, gerr.NewWithDetail(ctx, gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	//}
	//
	//if req.GetLabels() != nil {
	//	err = labelutil.SyncRuntimeLabels(ctx, runtimeId, req.GetLabels().GetValue())
	//	if err != nil {
	//		return nil, err
	//	}
	//}
	//
	//res := &pb.CreateRuntimeResponse{
	//	RuntimeId: pbutil.ToProtoString(runtimeId),
	//}
	//return res, nil
}

func (s *server) DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	log.Println("Hello,use function DescribeNfs at server end.")
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeNfs at server end. " }, nil
}

func (s *server) DescribeUserNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error) {
	log.Println("Hello,use function DescribeUserNfs at server end.")
	return &pb.DescribeNfsResponse{Message: "Hello,use function DescribeUserNfs at server end. " }, nil
}

