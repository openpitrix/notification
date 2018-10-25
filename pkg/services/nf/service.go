package nf

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/models"
)

//Contains all of the logic for the User model.
type nfservice struct {
	cfg         *config.Config
	db           *gorm.DB
}


// NewService initialization. db *gorm.DB
func NewServices() (Service, error) {
	var (
		err    error
		nfs=&nfservice{}
	)

	//log.Println("Configuring server Start...")
	nfs.cfg=config.NewConfig()
	//log.Println("Configuring server End...")
	log.Println("NewServices:call NewServices, set nfs.cfg")

	//log.Println("CreateDatabaseConnection Start...")
	nfs.db, err =nfs.createDatabaseConnection()
	//log.Println("CreateDatabaseConnection End...")
	log.Println("NewServices:call NewServices, nfs.db")

	if err != nil {
		return nil, err
	}
	return nfs, nil
}

// createDatabaseConn creates a new GORM database with the specified database configuration.
func (nfs *nfservice) createDatabaseConnection() (*gorm.DB, error) {

	var (
		db               *gorm.DB
		err              error
		dbCfg            = nfs.cfg.Db
		connectionString = fmt.Sprintf(
			"%v:%v@(%v:%v)/%v?charset=utf8&parseTime=True&loc=Local",
			dbCfg.User,
			dbCfg.Password,
			dbCfg.Host,
			dbCfg.Port,
			dbCfg.DatabaseName,
		)
	)

	db, err = gorm.Open("mysql", connectionString)


	if err != nil {
		log.Print(err)
		return nil, err
	}

	err = db.DB().Ping()

	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(10)
	db.LogMode(nfs.cfg.DBLogMode)

	 // 全局禁用表名复数
	db.SingularTable(true)

	return db, nil
}

func (nfs *nfservice) GetDataFromDB4Test() () {
	db, _ :=nfs.createDatabaseConnection()

	// 读取
	var product models.Product
	db.First(&product, 1) // 查询id为1的product
	db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
	fmt.Println(product)
}

func (nfs *nfservice) SayHello(str string) (string, error) {
	log.Println("Test SayHello in service..")
	log.Print(str)
	return "ss",nil
}


func (nfs *nfservice) CreateNfWaddrs(nfPostID string, nfPostType string, title string, content string, shortContent string, exporedDays int64, owner string) (error) {
	log.Println("Test CreateNfWaddrs..")

	return nil
}



func (nfs *nfservice) CreateNfWaddrs2(nf *models.NotificationCenterPost) error {
	err := nfs.db.Create(&nf).Error
	if err != nil {
		return err
	}
	return nil

}

func (nfs *nfservice) CreateNfWaddrs3(nf *models.NotificationCenterPost, job *models.Job, task *models.Task) error {
	panic("implement me")
	log.Print("Test CreateNfWaddrs3")

	tx := nfs.db.Begin()
	// 注意，一旦你在一个事务中，使用tx作为数据库句柄

	if err := tx.Create(&nf).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&job).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Create(&task).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
