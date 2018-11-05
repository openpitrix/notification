package task

import (
	"fmt"
	"log"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/emailutil"

	//	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"openpitrix.io/notification/pkg/util/etcdutil"

	//	"openpitrix.io/openpitrix/pkg/models"
	//	"openpitrix.io/openpitrix/pkg/pi"
	//	"openpitrix.io/openpitrix/pkg/plugins"
	//	"openpitrix.io/openpitrix/pkg/util/ctxutil"
	"time"
)

//Contains all of the logic for the User model.
type taskService struct {
	db    *gorm.DB
	queue *etcdutil.Queue
	runningTaskIds chan string
}

func NewService(db *gorm.DB, queue *etcdutil.Queue) Service {
	tasksc:=&taskService{db: db, queue: queue}
	tasksc.runningTaskIds=make(chan string, 10)
	return tasksc
}

func (sc *taskService) ExtractTasks() (error) {
	for {
		//taskId, err := sc.queue.Dequeue()
		taskId := time.Now().Format("2006-01-02 15:04:05")
		time.Sleep(1 * time.Second)
		//if err != nil {
		//	logger.Error(nil, "Failed to dequeue job from etcd queue: %+v", err)
		//	time.Sleep(3 * time.Second)
		//	continue
		//}
		//logger.Debug(nil, "Dequeue job [%s] from etcd queue success", taskId)
		fmt.Println("Dequeue from etcd queue success  "+taskId)
		sc.runningTaskIds <- taskId
	}
	return nil
}

func (sc *taskService) HandleTask(handlerNum string) (error) {
	for   {
		taskId := <-sc.runningTaskIds
		fmt.Println(time.Now().Format("2006-01-02 15:04:05")+" handlerNum:"+handlerNum+"  Receive:", taskId)

		taskWNfInfo, err :=sc.getTaskwithNfContentbyID(taskId)
		if err != nil {
			log.Println("something wrong")
		}
		addrsStr:=taskWNfInfo.AddrsStr
		log.Println(addrsStr)
		titel:=taskWNfInfo.Title
		content:=taskWNfInfo.Content
		emailutil.SendMail(addrsStr,titel,content)

	}
	return nil
}

func (sc *taskService) getTaskbyID(taskID string) ( *models.Task,error) {
	task := &models.Task{}
	err := sc.db.
		Where("task_id = ?", taskID).
		First(task).Error
	log.Println(task)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		return nil, nil
	}
	return task, nil
}


func (sc *taskService) getTaskwithNfContentbyID(taskID string) (*models.TaskWNfInfo,error) {
	//taskID="task-LBx4k82RMZOo"
	task := &models.TaskWNfInfo{}

	sc.db.Raw("SELECT  t3.title,t3.short_content,  t3.content,t1.task_id,t1.addrs_str "+
		"	FROM task t1,job t2,notification_center_post t3 where t1.job_id=t2.job_id and t2.nf_post_id=t3.nf_post_id  and t1.task_id=? ",taskID).Scan(&task)

	return task, nil
}
