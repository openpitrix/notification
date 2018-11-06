package task

import (
	"log"
	"openpitrix.io/notification/pkg/config"
)

type handler struct {
	tasksc Service
}

func NewHandler(tasksc Service) Handler {
	return &handler{
		tasksc: tasksc,
	}
}

func (h *handler) ExtractTasks() error {
	h.tasksc.ExtractTasks()
	return nil
}

func (h *handler) HandleTask(handlerNum string) error {
	h.tasksc.HandleTask(handlerNum)
	return nil
}

func (h *handler) ServeTask() error {
	log.Println("Call handlerImpl.ServeTask")
	go h.ExtractTasks()

	MaxWorkingTasks:=config.GetInstance().App.MaxWorkingTasks
	for i := 0; i < MaxWorkingTasks; i++ {
		go h.HandleTask(string(i))
	}
	return nil
}
