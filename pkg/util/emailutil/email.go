// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package emailutil

import (
	"context"
	"crypto/tls"
	"errors"

	gomail "gopkg.in/gomail.v2"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
)

func SendMail(ctx context.Context, emailAddr string, header string, body string, fmtType string) error {
	host := config.GetInstance().Email.EmailHost
	port := config.GetInstance().Email.Port
	email := config.GetInstance().Email.Email
	password := config.GetInstance().Email.Password
	displaySender := config.GetInstance().Email.DisplaySender
	sslEnable := config.GetInstance().Email.SSLEnable

	m := gomail.NewMessage()
	m.SetAddressHeader("From", email, displaySender)
	m.SetHeader("To", emailAddr)
	m.SetHeader("Subject", header)
	contentType := "text/html"
	if fmtType == "normal" {
		contentType = "text/plain"
	}
	m.SetBody(contentType, body)

	d := gomail.NewDialer(host, port, email, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !sslEnable}

	if err := d.DialAndSend(m); err != nil {
		if !sslEnable && err == errors.New("unencrypted connection") {
			d.Auth = &unencryptedPlainAuth{
				identity: "",
				username: d.Username,
				password: d.Password,
				host:     d.Host,
			}
			if err = d.DialAndSend(m); err != nil {
				logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)
				return err
			}
		} else {
			logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)
			return err
		}
	}
	return nil
}
