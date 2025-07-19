# OpenAPI Package

DB증권 Open API를 위한 Go 클라이언트 라이브러리입니다.

## 패키지 구조

```
openapi/
├── client/          # HTTP 클라이언트 레이어
├── domestic/        # 국내 주식 API 서비스
├── foreign/         # 해외 주식 API 서비스
├── futures/         # 선물 API 서비스 (예정)
├── models/          # 데이터 모델 정의
├── utils/           # 공통 유틸리티
├── errors/          # 에러 타입 정의
└── logger/          # 로깅 시스템
```

## 주요 기능

### 🏠 국내 주식 (Domestic)
- **종목 조회**: 국내 주식, ETF, ETN 종목 목록
- **현재가 조회**: 실시간 주가 정보 및 관련 지표

### 🌍 해외 주식 (Foreign)  
- **종목 조회**: 미국 주식 종목 목록 (뉴욕, 나스닥, 아멕스)
- **현재가 조회**: 실시간 미국 주식 가격 정보

## 빠른 시작

### 1. 클라이언트 초기화

```go
import (
    "stock-recommender/backend/config"
    "stock-recommender/backend/openapi/client"
    "stock-recommender/backend/openapi/domestic"
    "stock-recommender/backend/openapi/foreign"
)

// 설정 로드
cfg := config.LoadConfig()

// API 클라이언트 생성
apiClient := client.NewDBSecClient(cfg)

// 서비스 생성
domesticPrice := domestic.NewCurrentPriceService(apiClient)
foreignPrice := foreign.NewForeignCurrentPriceService(apiClient)
```

### 2. 국내 주식 현재가 조회

```go
// 삼성전자 현재가 조회
data, err := domesticPrice.GetStockPrice("005930")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("종목: %s\n", data.StockCode)
fmt.Printf("현재가: %.0f원\n", data.CurrentPrice)
fmt.Printf("전일대비: %.0f원 (%.2f%%)\n", data.PriceChange, data.PriceChangeRate)
```

### 3. 해외 주식 현재가 조회

```go
// 테슬라 현재가 조회
data, err := foreignPrice.GetNASDAQStockPrice("TSLA")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("종목: %s\n", data.StockCode)
fmt.Printf("현재가: $%.2f\n", data.CurrentPrice)
fmt.Printf("전일대비: $%.2f (%.2f%%)\n", data.PriceChange, data.PriceChangeRate)
```

## 아키텍처

### 🔧 Client Layer
- **인증 관리**: OAuth 2.0 토큰 자동 발급/갱신
- **Rate Limiting**: API 호출 제한 준수
- **에러 처리**: 자동 재시도 및 에러 타입 분류
- **로깅**: 구조화된 로깅 시스템

### 📊 Service Layer
- **Domain Services**: 비즈니스 로직 캡슐화
- **Data Transformation**: API 응답을 비즈니스 모델로 변환
- **Convenience Methods**: 자주 사용되는 조회 패턴 제공

### 📋 Models Layer
- **Request/Response Models**: API 스펙과 정확히 매핑
- **Business Models**: 사용하기 쉬운 구조화된 데이터
- **Constants**: API 경로, 코드 등 상수 정의

## 에러 처리

### 에러 타입

```go
import "stock-recommender/backend/openapi/errors"

// 에러 타입 확인
if errors.IsAuthError(err) {
    // 인증 에러 처리
} else if errors.IsRetryableError(err) {
    // 재시도 가능한 에러 처리
}
```

### 주요 에러 코드
- `AUTH_FAILED`: 인증 실패
- `TOKEN_EXPIRED`: 토큰 만료
- `RATE_LIMIT`: 호출 제한 초과
- `NETWORK_ERROR`: 네트워크 오류
- `PARSE_ERROR`: 응답 파싱 오류

## 로깅

### 로그 레벨 설정

```go
import "stock-recommender/backend/openapi/logger"

// 로그 레벨 설정
logger.SetDefaultLogger(logger.NewDefaultLogger(logger.DEBUG))

// 구조화된 로깅
logger.Info("API call completed", 
    logger.Field{Key: "stock_code", Value: "005930"},
    logger.Field{Key: "duration_ms", Value: 150})
```

## 유틸리티

### 데이터 파싱

