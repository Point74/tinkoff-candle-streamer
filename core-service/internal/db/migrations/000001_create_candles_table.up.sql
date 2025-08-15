CREATE TABLE IF NOT EXISTS candles (
    id SERIAL PRIMARY KEY,
    ticker VARCHAR(20) NOT NULL,
    high DOUBLE PRECISION NOT NULL,
    low DOUBLE PRECISION NOT NULL,
    open DOUBLE PRECISION NOT NULL,
    close DOUBLE PRECISION NOT NULL,
    last_trade_ts timestamptz NOT NULL
);

CREATE INDEX IF NOT EXISTS candle_ticker_ts_idx ON candles (ticker, last_trade_ts DESC);