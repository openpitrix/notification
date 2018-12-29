package task

import (
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/services/test"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {

	test.InitGlobelSetting()
	db := test.GetTestDB()
	q := test.GetEtcdQueue()

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
	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	tasksc := &taskService{db: db, queue: q}
	task, _ := tasksc.getTaskbyID("task-LBx4k82RMZOo")
	logger.Infof(nil, task.EmailAddr)
}

func TestGetTaskwithNfContentbyID(t *testing.T) {
	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	tasksc := &taskService{db: db, queue: q}
	task, _ := tasksc.getTaskwithNfContentbyID("task-LBx4k82RMZOo")
	logger.Infof(nil, task.EmailAddr)
}
