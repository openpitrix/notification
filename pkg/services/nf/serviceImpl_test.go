package nf

import (
	"log"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/dbutil"
	"os"
	"testing"
	"time"
)


func TestNewService(t *testing.T) {
	log.Println("Test NewServices")
	db := dbutil.GetInstance().GetMysqlDB()
	nfservice := NewService(db)
	nfservice.SayHello("ssss")
}


func TestSayHello(t *testing.T) {
	log.Println("Test NewServices")

	issucc :=  dbutil.GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}
	db := dbutil.GetInstance().GetMysqlDB()

	nfservice := NewService(db)
	nfservice.SayHello("ssss")
}

func TestCreateNfWaddrs(t *testing.T) {
	log.Println("Test CreateNfWaddrs")

	//set mysql db,init database pool
	issucc :=  dbutil.GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}
	db := dbutil.GetInstance().GetMysqlDB()
	nfservice := NewService(db)

	nf := &models.NotificationCenterPost{
		NfPostID:       CreatenfPostID(),
		NfPostType:  "Email",
		AddrsStr:"johuo@yunify.com;danma@yunify.com",
		Title:  "Title Test",
		Content: "Content",
		ShortContent :  "ShortContent",
		ExporedDays :5,
		Owner :  "Huojiao",
		Status:"New",
		CreatedAt:time.Now(),
		UpdatedAt:time.Now(),
		DeletedAt:time.Now(),
	}
	err:=nfservice.CreateNfWaddrs(nf)
	if err != nil {
		log.Println("something is wrong")
	}
}
