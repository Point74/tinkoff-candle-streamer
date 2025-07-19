package main

import (
	"api-connection-service/internal/api"
	"api-connection-service/internal/config"
	"api-connection-service/internal/logger"
	"os"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config", err)
		os.Exit(1)
	}

	client, err := api.NewClient(cfg, log)
	if err != nil {
		log.Error("Error creating client: ", err)
		os.Exit(1)
	}

	defer client.Close()

	log.Info("API Connection Service started!")

	if err = client.StartStream("BBG004730N88"); err != nil {
		log.Error("Error starting stream: ", err)
	}
}
