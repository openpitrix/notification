package emailutil

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
)

func SendMail(emailaddr string, header string, body string) error {
	logger.Debugf(nil, "emailaddr="+emailaddr)

	m := gomail.NewMessage()
	m.SetHeader("From", "openpitrix@app-center.cn")
	m.SetHeader("To", emailaddr)
	m.SetHeader("Subject", header)
	m.SetBody("text/plain", body)

	host := config.GetInstance().Email.EmailHost
	port := config.GetInstance().Email.EmailPort
	username := config.GetInstance().Email.EmailUsername
	password := config.GetInstance().Email.EmailPassword

	//logger.Debugf(nil, "host="+host)
	//logger.Debugf(nil, "%v",port)
	//logger.Debugf(nil, "username="+username)
	//logger.Debugf(nil, "password="+password)

	d := gomail.NewDialer(host, port , username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		logger.Warnf(nil, "%+v", err)
		return err
	}
	return nil
}
