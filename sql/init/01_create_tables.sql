-- Stock Recommender Database Schema
-- This file is automatically executed when PostgreSQL container starts

-- Create database extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Stocks table
CREATE TABLE IF NOT EXISTS stocks (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100),
    market VARCHAR(5) NOT NULL CHECK (market IN ('KR', 'US')),
    exchange VARCHAR(20),
    sector VARCHAR(50),
    industry VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Stock prices table (partitioned by month for performance)
CREATE TABLE IF NOT EXISTS stock_prices (
    id BIGSERIAL,
    symbol VARCHAR(20) NOT NULL,
    market VARCHAR(5) NOT NULL,
    open_price DECIMAL(12,4),
    high_price DECIMAL(12,4),
    low_price DECIMAL(12,4),
    close_price DECIMAL(12,4),
    volume BIGINT,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
) PARTITION BY RANGE (timestamp);

-- Create partitions for current and next months
CREATE TABLE IF NOT EXISTS stock_prices_2024_01 PARTITION OF stock_prices
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_02 PARTITION OF stock_prices
FOR VALUES FROM ('2024-02-01') TO ('2024-03-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_03 PARTITION OF stock_prices
FOR VALUES FROM ('2024-03-01') TO ('2024-04-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_04 PARTITION OF stock_prices
FOR VALUES FROM ('2024-04-01') TO ('2024-05-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_05 PARTITION OF stock_prices
FOR VALUES FROM ('2024-05-01') TO ('2024-06-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_06 PARTITION OF stock_prices
FOR VALUES FROM ('2024-06-01') TO ('2024-07-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_07 PARTITION OF stock_prices
FOR VALUES FROM ('2024-07-01') TO ('2024-08-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_08 PARTITION OF stock_prices
FOR VALUES FROM ('2024-08-01') TO ('2024-09-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_09 PARTITION OF stock_prices
FOR VALUES FROM ('2024-09-01') TO ('2024-10-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_10 PARTITION OF stock_prices
FOR VALUES FROM ('2024-10-01') TO ('2024-11-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_11 PARTITION OF stock_prices
FOR VALUES FROM ('2024-11-01') TO ('2024-12-01');

CREATE TABLE IF NOT EXISTS stock_prices_2024_12 PARTITION OF stock_prices
FOR VALUES FROM ('2024-12-01') TO ('2025-01-01');

-- Technical indicators table
CREATE TABLE IF NOT EXISTS technical_indicators (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    indicator_name VARCHAR(50) NOT NULL,
    indicator_value JSONB,
    calculated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Trading signals table
CREATE TABLE IF NOT EXISTS trading_signals (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    signal_type VARCHAR(10) NOT NULL CHECK (signal_type IN ('BUY', 'SELL', 'HOLD')),
    strength DECIMAL(3,2) CHECK (strength >= 0 AND strength <= 1),
    confidence DECIMAL(3,2) CHECK (confidence >= 0 AND confidence <= 1),
    reasons JSONB,
    source VARCHAR(20) DEFAULT 'AI' CHECK (source IN ('AI', 'RULE', 'MANUAL')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- News articles table
CREATE TABLE IF NOT EXISTS news_articles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT,
    url TEXT,
    source VARCHAR(100),
    sentiment_score DECIMAL(3,2) CHECK (sentiment_score >= -1 AND sentiment_score <= 1),
    keywords JSONB,
    related_symbols JSONB,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_stocks_symbol ON stocks(symbol);
CREATE INDEX IF NOT EXISTS idx_stocks_market_active ON stocks(market, is_active);

CREATE INDEX IF NOT EXISTS idx_stock_prices_symbol_timestamp ON stock_prices(symbol, timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_stock_prices_market_timestamp ON stock_prices(market, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_technical_indicators_symbol_calculated ON technical_indicators(symbol, calculated_at DESC);
CREATE INDEX IF NOT EXISTS idx_technical_indicators_name_calculated ON technical_indicators(indicator_name, calculated_at DESC);

CREATE INDEX IF NOT EXISTS idx_trading_signals_symbol_created ON trading_signals(symbol, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_trading_signals_type_created ON trading_signals(signal_type, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_news_articles_published ON news_articles(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_news_articles_sentiment ON news_articles(sentiment_score);

-- GIN indexes for JSONB columns
CREATE INDEX IF NOT EXISTS idx_technical_indicators_value_gin ON technical_indicators USING GIN(indicator_value);
CREATE INDEX IF NOT EXISTS idx_trading_signals_reasons_gin ON trading_signals USING GIN(reasons);
CREATE INDEX IF NOT EXISTS idx_news_articles_keywords_gin ON news_articles USING GIN(keywords);
CREATE INDEX IF NOT EXISTS idx_news_articles_symbols_gin ON news_articles USING GIN(related_symbols);