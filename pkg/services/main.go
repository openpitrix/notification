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
	"net"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/pb"
)

func Serve(cfg *config.Config) {

	config.GetInstance().PrintUsage()

	port := config.GetInstance().App.Port
	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Criticalf(nil, "failed to listen: %v", err)
	}
	nfserver, _ := NewServer()

	go nfserver.taskhandler.ServeTask()

	s := grpc.NewServer()
	pb.RegisterNotificationServer(s, nfserver)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		logger.Criticalf(nil, "failed to serve: %v", err)
	}
}
