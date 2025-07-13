-- Seed data for development and testing

-- Insert Korean stocks (KOSPI major companies)
INSERT INTO stocks (symbol, name, market, exchange, sector, industry) VALUES
('005930', '삼성전자', 'KR', 'KOSPI', 'Technology', 'Semiconductors'),
('000660', 'SK하이닉스', 'KR', 'KOSPI', 'Technology', 'Semiconductors'),
('035420', 'NAVER', 'KR', 'KOSPI', 'Technology', 'Internet'),
('051910', 'LG화학', 'KR', 'KOSPI', 'Materials', 'Chemicals'),
('006400', '삼성SDI', 'KR', 'KOSPI', 'Technology', 'Battery'),
('207940', '삼성바이오로직스', 'KR', 'KOSPI', 'Healthcare', 'Biotechnology'),
('068270', '셀트리온', 'KR', 'KOSPI', 'Healthcare', 'Biotechnology'),
('028260', '삼성물산', 'KR', 'KOSPI', 'Industrials', 'Construction'),
('066570', 'LG전자', 'KR', 'KOSPI', 'Technology', 'Electronics'),
('003670', 'POSCO홀딩스', 'KR', 'KOSPI', 'Materials', 'Steel')
ON CONFLICT (symbol) DO NOTHING;

-- Insert US stocks (Major tech companies)
INSERT INTO stocks (symbol, name, market, exchange, sector, industry) VALUES
('AAPL', 'Apple Inc.', 'US', 'NASDAQ', 'Technology', 'Consumer Electronics'),
('MSFT', 'Microsoft Corporation', 'US', 'NASDAQ', 'Technology', 'Software'),
('GOOGL', 'Alphabet Inc. Class A', 'US', 'NASDAQ', 'Technology', 'Internet'),
('AMZN', 'Amazon.com Inc.', 'US', 'NASDAQ', 'Consumer Discretionary', 'E-commerce'),
('TSLA', 'Tesla Inc.', 'US', 'NASDAQ', 'Consumer Discretionary', 'Electric Vehicles'),
('NVDA', 'NVIDIA Corporation', 'US', 'NASDAQ', 'Technology', 'Semiconductors'),
('META', 'Meta Platforms Inc.', 'US', 'NASDAQ', 'Technology', 'Social Media'),
('NFLX', 'Netflix Inc.', 'US', 'NASDAQ', 'Communication Services', 'Streaming'),
('AMD', 'Advanced Micro Devices', 'US', 'NASDAQ', 'Technology', 'Semiconductors'),
('CRM', 'Salesforce Inc.', 'US', 'NYSE', 'Technology', 'Software')
ON CONFLICT (symbol) DO NOTHING;

-- Insert sample stock price data for testing (last 30 days)
INSERT INTO stock_prices (symbol, market, open_price, high_price, low_price, close_price, volume, timestamp) VALUES
-- Samsung Electronics (005930) sample data
('005930', 'KR', 71000.00, 72500.00, 70500.00, 72000.00, 12000000, NOW() - INTERVAL '1 hour'),
('005930', 'KR', 70500.00, 71200.00, 70000.00, 71000.00, 11500000, NOW() - INTERVAL '1 day'),
('005930', 'KR', 69800.00, 70800.00, 69500.00, 70500.00, 13000000, NOW() - INTERVAL '2 days'),

-- Apple (AAPL) sample data
('AAPL', 'US', 185.50, 187.20, 184.80, 186.90, 45000000, NOW() - INTERVAL '1 hour'),
('AAPL', 'US', 184.20, 186.00, 183.50, 185.50, 42000000, NOW() - INTERVAL '1 day'),
('AAPL', 'US', 183.80, 185.40, 182.90, 184.20, 38000000, NOW() - INTERVAL '2 days'),

-- NVIDIA (NVDA) sample data
('NVDA', 'US', 445.20, 452.80, 440.15, 448.65, 25000000, NOW() - INTERVAL '1 hour'),
('NVDA', 'US', 442.30, 447.90, 439.80, 445.20, 28000000, NOW() - INTERVAL '1 day'),
('NVDA', 'US', 438.50, 444.70, 435.20, 442.30, 31000000, NOW() - INTERVAL '2 days');

-- Insert sample technical indicators
INSERT INTO technical_indicators (symbol, indicator_name, indicator_value, calculated_at) VALUES
('005930', 'RSI', '{"value": 65.4, "period": 14}', NOW() - INTERVAL '1 hour'),
('005930', 'MACD', '{"macd": 1250.5, "signal": 980.2, "histogram": 270.3}', NOW() - INTERVAL '1 hour'),
('005930', 'SMA_20', '{"value": 70800.0}', NOW() - INTERVAL '1 hour'),
('005930', 'EMA_12', '{"value": 71200.0}', NOW() - INTERVAL '1 hour'),

('AAPL', 'RSI', '{"value": 58.7, "period": 14}', NOW() - INTERVAL '1 hour'),
('AAPL', 'MACD', '{"macd": 2.45, "signal": 1.89, "histogram": 0.56}', NOW() - INTERVAL '1 hour'),
('AAPL', 'SMA_20', '{"value": 184.5}', NOW() - INTERVAL '1 hour'),
('AAPL', 'EMA_12', '{"value": 185.8}', NOW() - INTERVAL '1 hour'),

('NVDA', 'RSI', '{"value": 72.3, "period": 14}', NOW() - INTERVAL '1 hour'),
('NVDA', 'MACD', '{"macd": 8.92, "signal": 6.45, "histogram": 2.47}', NOW() - INTERVAL '1 hour'),
('NVDA', 'SMA_20', '{"value": 442.8}', NOW() - INTERVAL '1 hour'),
('NVDA', 'EMA_12', '{"value": 446.2}', NOW() - INTERVAL '1 hour');

-- Insert sample trading signals
INSERT INTO trading_signals (symbol, signal_type, strength, confidence, reasons, source) VALUES
('005930', 'HOLD', 0.6, 0.75, '["RSI in neutral zone", "Price above SMA20", "Market conditions stable"]', 'AI'),
('AAPL', 'BUY', 0.8, 0.85, '["Strong uptrend", "Positive MACD", "Good earnings forecast"]', 'AI'),
('NVDA', 'SELL', 0.7, 0.80, '["RSI overbought", "High volume selling", "Profit taking zone"]', 'AI');

-- Insert sample news articles
INSERT INTO news_articles (title, content, url, source, sentiment_score, keywords, related_symbols, published_at) VALUES
('삼성전자, 3분기 실적 예상치 상회', '삼성전자가 3분기 실적에서 시장 예상치를 상회하는 성과를 보였다...', 'https://example.com/news1', 'Naver Finance', 0.7, '["실적", "상회", "반도체"]', '["005930"]', NOW() - INTERVAL '2 hours'),
('Apple Reports Strong iPhone Sales', 'Apple Inc. reported stronger than expected iPhone sales in the latest quarter...', 'https://example.com/news2', 'Yahoo Finance', 0.8, '["iPhone", "sales", "strong"]', '["AAPL"]', NOW() - INTERVAL '3 hours'),
('NVIDIA Stock Reaches New High', 'NVIDIA Corporation shares hit a new all-time high amid AI boom...', 'https://example.com/news3', 'Reuters', 0.6, '["AI", "high", "nvidia"]', '["NVDA"]', NOW() - INTERVAL '1 hour');