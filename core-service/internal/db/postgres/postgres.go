package postgres

import (
	"context"
	"core-service/internal/db"
	"fmt"
	"github.com/Point74/tinkoff-candle-streamer/config"
	"github.com/Point74/tinkoff-candle-streamer/prometheus/metrics"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"time"
)

type Storage struct {
	database *pgx.Conn
	logger   *slog.Logger
}

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Storage, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	database, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	if err := database.Ping(ctx); err != nil {
		logger.Error("Ping to connect database failed", "error", err)
		database.Close(ctx)
		return nil, err
	}

	logger.Info("Connected to database")
	logger.Info("Running migrations succeeded", "url", dbUrl)

	return &Storage{
		database: database,
		logger:   logger,
	}, nil
}

func (s *Storage) Save(ctx context.Context, p *db.Page) error {
	start := time.Now()
	status := "success"

	defer func() {
		metrics.SendCandlesToDBTotal.WithLabelValues("Save", status).Inc()
		metrics.SendCandlesToDBDuration.WithLabelValues("Save", status).Observe(time.Since(start).Seconds())
	}()

	sql := `INSERT INTO candles (Ticker, High, Low, Open, Close, Last_trade_ts) VALUES ($1, $2, $3, $4, $5, $6)`
	if _, err := s.database.Exec(
		ctx,
		sql,
		p.Ticker,
		p.High,
		p.Low,
		p.Open,
		p.Close,
		p.LastTradeTs,
	); err != nil {
		status = "error"
		s.logger.Error("Unable to save candle", "error", err)
		return err
	}

	s.logger.Info("Candle saved successfully", "ticker", p.Ticker)

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.database != nil {
		return s.database.Close(ctx)
	}

	return nil
}
