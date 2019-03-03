// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package resource_control

import (
	"os"
	"strconv"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

func GetEmailServiceConfig() *pb.EmailServiceConfig {
	cfg := config.GetInstance()

	protocol := cfg.Email.Protocol
	emailHost := cfg.Email.EmailHost
	port := cfg.Email.Port
	displaySender := cfg.Email.DisplaySender
	email := cfg.Email.Email
	password := cfg.Email.Password
	sslEnable := cfg.Email.SSLEnable

	emailCfg := &pb.EmailServiceConfig{
		Protocol:      pbutil.ToProtoString(protocol),
		EmailHost:     pbutil.ToProtoString(emailHost),
		Port:          pbutil.ToProtoString(string(port)),
		DisplaySender: pbutil.ToProtoString(displaySender),
		Email:         pbutil.ToProtoString(email),
		Password:      pbutil.ToProtoString(password),
		SslEnable:     pbutil.ToProtoBool(sslEnable),
	}

	return emailCfg
}

func SetServiceConfig(req *pb.ServiceConfig) {
	protocol := req.GetEmailServiceConfig().GetProtocol().GetValue()
	emailHost := req.GetEmailServiceConfig().GetEmailHost().GetValue()
	port := req.GetEmailServiceConfig().GetPort().GetValue()
	displaySender := req.GetEmailServiceConfig().GetDisplaySender().GetValue()
	email := req.GetEmailServiceConfig().GetEmail().GetValue()
	password := req.GetEmailServiceConfig().GetPassword().GetValue()
	sslEnable := req.GetEmailServiceConfig().GetSslEnable().GetValue()

	os.Setenv("NOTIFICATION_EMAIL_PROTOCOL", protocol)
	os.Setenv("NOTIFICATION_EMAIL_EMAIL_HOST", emailHost)
	os.Setenv("NOTIFICATION_EMAIL_PORT", port)
	os.Setenv("NOTIFICATION_EMAIL_DISPLAY_SENDER", displaySender)
	os.Setenv("NOTIFICATION_EMAIL_EMAIL", email)
	os.Setenv("NOTIFICATION_EMAIL_PASSWORD", password)
	os.Setenv("NOTIFICATION_EMAIL_SSL_ENABLE", strconv.FormatBool(sslEnable))

	config.GetInstance().LoadConf()
	logger.Infof(nil, "Set ServiceConfig successfully, [%+v].", config.GetInstance().Email)

}
