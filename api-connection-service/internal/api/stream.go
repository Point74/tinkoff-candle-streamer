package api

import (
	"fmt"
	pb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/doc"
	"google.golang.org/grpc"
	"log/slog"
)

type Stream struct {
	client pb.MarketDataServiceClient
	logger *slog.Logger
	conn   *grpc.ClientConn
}

func NewStream(conn *grpc.ClientConn, logger *slog.Logger) (*Stream, error) {
	client := pb.NewMarketDataServiceClient(conn)
	if client == nil {
		logger.Error("cannot connect to MarketDataService")
		return nil, fmt.Errorf("cannot connect to MarketDataService")
	}

	return &Stream{
		client: client,
		logger: logger,
		conn:   conn,
	}, nil
}
