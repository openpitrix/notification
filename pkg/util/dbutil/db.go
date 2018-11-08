package dbutil

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"openpitrix.io/notification/pkg/config"
	"sync"
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
var err_db error

func GetInstance() *MysqlConnPool {
	once.Do(func() {
		instance = &MysqlConnPool{}
	})
	return instance
}



/*
* @fuc 初始化数据库连接(可在mail()适当位置调用)
*/
func (m *MysqlConnPool) InitDataPool() (issucc bool) {
	//db, err_db = gorm.Open("mysql", "root:password@tcp(192.168.0.10:13306)/notification?charset=utf8&parseTime=True&loc=Local")
	//fmt.Println(err_db)
	cfg:=config.GetInstance()

	var (
		dbCfg            = cfg.Db
		connectionString = fmt.Sprintf(
			"%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.DatabaseName,
		)
	)

	db, err_db = gorm.Open("mysql", connectionString)
	if err_db != nil {
		log.Print(err_db)
		return false
	}

	err_db = db.DB().Ping()

	if err_db != nil {
		return false
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(cfg.Db.DBLogMode)

	// 全局禁用表名复数
	db.SingularTable(true)


	if err_db != nil {
		log.Fatal(err_db)
		return false
	}
	//关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true
}

/*
* 对外获取数据库连接对象db
*/
func (m *MysqlConnPool) GetMysqlDB() (*gorm.DB) {
	return db
}

