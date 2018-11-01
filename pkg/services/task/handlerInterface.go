package task

type Handler interface {
	ExtractTasks() (error)
	HandleTasks() (error)

	ServeTask() (error)
}
