package main

import (
	"api-connection-service/internal/api"
	"api-connection-service/internal/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		log.Fatal("Error creating client: ", err)
	}

	defer client.Close()

	log.Printf("API Connection Service started!\n")
}
