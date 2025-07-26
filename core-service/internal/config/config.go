package config

type Config struct {
	KafkaBroker string
}

func LoadConfig() (*Config, error) {
	return &Config{
		KafkaBroker: "localhost:9092",
	}, nil
}
