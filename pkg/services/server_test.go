package services

import (
	"testing"
)

func TestNewServer(t *testing.T){
	server, _ :=NewServer()
	server.nfservice.SayHello("ssss")
	server.nfservice.GetDataFromDB4Test()
}
