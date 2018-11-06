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
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/pb"
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
	log.Println("step0:start********************************************")

	var (
		err    error
		server = &Server{}
	)

	log.Println("step1:set server.nfhandler**********************")
	log.Println("step1.1:create nfservice")
	log.Println("step1.1.1:create queue")
	cfg := config.GetInstance()
	endpoints := []string{cfg.Etcd.Endpoints}

	prefix := "test"
	nfetcd, err := etcdutil.Connect(endpoints, prefix)
	if err != nil {
		log.Fatal(err)
	}
	q := nfetcd.NewQueue("nf_task")

	log.Println("step1.1.2:get db")
	db := dbutil.GetInstance().GetMysqlDB()

	log.Println("step1.1:create new nfservice")
	nfservice := nf.NewService(db, q)
	log.Println("step1.2:create nfhandler")
	nfhandler := nf.NewHandler(nfservice)
	log.Println("step1.3:set server.nfhandler")
	server.nfhandler = nfhandler

	log.Println("step2:set server.taskhandler**********************")
	log.Println("step2.1:create taskservice")
	taskservice := task.NewService(db, q)
	log.Println("step2.2:create taskhandler")
	taskhandler := task.NewHandler(taskservice)
	log.Println("step2.3:set server.taskhandler")
	server.taskhandler = taskhandler

	if err != nil {
		return nil, err
	}
	log.Println("step0:end********************************************")
	return server, nil
}

func InitGlobelSetting() {
	log.Println("step0.1:初始化配置参数")
	config.GetInstance().InitCfg()

	log.Println("step0.2:初始化DB connection pool")
	issucc := dbutil.GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}

}
//**************************************************************************************************

func Serve() error {
	InitGlobelSetting()

	port := config.GetInstance().App.Port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	nfserver, _ := NewServer()

	go nfserver.taskhandler.ServeTask()

	s := grpc.NewServer()
	pb.RegisterNotificationServer(s, nfserver)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}



//**************************************************************************************************
