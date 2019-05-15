// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package websocket

import (
	"github.com/gorilla/websocket"

	"openpitrix.io/notification/pkg/util/idutil"
)

const ExpireTime = 60 // second

type userMessage struct {
	UserId      string
	MessageType string
	Message     Message
}

type Message struct {
	MessageId   string `json:"ws_message_id,omitempty"`
	UserId      string `json:"ws_user_id,omitempty"`
	MessageType string `json:"ws_message_type,omitempty"`
	Message     string `json:"ws_message,omitempty"`
}

type receiver struct {
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
