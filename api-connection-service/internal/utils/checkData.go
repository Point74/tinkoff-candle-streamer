package utils

import (
	"fmt"
	pb "github.com/Point74/tinkoff-candle-streamer/contracts/gen/doc"
)

func ValidCandle(candle *pb.Candle) error {
	if candle.Open == nil && candle.High == nil && candle.Low == nil && candle.Close == nil {
		return fmt.Errorf("invalid price of candle")
	}

	if candle.Volume < 0 {
		return fmt.Errorf("invalid volume of candle")
	}

	openPrice := toFloat64(candle.Open)
	highPrice := toFloat64(candle.High)
	lowPrice := toFloat64(candle.Low)
	closePrice := toFloat64(candle.Close)

	if openPrice <= 0 && highPrice <= 0 && lowPrice <= 0 && closePrice <= 0 {
		return fmt.Errorf("price of candle must be positive, bigger than 0")
	}

	return nil
}

func toFloat64(price *pb.Quotation) float64 {
	return float64(price.Units) + float64(price.Nano)
}
