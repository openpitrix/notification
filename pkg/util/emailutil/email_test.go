package emailutil

import (
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/services/test"
	"testing"
)

func TestSendMail(t *testing.T) {
	test.InitGlobelSetting4Test()
	emailaddr := "huojiao2006@163.com"
	header := "Subject-hello from Openpitrix notication"
	body := "Body-hello from Openpitrix notication"
	err := SendMail(emailaddr, header, body)

	if err != nil {
		logger.Warnf(nil, "%+v", err)
		logger.Errorf(nil, "send email failed, [%+v]", err)
	}
}
