package processor

import (
	"context"
	"core-service/internal/db"
	"core-service/internal/db/postgres"
	"github.com/Point74/tinkoff-candle-streamer/config"
	ownpb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/my"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"os"
)

type Processor struct {
	logger  *slog.Logger
	storage *postgres.Storage
}

func NewProcessor(ctx context.Context, cfg *config.Config, logger *slog.Logger) *Processor {
	storage, err := postgres.New(ctx, cfg, logger)
	if err != nil {
		os.Exit(1)
	}

	return &Processor{
		logger:  logger,
		storage: storage,
	}
}

func (p *Processor) Deserialization(ctx context.Context, serDataCandleChan chan []byte) {
	saveChan := make(chan *ownpb.CandleData, 100)
	var candle ownpb.CandleData
	go p.SaveToDB(ctx, saveChan)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Data processing is shutting down due to context", "data", ctx)
			close(saveChan)
			return

		case serData := <-serDataCandleChan:

			err := proto.Unmarshal(serData, &candle)
			if err != nil {
				p.logger.Error("Error unmarshaling data", "error", err)
				continue
			}

			p.logger.Info("Deserialized data", "candle", &candle)
			saveChan <- &candle
		}
	}
}

func (p *Processor) SaveToDB(ctx context.Context, saveChan chan *ownpb.CandleData) {
	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Data processing is shutting down due to context")
			return

		case data := <-saveChan:
			page := p.convertToPage(data)

			p.logger.Info("Coping data", "data", data)

			if err := p.storage.Save(ctx, &page); err != nil {
				p.logger.Error("Error saving data to database", "error", err)
			}
		}
	}

}

func (p *Processor) convertToPage(data *ownpb.CandleData) db.Page {
	page := &db.Page{
		Ticker:      data.Ticker,
		High:        p.convertPriceToFloat64(data.High),
		Low:         p.convertPriceToFloat64(data.Low),
		Open:        p.convertPriceToFloat64(data.Open),
		Close:       p.convertPriceToFloat64(data.Close),
		LastTradeTs: data.LastTradeTs.AsTime(),
	}

	return *page
}

func (p *Processor) convertPriceToFloat64(price *ownpb.PriceQuote) float64 {
	return float64(price.GetInteger()) + float64(price.GetFractional())/1e9
}
