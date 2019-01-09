package task

import (
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/services/test"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {

	test.InitGlobelSetting4Test()
	db, q := test.GetTestDBAndEtcd4Test()

	taskservice := NewService(db, q)

	go taskservice.ExtractTasks()
	go taskservice.HandleTask("A")
	go taskservice.HandleTask("B")
	//
	for {
		//println("...")
		time.Sleep(2 * time.Second)
	}
}

func TestGetTaskbyID(t *testing.T) {
	db, q := test.GetTestDBAndEtcd4Test()
	tasksc := &taskService{db: db, queue: q}
	task, _ := tasksc.getTaskbyID("task-QvQEG9n5BkZO")
	logger.Infof(nil, task.EmailAddr)
}

func TestGetTaskwithNfContentbyID(t *testing.T) {
	db, q := test.GetTestDBAndEtcd4Test()
	tasksc := &taskService{db: db, queue: q}
	task, _ := tasksc.GetTaskwithNfContentbyID("task-QvQEG9n5BkZO")
	logger.Infof(nil, task.EmailAddr)
}

func TestUpdateStatusById(t *testing.T) {
	db, q := test.GetTestDBAndEtcd4Test()
	tasksc := &taskService{db: db, queue: q}
	task, _ := tasksc.GetTaskwithNfContentbyID("task-QvQEG9n5BkZO")
	//status := "test_status"
	tasksc.UpdateStatus2SendingByIds(*task)
}
