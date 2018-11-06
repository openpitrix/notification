package nf

import (
	"fmt"
	"log"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/services/test"
	"openpitrix.io/notification/pkg/util/idutil"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	log.Println("Test NewServices")
	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	nfservice := NewService(db, q)
	nfservice.SayHello("ssss")
}

func TestSayHello(t *testing.T) {
	log.Println("Test NewServices")

	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	nfservice := NewService(db, q)
	nfservice.SayHello("ssss")
}

func TestCreateNfWaddrs(t *testing.T) {
	log.Println("Test CreateNfWaddrs")

	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	nfservice := NewService(db, q)

	nf := &models.NotificationCenterPost{
		NfPostID:     idutil.GetUuid(constants.NfPostIDPrifix),
		NfPostType:   "Email",
		//AddrsStr:     "johuo@yunify.com;danma@yunify.com",
		AddrsStr:     "johuo@yunify.com;johuo@yunify.com;johuo@yunify.com;johuo@yunify.com;johuo@yunify.com;huojiao2006@163.com;huojiao2006@163.com;huojiao2006@163.com;huojiao2006@163.com;huojiao2006@163.com",
		Title:        "Title Test2",
		Content:      "Content2",
		ShortContent: "ShortContent2",
		ExporedDays:  5,
		Owner:        "Huojiao",
		Status:       "New",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err := nfservice.CreateNfWaddrs(nf)
	if err != nil {
		log.Println("something is wrong")
	}
}

func TestDescribeNfs(t *testing.T) {
	nfID := "nf-KV4oN8ROJqPE"
	log.Println("TestDescribeNfs")
	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	nfservice := NewService(db, q)
	nf, err := nfservice.DescribeNfs(nfID)
	if err != nil {
		fmt.Println(nf)
	}
}
