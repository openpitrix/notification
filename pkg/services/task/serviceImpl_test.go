package task

import (
	"log"
	"openpitrix.io/notification/pkg/services/test"
	"testing"
)

func TestNewService(t *testing.T) {
	log.Println("Test NewServices")
	db := test.GetTestDB()

	q := test.GetEtcdQueue()
	taskservice := NewService(db, q)

	go taskservice.ExtractTasks()
	//go taskservice.HandleTask("A")
	//
	//for{
	//	//println("...")
	//	time.Sleep(2 * time.Second)
	//}
}

func TestGetTaskbyID(t *testing.T) {
	log.Println("Test TestgetTaskfromDBbyID")
	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	tasksc := &taskService{db: db, queue: q}
	task, _ := tasksc.getTaskbyID("task-LBx4k82RMZOo")
	log.Println(task.AddrsStr)
}


func TestGetTaskwithNfContentbyID(t *testing.T) {
	log.Println("Test TestGetTaskwithNfContentbyID")
	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	tasksc := &taskService{db: db, queue: q}
	task, _ := tasksc.getTaskwithNfContentbyID("task-LBx4k82RMZOo")
	log.Println(task.AddrsStr)
	log.Println(task)
}
