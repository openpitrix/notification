package nf

import (
	"log"
	"notification/pkg/models"
	"testing"
	"time"
)


func TestNewServices(t *testing.T) {
	log.Println("Test NewServices")
	nfs, _ := NewServices()
	nfs.GetDataFromDB4Test()
}


func TestSayHello(t *testing.T) {
	log.Println("Test SayHello")
	nfs, _ := NewServices()
	nfs.SayHello("TestSayHello")
}


func TestCreateNfWaddrs2(t *testing.T) {
	log.Println("Test CreateNfWaddrs2")
	nfs, _ := NewServices()
	nf := &models.NotificationCenterPost{
		NfPostID:        "2",
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
	err:=nfs.CreateNfWaddrs2(nf)
	if err != nil {
		log.Println("something is wrong")
	}
}
