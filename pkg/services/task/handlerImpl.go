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

func (h *handler) HandleTask(handlerNum string) (error) {
	h.tasksc.HandleTask(handlerNum)
	return nil
}


func (h *handler) ServeTask() (error) {
	log.Println("Call handlerImpl.ServeTask")
	go h.ExtractTasks()


	for i := 0; i < 10; i++ {
	go h.HandleTask("A")
	}

	return nil
}
