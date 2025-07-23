package api

import (
	"api-connection-service/internal/config"
	"api-connection-service/internal/kafka"
	tlsCred "api-connection-service/internal/tls"
	"context"
	"fmt"
	pb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/doc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"time"
)

type tokenAuth struct {
	token string
}

func (t *tokenAuth) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": "Bearer " + t.token,
	}, nil
}

func (t *tokenAuth) RequireTransportSecurity() bool {
	return true
}

type Client struct {
	conn     *grpc.ClientConn
	logger   *slog.Logger
	stream   *Stream
	producer *kafka.Producer
}

func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	creds, err := tlsCred.LoadTLSCredentials(cfg.TLS, logger)

	conn, err := grpc.DialContext(
		ctx,
		cfg.APIHost,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&tokenAuth{token: cfg.APIToken}),
		grpc.WithBlock(),
	)

	if err != nil {
		logger.Error("Failed to connect gRPC server: %v", err)
		return nil, err
	}

	state := conn.GetState()
	logger.Info("gRPC client connection state", "state", state.String())

	stream, err := NewStream(conn, logger)
	if err != nil {
		logger.Error("Failed to create stream", "error", err)
		return nil, err
	}

	producer, err := kafka.NewProducer(cfg.KafkaBroker, logger)

	return &Client{
		conn:     conn,
		logger:   logger,
		stream:   stream,
		producer: producer,
	}, nil
}

func (c *Client) GetInstrumentUIDFromTickerShare(ctx context.Context, ticker string) (string, error) {
	classCode := "TQBR"

	instrumentClient := pb.NewInstrumentsServiceClient(c.conn)
	req := &pb.InstrumentRequest{
		IdType:    pb.InstrumentIdType_INSTRUMENT_ID_TYPE_TICKER,
		ClassCode: &classCode,
		Id:        ticker,
	}

	resp, err := instrumentClient.GetInstrumentBy(ctx, req)
	if err != nil {
		c.logger.Error("Failed to get instrument by ticker", "error", err)
		return "", err
	}

	instrument := resp.GetInstrument()
	if instrument == nil {
		c.logger.Error("Failed to get instrument_uid by ticker", "error", err)
		return "", err
	}

	return instrument.GetUid(), nil
}

func (c *Client) StartStream(ctx context.Context, instrumentID string) (chan *pb.Candle, chan error, error) {
	if c.stream == nil {
		return nil, nil, fmt.Errorf("stream not initialized")
	}

	dataChan, errChan := c.stream.StartStream(ctx, instrumentID)

	go c.Serialization(ctx, dataChan)

	return dataChan, errChan, nil
}

func (c *Client) Serialization(ctx context.Context, dataChan chan *pb.Candle) {
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Serialization cancelled")
			return

		case candle, ok := <-dataChan:
			if !ok {
				c.logger.Info("Data channel closed")
				return
			}

			data, err := proto.Marshal(candle)
			if err != nil {
				c.logger.Error("Failed to serialize candle", "error", err)
				continue
			}

			c.producer.Send(ctx, data)
		}
	}
}

func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		c.logger.Error("Failed to close gRPC connection", "error", err)
	}

	if err := c.producer.Close; err != nil {
		c.logger.Error("Failed to close kafka producer", "error", err)
	}
}
