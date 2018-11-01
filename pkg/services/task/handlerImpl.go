package task

import "log"

type handler struct {
	tasksc Service
}

func NewHandler(tasksc Service) Handler {
	return &handler{
		tasksc: tasksc,
	}
}

func (h *handler) ExtractTasks() (error) {
	h.tasksc.ExtractTasks()
	return nil
}

func (h *handler) HandleTasks() (error) {
	h.tasksc.HandleTasks()
	return nil
}


func (h *handler) ServeTask() (error) {
	log.Println("Call handlerImpl.ServeTask")
	go h.ExtractTasks()
	go h.HandleTasks()

	return nil
}
