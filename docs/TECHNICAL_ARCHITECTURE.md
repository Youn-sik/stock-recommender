# 기술 아키텍처 설계

## 시스템 아키텍처 개요

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Backend       │
│   (React)       │◄──►│   (Gin)         │◄──►│   Services      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                       ┌─────────────────┐            │
                       │   Message Queue │◄───────────┘
                       │   (RabbitMQ)    │
                       └─────────────────┘
                                │
                                ▼
        ┌─────────────────┬─────────────────┬─────────────────┐
        │   Data Collector│   ML Pipeline   │   Crawler       │
        │   (DB증권 API)   │   (Python)      │   (News/Data)   │
        └─────────────────┴─────────────────┴─────────────────┘
                                │
                                ▼
                    ┌─────────────────────────────────┐
                    │          Database               │
                    │  ┌─────────────┬─────────────┐  │
                    │  │ PostgreSQL  │   Redis     │  │
                    │  │ (시계열)     │  (캐시)      │  │
                    │  └─────────────┴─────────────┘  │
                    └─────────────────────────────────┘
```

## 마이크로서비스 구성

### 1. API Gateway Service (Go)
**책임**: 라우팅, 인증, Rate Limiting
```go
// main.go
func main() {
    r := gin.Default()
    
    // 미들웨어
    r.Use(middleware.Auth())
    r.Use(middleware.RateLimit())
    r.Use(middleware.CORS())
    
    // 라우트 그룹
    api := r.Group("/api/v1")
    {
        api.GET("/stocks/:symbol", stockHandler.GetStock)
        api.GET("/signals", signalHandler.GetSignals)
        api.POST("/alerts", alertHandler.CreateAlert)
    }
    
    r.Run(":8080")
}
```

### 2. Data Collection Service (Go)
**책임**: DB증권 API 연동, 실시간 데이터 수집
```go
type DataCollector struct {
    client     *http.Client
    apiKey     string
    symbols    []string
    publisher  MessagePublisher
}

func (dc *DataCollector) CollectRealTimeData() {
    for _, symbol := range dc.symbols {
        data := dc.fetchStockData(symbol)
        dc.publisher.Publish("stock.data", data)
    }
}
```

### 3. Technical Indicator Service (Go + Python)
**책임**: 기술지표 계산, 신호 생성
```python
# ml/indicators/calculator.py
class IndicatorCalculator:
    def __init__(self):
        self.indicators = {
            'rsi': self.calculate_rsi,
            'macd': self.calculate_macd,
            'bollinger': self.calculate_bollinger_bands
        }
    
    def calculate_all_indicators(self, data):
        results = {}
        for name, func in self.indicators.items():
            results[name] = func(data)
        return results
```

### 4. ML Prediction Service (Python)
**책임**: 주가 예측, 모델 학습
```python
# ml/models/predictor.py
class StockPredictor:
    def __init__(self):
        self.models = {
            'lstm': LSTMModel(),
            'random_forest': RandomForestModel(),
            'ensemble': EnsembleModel()
        }
    
    def predict(self, symbol, features):
        predictions = {}
        for name, model in self.models.items():
            predictions[name] = model.predict(features)
        
        return self.ensemble_predict(predictions)
```

### 5. News Crawler Service (Python)
**책임**: 뉴스 크롤링, 감성 분석
```python
# crawler/news_crawler.py
class NewsCrawler:
    def __init__(self):
        self.sources = [
            NaverFinanceSource(),
            DaumFinanceSource(),
            InvestingSource()
        ]
    
    async def crawl_news(self, keywords):
        tasks = []
        for source in self.sources:
            tasks.append(source.fetch_news(keywords))
        
        results = await asyncio.gather(*tasks)
        return self.merge_results(results)
