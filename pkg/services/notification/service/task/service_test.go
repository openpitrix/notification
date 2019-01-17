// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"testing"

	"openpitrix.io/logger"
)

//func TestNewService(t *testing.T) {
//
//	//taskservice := NewService()
//
//	//go taskservice.ExtractTasks()
//	//go taskservice.HandleTask("A")
//	//go taskservice.HandleTask("B")
//	////
//	//for {
//	//	//println("...")
//	//	time.Sleep(2 * time.Second)
//	//}
//}

func TestGetTaskByID(t *testing.T) {
	taskService := &taskService{}
	task, _ := taskService.getTaskbyID("task-6J1BEDx9wJ94")
	logger.Infof(nil, task.EmailAddr)
}

func TestGetTaskWithNfContentByID(t *testing.T) {
	taskService := &taskService{}
	task, _ := taskService.GetTaskWithNfContentByID("task-mqY0kxG9yl98")
	logger.Infof(nil, "EmailAddr=[%+s]", task.EmailAddr)
}

func TestUpdateStatusById(t *testing.T) {
	tasksc := &taskService{}
	task, _ := tasksc.GetTaskWithNfContentByID("task-QvQEG9n5BkZO")
	tasksc.UpdateStatus2SendingByIds(*task)
}
