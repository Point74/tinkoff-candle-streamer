package api

import (
	"api-connection-service/internal/config"
	tlsCred "api-connection-service/internal/tls"
	"context"
	"fmt"
	"google.golang.org/grpc"
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
	conn   *grpc.ClientConn
	logger *slog.Logger
	stream *Stream
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

	return &Client{
		conn:   conn,
		logger: logger,
		stream: stream,
	}, nil
}

func (c *Client) StartStream(instrumentID string) error {
	if c.stream == nil {
		return fmt.Errorf("stream not initialized")
	}

	return c.stream.StartStream(instrumentID)
}

func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		c.logger.Error("Failed to close gRPC connection: %v", err)
	}
}
