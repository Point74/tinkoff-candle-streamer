package api

import (
	"api-connection-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(cfg *config.Config) (*Client, error) {
	creds := credentials.NewClientTLSFromCert(nil, "")

	conn, err := grpc.NewClient(cfg.APIHost, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Printf("Failed to connect gRPC server: %v", err)
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}

}
