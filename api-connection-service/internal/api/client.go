package api

import (
	"api-connection-service/internal/config"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log/slog"
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
}

func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	creds := credentials.NewClientTLSFromCert(nil, "")

	conn, err := grpc.NewClient(
		cfg.APIHost,
		grpc.WithTransportCredentials(creds),
	)

	if err != nil {
		logger.Error("Failed to connect gRPC server: %v", err)
		return nil, err
	}

	return &Client{
		conn:   conn,
		logger: logger,
	}, nil
}

func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		c.logger.Error("Failed to close gRPC connection: %v", err)
	}
}
