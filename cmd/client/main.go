// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"openpitrix.io/notification/pkg/client"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
	"strconv"
	"time"
)

func main() {
	client.Serve()
}

func Serve() {

	config.GetInstance().LoadConf()
	host := constants.NotificationManagerHost
	port := constants.NotificationManagerPort
	address := host + strconv.Itoa(port)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNotificationClient(conn)

	// Contact the server and print out its response.

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	testAddrsStr := "huojiao2006@163.com;513590612@qq.com"
	r, err := c.CreateNfWithAddrs(ctx, &pb.CreateNfWithAddrsRequest{
		ContentType:  pbutil.ToProtoString("Information"),
		SentType:     pbutil.ToProtoString("Email"),
		AddrsStr:     pbutil.ToProtoString(testAddrsStr),
		Title:        pbutil.ToProtoString("Run case"),
		Content:      pbutil.ToProtoString("Run case content"),
		ShortContent: pbutil.ToProtoString("Run case ShortContent"),
		ExpiredDays:  pbutil.ToProtoString("7"),
		Owner:        pbutil.ToProtoString("HuoJiao"),
		Status:       pbutil.ToProtoString("New"),
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting output client - SayHello: %s", r.NotificationId)

	c.CreateNfWithAddrs(ctx, &pb.CreateNfWithAddrsRequest{})

}
