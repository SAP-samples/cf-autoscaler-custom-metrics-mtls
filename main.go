package main

import (
	"github.com/asalan316/golang-autoscaler-custom-metrics/app/app_config"
	"github.com/asalan316/golang-autoscaler-custom-metrics/app/server"
	"github.com/cloudfoundry-community/go-cfenv"
	"log"
)

const serverAddress = "0.0.0.0:8080"

func main() {
	appConfig, err := app_config.GetAppEnv(cfenv.Current)
	if err != nil {
		log.Fatal("app configuration not found %w", err)
	}
	server := server.NewServer(appConfig)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

//TODO
// use this for better structuring the code
//https://dev.to/techschoolguru/implement-restful-http-api-in-go-using-gin-4ap1
