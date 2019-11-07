// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package emailutil

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"regexp"
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
	//In fact email is the SMTP user which used to validate access to SMTP Server with password.
	//for some SMTP server, SMTP user is an email address, for other SMPT server, smtp user is not an email address.
	usernameOfSMTP := config.GetInstance().Email.Email
	password := config.GetInstance().Email.Password
	displaySender := config.GetInstance().Email.DisplaySender
	//sslEnable := config.GetInstance().Email.SSLEnable
	fromEmailAddr := config.GetInstance().Email.FromEmailAddr

	if fromEmailAddr == "" && VerifyEmailFmt(ctx, usernameOfSMTP) {
		fromEmailAddr = usernameOfSMTP
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", fromEmailAddr, displaySender)

	m.SetHeader("To", emailAddr)
	m.SetHeader("Subject", header)
	contentType := "text/html"
	if fmtType == "normal" {
		contentType = "text/plain"
	}
	m.SetBody(contentType, body)

	d := gomail.NewDialer(host, port, usernameOfSMTP, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)

		//Attention!!! As Gomail does not support using plainauth without STARTTLS,
		//so if the email server is without TSL setting, the mail can not be sent.
		//Here is to add noStartTLSPlainAuth to support this scenario.
		if err.Error() == errors.New("unencrypted connection").Error() {
			err = dialAndSendByNoStartTLSPlainAuth(ctx, emailAddr, d, m)
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

func getDefaultMessage(iconstr string, language string) string {
	notifytemplate := constants.ValidationEmailNotifyTemplate
	if language == "en" {
		notifytemplate = constants.ValidationEmailNotifyTemplateEn
	}
	t, _ := template.New("validationEmail").Parse(notifytemplate)

	b := bytes.NewBuffer([]byte{})
	emailIcon := &EmailIcon{
		Icon: iconstr,
	}

	t.Execute(b, emailIcon)
	return b.String()
}

func SendMail4ValidateEmailService(ctx context.Context, emailServiceConfig *pb.EmailServiceConfig, testEmailRecipient string, language string) error {
	host := emailServiceConfig.GetEmailHost().GetValue()
	port := emailServiceConfig.GetPort().GetValue()
	usernameOfSMTP := emailServiceConfig.GetEmail().GetValue() //smtp user
	password := emailServiceConfig.GetPassword().GetValue()
	displaySender := emailServiceConfig.GetDisplaySender().GetValue()
	icon := emailServiceConfig.GetValidationIcon().GetValue()
	title := emailServiceConfig.GetValidationTitle().GetValue()
	fromEmailAddr := emailServiceConfig.FromEmailAddr.GetValue()

	if fromEmailAddr == "" && VerifyEmailFmt(ctx, usernameOfSMTP) {
		fromEmailAddr = usernameOfSMTP
	}

	if testEmailRecipient == "" {
		testEmailRecipient = fromEmailAddr
	}

	body := getDefaultMessage(icon, language)

	m := gomail.NewMessage()
	m.SetAddressHeader("From", fromEmailAddr, displaySender)
	m.SetHeader("To", testEmailRecipient)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, int(port), usernameOfSMTP, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", testEmailRecipient, err)

		//Attention!!! As Gomail does not support using plainauth without STARTTLS,
		//so if the email server is without TSL setting, the mail can not be sent.
		//Here is to add noStartTLSPlainAuth to support this scenario.
		if err.Error() == errors.New("unencrypted connection").Error() {
			err = dialAndSendByNoStartTLSPlainAuth(ctx, testEmailRecipient, d, m)
		} else {
			logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", testEmailRecipient, err)
			return err
		}
	}

	return nil
}

func dialAndSendByNoStartTLSPlainAuth(ctx context.Context, emailAddr string, d *gomail.Dialer, m *gomail.Message) error {
	logger.Debugf(ctx, "Try to use noStartTLSPlainAuth to send mail.")
	d.Auth = &noStartTLSPlainAuth{
		identity: "",
		username: d.Username,
		password: d.Password,
		host:     d.Host,
	}
	if err := d.DialAndSend(m); err != nil {
		logger.Errorf(ctx, "Send email to [%s] failed, [%+v]", emailAddr, err)
		return err
	}
	logger.Debugf(ctx, "Send out mail by noStartTLSPlainAuth successfully.")
	return nil
}

func VerifyEmailFmt(ctx context.Context, emailStr string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	result := reg.MatchString(emailStr)
	if result {
		return true
	} else {
		return false
	}
}
