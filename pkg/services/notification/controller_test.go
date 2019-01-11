package notification

import (
	"openpitrix.io/notification/pkg/services/notification/service/notification"
	"openpitrix.io/notification/pkg/services/notification/service/task"
	"testing"
	"time"
)

func TestServe(t *testing.T) {
	nfservice := notification.NewService()
	taskservice := task.NewService()

	c := NewController(nfservice, taskservice)

	go c.ExtractTasks()
	go c.HandleTask("A")
	go c.HandleTask("B")
	//
	for {
		//println("...")
		time.Sleep(2 * time.Second)
	}
}
