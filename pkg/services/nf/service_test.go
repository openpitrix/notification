package nf

import (
	"log"
	"testing"
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
