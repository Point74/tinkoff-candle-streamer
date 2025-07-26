package main

import (
	"core-service/internal/config"
	"core-service/internal/kafka"
	"core-service/internal/logger"
	"os"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config", err)
		os.Exit(1)
	}

	consumer, err := kafka.NewConsumer(cfg.KafkaBroker, log)
	if err != nil {
		log.Error("Error creating consumer", "error", err)
		os.Exit(1)
	}

	defer consumer.Close()
}
