package db

import (
	"context"
	"time"
)

type Storage interface {
	Save(ctx context.Context, p *Page)
}

type Page struct {
	Ticker      string
	High        float64
	Low         float64
	Open        float64
	Close       float64
	LastTradeTs time.Time
}
