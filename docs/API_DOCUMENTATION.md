# 📡 API 문서

## 개요

주식 투자 추천 시스템의 REST API 엔드포인트 상세 문서입니다.

## 기본 정보

- **Base URL**: `http://localhost:8080`
- **Content-Type**: `application/json`
- **응답 형식**: JSON

## 🏥 헬스 체크

### GET /health

시스템 상태를 확인합니다.

**응답 예시:**
```json
{
  "status": "ok",
  "timestamp": "2024-07-13T15:30:00Z",
  "database": "connected",
  "version": "1.0.0"
}
```

## 📈 주식 정보 API

### GET /api/v1/stocks

전체 종목 목록을 조회합니다.

**쿼리 파라미터:**
- `market` (선택): KR, US, INDEX
- `active` (선택): true, false
- `limit` (선택): 결과 개수 제한
- `offset` (선택): 페이지네이션

**응답 예시:**
```json
{
  "stocks": [
    {
      "id": 1,
      "symbol": "005930",
      "name": "삼성전자",
      "market": "KR",
      "exchange": "KOSPI",
      "sector": "Technology",
      "is_active": true,
      "created_at": "2024-07-13T10:00:00Z"
    }
  ],
  "total": 50,
  "page": 1
}
```

### GET /api/v1/stocks/{symbol}

특정 종목의 상세 정보를 조회합니다.

**경로 파라미터:**
- `symbol`: 종목 코드 (예: 005930, AAPL)

**응답 예시:**
```json
{
  "stock": {
    "id": 1,
    "symbol": "005930",
    "name": "삼성전자",
    "market": "KR",
    "exchange": "KOSPI",
    "sector": "Technology",
    "industry": "Semiconductors",
    "is_active": true
  }
}
```

### GET /api/v1/stocks/{symbol}/price

특정 종목의 최신 주가 정보를 조회합니다.

**응답 예시:**
```json
{
  "price": {
    "symbol": "005930",
    "current_price": 70500.0,
    "open_price": 70000.0,
    "high_price": 71000.0,
    "low_price": 69500.0,
    "prev_close_price": 70200.0,
    "change": 300.0,
    "change_rate": 0.43,
    "volume": 12345678,
    "trade_amount": 870000000000,
    "timestamp": "2024-07-13T15:30:00Z"
  }
}
```

### GET /api/v1/stocks/{symbol}/indicators

특정 종목의 기술지표를 조회합니다.

**응답 예시:**
```json
{
  "indicators": {
    "symbol": "005930",
    "rsi": 65.4,
    "macd": {
      "macd": 1250.5,
      "signal": 1180.2,
      "histogram": 70.3
    },
    "bollinger_bands": {
      "upper": 72000.0,
      "middle": 70000.0,
      "lower": 68000.0
    },
    "sma_20": 70800.0,
    "sma_50": 69500.0,
    "ema_12": 70900.0,
    "ema_26": 70200.0,
    "stochastic": {
      "k": 75.2,
      "d": 72.8
    },
    "williams_r": -24.6,
    "atr": 1250.0,
    "obv": 15000000,
    "calculated_at": "2024-07-13T15:30:00Z"
  }
}
```

## 🎯 매매 신호 API

### GET /api/v1/signals

전체 매매 신호를 조회합니다.

**쿼리 파라미터:**
- `signal_type` (선택): BUY, SELL, HOLD
- `market` (선택): KR, US
- `limit` (선택): 결과 개수 제한

**응답 예시:**
```json
{
  "signals": [
    {
      "id": 1,
      "symbol": "005930",
      "signal_type": "BUY",
      "strength": 0.85,
      "confidence": 0.78,
      "reasons": [
        "RSI oversold condition",
        "MACD bullish crossover",
        "Positive news sentiment"
      ],
      "source": "AI",
      "created_at": "2024-07-13T15:30:00Z"
    }
  ],
  "total": 25
}
```

### GET /api/v1/signals/{symbol}

특정 종목의 매매 신호를 조회합니다.

**쿼리 파라미터:**
- `limit` (선택): 최근 신호 개수 (기본값: 10)

**응답 예시:**
```json
{
  "signals": [
    {
      "symbol": "005930",
      "signal_type": "BUY",
      "strength": 0.85,
      "confidence": 0.78,
      "reasons": [
        "강한 상승 모멘텀",
        "거래량 급증",
        "긍정적 뉴스 감성"
      ],
      "source": "AI",
      "created_at": "2024-07-13T15:30:00Z"
    }
  ]
}
```

## 🔧 관리자 API

### 종목 관리

#### POST /api/v1/admin/stocks

새로운 종목을 등록합니다.

**요청 본문:**
```json
{
  "symbol": "TEST001",
  "name": "테스트 종목",
  "market": "KR",
  "exchange": "KOSPI",
  "sector": "Technology",
  "industry": "Software"
}
```