```go
import "stock-recommender/backend/openapi/utils"

// 안전한 문자열 파싱
price := utils.ParseFloat("55550.00")    // 55550.0
volume := utils.ParseInt("7240324")      // 7240324
date := utils.ParseDate("20231225")      // time.Time
```

### 페이지네이션

```go
// 페이지네이션 헬퍼 사용
pagination := utils.NewPaginationHelper()
for {
    response, nextKey, err := service.GetStockTickers(marketDiv, pagination.ContKey)
    if err != nil {
        break
    }
    
    // 데이터 처리
    processData(response.Out)
    
    pagination.SetNextKey(nextKey)
    if !pagination.HasNext() {
        break
    }
}
```

## 테스트

### 단위 테스트

```go
import "stock-recommender/backend/openapi/utils"

func TestMyService(t *testing.T) {
    // 테스트 클라이언트 생성
    client, cleanup := utils.CreateTestClient(t)
    defer cleanup()
    
    // 자격증명 확인
    utils.SkipIfNoCredentials(t, client)
    
    // 테스트 실행
    service := NewMyService(client)
    result, err := service.DoSomething()
    
    // 검증
    utils.AssertStringEqual(t, "expected", result.Value, "Result value")
}
```

### 모의 테스트

```go
// 모의 서버 생성
handler := utils.CreateCurrentPriceMockHandler(t, "005930", mockData)
mockServer := utils.NewMockServer(t, handler)
defer mockServer.Close()
```

## 설정

### 환경 변수

```bash
# DB증권 API 키 설정
export DBSEC_APP_KEY="your_app_key"
export DBSEC_APP_SECRET="your_app_secret"
```

### 설정 파일

```yaml
# config.yaml
api:
  dbsec_app_key: "your_app_key"
  dbsec_app_secret: "your_app_secret"
  timeout: 30s
  rate_limit: 20
```

## API 제한사항

### 호출 제한
- **국내 주식 종목 조회**: 초당 3건
- **국내 주식 현재가 조회**: 초당 5건  
- **해외 주식 종목 조회**: 초당 2건
- **해외 주식 현재가 조회**: 초당 2건

### 데이터 제한
- **페이지네이션**: 대량 데이터는 연속키를 통한 분할 조회
- **실시간성**: 현재가는 실시간, 종목 목록은 정기 업데이트
- **통화**: 국내 주식은 KRW, 해외 주식은 USD

## 모범 사례

### 1. 에러 처리

```go
data, err := service.GetStockPrice("005930")
if err != nil {
    if errors.IsAuthError(err) {
        // 인증 재시도 로직
        return handleAuthError(err)
    } else if errors.IsRetryableError(err) {
        // 재시도 로직
        return retryWithBackoff(func() error {
            data, err = service.GetStockPrice("005930")
            return err
        })
    }
    return err
}
```

### 2. 대량 조회 최적화

```go
// 병렬 처리로 성능 향상
codes := []string{"005930", "000660", "035720"}
results := make(chan *models.CurrentPriceData, len(codes))

for _, code := range codes {
    go func(stockCode string) {
        data, err := service.GetStockPrice(stockCode)
        if err == nil {
            results <- data
        }
    }(code)
}
```

### 3. 캐싱 활용

```go
// 종목 목록은 하루 단위로 캐싱
type CachedStockService struct {
    service domestic.StockTickerService
    cache   map[string][]models.StockTickerOutput
    lastUpdate time.Time
}

func (c *CachedStockService) GetStocks() ([]models.StockTickerOutput, error) {
    if time.Since(c.lastUpdate) > 24*time.Hour {
        stocks, err := c.service.GetStocks()
        if err != nil {
            return nil, err
        }
        c.cache["stocks"] = stocks
        c.lastUpdate = time.Now()
    }
    return c.cache["stocks"], nil
}
```

## 기여하기

1. 이슈 생성 또는 기능 요청
2. 브랜치 생성: `git checkout -b feature/amazing-feature`
3. 변경사항 커밋: `git commit -m 'Add amazing feature'`
4. 브랜치 푸시: `git push origin feature/amazing-feature`
5. Pull Request 생성

## 라이센스

이 프로젝트는 MIT 라이센스 하에 배포됩니다.

## 지원

- 📚 [API 문서](./domestic/README.md)
- 🐛 [이슈 리포트](https://github.com/your-repo/issues)
- 💬 [디스커션](https://github.com/your-repo/discussions)