package task

import (
	"openpitrix.io/notification/pkg/services/test"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	db := test.GetTestDB()
	q := test.GetEtcdQueue()

	taskservice := NewService(db, q)
	handler := NewHandler(taskservice)

//	go handler.ExtractTasks()
//	go handler.HandleTask("1")
	go handler.ServeTask()

	for{
		//println("...")
		time.Sleep(2 * time.Second)
	}
}
