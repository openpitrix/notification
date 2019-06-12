// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package ws_message_redis

import (
	"context"
	"fmt"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	wstypes "openpitrix.io/notification/pkg/services/websocket/types"
	"testing"
	"time"

	r "github.com/go-redis/redis"
	pkg "openpitrix.io/notification/pkg"
)

func TestRedisPubSubSample(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	options, err := r.ParseURL("redis://192.168.0.6:6379")
	if err != nil {
		logger.Errorf(nil, "err=[%+v]", err)
	}
	rClient := r.NewClient(options)
	pubsub := rClient.Subscribe("mychannel")

	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive()
	if err != nil {
		panic(err)
	}
	var chSample = make(<-chan *r.Message)
	chSample = pubsub.Channel()
	// Publish a message.
	err = rClient.Publish("mychannel", "hello").Err()
	if err != nil {
		panic(err)
	}

	err = rClient.Publish("mychannel", "hello2").Err()
	if err != nil {
		panic(err)
	}

	err = rClient.Publish("mychannel", "hello3").Err()
	if err != nil {
		panic(err)
	}

	time.AfterFunc(time.Second, func() {
		// When pubsub is closed channel is closed too.
		_ = pubsub.Close()
	})

	// Consume messages.
	for msg := range chSample {
		fmt.Println(msg.Channel, msg.Payload)
	}

}

func TestPublishWsMessage(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	//pubsubConnStr := cfg.PubSub.Addr
	//pubsubType := cfg.PubSub.Type
	pubsubConnStr := "redis://192.168.0.6:6379"
	pubsubType := "redis"
	pubsubConfigMap := map[string]interface{}{
		"connStr": pubsubConnStr,
	}

	psClient, err := wstypes.New(pubsubType, pubsubConfigMap)
	if err != nil {
		logger.Errorf(nil, "Failed to connect pubsub server: %+v.", err)
	}
	rClient := (psClient).(wstypes.RedisClient)

	testUserid := "system"

	msgDetail := wstypes.MessageDetail{
		MessageId:      wstypes.NewWsMessageId(),
		UserId:         testUserid,
		Service:        "op",
		MessageType:    "event",
		MessageContent: "Message_content_test",
	}

	userMsg0 := wstypes.UserMessage{
		UserId:        testUserid,
		Service:       "op",
		MessageType:   "event",
		MessageDetail: msgDetail,
	}

	//channel := formatRedisChannel(clientService, wsMessageType)
	//pubsub := rClient.PSubscribe("op/event/*")
	//pubsub := rClient.PSubscribe(channel)

	wsMessageRedis := new(WsMessageRedis)

	err = wsMessageRedis.PublishWsMessage(context.Background(), &rClient, &userMsg0)

	//time.AfterFunc(time.Second, func() {
	//	// When pubsub is closed channel is closed too.
	//	_ = pubsub.Close()
	//})

	//getWsMessagesByRredis(pubsub)
}

func TestWatchWsMessages(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	config.GetInstance().LoadConf()
	//pubsubConnStr := cfg.PubSub.Addr
	//pubsubType := cfg.PubSub.Type
	pubsubConnStr := "redis://192.168.0.6:6379"
	pubsubType := "redis"
	pubsubConfigMap := map[string]interface{}{
		"connStr": pubsubConnStr,
	}

	psClient, err := wstypes.New(pubsubType, pubsubConfigMap)
	if err != nil {
		logger.Errorf(nil, "Failed to connect pubsub server: %+v.", err)
	}
	rClient := (psClient).(wstypes.RedisClient)

	//channel := formatRedisChannel(clientService, wsMessageType)
	//pubsub := rClient.PSubscribe("op/event/*")
	//pubsub := rClient.PSubscribe(channel)

	wsMessageRedis := new(WsMessageRedis)
	c := wsMessageRedis.WatchWsMessages(&rClient)
	time.Sleep(1 * time.Second)

	//============================================================================================================================
	service := "op"
	messageType := "nf"
	testUserid := "system"

	msgDetail := wstypes.MessageDetail{
		MessageId:      wstypes.NewWsMessageId(),
		UserId:         testUserid,
		Service:        service,
		MessageType:    messageType,
		MessageContent: "Message_content_test",
	}

	userMsg := wstypes.UserMessage{
		UserId:        testUserid,
		Service:       service,
		MessageType:   messageType,
		MessageDetail: msgDetail,
	}

	err = wsMessageRedis.PublishWsMessage(context.Background(), &rClient, &userMsg)
	outmsg := <-c
	t.Log(outmsg)
	//============================================================================================================================

	service1 := "op"
	messageType1 := "nf"
	testUserid1 := "system"

	msgDetail1 := wstypes.MessageDetail{
		MessageId:      wstypes.NewWsMessageId(),
		UserId:         testUserid1,
		Service:        service1,
		MessageType:    messageType1,
		MessageContent: "Message_content_test",
	}

	userMsg1 := wstypes.UserMessage{
		UserId:        testUserid1,
		Service:       service1,
		MessageType:   messageType1,
		MessageDetail: msgDetail1,
	}

	err = wsMessageRedis.PublishWsMessage(context.Background(), &rClient, &userMsg1)
	outmsg1 := <-c
	t.Log(outmsg1)
	//time.AfterFunc(time.Second, func() {
	//	// When pubsub is closed channel is closed too.
	//	_ = pubsub.Close()
	//})

	//getWsMessagesByRredis(pubsub)
	//wsMessageRedis.WatchWsMessages(&rClient)

}
