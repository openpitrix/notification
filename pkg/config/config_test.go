package config

import (
	"openpitrix.io/logger"
	"testing"
)

func TestGetCfg(t *testing.T) {
	logger.SetLevelByString("debug")
	cfg := GetInstance()
	cfg.InitCfg()
	cfg.Print()
	cfg.Validate()
}


func TestLoadConf(t *testing.T) {
	logger.SetLevelByString("debug")
	cfg := GetInstance()
 	cfg.LoadConf()
	cfg.Print()

}
