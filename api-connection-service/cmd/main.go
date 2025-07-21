package main

import (
	"api-connection-service/internal/api"
	"api-connection-service/internal/config"
	"api-connection-service/internal/logger"
	"context"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("API Connection Service started!")

	dataChan, errChan, err := client.StartStream(ctx, "BBG004730N88")
	if err != nil {
		log.Error("Error starting stream: ", err)
		os.Exit(1)
	}

	for {
		select {
		case err, ok := <-errChan:
			if !ok {
				log.Info("Error channel closed")
				return
			}

			log.Info("Stream error: %v", err)

		case data, ok := <-dataChan:
			if !ok {
				log.Info("Data channel closed")
				return
			}

			log.Info("Received candle data", "data", data)

		case <-ctx.Done():
			log.Info("Context done")
		}
	}
}
