package postgres

import (
	"context"
	"fmt"
	"github.com/Point74/tinkoff-candle-streamer/config"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

type Storage struct {
	db *pgx.Conn
}

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Storage, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	db, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		logger.Error("Unable to connect to database", "error", err)
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		logger.Error("Ping to connect database failed", "error", err)
		db.Close(ctx)
		return nil, err
	}

	logger.Info("Connected to database")

	return &Storage{db: db}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	if s.db != nil {
		return s.db.Close(ctx)
	}

	return nil
}
