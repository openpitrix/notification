// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package ws_message_redis

import (
	"context"
	"fmt"
	"strings"

	r "github.com/go-redis/redis"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	wstypes "openpitrix.io/notification/pkg/services/websocket/types"
	"openpitrix.io/notification/pkg/util/jsonutil"
)

type WsMessageRedis struct {
}

func (m WsMessageRedis) PublishWsMessage(ctx context.Context, rClient *wstypes.RedisClient, userMsg *wstypes.UserMessage) error {
	service := userMsg.Service
	messageType := userMsg.MessageType
	msgDetail := userMsg.MessageDetail

	channel := formatRedisChannel(service, messageType)
	msgValue, err := jsonutil.Encode(msgDetail)
	if err != nil {
		logger.Errorf(ctx, "Encode message [%+v] to json failed", msgDetail)
		return err
	}
	// Publish a message.
	err = rClient.Publish(channel, msgValue).Err()
	if err != nil {
		panic(err)
	}
	return nil
}

func (m WsMessageRedis) WatchWsMessages(rClient *wstypes.RedisClient) chan wstypes.UserMessage {
	wscfg := config.GetInstance().Websocket.ServiceMessageTypes
	serviceMessageTypes := strings.Split(wscfg, ",")
	var c = make(chan wstypes.UserMessage, 255)

	for _, serviceMessageType := range serviceMessageTypes {
		clientService, messageType := parseRedisChannel(serviceMessageType)
		channel := formatRedisChannel(clientService, messageType)

		pubsub := rClient.PSubscribe(channel)

		go getWsMessagesByRredis(pubsub, c)
	}
	return c
}

func getWsMessagesByRredis(ps *r.PubSub, c chan wstypes.UserMessage) {
	ch := ps.Channel()
	for msg := range ch {
		logger.Debugf(nil, "msg.Channel=[%s],msg.Payload=[%s]. ", msg.Channel, msg.Payload)
		var msgDetail wstypes.MessageDetail
		var data []byte = []byte(msg.Payload)
		err := jsonutil.Decode(data, &msgDetail)
		if err != nil {
			logger.Errorf(nil, "Decode userMessage[%s] failed: %+v", msg.Payload, err)
		} else {
			logger.Debugf(nil, "Get userMessage=[%+v].]", msgDetail)
			c <- wstypes.UserMessage{
				UserId:        msgDetail.UserId,
				Service:       msgDetail.Service,
				MessageType:   msgDetail.MessageType,
				MessageDetail: msgDetail,
			}
		}
	}
}

func parseRedisChannel(topic string) (wsService string, wsMessageType string) {
	t := strings.Split(topic, "/")
	wsService = t[0]
	wsMessageType = t[1]
	return
}

func formatRedisChannel(service string, wsMessageType string) string {
	return fmt.Sprintf("%s/%s/*", service, wsMessageType)
}
