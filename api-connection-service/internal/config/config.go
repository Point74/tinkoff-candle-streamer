package config

type Config struct {
	APIHost     string
	APIToken    string
	TLS         string
	KafkaBroker string
}

func LoadConfig() (*Config, error) {
	return &Config{
		APIHost:     "invest-public-api.tinkoff.ru:443",
		APIToken:    "your-api-token",
		TLS:         "certificate/tinkoff.pem",
		KafkaBroker: "localhost:9092",
	}, nil
}
