package main

import (
	"context"
	"core-service/internal/db/postgres"
	"core-service/internal/kafka"
	"core-service/internal/logger"
	"github.com/Point74/tinkoff-candle-streamer/config"
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
		log.Info("Starting metrics server core-service on :9090")
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Error("Error starting metrics server", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer, err := kafka.NewConsumer(cfg, log)
	if err != nil {
		log.Error("Error creating consumer", "error", err)
		os.Exit(1)
	}

	defer consumer.Close()

	storage, err := postgres.New(ctx, cfg, log)
	if err != nil {
		os.Exit(1)
	}

	defer storage.Close(ctx)

	consumer.Get(ctx)
}
