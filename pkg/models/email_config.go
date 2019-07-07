// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.
package models

import (
	"time"

	"openpitrix.io/notification/pkg/pb"
	"openpitrix.io/notification/pkg/util/pbutil"
)

type EmailConfig struct {
	Protocol      string    `gorm:"column:protocol"`
	EmailHost     string    `gorm:"column:email_host"`
	Port          uint32    `gorm:"column:port"`
	DisplaySender string    `gorm:"column:display_sender"`
	Email         string    `gorm:"column:email"`
	Password      string    `gorm:"column:password"`
	SSLEnable     bool      `gorm:"column:ssl_enable"`
	CreateTime    time.Time `gorm:"column:create_time"`
	StatusTime    time.Time `gorm:"column:status_time"`
}

//table name
const (
	TableEmailConfig = "email_config"
)

//field name
const (
	EmailCfgColProtocol      = "protocol"
	EmailCfgColEmailHost     = "email_host"
	EmailCfgColPort          = "port"
	EmailCfgColDisplaySender = "display_sender"
	EmailCfgColEmail         = "email"
	EmailCfgColPassword      = "password"
	EmailCfgColSSLEnable     = "ssl_enable"
	EmailCfgColCreateTime    = "create_time"
	EmailCfgColStatusTime    = "status_time"
)

func NewEmailConfig(req *pb.ServiceConfig) *EmailConfig {
	emailCfg := &EmailConfig{
		Protocol:      req.GetEmailServiceConfig().GetProtocol().GetValue(),
		EmailHost:     req.GetEmailServiceConfig().GetEmailHost().GetValue(),
		Port:          req.GetEmailServiceConfig().GetPort().GetValue(),
		DisplaySender: req.GetEmailServiceConfig().GetDisplaySender().GetValue(),
		Email:         req.GetEmailServiceConfig().GetEmail().GetValue(),
		Password:      req.GetEmailServiceConfig().GetPassword().GetValue(),
		SSLEnable:     req.GetEmailServiceConfig().GetSslEnable().GetValue(),
		CreateTime:    time.Now(),
		StatusTime:    time.Now(),
	}
	return emailCfg
}

func EmailConfigToPb(emailConfig *EmailConfig) *pb.EmailServiceConfig {
	pbEmailConfig := pb.EmailServiceConfig{}
	pbEmailConfig.Protocol = pbutil.ToProtoString(emailConfig.Protocol)
	pbEmailConfig.EmailHost = pbutil.ToProtoString(emailConfig.EmailHost)
	pbEmailConfig.Port = pbutil.ToProtoUInt32(emailConfig.Port)
	pbEmailConfig.DisplaySender = pbutil.ToProtoString(emailConfig.DisplaySender)
	pbEmailConfig.Email = pbutil.ToProtoString(emailConfig.Email)
	pbEmailConfig.Password = pbutil.ToProtoString(emailConfig.Password)
	pbEmailConfig.SslEnable = pbutil.ToProtoBool(emailConfig.SSLEnable)
	return &pbEmailConfig
}
