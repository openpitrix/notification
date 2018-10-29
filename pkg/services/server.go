/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.pb

package services

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/services/nf"
	"openpitrix.io/notification/pkg/util/dbutil"
	"os"
)

const (
	port = ":50051"
)


// Server is used to implement nf.RegisterNotificationServer.
type Server struct{
	cfg         *config.Config
	db          *gorm.DB
	nfhandler   nf.Handler
}

// NewServer initializes a new Server instance.
func NewServer() (*Server, error) {
	var (
		err    error
		server = &Server{}
	)

	server.cfg=config.NewConfig()

	//set mysql db,init database pool
	log.Println("step1:Set db")
	issucc :=  dbutil.GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}
	server.db = dbutil.GetInstance().GetMysqlDB()

	log.Println("step2:create new nfservice")
	nfservice :=nf.NewService(server.db)

	log.Println("step3:create new nfhandler")
	nfhandler:=nf.NewHandler(nfservice)

	log.Println("step4:set server.nfhandler")
	//set nfhandler
	server.nfhandler=nfhandler

	if err != nil {
		return nil, err
	}

	return server, nil
}


func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ss := grpc.NewServer()
	nfserver, _ :=NewServer()
	pb.RegisterNotificationServer(ss,nfserver)

	// Register reflection service on gRPC server.
	reflection.Register(ss)
	if err := ss.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}



// SayHello implements nf.RegisterNotificationServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Print("step5:call s.nfhandler.SayHello")
	s.nfhandler.SayHello(ctx,in)
	return &pb.HelloReply{Message: "Hello,use function SayHello at server end. " + in.Name}, nil
}

func (s *Server) CreateNfWaddrs(ctx context.Context, in *pb.CreateNfWaddrsRequest) (*pb.CreateNfResponse, error) {
	log.Println("Hello,use function CreateNfWaddrs at server end.")
	s.nfhandler.CreateNfWaddrs(ctx,in)
	return &pb.CreateNfResponse{Message: "Hello,use function CreateNfWaddrs at server end. " }, nil
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


