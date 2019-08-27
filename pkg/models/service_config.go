// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

const (
	ServiceType = "service_type"
)

const (
	ServiceCfgProtocol     = "protocol"
	ServiceCfgEmailHost    = "email_host"
	ServiceCfgPort         = "port"
	ServiceCfgDisplayEmail = "display_email"
	ServiceCfgEmail        = "email"
	ServiceCfgPassword     = "password"
)

const (
	TestEmailRecipient = "test_email_recipient"
)

const (
	ProtocolTypeSMTP = "SMTP"
	ProtocolTypePOP3 = "POP3"
	ProtocolTypeIMAP = "IMAP"
)

var ProtocolTypes = []string{
	ProtocolTypeSMTP,
	ProtocolTypePOP3,
	ProtocolTypeIMAP,
}
