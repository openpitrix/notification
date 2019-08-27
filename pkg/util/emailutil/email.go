// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package emailutil

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"text/template"

	gomail "gopkg.in/gomail.v2"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/pb"
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

	//note:if the smtp server supports certificate,should not skip InsecureSkipVerify
	//if just internal smtp server and without certification, should skip InsecureSkipVerify,otherwise email can not be sent out by stmp.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !sslEnable}

	if err := d.DialAndSend(m); err != nil {
		logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)

		//Attention!!! As Gomail does not support using plainauth without STARTTLS,
		//so if the email server is without TSL setting, the mail can not be sent.
		//Here is to add noStartTLSPlainAuth to support this scenario.
		if !sslEnable && err.Error() == errors.New("unencrypted connection").Error() {
			logger.Debugf(ctx, "Try to use noStartTLSPlainAuth to send mail.")
			d.Auth = &noStartTLSPlainAuth{
				identity: "",
				username: d.Username,
				password: d.Password,
				host:     d.Host,
			}
			if err = d.DialAndSend(m); err != nil {
				logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)
				return err
			}
			logger.Debugf(ctx, "Send out mail by noStartTLSPlainAuth successfully.")
		} else {
			logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)
			return err
		}
	}

	return nil
}

type EmailIcon struct {
	Icon string
}

func getDefaultMessage(iconstr string) string {
	t, _ := template.New("validationEmail").Parse(constants.ValidationEmailNotifyTemplate)

	b := bytes.NewBuffer([]byte{})
	emailIcon := &EmailIcon{
		Icon: iconstr,
	}

	t.Execute(b, emailIcon)
	return b.String()
}

func SendMail4ValidateEmailService(ctx context.Context, emailServiceConfig *pb.EmailServiceConfig, testEmailRecipient string) error {
	host := emailServiceConfig.GetEmailHost().GetValue()
	port := emailServiceConfig.GetPort().GetValue()
	email := emailServiceConfig.GetEmail().GetValue() //smtp user
	password := emailServiceConfig.GetPassword().GetValue()
	displaySender := emailServiceConfig.GetDisplaySender().GetValue()
	sslEnable := emailServiceConfig.GetSslEnable().GetValue()
	icon := emailServiceConfig.GetValidationIcon().GetValue()
	title := emailServiceConfig.GetValidationTitle().GetValue()
	if testEmailRecipient == "" {
		testEmailRecipient = email
	}

	body := getDefaultMessage(icon)

	m := gomail.NewMessage()
	m.SetAddressHeader("From", email, displaySender)
	m.SetHeader("To", testEmailRecipient)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, int(port), email, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: !sslEnable}
	if err := d.DialAndSend(m); err != nil {
		logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", testEmailRecipient, err)

		//Attention!!! As Gomail does not support using plainauth without STARTTLS,
		//so if the email server is without TSL setting, the mail can not be sent.
		//Here is to add noStartTLSPlainAuth to support this scenario.
		if !sslEnable && err.Error() == errors.New("unencrypted connection").Error() {
			logger.Debugf(ctx, "Try to use noStartTLSPlainAuth to send mail.")
			d.Auth = &noStartTLSPlainAuth{
				identity: "",
				username: d.Username,
				password: d.Password,
				host:     d.Host,
			}
			if err = d.DialAndSend(m); err != nil {
				logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", testEmailRecipient, err)
				return err
			}
			logger.Debugf(ctx, "Send out mail by noStartTLSPlainAuth successfully.")
		} else {
			logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", testEmailRecipient, err)
			return err
		}
	}

	return nil
}
