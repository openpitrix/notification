// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package ws_message_etcd

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	wstypes "openpitrix.io/notification/pkg/services/websocket/types"
	"openpitrix.io/notification/pkg/util/idutil"
	"openpitrix.io/notification/pkg/util/jsonutil"
)

type WsMessageEtcd struct {
}

func (m WsMessageEtcd) PushWsMessage(ctx context.Context, eClient *wstypes.EtcdClient, userMsg *wstypes.UserMessage) error {
	service := userMsg.Service
	messageType := userMsg.MessageType
	uid := userMsg.UserId
	msgDetail := userMsg.MessageDetail

	var messageId = idutil.GetIntId()

	var key = formatTopic4Etcd(service, messageType, uid, messageId)
	msgValue, err := jsonutil.Encode(msgDetail)
	if err != nil {
		logger.Errorf(ctx, "Encode message[%+v] to json failed", msgDetail)
		return err
	}

	resp, err := eClient.Grant(ctx, wstypes.ExpireTime)
	if err != nil {
		logger.Errorf(ctx, "Grant ttl from etcd failed: %+v", err)
		return err
	}

	_, err = eClient.Put(ctx, key, string(msgValue), clientv3.WithLease(resp.ID))
	if err != nil {
		logger.Errorf(ctx, "Push message[%+v] to etcd failed: %+v", msgDetail, err)
		return err
	}
	logger.Debugf(ctx, "Push message[%+v] to etcd successfully.", msgDetail)

	return nil
}

func (m WsMessageEtcd) WatchWsMessages(e *wstypes.EtcdClient) chan wstypes.UserMessage {
	serviceMessageTypes := strings.Split(config.GetInstance().Websocket.ServiceMessageTypes, ",")
	var c = make(chan wstypes.UserMessage, 255)

	for _, serviceMessageType := range serviceMessageTypes {
		go getMessageByWatches4Etcd(e, serviceMessageType, c)
	}
	return c
}

func getMessageByWatches4Etcd(e *wstypes.EtcdClient, serviceMessageType string, c chan wstypes.UserMessage) {
	watchRes := e.Watch(context.Background(), serviceMessageType+"/", clientv3.WithPrefix())
	for res := range watchRes {
		for _, ev := range res.Events {
			if ev.Type == mvccpb.PUT {
				var msgDetail wstypes.MessageDetail
				value := ev.Kv.Value
				err := jsonutil.Decode(value, &msgDetail)
				if err != nil {
					logger.Errorf(nil, "Decode ws_message[%s] failed: %+v", string(value), err)
				} else {
					c <- wstypes.UserMessage{
						UserId:        msgDetail.UserId,
						Service:       msgDetail.Service,
						MessageType:   msgDetail.MessageType,
						MessageDetail: msgDetail,
					}
				}
			}
		}
	}
}

func formatTopic4Etcd(service string, wsMessageType string, uid string, messageId uint64) string {
	return fmt.Sprintf("%s/%s/%s/%d", service, wsMessageType, uid, messageId)
}
