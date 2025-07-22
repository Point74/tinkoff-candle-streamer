module api-connection-service

go 1.24.3

require (
	github.com/Point74/tinkoff-candle-streamer/contracts v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.73.0
)

require (
	github.com/twmb/franz-go v1.19.5 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/Point74/tinkoff-candle-streamer/contracts => ../contracts
