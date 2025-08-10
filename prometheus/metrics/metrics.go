package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	GrpcRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_request_total",
			Help: "The total number of grpc requests",
		},
		[]string{"method", "status"},
	)

	GrpcRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	SerializedCandlesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "serialized_candles_total",
			Help: "The total number of serialized candles",
		},
		[]string{"method", "status"},
	)

	SerializedCandlesDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "serialized_candles_duration",
			Help: "Duration of serialization candles",
		},
		[]string{"method", "status"},
	)

	KafkaSendDataTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_send_data_total",
			Help: "Total number of messages send to Kafka",
		},
		[]string{"method", "status"},
	)

	KafkaGetDataTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_get_data_total",
			Help: "Total number of messages get to Kafka",
		},
		[]string{"method", "status"},
	)

	KafkaRecordProcessingDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "kafka_get_data_duration",
			Help: "Duration of getting data from Kafka",
		},
		[]string{"method", "status"},
	)

	KafkaFetchDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "kafka_fetch_duration",
			Help: "Duration of fetch data from Kafka",
		},
		[]string{"method", "status"},
	)

	SendCandlesToDBTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "send_candles_to_db_total",
			Help: "Total number of record in db",
		},
		[]string{"method", "status"},
	)

	SendCandlesToDBDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "send_candles_to_db_duration",
			Help: "Duration of record in db",
		},
		[]string{"method", "status"},
	)
)
