// Copyright 2019 The OpenPitrix Authors. All rights reserved.
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

	emailaddr := "openpitrix@163.com"
	header := "email_test.go sends an email."
	body := "<p>Content:email_test.go sends an email!</p>"
	err := SendMail(nil, emailaddr, header, body)

	if err != nil {
		logger.Errorf(nil, "send email failed, [%+v]", err)
	}
}
