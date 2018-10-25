package main

import (
	"log"
	"openpitrix.io/notification/pkg/services"
)

func main() {

	var (
		err error
		s   *services.Server
	)

	log.Println("Starting server...")

	s, _ = services.NewServer()
	log.Println("Serveing ...")

	err = s.Serve()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Server shuting down...")

}
