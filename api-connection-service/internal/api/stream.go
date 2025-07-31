package api

import (
	"context"
	"fmt"
	pb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/doc"
	"google.golang.org/grpc"
	"log/slog"
	"time"
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

	logger.Info("connect to MarketDataService")

	return &Stream{
		client: client,
		logger: logger,
		conn:   conn,
	}, nil
}

func (s *Stream) StartStream(ctx context.Context, instrumentID string) (chan *pb.Candle, chan error) {
	dataChan := make(chan *pb.Candle, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(dataChan)
		defer close(errChan)

		retry := 5 * time.Second
		maxRetry := 60 * time.Second
		retryMultiplier := 2.0

		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Stream stopped")
				return

			default:
				stream, err := s.client.MarketDataStream(ctx)

				if err != nil {
					s.logger.Error("cannot connect to MarketDataService", "error", err)
					errChan <- fmt.Errorf("cannot connect to MarketDataService")

					time.Sleep(retry)
					retry = time.Duration(float64(retry) * retryMultiplier)

					if retry > maxRetry {
						retry = maxRetry
					}

					continue
				}

				candleSource := pb.GetCandlesRequest_CANDLE_SOURCE_UNSPECIFIED
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
							CandleSourceType: &candleSource,
							WaitingClose:     true,
						},
					},
				}

				if err := stream.Send(request); err != nil {
					s.logger.Error("cannot send request to MarketDataService", "error", err)
					errChan <- fmt.Errorf("cannot send request to MarketDataService")

					time.Sleep(retry)
					retry = time.Duration(float64(retry) * retryMultiplier)

					if retry > maxRetry {
						retry = maxRetry
					}

					continue
				}

				s.logger.Info("starting stream", "instrumentID", instrumentID)

				retry = 5 * time.Second

				for {
					resp, err := stream.Recv()

					if err != nil {
						s.logger.Error("Error receiving stream data", "error", err)
						errChan <- fmt.Errorf("error receiving stream data")
						break
					}

					if candle := resp.GetCandle(); candle != nil {
						dataChan <- candle
						s.logger.Info("received candle data", "candle", candle)
					}
				}
			}
		}
	}()

	return dataChan, errChan
}
