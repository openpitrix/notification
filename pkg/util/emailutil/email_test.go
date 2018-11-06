package emailutil

import "testing"

func TestSendMail(t *testing.T) {

	emailaddr:="huojiao2006@163.com"
	header:="Subject-hello from Openpitrix notication"
	body:="Body-hello from Openpitrix notication"
	SendMail(emailaddr,header,body)
}
