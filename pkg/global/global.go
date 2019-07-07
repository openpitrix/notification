// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package global

import (
	"os"
	"sync"

	"github.com/google/gops/agent"
	"github.com/jinzhu/gorm"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	nfdb "openpitrix.io/notification/pkg/db"
	wstypes "openpitrix.io/notification/pkg/services/websocket/types"
)

type GlobalCfg struct {
	cfg          *config.Config
	database     *gorm.DB
	pubsubClient *wstypes.PubsubClient
}

var instance *GlobalCfg
var once sync.Once

func GetInstance() *GlobalCfg {
	once.Do(func() {
		instance = newGlobalCfg()
	})
	return instance
}

func newGlobalCfg() *GlobalCfg {
	cfg := config.GetInstance().LoadConf()
	g := &GlobalCfg{cfg: cfg}

	g.setLoggerLevel()
	g.openDatabase()
	g.setPubSubClient()

	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true,
	}); err != nil {
		logger.Criticalf(nil, "Failed to start gops agent")
	}
	return g
}

func (g *GlobalCfg) openDatabase() *GlobalCfg {
	if g.cfg.Mysql.Disable {
		logger.Debugf(nil, "%+s", "Database setting for Mysql.Disable is true.")
		return g
	}
	isSucc := nfdb.GetInstance().InitDataPool()

	if !isSucc {
		logger.Criticalf(nil, "%+s", "Init database pool failure...")
		os.Exit(1)
	}
	logger.Debugf(nil, "%+s", "Init database pool successfully.")

	db := nfdb.GetInstance().GetMysqlDB()
	g.database = db
	logger.Debugf(nil, "%+s", "Set globalcfg database value.")
	return g
}

func (g *GlobalCfg) setLoggerLevel() *GlobalCfg {
	AppLogMode := config.GetInstance().Log.Level
	logger.SetLevelByString(AppLogMode)
	logger.Infof(nil, "Set app log level to %+s", AppLogMode)
	return g
}

func (g *GlobalCfg) setPubSubClient() *GlobalCfg {
	pubsubConnStr := g.cfg.PubSub.Addr
	pubsubType := g.cfg.PubSub.Type

	pubsubConfigMap := map[string]interface{}{
		"connStr": pubsubConnStr,
	}

	psClient, err := wstypes.New(pubsubType, pubsubConfigMap)
	if err != nil {
		logger.Errorf(nil, "Failed to connect pubsub server: %+v.", err)
	}

	g.pubsubClient = &psClient
	return g
}

func (g *GlobalCfg) GetDB() *gorm.DB {
	return g.database
}
func (g *GlobalCfg) GetPubSubClient() *wstypes.PubsubClient {
	return g.pubsubClient
}
