package emailutil

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"openpitrix.io/logger"
)

func SendMail(emailaddr string,header string,body string)  error   {
	m := gomail.NewMessage()
	m.SetHeader("From", "openpitrix@app-center.cn")
	m.SetHeader("To", emailaddr)
	m.SetHeader("Subject", header)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("mail.app-center.cn", 25, "openpitrix@app-center.cn", "openpitrix")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		logger.Warnf(nil, "%+v", err)
		return  err
	}
	return nil
}