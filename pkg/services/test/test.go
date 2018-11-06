package test

import (
	"github.com/jinzhu/gorm"
	"log"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/util/dbutil"
	"openpitrix.io/notification/pkg/util/etcdutil"
	"os"
)

func GetTestDB() *gorm.DB {
	InitGlobelSetting()
	db := dbutil.GetInstance().GetMysqlDB()
	db.LogMode(true)
	return db
}

func GetEtcdQueue() *etcdutil.Queue {
	endpoints := []string{"192.168.0.7:2379"}
	prefix := "test"
	nfetcd, err := etcdutil.Connect(endpoints, prefix)
	if err != nil {
		log.Fatal(err)
	}
	q := nfetcd.NewQueue("nf_task")
	return q
}


func InitGlobelSetting() {
	logger.Debugf(nil,"step0.1:初始化配置参数11111")
	config.GetInstance().InitCfg()

	logger.Debugf(nil,"step0.2:初始化DB connection pool")
	issucc := dbutil.GetInstance().InitDataPool()
	if !issucc {
		logger.Debugf(nil,"init database pool failure...")
		os.Exit(1)
	}

	AppLogMode:=config.GetInstance().AppLogMode
	logger.SetLevelByString(AppLogMode)

}