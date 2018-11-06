package task

import (
	"log"
	"openpitrix.io/notification/pkg/services/test"
	"testing"
)

func TestNewHandler(t *testing.T) {
	log.Println("Test func NewHandler")

	db := test.GetTestDB()
	q := test.GetEtcdQueue()

	taskservice := NewService(db, q)
	handler := NewHandler(taskservice)

	handler.ExtractTasks()
}
