package test

import (
	"github.com/jinzhu/gorm"
	"log"
	"openpitrix.io/notification/pkg/util/dbutil"
	"openpitrix.io/notification/pkg/util/etcdutil"
	"os"

)

func GetTestDB()  *gorm.DB {
	issucc := dbutil.GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}
	db := dbutil.GetInstance().GetMysqlDB()
	return db
}


func GetEtcdQueue()   *etcdutil.Queue {
	endpoints := []string{"192.168.0.7:2379"}
	prefix := "test"
	nfetcd, err := etcdutil.Connect(endpoints, prefix)
	log.Println(nfetcd)
	if err != nil {
		log.Fatal(err)
	}
	q := nfetcd.NewQueue("nf_task")
	return q
}