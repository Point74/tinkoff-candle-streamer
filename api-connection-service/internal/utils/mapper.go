package utils

import (
	pb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/doc"
	ownpb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/my"
	"github.com/jinzhu/copier"
	"log/slog"
)

func quotationToPriceQuote(price *pb.Quotation) *ownpb.PriceQuote {
	return &ownpb.PriceQuote{
		Integer:    price.Units,
		Fractional: price.Nano,
	}
}

func MapCandle(logger *slog.Logger, sourceData *pb.Candle, tickerShare string) (*ownpb.CandleData, error) {
	if sourceData == nil {
		logger.Error("source data is nil", "source", sourceData)
		return nil, nil
	}

	target := &ownpb.CandleData{}
	if err := copier.Copy(&target, &sourceData); err != nil {
		logger.Warn("failed to copy data candle", "err", err)
		return nil, err
	}

	target.Ticker = tickerShare
	target.High = quotationToPriceQuote(sourceData.High)
	target.Low = quotationToPriceQuote(sourceData.Low)
	target.Open = quotationToPriceQuote(sourceData.Open)
	target.Close = quotationToPriceQuote(sourceData.Close)

	return target, nil
}
