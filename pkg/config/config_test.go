package config

import (
	"log"
	"testing"
)

func TestGetCfg(t *testing.T) {
	cfg:=GetInstance()
	cfg.InitCfg()
	log.Println(cfg.Etcd)
	cfg.Print()
	cfg.Validate()
}