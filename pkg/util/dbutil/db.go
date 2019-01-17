// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package dbutil

import (
	"fmt"
	"log"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"openpitrix.io/notification/pkg/config"
)

/*
* MysqlConnPool
* 数据库连接操作库
* 基于gorm封装开发
 */
type MysqlConnPool struct {
}

var instance *MysqlConnPool
var once sync.Once

var db *gorm.DB
var err error

func GetInstance() *MysqlConnPool {
	once.Do(func() {
		instance = &MysqlConnPool{}
	})
	return instance
}

/*
* @fuc 初始化数据库连接
 */
func (m *MysqlConnPool) InitDataPool() (isSucc bool) {
	cfg := config.GetInstance()

	var (
		dbCfg            = cfg.Mysql
		connectionString = fmt.Sprintf(
			"%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.Database,
		)
	)

	db, err = gorm.Open("mysql", connectionString)
	if err != nil {
		log.Print(err)
		return false
	}

	err = db.DB().Ping()

	if err != nil {
		return false
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(cfg.Mysql.LogMode)

	// 全局禁用表名复数
	db.SingularTable(true)

	if err != nil {
		log.Fatal(err)
		return false
	}
	//关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true
}

/*
* 对外获取数据库连接对象db
 */
func (m *MysqlConnPool) GetMysqlDB() *gorm.DB {
	return db
}
