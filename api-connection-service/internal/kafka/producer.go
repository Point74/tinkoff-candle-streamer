package kafka

import (
	"context"
	"fmt"
	"github.com/twmb/franz-go/pkg/kgo"
	"log/slog"
	"time"
)

type Producer struct {
	client *kgo.Client
	logger *slog.Logger
	topic  string
}

func NewProducer(brokers string, logger *slog.Logger) (*Producer, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no brokers provided")
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers),
		kgo.DefaultProduceTopic("tinkoff-candle"),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		logger.Error("Error creating Kafka Producer", "error", err, "brokers", brokers)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		logger.Error("Error connect to Kafka Producer", "error", err, "brokers", brokers)
		client.Close()

		return nil, err
	}

	return &Producer{
		client: client,
		logger: logger,
		topic:  "tinkoff-candle",
	}, nil
}

func (p *Producer) Send(ctx context.Context, dataChan chan []byte) {
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Producer is shutting down due to context", "data", ctx)
			return
		case data, ok := <-dataChan:
			if !ok {
				p.logger.Info("Producer data chanel is closed", "data", ok)
				return
			}

			record := &kgo.Record{
				Value: data,
				Topic: p.topic,
			}

			p.client.Produce(ctx, record, func(_ *kgo.Record, err error) {
				if err != nil {
					fmt.Printf("record had a produce error: %v\n", err)
				}
				if err == nil {
					fmt.Println("record sent")
				}
			})
		}
	}
}

func (p *Producer) Close() {
	if p.client != nil {
		p.client.Close()
	}
}
