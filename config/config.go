package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	APIHost     string `env:"API_HOST" envDefault:"invest-public-api.tinkoff.ru:443"`
	APIToken    string `env:"API_TOKEN" envDefault:""`
	KafkaBroker string `env:"KAFKA_BROKER" envDefault:"localhost:9092"`
	PostgresUrl string `env:"POSTGRES_URL"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load("../docker/.env"); err != nil {
		return nil, err
	}

	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
