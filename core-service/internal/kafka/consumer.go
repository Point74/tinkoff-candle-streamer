package kafka

import (
	"context"
	"core-service/internal/processor"
	"fmt"
	"github.com/Point74/tinkoff-candle-streamer/config"
	"github.com/twmb/franz-go/pkg/kgo"
	"log/slog"
	"time"
)

type Consumer struct {
	client    *kgo.Client
	logger    *slog.Logger
	topic     string
	group     string
	processor processor.Processor
}

func NewConsumer(cfg *config.Config, logger *slog.Logger) (*Consumer, error) {
	if len(cfg.KafkaBroker) == 0 {
		return nil, fmt.Errorf("no brokers provided")
	}

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.KafkaBroker),
		kgo.ConsumeTopics("tinkoff-candle"),
		kgo.ConsumerGroup("consumer-group"),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		logger.Error("Error creating Kafka Consumer", "error", err, "brokers", cfg.KafkaBroker)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		logger.Error("Error connect to Kafka Consumer", "error", err, "brokers", cfg.KafkaBroker)
		client.Close()

		return nil, err
	}

	logger.Info("Connected to Kafka Consumer", "brokers", cfg.KafkaBroker)

	proc := processor.NewProcessor(context.Background(), cfg, logger)

	return &Consumer{
		client:    client,
		logger:    logger,
		topic:     "tinkoff-candle",
		group:     "consumer-group",
		processor: *proc,
	}, nil
}

func (c *Consumer) Get(ctx context.Context) {
	recordChan := make(chan []byte, 100)
	go c.processor.Deserialization(ctx, recordChan)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Consumer is shutting down dou to context", "data", ctx)
			close(recordChan)
			return
		default:
			fetches := c.client.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				c.logger.Error("Error polling fetches", "error", errs)
				continue
			}

			iter := fetches.RecordIter()
			for !iter.Done() {
				record := iter.Next()
				recordChan <- record.Value
			}
		}
	}
}

func (c *Consumer) Close() {
	if c.client != nil {
		c.client.Close()
	}
}
