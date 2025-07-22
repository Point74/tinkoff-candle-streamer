package kafka

import (
	"api-connection-service/internal/config"
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"log/slog"
)

type Producer struct {
	client *kgo.Client
	logger *slog.Logger
}

func NewProducer(cfg *config.Config, brokers []string, logger *slog.Logger) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no brokers provided")
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		logger.Error("Error creating Kafka Producer", "error", err, "brokers", brokers)
		return nil, err
	}

	return &Producer{
		client: client,
		logger: logger,
	}, nil
}

func (p *Producer) Send(ctx context.Context, data []byte) {

}

func (p *Producer) Close() {
	if p.client != nil {
		p.client.Close()
	}
}
