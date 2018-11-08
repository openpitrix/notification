package emailutil

import (
	"openpitrix.io/notification/pkg/services/test"
	"testing"
)

func TestSendMail(t *testing.T) {
	test.InitGlobelSetting()
	emailaddr:="huojiao2006@163.com"
	header:="Subject-hello from Openpitrix notication"
	body:="Body-hello from Openpitrix notication"
	SendMail(emailaddr,header,body)
}
