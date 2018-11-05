package task

// Service interface describes all functions that must be implemented.
type Service interface {
	ExtractTasks() (error)
	HandleTask(handlerNum string) (error)
}
