// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package websocket

import (
	"context"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"openpitrix.io/logger"

	nfclient "openpitrix.io/notification/pkg/client/notification"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type Receiver struct {
	Service     string
	MessageType string
	UserId      string
	Conn        *websocket.Conn
}

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
	accessServiceName string
	messageTypes      []string
	receiverMap       map[string]receiversT
	addReceiver       chan Receiver
	delReceiver       chan Receiver
	msgChan           chan models.UserMessage
}

func NewWsManager(accessServiceName string, messageTypes []string) (*wsManager, error) {
	var wsm wsManager
	wsm.accessServiceName = accessServiceName
	wsm.messageTypes = messageTypes
	wsm.addReceiver = make(chan Receiver, 255)
	wsm.delReceiver = make(chan Receiver, 255)
	wsm.receiverMap = make(map[string]receiversT)

	client, err := nfclient.NewClient()
	if err != nil {
		logger.Errorf(nil, "failed to create nfclient,err=%+v", err)
		return nil, err
	}

	reqstreamData := &pb.StreamReqData{
		Service: pbutil.ToProtoString(wsm.accessServiceName),
	}
	ctx := context.Background()
	channelClient, err := client.CreateNotificationChannel(ctx, reqstreamData)
	if err != nil {
		logger.Errorf(nil, "failed to get msgs by grpc stream,err=%+v", err)
		return nil, err
	}

	wsm.msgChan = wsm.ReceiveMsg(client, channelClient, ctx)
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

			serviceMsgTypeFromData := service + "/" + userMsgType

			for _, msgType := range wsm.messageTypes {
				s := wsm.accessServiceName + "/" + msgType
				if s == serviceMsgTypeFromData {
					receivers := wsm.getReceivers(service, userMsgType, userId)

					for r, mutex := range receivers {
						go wsm.writeMessage(r, service, userMsgType, mutex, userMsg)
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

func (wsm *wsManager) writeMessage(conn *websocket.Conn, service string, messageType string, mutex *sync.Mutex, userMsg models.UserMessage) {
	mutex.Lock()
	defer mutex.Unlock()

	if service == userMsg.Service && messageType == userMsg.MessageType {
		var err error
		err = conn.WriteJSON(userMsg)

		if err != nil {
			logger.Errorf(nil, "Failed to send message [%+v] to [%+v], error: %+v", userMsg, conn, err)
		}
		logger.Debugf(nil, "Message sent [%+v]", userMsg)
	}
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
			logger.Errorf(nil, "Upgrade websocket_manager request failed: %+v", err)
			return
		}

		receiver := Receiver{
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

func (wsm *wsManager) ReceiveMsg(client *nfclient.Client, channelClient pb.Notification_CreateNotificationChannelClient, ctx context.Context) chan models.UserMessage {
	var msgChan = make(chan models.UserMessage, 255)
	go wsm.getMsgsByGrpcStream(client, channelClient, ctx, msgChan)
	return msgChan
}

func (wsm *wsManager) getMsgsByGrpcStream(client *nfclient.Client, channelClient pb.Notification_CreateNotificationChannelClient, ctx context.Context, msgChan chan models.UserMessage) {
	for {
		userWsMsgStreamData, err := channelClient.Recv()
		if err != nil {
			logger.Errorf(nil, "failed to recv: %+v", err)
			continue
		}

		userMsg := models.PbToUserMessage(userWsMsgStreamData.UserMsg)
		if userMsg.Service == wsm.accessServiceName {
			msgChan <- *userMsg
		}
		if err == io.EOF {
			break
		}
	}

}
