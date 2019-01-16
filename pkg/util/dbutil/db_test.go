// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package dbutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/jinzhu/gorm"

	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestGetMysqlDB(t *testing.T) {

	logger.Debugf(nil, "step0.1:初始化配置参数")
	config.GetInstance().LoadConf()

	logger.Debugf(nil, "step0.2:初始化DB connection pool")
	issucc := GetInstance().InitDataPool()
	if !issucc {
		logger.Criticalf(nil, "init database pool failure...")
		os.Exit(1)
	}

	db = GetInstance().GetMysqlDB()

	// 读取
	var product Product
	db.First(&product, 1)                   // 查询id为1的product
	db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
	fmt.Println(product)

}
