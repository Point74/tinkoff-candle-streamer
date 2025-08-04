module core-service

go 1.24.3

require (
	github.com/Point74/tinkoff-candle-streamer/config v0.0.0-20250731202554-1b85203d13a4
	github.com/Point74/tinkoff-candle-streamer/contracts v0.0.0-00010101000000-000000000000
	github.com/golang-migrate/migrate/v4 v4.18.3
	github.com/jackc/pgx/v5 v5.7.5
	github.com/twmb/franz-go v1.19.5
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/caarlos0/env/v11 v11.3.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/twmb/franz-go/pkg/kmsg v1.11.2 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/text v0.25.0 // indirect
)

replace github.com/Point74/tinkoff-candle-streamer/contracts => ../contracts

replace github.com/Point74/tinkoff-candle-streamer/config => ../config