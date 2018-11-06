package nf

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/etcdutil"
)

type nfService struct {
	db    *gorm.DB
	queue *etcdutil.Queue
}

func NewService(db *gorm.DB, q *etcdutil.Queue) Service {
	return &nfService{db: db, queue: q}
}

func (sc *nfService) SayHello(str string) (string, error) {
	logger.Debugf(nil,"Step 7: deep func" + str)
	sc.GetDataFromDB4Test()
	return str, nil
}

func (sc *nfService) CreateNfWaddrs(nf *models.NotificationCenterPost) error {
	var err error
	var job *models.Job

	tx := sc.db.Begin()

	if err = tx.Create(&nf).Error; err != nil {
		tx.Rollback()
		logger.Criticalf(nil, "Cannot insert data to db:%+v", err)
		return err
	}

	parser := &NfHandlerModelParser{}
	job, err = parser.GenJobfromNf(nf)
	if err := tx.Create(&job).Error; err != nil {
		tx.Rollback()
		logger.Criticalf(nil, "Cannot insert data to db:%+v", err)
		return err
	}

	tasks, err := parser.GenTasksfromJob(job)

	for _, task := range tasks {
		if err := tx.Create(&task).Error; err != nil {
			tx.Rollback()
			logger.Criticalf(nil, "Cannot insert data to db:%+v", err)
			return err
		}
		err = sc.queue.Enqueue(task.TaskID)
	}

	tx.Commit()
	return nil
}

func (sc *nfService) GetDataFromDB4Test() {
	type Product struct {
		gorm.Model
		Code  string
		Price uint
	}
	db := sc.db
	// 读取
	var product Product
	db.First(&product, 1) // 查询id为1的product
	//db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
	logger.SetLevelByString("debug")
	logger.Debugf(nil, "%+v", product)
	logger.Infof(nil, "%+v", product)
}

func (sc *nfService) DescribeNfs(nfID string) (*models.NotificationCenterPost, error) {
	nf := &models.NotificationCenterPost{}

	err := sc.db.
		Where("nf_post_id = ?", nfID).
		First(nf).Error

	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}

		return nil, nil
	}

	return nf, nil
}
