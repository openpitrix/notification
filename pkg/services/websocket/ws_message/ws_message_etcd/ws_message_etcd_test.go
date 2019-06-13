// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package ws_message_etcd

import (
	"context"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"testing"
	"time"

	pkg "openpitrix.io/notification/pkg"
	wstypes "openpitrix.io/notification/pkg/services/websocket/types"
)

func TestPushWsMessage(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	//cfg := config.GetInstance().LoadConf()
	//pubsubConnStr := cfg.PubSub.Addr
	//pubsubType := cfg.PubSub.Type

	pubsubConnStr := "192.168.0.6:12379"
	pubsubType := "etcd"

	pubsubConfigMap := map[string]interface{}{
		"connStr": pubsubConnStr,
	}

	psClient, err := wstypes.New(pubsubType, pubsubConfigMap)
	if err != nil {
		logger.Errorf(nil, "Failed to connect pubsub server: %+v.", err)
	}
	eClient := (psClient).(wstypes.EtcdClient)

	service := "op"
	messageType := "event"
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

	wsMessageEtcd := new(WsMessageEtcd)
	err = wsMessageEtcd.PushWsMessage(context.Background(), &eClient, &userMsg)

}

func TestWatchWsMessages(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	config.GetInstance().LoadConf()
	//pubsubConnStr := cfg.PubSub.Addr
	//pubsubType := cfg.PubSub.Type
	pubsubConnStr := "192.168.0.6:12379"
	pubsubType := "etcd"

	pubsubConfigMap := map[string]interface{}{
		"connStr": pubsubConnStr,
	}

	psClient, err := wstypes.New(pubsubType, pubsubConfigMap)
	if err != nil {
		logger.Errorf(nil, "Failed to connect pubsub server: %+v.", err)
	}
	eClient := (psClient).(wstypes.EtcdClient)

	wsMessageEtcd := new(WsMessageEtcd)
	c := wsMessageEtcd.WatchWsMessages(&eClient)
	time.Sleep(1 * time.Second)

	//====================================================================================================================================
	service := "op"
	messageType := "event"
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
	err = wsMessageEtcd.PushWsMessage(context.Background(), &eClient, &userMsg)
	outmsg := <-c
	//require.Equal(t, testUserid, event.UserId)
	//require.Equal(t, "MessageId_test", event.Message.MessageId)
	//require.Equal(t, "event", event.Message.MessageType)
	//require.Equal(t, "Message_content_test", event.Message.Message)
	t.Log(outmsg)
	//====================================================================================================================================
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
	err = wsMessageEtcd.PushWsMessage(context.Background(), &eClient, &userMsg1)
	outmsg1 := <-c
	//require.Equal(t, testUserid, event.UserId)
	//require.Equal(t, "MessageId_test", event.Message.MessageId)
	//require.Equal(t, "event", event.Message.MessageType)
	//require.Equal(t, "Message_content_test", event.Message.Message)
	t.Log(outmsg1)

}
