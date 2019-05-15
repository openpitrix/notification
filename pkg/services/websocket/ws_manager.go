// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package websocket

import (
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/etcd"
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
	*etcd.Etcd
	receiverMap map[string]receiversT
	addReceiver chan receiver
	delReceiver chan receiver
	msgChan     chan userMessage
}

func NewWsManager(e *etcd.Etcd) *wsManager {
	var wsm wsManager
	wsm.Etcd = e
	wsm.msgChan = watchWsMessages(e)
	wsm.addReceiver = make(chan receiver, 255)
	wsm.delReceiver = make(chan receiver, 255)
	wsm.receiverMap = make(map[string]receiversT)
	return &wsm
}

func (wsm *wsManager) Run() {
	for {
		select {
		case receiver := <-wsm.addReceiver:
			receivers := wsm.getReceivers(receiver.MessageType, receiver.UserId)
			receivers[receiver.Conn] = &sync.Mutex{}

		case receiver := <-wsm.delReceiver:
			receivers := wsm.getReceivers(receiver.MessageType, receiver.UserId)
			delete(receivers, receiver.Conn)
			if len(receivers) == 0 {
				delete(wsm.receiverMap, receiver.UserId)
			}
			go receiver.Conn.Close()

		case userMsg := <-wsm.msgChan:
			userMsgType := userMsg.MessageType

			mycfg := config.GetInstance()
			msgTypes := strings.Split(mycfg.Websocket.MessageTypes, ",")

			for _, msgType := range msgTypes {
				if msgType == userMsgType {
					receivers := wsm.getReceivers(userMsgType, userMsg.UserId)
					for c, mutex := range receivers {
						go writeMessage(c, userMsgType, mutex, userMsg)
					}
				}
			}
		}
	}
}

func (wsm *wsManager) getReceivers(messageType string, userId string) receiversT {
	key := messageType + "_" + userId
	receiverT, ok := wsm.receiverMap[key]
	if !ok {
		receiverT = make(receiversT)
		wsm.receiverMap[key] = receiverT
	}
	return receiverT
}

func writeMessage(conn *websocket.Conn, messageType string, mutex *sync.Mutex, userMsg userMessage) {
	mutex.Lock()
	defer mutex.Unlock()

	var err error
	if messageType == userMsg.MessageType {
		err = conn.WriteJSON(userMsg.Message)
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

		messageType := r.URL.Query().Get("message_type")
		if messageType == "" {
			http.Error(w, "Unauthorized: [message_type] is required.", http.StatusUnauthorized)
			return
		}

		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Infof(nil, "Upgrade websocket request failed: %+v", err)
			return
		}

		receiver := receiver{
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
