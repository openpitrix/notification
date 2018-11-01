package task

import (
	"log"
	"openpitrix.io/notification/pkg/services/test"
	"testing"
)

func TestNewService(t *testing.T) {
	log.Println("Test NewServices")
	db := test.GetTestDB()
	q := test.GetEtcdQueue()
	taskservice := NewService(db, q)
	taskservice.ExtractTasks()
}
