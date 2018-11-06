package task

type Handler interface {
	ExtractTasks() error
	HandleTask(handlerNum string) error

	ServeTask() error
}
