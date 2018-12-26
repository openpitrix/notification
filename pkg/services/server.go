package services

import (
	"golang.org/x/net/context"
	"log"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
	"openpitrix.io/notification/pkg/services/nf"
	"openpitrix.io/notification/pkg/services/task"
	"openpitrix.io/notification/pkg/util/dbutil"
	"openpitrix.io/notification/pkg/util/etcdutil"
	"os"
)



// Server is used to implement nf.RegisterNotificationServer.
type Server struct {
	nfhandler   nf.Handler
	taskhandler task.Handler
}



// NewServer initializes a new Server instance.
func NewServer() (*Server, error) {
	logger.Debugf(nil,"step0:start********************************************")

	var (
		err    error
		server = &Server{}
	)

	logger.Debugf(nil,"step1:set server.nfhandler**********************")
	logger.Debugf(nil,"step1.1:create nfservice")
	logger.Debugf(nil,"step1.1.1:create queue")
	cfg := config.GetInstance()
	endpoints := []string{cfg.Etcd.Endpoints}

	prefix:=cfg.Etcd.Prefix
	nfetcd, err := etcdutil.Connect(endpoints, prefix)
	if err != nil {
		logger.Criticalf(nil,"%+v",err)
	}

	topic:=cfg.Etcd.Topic
	q := nfetcd.NewQueue(topic)

	logger.Debugf(nil,"step1.1.2:get db")
	db := dbutil.GetInstance().GetMysqlDB()

	logger.Debugf(nil,"step1.1:create new nfservice")
	nfservice := nf.NewService(db, q)
	logger.Debugf(nil,"step1.2:create nfhandler")
	nfhandler := nf.NewHandler(nfservice)
	logger.Debugf(nil,"step1.3:set server.nfhandler")
	server.nfhandler = nfhandler

	logger.Debugf(nil,"step2:set server.taskhandler**********************")
	logger.Debugf(nil,"step2.1:create taskservice")
	taskservice := task.NewService(db, q)
	logger.Debugf(nil,"step2.2:create taskhandler")
	taskhandler := task.NewHandler(taskservice)
	logger.Debugf(nil,"step2.3:set server.taskhandler")
	server.taskhandler = taskhandler

	if err != nil {
		logger.Criticalf(nil,"%+v",err)
		return nil, err
	}
	logger.Debugf(nil,"step0:end********************************************")
	return server, nil
}

func InitGlobelSetting() {
	logger.Debugf(nil,"step0.1:初始化配置参数")
	config.GetInstance().LoadConf()

	logger.Debugf(nil,"step0.2:初始化DB connection pool")
	issucc := dbutil.GetInstance().InitDataPool()
	if !issucc {
		logger.Criticalf(nil,"init database pool failure...")
		os.Exit(1)
	}

	AppLogMode:=config.GetInstance().App.Applogmode
	logger.SetLevelByString(AppLogMode)
}


// SayHello implements nf.RegisterNotificationServer
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	logger.Debugf(nil,"step5:call s.nfhandler.SayHello")
	s.nfhandler.SayHello(ctx, in)
	return &pb.HelloReply{Message: "Hello,use function SayHello at server end. " + in.Name}, nil
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
//func (s *Server) CreateNfWithAddrs(context.Context, *pb.CreateNfWithAddrsRequest) (*pb.CreateNfResponse, error) {
//	panic("implement me")
//}


func (s *Server) CreateNfWithAddrs(ctx context.Context, in *pb.CreateNfWithAddrsRequest) (*pb.CreateNfResponse, error) {
	res,err:= s.nfhandler.CreateNfWithAddrs(ctx, in)
	return res, err
}
