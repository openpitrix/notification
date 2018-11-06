package task

import (
	"openpitrix.io/notification/pkg/services/test"
	"testing"
)

func TestNewHandler(t *testing.T) {
	db := test.GetTestDB()
	q := test.GetEtcdQueue()

	taskservice := NewService(db, q)
	handler := NewHandler(taskservice)

	handler.ExtractTasks()
}