**응답 예시:**
```json
{
  "message": "Stock created successfully",
  "stock": {
    "id": 100,
    "symbol": "TEST001",
    "name": "테스트 종목",
    "market": "KR",
    "exchange": "KOSPI",
    "sector": "Technology",
    "industry": "Software",
    "is_active": true,
    "created_at": "2024-07-13T15:30:00Z"
  }
}
```

#### GET /api/v1/admin/stocks

관리자용 전체 종목 목록을 조회합니다.

**응답 예시:**
```json
{
  "stocks": [...],
  "count": 50
}
```

#### PUT /api/v1/admin/stocks/{symbol}/status

종목의 활성화 상태를 변경합니다.

**요청 본문:**
```json
{
  "is_active": false
}
```

#### DELETE /api/v1/admin/stocks/{symbol}

종목을 삭제합니다.

### 데이터 수집 관리

#### POST /api/v1/admin/collect/{symbol}

특정 종목의 데이터 수집을 트리거합니다.

**응답 예시:**
```json
{
  "message": "Data collection triggered successfully",
  "symbol": "005930"
}
```

#### POST /api/v1/admin/collect/all

전체 종목의 데이터 수집을 트리거합니다.

**응답 예시:**
```json
{
  "message": "Batch data collection started"
}
```

#### POST /api/v1/admin/initialize/major-stocks

주요 종목을 자동으로 등록합니다.

**응답 예시:**
```json
{
  "message": "Major stocks initialized successfully"
}
```

### 시스템 모니터링

#### GET /api/v1/admin/api-status

DB증권 API 연결 상태를 확인합니다.

**응답 예시:**
```json
{
  "api_status": {
    "authenticated": true,
    "base_url": "https://openapi.dbsec.co.kr",
    "rate_limit": "20 requests/second",
    "api_available": true,
    "timestamp": "2024-07-13T15:30:00Z"
  }
}
```

#### GET /api/v1/admin/database/stats

데이터베이스 통계를 조회합니다.

**응답 예시:**
```json
{
  "total_stocks": 50,
  "active_stocks": 45,
  "total_price_points": 125000,
  "total_signals": 850,
  "last_update": "2024-07-13 15:30:00"
}
```

## 🚨 오류 응답

모든 API는 표준화된 오류 형식을 사용합니다.

### 4xx 클라이언트 오류

```json
{
  "error": "Bad Request",
  "message": "Invalid symbol format",
  "code": 400,
  "timestamp": "2024-07-13T15:30:00Z"
}
```

### 5xx 서버 오류

```json
{
  "error": "Internal Server Error",
  "message": "Database connection failed",
  "code": 500,
  "timestamp": "2024-07-13T15:30:00Z"
}
```

## 📊 상태 코드

| 코드 | 설명 |
|------|------|
| 200 | 요청 성공 |
| 201 | 리소스 생성 성공 |
| 202 | 요청 접수 (비동기 처리) |
| 400 | 잘못된 요청 |
| 401 | 인증 필요 |
| 403 | 권한 부족 |
| 404 | 리소스 없음 |
| 409 | 리소스 충돌 |
| 429 | 요청 한도 초과 |
| 500 | 서버 내부 오류 |
| 502 | 게이트웨이 오류 |
| 503 | 서비스 이용 불가 |

## 🔒 인증 및 보안

현재 버전은 개발/테스트 목적으로 인증이 없지만, 프로덕션 환경에서는 다음 보안 기능이 추가될 예정입니다:

- **JWT 토큰 인증**
- **API 키 기반 접근 제어**
- **Rate Limiting**
- **HTTPS 강제**
- **요청 로깅 및 모니터링**

## 📝 사용 예시

### cURL 예시

```bash
# 헬스 체크
curl http://localhost:8080/health

# 삼성전자 주가 조회
curl http://localhost:8080/api/v1/stocks/005930/price

# 매매 신호 조회
curl http://localhost:8080/api/v1/signals/005930

# 데이터 수집 트리거
curl -X POST http://localhost:8080/api/v1/admin/collect/005930
```

### JavaScript 예시

```javascript
// 주가 정보 조회
const response = await fetch('http://localhost:8080/api/v1/stocks/005930/price');
const data = await response.json();
console.log(data.price.current_price);

// 매매 신호 조회
const signals = await fetch('http://localhost:8080/api/v1/signals/005930');
const signalData = await signals.json();
console.log(signalData.signals[0].signal_type);
```

## 📈 실시간 데이터

- **주가 데이터**: 5분마다 자동 업데이트
- **기술지표**: 주가 업데이트 시 자동 계산
- **매매 신호**: 지표 계산 후 AI 분석 실행
- **뉴스 데이터**: 1시간마다 수집 및 감성 분석

---

**📅 문서 업데이트**: 2024년 7월 13일  
**🔗 관련 문서**: [시스템 아키텍처](TECHNICAL_ARCHITECTURE.md), [배포 가이드](DEPLOYMENT_GUIDE.md)