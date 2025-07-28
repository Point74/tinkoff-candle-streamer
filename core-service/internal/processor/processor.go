package processor

import (
	"context"
	ownpb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/my"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type Processor struct {
	logger *slog.Logger
}

type Candle struct {
	ticker string
	open   float64
	high   float64
	low    float64
	close  float64
	time   timestamppb.Timestamp
}

func NewProcessor(logger *slog.Logger) *Processor {
	return &Processor{
		logger: logger,
	}
}

func (p *Processor) Deserialization(ctx context.Context, serDataCandleChan chan []byte) {
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Data processing is shutting down due to context", "data", ctx)
			return
		case serData := <-serDataCandleChan:
			var candle ownpb.CandleData
			err := proto.Unmarshal(serData, &candle)
			if err != nil {
				p.logger.Error("Error unmarshaling data", "error", err)
				return
			}

			p.logger.Info("Deserialized data", "candle", &candle)
		}
	}
}
