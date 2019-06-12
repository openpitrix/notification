// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package types

import (
	"github.com/gorilla/websocket"
	"openpitrix.io/notification/pkg/util/idutil"
)

const ExpireTime = 60 // second

type UserMessage struct {
	UserId        string
	Service       string
	MessageType   string
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
