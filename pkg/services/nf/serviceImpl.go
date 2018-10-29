package nf

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"openpitrix.io/notification/pkg/models"
)

//Contains all of the logic for the User model.
type nfService struct {
	db           *gorm.DB
}

func (nfs *nfService) CreateNfWaddrs(nf *models.NotificationCenterPost) error {
	err := nfs.db.Create(&nf).Error
	if err != nil {
		return err
	}
	return nil
}



func (nfs *nfService) GetDataFromDB4Test() {
		db :=nfs.db
		// 读取
		var product models.Product
		db.First(&product, 1) // 查询id为1的product
		//db.First(&product, "code = ?", "L1212") // 查询code为l1212的product
		fmt.Println(product)
}

func NewService(db *gorm.DB) Service {
	return &nfService{db: db}
}


func (nfs  *nfService) SayHello(str string) (string, error) {
	log.Println("Step 7: deep func"+str)
	nfs.GetDataFromDB4Test()
	return str,nil
}


//func (nfs *nfService) CreateNfWaddrs(nf *models.NotificationCenterPost) error {
//	err := nfs.db.Create(&nf).Error
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//