package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GetCandlesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "get_candles_total",
			Help: "The total number of getting candles",
		})

	ProcessedCandlesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "processed_candles_total",
			Help: "The total number of processed candles",
		})
)
