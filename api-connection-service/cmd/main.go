package main

import (
	"api-connection-service/internal/api"
	"api-connection-service/internal/logger"
	"context"
	"github.com/Point74/tinkoff-candle-streamer/config"
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

func main() {
	log := logger.NewLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Error loading config", err)
		os.Exit(1)
	}

	go func() {
		log.Info("Starting metrics server api-connection-service on :9090")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Error("Error starting metrics server", err)
		}
	}()

	client, err := api.NewClient(cfg, log)
	if err != nil {
		log.Error("Error creating client", "error", err)
		os.Exit(1)
	}

	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Info("API Connection Service started!")

	shareTicker := "T"

	shareUid, err := client.GetInstrumentUIDFromTickerShare(ctx, shareTicker)
	if err != nil {
		os.Exit(1)
	}

	_, errChan, err := client.StartStream(ctx, shareUid, shareTicker)
	if err != nil {
		log.Error("Error starting stream", "error", err)
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

		case <-ctx.Done():
			log.Info("Context done")
		}
	}
}
