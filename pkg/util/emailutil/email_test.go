// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package emailutil

import (
	"testing"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
)

func TestSendMail(t *testing.T) {
	config.GetInstance().LoadConf()
	emailaddr := "huojiao2006@163.com"
	header := "Subject-hello from Openpitrix notication"
	body := "Body-hello from Openpitrix notication"
	err := SendMail(emailaddr, header, body)

	if err != nil {
		logger.Warnf(nil, "%+v", err)
		logger.Errorf(nil, "send email failed, [%+v]", err)
	}
}
