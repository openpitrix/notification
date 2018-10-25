package main

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

func main() {
	m := gomail.NewMessage()
	m.SetHeader("From", "openpitrix@app-center.cn")
	m.SetHeader("To", "huojiao2006@163.com")
	//m.SetHeader("To", "13009254@qq.com")
	m.SetHeader("Subject", "nf from goland,huojiao!")
	m.SetBody("text/plain", "终于搞好了！!")

	d := gomail.NewDialer("mail.app-center.cn", 25, "openpitrix@app-center.cn", "openpitrix")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}

//dGVzdEB0ZXN0LmNvbQ==
//dGVzdA==