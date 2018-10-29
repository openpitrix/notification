// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package dbutil

import (
	"fmt"
	"log"
	"openpitrix.io/notification/pkg/models"
	"os"
	"testing"
)

func TestGetMysqlDB(t *testing.T) {
	// init database pool
	issucc := GetInstance().InitDataPool()
	if !issucc {
		log.Println("init database pool failure...")
		os.Exit(1)
	}

	db = GetInstance().GetMysqlDB()

	// 读取
	var product models.Product
	db.First(&product, 1) // 查询id为1的product
	db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
	fmt.Println(product)

}
