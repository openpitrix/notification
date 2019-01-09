package main

import (
	"log"
	"openpitrix.io/notification/pkg/config"
	"openpitrix.io/notification/pkg/services"
)

func main() {
	log.Println("Starting server...")

	cfg := config.GetInstance()
	cfg.LoadConf()

	services.InitGlobelSetting()

	services.Serve(cfg)

	log.Println("Server shuting down...")

}
