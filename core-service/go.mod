module core-service

go 1.24.3

require (
	github.com/Point74/tinkoff-candle-streamer/contracts v0.0.0-00010101000000-000000000000
	github.com/twmb/franz-go v1.19.5
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.11.2 // indirect
)

replace github.com/Point74/tinkoff-candle-streamer/contracts => ../contracts
