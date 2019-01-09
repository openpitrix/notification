package notification

import (
	"log"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/services/test"
	"openpitrix.io/notification/pkg/util/idutil"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	test.InitGlobelSetting4Test()
	db, q := test.GetTestDBAndEtcd4Test()
	nfservice := NewService(db, q)
	nfservice.SayHello("ssss")
}

func TestSayHello(t *testing.T) {
	db, q := test.GetTestDBAndEtcd4Test()
	nfservice := NewService(db, q)
	nfservice.SayHello("ssss")
}

func TestCreateNfWaddrs(t *testing.T) {
	db, q := test.GetTestDBAndEtcd4Test()
	nfservice := NewService(db, q)
	testAddrsStr := "johuo@yunify.com;513590612@qq.com"

	nf := &models.Notification{
		NotificationId: idutil.GetUuid(constants.NfPostIDPrifix),
		ContentType:    "Email",
		AddrsStr:       testAddrsStr,
		Title:          "Title Test2",
		Content:        "Content2",
		ShortContent:   "ShortContent2",
		ExporedDays:    5,
		Owner:          "Huojiao",
		Status:         "New",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	nfId, err := nfservice.CreateNfWithAddrs(nf)
	if err != nil {
		logger.Criticalf(nil, "Cannot create NfWaddrs:%+v", err)
	}
	logger.Debugf(nil, nfId)
}

func TestDescribeNfs(t *testing.T) {
	nfID := "nf-KV4oN8ROJqPE"
	log.Println("TestDescribeNfs")
	db, q := test.GetTestDBAndEtcd4Test()
	nfservice := NewService(db, q)
	nf, err := nfservice.DescribeNfs(nfID)
	logger.Infof(nil, "%+v", nf)

	if err != nil {
		logger.Warnf(nil, "%+v", err)
	}
}
