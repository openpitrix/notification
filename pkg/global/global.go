// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package global

import (
	"errors"
	"os"
	"sync"

	"github.com/google/gops/agent"
	"github.com/jinzhu/gorm"
	i "openpitrix.io/libqueue"
	qetcd "openpitrix.io/libqueue/etcd"
	q "openpitrix.io/libqueue/queue"
	qredis "openpitrix.io/libqueue/redis"
	"openpitrix.io/logger"

	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	nfdb "openpitrix.io/notification/pkg/db"
)

type GlobalCfg struct {
	cfg         *config.Config
	database    *gorm.DB
	queueClient *i.IClient
	pubsub      *i.IPubSub
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
	g.setQueueClient()
	if config.GetInstance().Websocket.Service != "none" {
		_, err := g.setPubSub()
		if err != nil {
			logger.Errorf(nil, "Failed to set pubsub,err=%+v", err)
		}
	}

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

func (g *GlobalCfg) GetDB() *gorm.DB {
	return g.database
}

func (g *GlobalCfg) setQueueClient() *GlobalCfg {
	pubsubConnStr := g.cfg.Queue.Addr
	pubsubType := g.cfg.Queue.Type

	pubsubConfigMap := map[string]interface{}{
		"connStr": pubsubConnStr,
	}

	qClient, err := q.NewIClient(pubsubType, pubsubConfigMap)
	if err != nil {
		logger.Errorf(nil, "Failed to connect %s pubsub server: %+v.", pubsubType, err)
	}

	g.queueClient = &qClient
	return g
}

func (g *GlobalCfg) GetQueueClient() *i.IClient {
	return g.queueClient
}

func (g *GlobalCfg) setPubSub() (*GlobalCfg, error) {
	queueType := config.GetInstance().Queue.Type
	var ipubsub i.IPubSub
	if queueType == constants.QueueTypeRedis {
		redisPubSub := qredis.RedisPubSub{}
		ipubsub = &redisPubSub
	} else if queueType == constants.QueueTypeEtcd {
		etcdPubSub := qetcd.EtcdPubSub{}
		ipubsub = &etcdPubSub
	} else {
		return nil, errors.New("Unsupport queue type, currently support redis and etcd.")
	}

	ipubsub.SetClient(g.queueClient)
	g.pubsub = &ipubsub
	return g, nil
}

func (g *GlobalCfg) GetPubSub() *i.IPubSub {
	return g.pubsub
}
