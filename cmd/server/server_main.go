package main

import (
	"log"
	"openpitrix.io/notification/pkg/services"
)

func main() {
	log.Println("Starting server...")

	err := services.Serve()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Server shuting down...")

}
