// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package websocket

import (
	"context"
	"fmt"
	"openpitrix.io/notification/pkg/constants"
	"strconv"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/etcd"
	"openpitrix.io/notification/pkg/util/idutil"
	"openpitrix.io/notification/pkg/util/jsonutil"
)

func parseTopic(topic string) (uid string, messageId uint64) {
	t := strings.Split(topic, "/")
	uid = t[1]
	mid, _ := strconv.Atoi(t[2])
	messageId = uint64(mid)
	return
}

func FormatTopic(wsMessageType string, uid string, messageId uint64) string {
	return fmt.Sprintf("%s/%s/%d", wsMessageType, uid, messageId)
}

func PushWsMessage(ctx context.Context, wsEtcd *etcd.Etcd, wsMessageType string, uid string, msg Message) error {
	wsMessageType = constants.WsMessagePrefix + wsMessageType
	var messageId = idutil.GetIntId()
	var key = FormatTopic(wsMessageType, uid, messageId)
	msgValue, err := jsonutil.Encode(msg)
	if err != nil {
		logger.Errorf(ctx, "Encode message [%+v] to json failed", msg)
		return err
	}

	resp, err := wsEtcd.Grant(ctx, ExpireTime)
	if err != nil {
		logger.Errorf(ctx, "Grant ttl from etcd failed: %+v", err)
		return err
	}

	_, err = wsEtcd.Put(ctx, key, string(msgValue), clientv3.WithLease(resp.ID))
	if err != nil {
		logger.Errorf(ctx, "Push user [%s] message [%d] [%s] to etcd failed: %+v", uid, messageId, string(msgValue), err)
		return err
	}
	return nil
}

func watchWsMessages(e *etcd.Etcd) chan userMessage {
	mycfg := config.GetInstance()
	mycfg.LoadConf()

	messageTypes := strings.Split(mycfg.Websocket.MessageTypes, ",")
	var c = make(chan userMessage, 255)

	for _, messageType := range messageTypes {
		go getMessageByWatches(e, messageType, c)
	}
	return c
}

func getMessageByWatches(e *etcd.Etcd, wsMessageType string, c chan userMessage) {
	wsMessageType = constants.WsMessagePrefix + wsMessageType
	watchRes := e.Watch(context.Background(), wsMessageType+"/", clientv3.WithPrefix())
	for res := range watchRes {
		for _, ev := range res.Events {
			if ev.Type == mvccpb.PUT {
				var message Message
				key := string(ev.Kv.Key)
				userId, msgId := parseTopic(key)
				err := jsonutil.Decode(ev.Kv.Value, &message)
				if err != nil {
					logger.Errorf(nil, "Decode ws message [%s] [%d] [%s] failed: %+v", userId, msgId, string(ev.Kv.Value), err)
				} else {
					logger.Infof(nil, "Get ws message [%s] [%d] [%s]", userId, msgId, string(ev.Kv.Value))
					c <- userMessage{
						UserId:      userId,
						MessageType: message.MessageType,
						Message:     message,
					}
				}
			}
		}
	}
}
