// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package websocket

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/etcd"
)

func TestFormat(t *testing.T) {
	testUid := "uid"
	testMessageId := uint64(1111111)
	topic := "test_topic"
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	topic = FormatTopic(topic, testUid, testMessageId)
	uid, messageId := parseTopic(topic)
	require.Equal(t, testUid, uid)
	require.Equal(t, testMessageId, messageId)
}

func TestPushWsTask(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}

	mycfg := config.GetInstance()
	mycfg.LoadConf()

	endpoints := strings.Split(mycfg.Etcd.Endpoints, ",")
	nfetcd, err := etcd.Connect(endpoints, constants.EtcdPrefix)
	require.NoError(t, err)

	wsMessageType := "ws_op_event"
	testUid := "system"

	err = PushWsMessage(context.Background(), nfetcd, wsMessageType, testUid, Message{
		MessageId:   "test_message_id",
		MessageType: "event",
		Message:     "test_message",
		UserId:      "system",
	})

}

func TestWatchWsTasks(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("Local Dev testing env disabled.")
	}
	mycfg := config.GetInstance()
	mycfg.LoadConf()
	endpoints := strings.Split(mycfg.Etcd.Endpoints, ",")

	wsEtcd, err := etcd.Connect(endpoints, constants.EtcdPrefix)
	require.NoError(t, err)

	c := watchWsMessages(wsEtcd)
	time.Sleep(1 * time.Second)

	//============================================================================================

	testUid := "system"
	err = PushWsMessage(context.Background(), wsEtcd, "ws_op_event", testUid, Message{
		MessageId:   "ws_op_event",
		MessageType: "ws_op_event",
		Message:     "ws_op_event",
		UserId:      testUid,
	})
	event := <-c
	//require.Equal(t, testUid, event.UserId)
	//require.Equal(t, "ws_op_event", event.Message.MessageId)
	//require.Equal(t, "ws_op_event", event.Message.MessageType)
	//require.Equal(t, "ws_op_event", event.Message.Message)
	t.Log(event)

	err = PushWsMessage(context.Background(), wsEtcd, "ws_op_nf", testUid, Message{
		MessageId:   "ws_op_nf",
		MessageType: "ws_op_nf",
		Message:     "ws_op_nf",
		UserId:      testUid,
	})
	event = <-c

	err = PushWsMessage(context.Background(), wsEtcd, "ws_ks_nf", testUid, Message{
		MessageId:   "ws_ks_nf",
		MessageType: "ws_ks_nf",
		Message:     "ws_ks_nf",
		UserId:      testUid,
	})
	event = <-c

	err = PushWsMessage(context.Background(), wsEtcd, "ws_ks_event", testUid, Message{
		MessageId:   "ws_ks_event",
		MessageType: "ws_ks_event",
		Message:     "ws_ks_event",
		UserId:      testUid,
	})
	event = <-c
	//============================================================================================
	testUid2 := "huojiao"
	err = PushWsMessage(context.Background(), wsEtcd, "ws_op_event", testUid2, Message{
		MessageId:   "ws_op_event",
		MessageType: "ws_op_event",
		Message:     "ws_op_event",
		UserId:      testUid2,
	})
	event = <-c
	err = PushWsMessage(context.Background(), wsEtcd, "ws_op_nf", testUid2, Message{
		MessageId:   "ws_op_nf",
		MessageType: "ws_op_nf",
		Message:     "ws_op_nf",
		UserId:      testUid2,
	})
	event = <-c

	err = PushWsMessage(context.Background(), wsEtcd, "ws_ks_nf", testUid2, Message{
		MessageId:   "ws_ks_nf",
		MessageType: "ws_ks_nf",
		Message:     "ws_ks_nf",
		UserId:      testUid2,
	})
	event = <-c

	err = PushWsMessage(context.Background(), wsEtcd, "ws_ks_event", testUid2, Message{
		MessageId:   "ws_ks_event",
		MessageType: "ws_ks_event",
		Message:     "ws_ks_event",
		UserId:      testUid2,
	})
	event = <-c

}
