package emailutil

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func SendMail(emailaddr string)    {
	m := gomail.NewMessage()
	m.SetHeader("From", "openpitrix@app-center.cn")
	m.SetHeader("To", emailaddr)
	m.SetHeader("Subject", "Subject-hello from Openpitrix notication")
	m.SetBody("text/plain", "Body-hello from Openpitrix notication")

	d := gomail.NewDialer("mail.app-center.cn", 25, "openpitrix@app-center.cn", "openpitrix")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}