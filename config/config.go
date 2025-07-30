package config

type Config struct {
	APIHost     string `env:"API_HOST" envDefault:"invest-public-api.tinkoff.ru:443"`
	APIToken    string `env:"API_TOKEN" envDefault:""`
	KafkaBroker string `env:"KAFKA_BROKER" envDefault:"localhost:9092"`
}
