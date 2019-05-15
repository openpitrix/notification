// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"openpitrix.io/logger"

	pkg "openpitrix.io/notification/pkg"
	"openpitrix.io/notification/pkg/config"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestGetMysqlDB(t *testing.T) {
	if !*pkg.LocalDevEnvEnabled {
		t.Skip("LocalDevEnv disabled")
	}

	logger.Debugf(nil, "step0.1:init params")
	config.GetInstance().LoadConf()

	logger.Debugf(nil, "step0.2:init db connection pool")
	isSucc := GetInstance().InitDataPool()
	if !isSucc {
		logger.Criticalf(nil, "init database pool failure...")
		os.Exit(1)
	}

	db = GetInstance().GetMysqlDB()

}
