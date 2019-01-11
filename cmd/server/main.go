package main

import (
	"log"
	"openpitrix.io/notification/pkg/services/notification"
)

func main() {
	log.Println("Starting server...")

	notification.Serve()

	log.Println("Server shuting down...")

}
