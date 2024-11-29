package main

import (
	"github.com/asalan316/golang-autoscaler-custom-metrics/app/app_config"
	"github.com/asalan316/golang-autoscaler-custom-metrics/app/server"
	"github.com/cloudfoundry-community/go-cfenv"
	"log"
)

const serverAddress = "0.0.0.0:8080"

func main() {
	var appConfig *cfenv.App
	appConfig, err := app_config.GetAppEnv(cfenv.Current)
	if err != nil {
		log.Fatal("app configuration not found: ", err)
	}
	httpServer := server.NewServer(appConfig)
	err = httpServer.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Printf("server started on %s\n", serverAddress)
}
