package config

import (
	"log"
	"testing"
	"openpitrix.io/logger"
)

func TestGetCfg(t *testing.T) {
    logger.SetLevelByString("debug")
	cfg:=GetInstance()
	cfg.InitCfg()
	log.Println(cfg.Etcd)
	cfg.Print()
	cfg.Validate()
}