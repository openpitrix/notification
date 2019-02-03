// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package emailutil

import (
	"context"
	"crypto/tls"

	gomail "gopkg.in/gomail.v2"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
)

func SendMail(ctx context.Context, emailAddr string, header string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "openpitrix@app-center.cn")
	m.SetHeader("To", emailAddr)
	m.SetHeader("Subject", header)
	m.SetBody("text/html", body)

	host := config.GetInstance().Email.EmailHost
	port := config.GetInstance().Email.Port
	email := config.GetInstance().Email.Email
	password := config.GetInstance().Email.Password

	d := gomail.NewDialer(host, port, email, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)
		return err
	}
	return nil
}
