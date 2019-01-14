package globalcfg

import (
	"github.com/google/gops/agent"
	"github.com/jinzhu/gorm"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/util/dbutil"
	"openpitrix.io/openpitrix/pkg/etcd"
	"os"
	"strings"
	"sync"
)

type GlobalCfg struct {
	cfg      *config.Config
	database *gorm.DB
	etcd     *etcd.Etcd
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
	g.openEtcd()

	if err := agent.Listen(agent.Options{
		ShutdownCleanup: true,
	}); err != nil {
		logger.Criticalf(nil, "failed to start gops agent")
	}
	return g
}

func (g *GlobalCfg) openDatabase() *GlobalCfg {
	if g.cfg.Mysql.Disable {
		logger.Debugf(nil, "%+s", "Database setting for Mysql.Disable is true.")
		return g
	}
	issucc := dbutil.GetInstance().InitDataPool()
	logger.Debugf(nil, "%+s", "Init database pool successfully.")
	if !issucc {
		logger.Criticalf(nil, "%+s", "init database pool failure...")
		os.Exit(1)
	}

	db := dbutil.GetInstance().GetMysqlDB()
	g.database = db
	logger.Debugf(nil, "%+s", "Set globalcfg database value.")
	return g
}

func (g *GlobalCfg) openEtcd() *GlobalCfg {
	endpoints := strings.Split(g.cfg.Etcd.Endpoints, ",")
	e, err := etcd.Connect(endpoints, constants.EtcdPrefix)
	if err != nil {
		logger.Criticalf(nil, "%+s", "failed to connect etcd...")
		panic(err)
	}
	logger.Debugf(nil, "%+s", "Connect to etcd succesfully.")
	g.etcd = e
	logger.Debugf(nil, "%+s", "Set globalcfg etcd value.")
	return g
}

func (g *GlobalCfg) setLoggerLevel() *GlobalCfg {
	AppLogMode := config.GetInstance().Log.Level
	logger.SetLevelByString(AppLogMode)
	logger.Debugf(nil, "Set app log level to %+s", AppLogMode)
	return g
}

func (g *GlobalCfg) GetEtcd() *etcd.Etcd {
	return g.etcd
}

func (g *GlobalCfg) GetDB() *gorm.DB {
	return g.database
}