```

## 데이터베이스 설계

### PostgreSQL 스키마
```sql
-- 주가 데이터 (시계열)
CREATE TABLE stock_prices (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    market VARCHAR(5) NOT NULL,
    open_price DECIMAL(12,4),
    high_price DECIMAL(12,4),
    low_price DECIMAL(12,4),
    close_price DECIMAL(12,4),
    volume BIGINT,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 파티셔닝 (월별)
CREATE TABLE stock_prices_y2024m01 PARTITION OF stock_prices
FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- 기술지표 결과
CREATE TABLE technical_indicators (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    indicator_name VARCHAR(50) NOT NULL,
    indicator_value JSONB,
    calculated_at TIMESTAMPTZ NOT NULL
);

-- 매매 신호
CREATE TABLE trading_signals (
    id BIGSERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    signal_type VARCHAR(10) NOT NULL, -- BUY/SELL/HOLD
    strength DECIMAL(3,2),
    reasons JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 뉴스 데이터
CREATE TABLE news_articles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT,
    sentiment_score DECIMAL(3,2),
    keywords JSONB,
    source VARCHAR(100),
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 인덱스
CREATE INDEX idx_stock_prices_symbol_timestamp 
ON stock_prices (symbol, timestamp DESC);

CREATE INDEX idx_trading_signals_symbol_created 
ON trading_signals (symbol, created_at DESC);
```

### Redis 캐싱 전략
```go
// cache/strategy.go
type CacheStrategy struct {
    client redis.Client
}

// 실시간 주가 캐싱 (TTL: 1분)
func (c *CacheStrategy) CacheStockPrice(symbol string, price StockPrice) {
    key := fmt.Sprintf("stock:price:%s", symbol)
    c.client.Set(key, price, time.Minute)
}

// 기술지표 캐싱 (TTL: 5분)
func (c *CacheStrategy) CacheIndicators(symbol string, indicators map[string]float64) {
    key := fmt.Sprintf("indicators:%s", symbol)
    c.client.HMSet(key, indicators)
    c.client.Expire(key, 5*time.Minute)
}
```

## 메시지 큐 설계

### RabbitMQ Exchange/Queue 구조
```
Exchange: stock.data
├── Queue: price.updates (실시간 가격 업데이트)
├── Queue: indicator.calculation (지표 계산 요청)
└── Queue: ml.prediction (예측 모델 실행)

Exchange: trading.signals
├── Queue: signal.generation (신호 생성)
└── Queue: alert.notification (알림 발송)

Exchange: news.analysis
├── Queue: crawl.requests (크롤링 요청)
└── Queue: sentiment.analysis (감성 분석)
```

### 메시지 처리 워커
```go
// workers/price_worker.go
type PriceWorker struct {
    consumer MessageConsumer
    processor PriceProcessor
}

func (w *PriceWorker) Start() {
    w.consumer.Subscribe("price.updates", func(msg Message) {
        var priceData StockPrice
        json.Unmarshal(msg.Body, &priceData)
        
        // 기술지표 계산 트리거
        w.processor.TriggerIndicatorCalculation(priceData.Symbol)
        
        // 실시간 알림 확인
        w.processor.CheckAlerts(priceData)
    })
}
```

## API 설계

### RESTful API Endpoints
```yaml
# OpenAPI 3.0 스펙
paths:
  /api/v1/stocks/{symbol}:
    get:
      summary: 종목 상세 정보 조회
      parameters:
        - name: symbol
          in: path
          required: true
          schema:
            type: string
      responses:
        200:
          description: 성공
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/StockDetail'

  /api/v1/signals:
    get:
      summary: 매매 신호 목록
      parameters:
        - name: market
          in: query
          schema:
            type: string
            enum: [KR, US]
        - name: signal_type
          in: query
          schema:
            type: string
            enum: [BUY, SELL, HOLD]
      responses:
        200:
          description: 성공
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/TradingSignal'

components:
  schemas:
    StockDetail:
      type: object
      properties:
        symbol:
          type: string
        current_price:
          type: number
        indicators:
          type: object
        signals:
          type: array
          items:
            $ref: '#/components/schemas/TradingSignal'
    
    TradingSignal:
      type: object
      properties:
        symbol:
          type: string
        signal_type:
          type: string
          enum: [BUY, SELL, HOLD]
        strength:
          type: number
          minimum: 0
          maximum: 1
        reasons:
          type: array
          items:
            type: string
        created_at:
          type: string
          format: date-time
```

### WebSocket API (실시간)
```go
// websocket/hub.go
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
    
    h.register <- client
    
    go client.readPump()
    go client.writePump()
}

// 실시간 데이터 브로드캐스트
func (h *Hub) BroadcastPriceUpdate(symbol string, price float64) {
    message := map[string]interface{}{
        "type": "price_update",
        "symbol": symbol,
        "price": price,
        "timestamp": time.Now(),
    }
    
    data, _ := json.Marshal(message)
    h.broadcast <- data
}
```

## 보안 설계

### 인증/인가
```go
// middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        claims, err := jwt.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

### API Rate Limiting
```go
// middleware/rate_limit.go
func RateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(time.Second), 100) // 초당 100 요청
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## 모니터링 및 로깅

### Prometheus 메트릭
```go
// metrics/collector.go
var (
    apiRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "api_requests_total",
            Help: "Total number of API requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    stockPriceUpdates = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "stock_price_updates_total",
            Help: "Total number of stock price updates",
        },
        []string{"symbol", "market"},
    )
)
```

### 구조화된 로깅
```go
// logger/logger.go
func LogTradingSignal(signal TradingSignal) {
    log.WithFields(logrus.Fields{
        "symbol":      signal.Symbol,
        "signal_type": signal.SignalType,
        "strength":    signal.Strength,
        "timestamp":   signal.CreatedAt,
    }).Info("Trading signal generated")
}
```

## 배포 및 운영

### Docker Compose 구성
```yaml
# docker-compose.yml
version: '3.8'
services:
  api-gateway:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis
  
  data-collector:
    build: ./backend
    command: ["./data-collector"]
    environment:
      - DBSEC_APP_KEY=${DBSEC_APP_KEY}
    depends_on:
      - postgres
      - rabbitmq
  
  ml-service:
    build: ./ml
    ports:
      - "8001:8001"
    volumes:
      - ./models:/app/models
    depends_on:
      - postgres
  
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: stockdb
      POSTGRES_USER: stockuser
      POSTGRES_PASSWORD: stockpass
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    image: redis:6-alpine
    command: redis-server --maxmemory 256mb --maxmemory-policy allkeys-lru
  
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"

volumes:
  postgres_data:
```

### Kubernetes 배포 (운영환경)
```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: stock-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: stock-api
  template:
    metadata:
      labels:
        app: stock-api
    spec:
      containers:
      - name: api
        image: stock-recommender/api:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          value: "postgres-service"
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

이 아키텍처는 고가용성, 확장성, 보안을 모두 고려한 설계로, 실제 운영 환경에서 안정적으로 동작할 수 있도록 구성되었습니다.