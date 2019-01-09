package task

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	nfconstants "openpitrix.io/notification/pkg/constants"
	"openpitrix.io/notification/pkg/models"
	"openpitrix.io/notification/pkg/util/emailutil"
	"openpitrix.io/notification/pkg/util/etcdutil"
	"strings"
	"time"
)

// Service interface describes all functions that must be implemented.
type Service interface {
	ExtractTasks() error
	HandleTask(handlerNum string) error
	GetTaskwithNfContentbyID(taskID string) (*models.TaskWithNfInfo, error)
	UpdateStatus2SendingByIds(taskWithNfInfo models.TaskWithNfInfo) (bool, error)
}

//Contains all of the logic for the User model.
type taskService struct {
	db             *gorm.DB
	queue          *etcdutil.Queue
	runningTaskIds chan string
	nfIdLast       string
}

func NewService(db *gorm.DB, queue *etcdutil.Queue) Service {
	tasksc := &taskService{db: db, queue: queue}
	MaxTasks := config.GetInstance().App.Maxtasks
	tasksc.runningTaskIds = make(chan string, MaxTasks)
	return tasksc
}

func (sc *taskService) ExtractTasks() error {
	for {
		nfTaskIdsStr, err := sc.queue.Dequeue()
		//taskId := time.Now().Format("2006-01-02 15:04:05")
		//time.Sleep(1 * time.Second)
		if err != nil {
			logger.Errorf(nil, "Failed to dequeue job from etcd queue: %+v", err)
			time.Sleep(3 * time.Second)
			continue
		}

		logger.Infof(nil, "%+v", "Dequeue from etcd queue success,  "+nfTaskIdsStr)
		sc.runningTaskIds <- nfTaskIdsStr
	}
	return nil
}

func (sc *taskService) HandleTask(handlerNum string) error {
	sc.nfIdLast = ""
	for {
		nfTaskIdsStr := <-sc.runningTaskIds
		logger.Debugf(nil, time.Now().Format("2006-01-02 15:04:05")+" handlerNum:"+handlerNum+"  Receive:", nfTaskIdsStr)

		ids := strings.Split(nfTaskIdsStr, ",")
		taskId := ids[1]
		logger.Debugf(nil, "test=======handlerNums%d", handlerNum)
		logger.Debugf(nil, "test=======taskId=s%", taskId)
		//	nfId := ids[1]
		taskWithNfInfo, err := sc.GetTaskwithNfContentbyID(taskId)
		if err != nil {
			logger.Errorf(nil, "got TaskwithNfContentbyID failed, [%+v]", err)
			return err
		}
		logger.Debugf(nil, "got TaskwithNfContentbyID successed, : [%+v]", taskWithNfInfo)

		emailAddr := taskWithNfInfo.EmailAddr
		titel := taskWithNfInfo.Title
		content := taskWithNfInfo.Content
		err = emailutil.SendMail(emailAddr, titel, content)
		if err != nil {
			logger.Warnf(nil, "send email failed, [%+v]", err)
			return err
		}
		//if send successfully,need to update notification, job and task status.
		_, err = sc.UpdateStatus2FinishedByIds(*taskWithNfInfo)
		if err != nil {
			logger.Errorf(nil, "update job and task status  to finished failed, [%+v]", err)
			return err
		}
		logger.Debugf(nil, "update job and task status to finished: [%+v]", taskWithNfInfo)

		//if the nfId is different from nfIdLast,that means the nf including all the tasks is finished.
		//update notification status to finished
		//if sc.nfIdLast != nfId && sc.nfIdLast != "" {
		//	nfsc := nfsc.NewService(sc.db, sc.queue)
		//	_, err = nfsc.UpdateStatus2FinishedById(sc.nfIdLast)
		//	if err != nil {
		//		logger.Errorf(nil, "update notification status to finished failed, [%+v]", err)
		//		return err
		//	}
		//}
		//sc.nfIdLast = nfId

	}
	return nil
}

func (sc *taskService) getTaskbyID(taskID string) (*models.Task, error) {
	task := &models.Task{}
	err := sc.db.
		Where("task_id = ?", taskID).
		First(task).Error
	if err != nil {
		//if err != gorm.ErrRecordNotFound {
		//	return nil, err
		//}
		return nil, err
	}
	return task, nil
}

func (sc *taskService) GetTaskwithNfContentbyID(taskID string) (*models.TaskWithNfInfo, error) {
	logger.Debugf(nil, "test========taskID=%s", taskID)
	taskWithNfInfo := &models.TaskWithNfInfo{}
	sql := models.GetTaskwithNfContentbyIDSQL
	sc.db.Raw(sql, taskID).Scan(&taskWithNfInfo)
	logger.Debugf(nil, "getTaskwithNfContentbyID got a task: [%+v]", taskWithNfInfo)
	return taskWithNfInfo, nil
}

func (sc *taskService) UpdateStatus2SendingByIds(taskWithNfInfo models.TaskWithNfInfo) (bool, error) {
	jobId := taskWithNfInfo.JobID
	taskId := taskWithNfInfo.TaskID
	nfId := taskWithNfInfo.NotificationId

	job := &models.Job{
		JobID: jobId,
	}
	task := &models.Task{
		TaskID: taskId,
	}
	nf := &models.Notification{
		NotificationId: nfId,
	}

	tx := sc.db.Begin()
	status := nfconstants.StatusSending
	err := sc.db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Error
	err = sc.db.Model(&job).Where("job_id = ?", jobId).Update("status", status).Error
	err = sc.db.Model(&nf).Where("notification_id = ?", nfId).Update("status", status).Error
	if err != nil {
		logger.Errorf(nil, "%+v", err)
		return false, err
	}
	tx.Commit()

	return true, nil
}

func (sc *taskService) UpdateStatus2FinishedByIds(taskWithNfInfo models.TaskWithNfInfo) (bool, error) {

	jobId := taskWithNfInfo.JobID
	taskId := taskWithNfInfo.TaskID

	job := &models.Job{
		JobID: jobId,
	}
	task := &models.Task{
		TaskID: taskId,
	}

	tx := sc.db.Begin()
	status := nfconstants.StatusFinished
	err := sc.db.Model(&task).Where("task_id = ?", taskId).Update("status", status).Error
	err = sc.db.Model(&job).Where("job_id = ?", jobId).Update("status", status).Error
	if err != nil {
		logger.Errorf(nil, "%+v", err)
		return false, err
	}
	tx.Commit()

	return true, nil
}
