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

package client

/*

import (
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/pb"
)

const (
	defaultName = "world"
)

func Serve() {
	log.Println("Start to run client. Step 1")

	config.GetInstance().LoadConf()
	host := config.GetInstance().App.Host
	port := config.GetInstance().App.Port
	address := host + port

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNotificationClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting output client - SayHello: %s", r.Message)

}
*/
