// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package config

import (
	"os"
	"testing"

	"openpitrix.io/logger"
)

func TestLoadConf(t *testing.T) {

	logger.SetLevelByString("debug")

	os.Setenv("NOTIFICATION_LOG_LEVEL", "debug")
	os.Setenv("NOTIFICATION_GRPC_SHOW_ERROR_CAUSE", "false")

	os.Setenv("NOTIFICATION_MYSQL_HOST", "MYSQL_HOST_test")
	os.Setenv("NOTIFICATION_MYSQL_PORT", "13306")

	os.Setenv("NOTIFICATION_APP_API_HOST", "TESTAPP_API_HOST")

	os.Setenv("NOTIFICATION_APP_MAX_WORKING_NOTIFICATIONS", "11")
	os.Setenv("NOTIFICATION_APP_MAX_WORKING_TASKS", "11")

	mycfg := GetInstance()
	mycfg.LoadConf()

	logger.Debugf(nil, "Other=========================================")
	logger.Debugf(nil, "NOTIFICATION_LOG_LEVEL : %+v", mycfg.Log.Level)
	logger.Debugf(nil, "NOTIFICATION_GRPC_SHOW_ERROR_CAUSE : %+v", mycfg.Grpc.ShowErrorCause)
	logger.Debugf(nil, "")

	logger.Debugf(nil, "Mysql=========================================")
	logger.Debugf(nil, "NOTIFICATION_MYSQL_HOST : %+v", mycfg.Mysql.Host)
	logger.Debugf(nil, "NOTIFICATION_MYSQL_PORT : %+v", mycfg.Mysql.Port)
	logger.Debugf(nil, "NOTIFICATION_MYSQL_USER : %+v", mycfg.Mysql.User)
	logger.Debugf(nil, "NOTIFICATION_MYSQL_PASSWORD : %+v", mycfg.Mysql.Password)
	logger.Debugf(nil, "NOTIFICATION_MYSQL_DATABASE : %+v", mycfg.Mysql.Database)
	logger.Debugf(nil, "NOTIFICATION_MYSQL_DISABLE : %+v", mycfg.Mysql.Disable)
	logger.Debugf(nil, "NOTIFICATION_MYSQL_LOG_MODE : %+v", mycfg.Mysql.LogMode)
	logger.Debugf(nil, "")

	logger.Debugf(nil, "ETCD=========================================")
	logger.Debugf(nil, "NOTIFICATION_ETCD_ENDPOINTS : %+v", mycfg.Etcd.Endpoints)
	logger.Debugf(nil, "")

	logger.Debugf(nil, "Email=========================================")
	logger.Debugf(nil, "NOTIFICATION_EMAIL_PROTOCOL : %+v", mycfg.Email.Protocol)
	logger.Debugf(nil, "NOTIFICATION_EMAIL_EMAIL_HOST : %+v", mycfg.Email.EmailHost)
	logger.Debugf(nil, "NOTIFICATION_EMAIL_PORT : %+v", mycfg.Email.Port)
	logger.Debugf(nil, "NOTIFICATION_EMAIL_DISPLAY_SENDER : %+v", mycfg.Email.DisplaySender)
	logger.Debugf(nil, "NOTIFICATION_EMAIL_EMAIL : %+v", mycfg.Email.Email)
	logger.Debugf(nil, "NOTIFICATION_EMAIL_PASSWORD : %+v", mycfg.Email.Password)
	logger.Debugf(nil, "NOTIFICATION_EMAIL_SSL_ENABLE : %+v", mycfg.Email.SSLEnable)
	logger.Debugf(nil, "")

	logger.Debugf(nil, "App=========================================")
	logger.Debugf(nil, "NOTIFICATION_APP_HOST : %+v", mycfg.App.Host)
	logger.Debugf(nil, "NOTIFICATION_APP_PORT : %+v", mycfg.App.Port)
	logger.Debugf(nil, "NOTIFICATION_APP_API_HOST : %+v", mycfg.App.ApiHost)
	logger.Debugf(nil, "NOTIFICATION_APP_API_PORT : %+v", mycfg.App.ApiPort)
	logger.Debugf(nil, "NOTIFICATION_APP_MAX_WORKING_NOTIFICATIONS : %+v", mycfg.App.MaxWorkingNotifications)
	logger.Debugf(nil, "NOTIFICATION_APP_MAX_WORKING_TASKS : %+v", mycfg.App.MaxWorkingTasks)
	logger.Debugf(nil, "")

	mycfg.PrintUsage()

}
