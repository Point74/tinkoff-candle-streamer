package api

import (
	"context"
	"fmt"
	pb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/doc"
	"google.golang.org/grpc"
	"log/slog"
	"sync"
)

type Stream struct {
	client pb.MarketDataStreamServiceClient
	logger *slog.Logger
	conn   *grpc.ClientConn
}

func NewStream(conn *grpc.ClientConn, logger *slog.Logger) (*Stream, error) {
	client := pb.NewMarketDataStreamServiceClient(conn)
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

func (s *Stream) StartStream(instrumentID string) error {
	s.logger.Info("starting stream", "instrumentID", instrumentID)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	request := &pb.MarketDataRequest{
		Payload: &pb.MarketDataRequest_SubscribeCandlesRequest{
			SubscribeCandlesRequest: &pb.SubscribeCandlesRequest{
				SubscriptionAction: pb.SubscriptionAction_SUBSCRIPTION_ACTION_SUBSCRIBE,
				Instruments: []*pb.CandleInstrument{
					{
						InstrumentId: instrumentID,
						Interval:     pb.SubscriptionInterval_SUBSCRIPTION_INTERVAL_ONE_MINUTE,
					},
				},
			},
		},
	}

	stream, err := s.client.MarketDataStream(ctx)
	if err != nil {
		s.logger.Error("cannot connect to MarketDataService", "error", err)
		return err
	}

	if err := stream.Send(request); err != nil {
		s.logger.Error("cannot send request to MarketDataService", "error", err)
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			resp, err := stream.Recv()
			if err != nil {
				s.logger.Error("Error receiving stream data", "error", err)
				cancel()
				return
			}

			s.logger.Info("received candle data", "data", resp)
		}
	}()

	wg.Wait()

	return nil
}
