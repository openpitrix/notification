// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package websocket

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	wstypes "openpitrix.io/notification/pkg/services/websocket/types"
	"openpitrix.io/notification/pkg/services/websocket/ws_message/ws_message_etcd"
	"openpitrix.io/notification/pkg/services/websocket/ws_message/ws_message_redis"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: get allowed host from global config
		return true
	},
}

type receiversT map[*websocket.Conn]*sync.Mutex

type wsManager struct {
	psClient    *wstypes.PubsubClient
	receiverMap map[string]receiversT
	addReceiver chan wstypes.Receiver
	delReceiver chan wstypes.Receiver
	msgChan     chan wstypes.UserMessage
}

func NewWsManager(pubsubType string, psClient *wstypes.PubsubClient) (*wsManager, error) {
	var wsm wsManager
	wsm.psClient = psClient

	if pubsubType == "etcd" {
		m := ws_message_etcd.WsMessageEtcd{}
		eClient := (*psClient).(wstypes.EtcdClient)
		wsm.msgChan = m.WatchWsMessages(&eClient)
	} else if pubsubType == "redis" {
		m := ws_message_redis.WsMessageRedis{}
		rClient := (*psClient).(wstypes.RedisClient)
		wsm.msgChan = m.WatchWsMessages(&rClient)
	} else {
		return nil, fmt.Errorf("unsupported queueType [%s]", pubsubType)
	}

	wsm.addReceiver = make(chan wstypes.Receiver, 255)
	wsm.delReceiver = make(chan wstypes.Receiver, 255)
	wsm.receiverMap = make(map[string]receiversT)
	return &wsm, nil
}

func (wsm *wsManager) Run() {
	for {
		select {
		case receiver := <-wsm.addReceiver:
			receivers := wsm.getReceivers(receiver.Service, receiver.MessageType, receiver.UserId)
			receivers[receiver.Conn] = &sync.Mutex{}

		case receiver := <-wsm.delReceiver:
			receivers := wsm.getReceivers(receiver.Service, receiver.MessageType, receiver.UserId)
			delete(receivers, receiver.Conn)
			if len(receivers) == 0 {
				delete(wsm.receiverMap, receiver.UserId)
			}
			go receiver.Conn.Close()

		case userMsg := <-wsm.msgChan:
			service := userMsg.Service
			userMsgType := userMsg.MessageType
			userId := userMsg.UserId

			userServiceMessageType := service + "/" + userMsgType
			serviceMessageTypes := strings.Split(config.GetInstance().Websocket.ServiceMessageTypes, ",")

			for _, serviceMessageType := range serviceMessageTypes {
				if serviceMessageType == userServiceMessageType {
					receivers := wsm.getReceivers(service, userMsgType, userId)

					for r, mutex := range receivers {
						go writeMessage(r, service, userMsgType, mutex, userMsg)
					}
				}
			}
		}
	}
}

func (wsm *wsManager) getReceivers(service string, messageType string, userId string) receiversT {
	key := service + "/" + messageType + "/" + userId
	receiverT, ok := wsm.receiverMap[key]
	if !ok {
		receiverT = make(receiversT)
		wsm.receiverMap[key] = receiverT
	}
	return receiverT
}

func writeMessage(conn *websocket.Conn, service string, messageType string, mutex *sync.Mutex, userMsg wstypes.UserMessage) {
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	if service == userMsg.Service && messageType == userMsg.MessageType {
		err = conn.WriteJSON(userMsg.MessageDetail)
	}

	if err != nil {
		logger.Errorf(nil, "Failed to send message [%+v] to [%+v], error: %+v", userMsg, conn, err)
	}
	logger.Debugf(nil, "Message sent [%+v]", userMsg)
}

func (wsm *wsManager) HandleWsTask() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("user_id")
		if userId == "" {
			http.Error(w, "Unauthorized: [user_id] is required.", http.StatusUnauthorized)
			return
		}

		service := r.URL.Query().Get("service")
		if service == "" {
			http.Error(w, "Unauthorized: [service] is required.", http.StatusUnauthorized)
			return
		}

		messageType := r.URL.Query().Get("message_type")
		if messageType == "" {
			http.Error(w, "Unauthorized: [message_type] is required.", http.StatusUnauthorized)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Infof(nil, "Upgrade websocket_manager request failed: %+v", err)
			return
		}

		receiver := wstypes.Receiver{
			Service:     service,
			MessageType: messageType,
			UserId:      userId,
			Conn:        c,
		}
		wsm.addReceiver <- receiver
		for {
			_, _, err := receiver.Conn.ReadMessage()
			if err != nil {
				wsm.delReceiver <- receiver
				logger.Errorf(nil, "Connection [%p] error: %+v", receiver.Conn, err)
				return
			}
		}
	}
}

func PushWsMessage(ctx context.Context, pubsubType string, psClient *wstypes.PubsubClient, userMsg *wstypes.UserMessage) error {
	if pubsubType == "etcd" {
		wsMessageEtcd := new(ws_message_etcd.WsMessageEtcd)
		eClient := (*psClient).(wstypes.EtcdClient)
		err := wsMessageEtcd.PushWsMessage(ctx, &eClient, userMsg)
		if err != nil {
			return err
		}
	} else if pubsubType == "redis" {
		wsMessageRedis := new(ws_message_redis.WsMessageRedis)
		rClient := (*psClient).(wstypes.RedisClient)

		err := wsMessageRedis.PublishWsMessage(ctx, &rClient, userMsg)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unsupported queueType [%s]", pubsubType)
	}

	return nil
}
