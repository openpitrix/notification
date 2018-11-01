package task

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"openpitrix.io/notification/pkg/util/etcdutil"
	"time"
)

//Contains all of the logic for the User model.
type taskService struct {
	db    *gorm.DB
	queue *etcdutil.Queue
}

func NewService(db *gorm.DB, queue *etcdutil.Queue) Service {
	return &taskService{db: db, queue: queue}
}

func (sc *taskService) ExtractTasks() (error) {
	for true {
		fmt.Printf("从ETCD 队列上提取tasks\n");
		time.Sleep(3 * time.Second)
		continue
	}
	return nil

}

func (sc *taskService) HandleTasks() (error) {
	for true {
		fmt.Printf("处理取到的tasks\n");
		time.Sleep(3 * time.Second)
		continue
	}
	return nil
}
