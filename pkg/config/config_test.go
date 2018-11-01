package config

import (
	"log"
	"testing"
)

func TestNewConfig(t *testing.T){
	cfg := NewConfig()
	println("cfg.SessionLifeTime.String():"+cfg.SessionLifeTime.String())
	println("cfg.App.AppName:"+cfg.App.AppName)
	println("cfg.Db.Host："+cfg.Db.Host)
}

func TestSayHello(t *testing.T){
	cfg := NewConfig()
	cfg.SayHello("test say hello")
}


func TestValidate(t *testing.T) {
	cfg := NewConfig()
	cfg.Validate()
	log.Println(cfg.App.AppName)
}

func TestPrint(t *testing.T) {
	cfg := NewConfig()
	cfg.Print()
	s :=  cfg.App.AppName
	println(s)

}