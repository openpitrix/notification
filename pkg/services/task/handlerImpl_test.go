package task

import (
	"fmt"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/services/test"
	"strconv"
	"testing"
	"time"
)

func TestNewHandler(t *testing.T) {
	test.InitGlobelSetting()
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


func TestDescribeNfs(t *testing.T){
	test.InitGlobelSetting()

	MaxWorkingTasks:=config.GetInstance().App.MaxWorkingTasks
	for a := 0; a < MaxWorkingTasks; a++ {
		fmt.Printf("a 的值为: %d\n", a)
		ss := strconv.Itoa(a)
		logger.Infof(nil,ss)
	}




}