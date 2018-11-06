package emailutil

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func SendMail(emailaddr string,header string,body string)    {
	m := gomail.NewMessage()
	m.SetHeader("From", "openpitrix@app-center.cn")
	m.SetHeader("To", emailaddr)
	//m.SetHeader("Cc",
	//	m.FormatAddress("513590612@qq.com", "收件人")) //抄送
	m.SetHeader("Subject", header)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("mail.app-center.cn", 25, "openpitrix@app-center.cn", "openpitrix")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}