// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"github.com/gorilla/websocket"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/idutil"
	"openpitrix.io/notification/pkg/util/jsonutil"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type UserMessage struct {
	UserId        string `json:"user_id,omitempty"`
	Service       string `json:"service,omitempty"`
	MessageType   string `json:"message_type,omitempty"`
	MessageDetail MessageDetail
}

type MessageDetail struct {
	MessageId      string `json:"ws_message_id,omitempty"`
	UserId         string `json:"ws_user_id,omitempty"`
	Service        string `json:"ws_service,omitempty"`
	MessageType    string `json:"ws_message_type,omitempty"`
	MessageContent string `json:"ws_message,omitempty"`
}

type Receiver struct {
	Service     string
	MessageType string
	UserId      string
	Conn        *websocket.Conn
}

const (
	WsMessageIdPrefix = "msg-"
)

func NewWsMessageId() string {
	return idutil.GetUuid(WsMessageIdPrefix)
}

func MessageDetailToPb(userMsgDetail *MessageDetail) *pb.MessageDetail {
	pbUserMsgDetail := pb.MessageDetail{}
	pbUserMsgDetail.MessageId = pbutil.ToProtoString(userMsgDetail.MessageId)
	pbUserMsgDetail.Service = pbutil.ToProtoString(userMsgDetail.Service)
	pbUserMsgDetail.MessageType = pbutil.ToProtoString(userMsgDetail.MessageType)
	pbUserMsgDetail.MessageContent = pbutil.ToProtoString(userMsgDetail.MessageContent)
	return &pbUserMsgDetail
}

func UserMessageToPb(userMsg *UserMessage) *pb.UserMessage {
	pbUserMsg := pb.UserMessage{}
	pbUserMsg.UserId = pbutil.ToProtoString(userMsg.UserId)
	pbUserMsg.MessageType = pbutil.ToProtoString(userMsg.MessageType)
	pbUserMsg.Service = pbutil.ToProtoString(userMsg.Service)
	pbUserMsg.MsgDetail = MessageDetailToPb(&(userMsg.MessageDetail))
	return &pbUserMsg
}

func UseMsgStringToPb(dataStr string) (*pb.UserMessage, error) {
	userMsg := new(UserMessage)
	err := jsonutil.Decode([]byte(dataStr), userMsg)
	if err != nil {
		logger.Errorf(nil, "Decode [%s] into models.UserMessage failed: %+v", dataStr, err)
	}
	pbUserMsg := UserMessageToPb(userMsg)
	return pbUserMsg, err
}

func PbToUserMessage(pbUserMsg *pb.UserMessage) *UserMessage {
	userMsg := UserMessage{}
	userMsg.MessageType = pbUserMsg.MessageType.GetValue()
	userMsg.Service = pbUserMsg.Service.GetValue()
	userMsg.UserId = pbUserMsg.UserId.GetValue()
	pbMsgDetail := pbUserMsg.MsgDetail
	MsgDetail := PbToMessageDetail(pbMsgDetail)
	userMsg.MessageDetail = *MsgDetail
	return &userMsg
}

func PbToMessageDetail(pbMsgDetail *pb.MessageDetail) *MessageDetail {
	msgDetail := MessageDetail{}
	msgDetail.UserId = pbMsgDetail.UserId.GetValue()
	msgDetail.Service = pbMsgDetail.Service.GetValue()
	msgDetail.MessageType = pbMsgDetail.MessageType.GetValue()
	msgDetail.MessageContent = pbMsgDetail.MessageContent.GetValue()
	msgDetail.MessageId = pbMsgDetail.MessageId.GetValue()
	return &msgDetail
}
