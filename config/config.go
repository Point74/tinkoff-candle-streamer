package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	APIHost          string `env:"API_HOST" envDefault:"invest-public-api.tinkoff.ru:443"`
	APIToken         string `env:"API_TOKEN" envDefault:""`
	KafkaBroker      string `env:"KAFKA_BROKER" envDefault:"localhost:9092"`
	PostgresUser     string `env:"POSTGRES_USER" envDefault:"user"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	PostgresDB       string `env:"POSTGRES_DB" envDefault:"postgres"`
	PostgresPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	MigrationsPath   string `env:"MIGRATIONS_PATH" envDefault:"/app/core-service/internal/db"`
}

func LoadConfig() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
