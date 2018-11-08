package nf

import (
	"golang.org/x/net/context"
	"openpitrix.io/notification/pkg/pb"
)

type Handler interface {
	SayHello(ctx context.Context, in *pb.HelloRequest) error
	CreateNfWaddrs(ctx context.Context, in *pb.CreateNfWaddrsRequest) (*pb.CreateNfResponse, error)
	DescribeNfs(ctx context.Context, in *pb.DescribeNfsRequest) (*pb.DescribeNfsResponse, error)
}
