# DB증권 Open API 통합 가이드

## 📋 개요

본 문서는 주식 투자 추천 프로그램에서 DB증권 Open API를 통합하여 실시간 주식 데이터를 수집하는 방법을 설명합니다.

## 🔗 API 연동 현황

### ✅ 구현 완료된 기능

1. **인증 시스템**
   - OAuth 2.0 기반 토큰 인증
   - 자동 토큰 갱신
   - Rate limiting (초당 20요청)

2. **국내주식 API**
   - 현재가 조회 (`/uapi/domestic-stock/v1/quotations/inquire-price`)
   - 호가 정보 조회 (`/uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn`)
   - 일봉차트 조회 (`/uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice`)
   - 종목 리스트 조회 (`/uapi/domestic-stock/v1/quotations/search-stock-info`)

3. **해외주식 API**
   - 현재가 조회 (`/uapi/overseas-price/v1/quotations/price`)
   - 일봉차트 조회 (`/uapi/overseas-price/v1/quotations/dailyprice`)

4. **지수 API**
   - 지수 현재가 조회 (`/uapi/domestic-stock/v1/quotations/inquire-index-price`)

5. **데이터 수집 시스템**
   - 자동 정기 수집 (5분 간격)
   - 주요 종목 자동 등록
   - Mock 데이터 fallback (개발용)

## 🔧 API 키 설정

### 1. DB증권 API 키 발급

1. [DB증권 Open API 포털](https://openapi.dbsec.co.kr/) 방문
2. 회원가입 및 로그인
3. 앱 등록 후 API 키 발급:
   - `APP_KEY`: 애플리케이션 키
   - `APP_SECRET`: 애플리케이션 시크릿

### 2. 환경변수 설정

```bash
# .env 파일에 추가
DBSEC_APP_KEY=your_actual_app_key_here
DBSEC_APP_SECRET=your_actual_app_secret_here
```

## 📊 지원하는 데이터 타입

### 국내주식 데이터
- **주가 정보**: 시가, 고가, 저가, 현재가, 거래량
- **호가 정보**: 매수/매도 호가 1~5단계, 잔량 정보
- **일봉 데이터**: 과거 주가 데이터 (최대 100일)
- **시장**: KOSPI, KOSDAQ, KONEX

### 해외주식 데이터
- **주가 정보**: 현재가, 변동률, 거래량
- **거래소**: NASDAQ, NYSE, AMEX
- **일봉 데이터**: 과거 주가 데이터

### 지수 데이터
- **KOSPI**: 0001
- **KOSDAQ**: 1001  
- **KOSPI200**: 1028

## 🚀 사용 방법

### 1. 서비스 실행

```bash
# Docker Compose로 전체 시스템 실행
docker-compose up -d

# 또는 개발 모드 실행
go run main.go
```

### 2. 주요 종목 초기화

```bash
# 주요 종목 자동 등록
curl -X POST http://localhost:8080/api/v1/admin/initialize/major-stocks
```

### 3. 데이터 수집 트리거

```bash
# 특정 종목 데이터 수집
curl -X POST http://localhost:8080/api/v1/admin/collect/005930

# 전체 종목 데이터 수집
curl -X POST http://localhost:8080/api/v1/admin/collect/all
```

### 4. API 상태 확인

```bash
# DB증권 API 연결 상태 확인
curl http://localhost:8080/api/v1/admin/api-status
```

## 📋 API 엔드포인트

### 관리자 API

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET | `/api/v1/admin/api-status` | API 연결 상태 확인 |
| POST | `/api/v1/admin/initialize/major-stocks` | 주요 종목 초기화 |
| POST | `/api/v1/admin/collect/:symbol` | 특정 종목 데이터 수집 |
| POST | `/api/v1/admin/collect/all` | 전체 종목 데이터 수집 |
| GET | `/api/v1/admin/database/stats` | 데이터베이스 통계 |

### 일반 API

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET | `/api/v1/stocks` | 전체 종목 목록 |
| GET | `/api/v1/stocks/:symbol` | 특정 종목 정보 |
| GET | `/api/v1/stocks/:symbol/price` | 최신 주가 정보 |
| GET | `/api/v1/stocks/:symbol/indicators` | 기술지표 정보 |

## 🎯 주요 종목

### 국내주식 (KR)
- **삼성전자** (005930)
- **SK하이닉스** (000660)
- **LG화학** (051910)
- **NAVER** (035420)
- **삼성SDI** (006400)
- **카카오** (035720)

### 미국주식 (US)
- **Apple** (AAPL)
- **Microsoft** (MSFT)
- **Google** (GOOGL)
- **Amazon** (AMZN)
- **Tesla** (TSLA)
- **Meta** (META)

## 🔄 데이터 플로우

```
1. DB증권 API → 실시간 데이터 수집
2. 데이터 검증 및 파싱
3. PostgreSQL 저장 (파티셔닝)
4. Redis 캐싱 (성능 최적화)
5. 기술지표 계산 (Go 엔진)
6. AI 분석 요청 (Python 서비스)
7. 매매 신호 생성
8. RabbitMQ 알림 전송
```

## ⚠️ 제한사항 및 주의사항

### API 제한
- **Rate Limit**: 초당 20요청
- **일일 한도**: 10,000요청 (무료 계정 기준)
- **토큰 유효기간**: 24시간

### 데이터 수집 주기
- **실시간 데이터**: 5분마다 자동 수집
- **일봉 데이터**: 일 1회 (장 마감 후)
- **주요 지수**: 1분마다 수집

### 오류 처리
- API 키 없음: Mock 데이터 사용
- Rate limit 초과: 자동 대기 후 재시도
- 토큰 만료: 자동 재인증
- 네트워크 오류: 3회 재시도 후 실패 처리

## 🛠 개발자 가이드

### 새로운 종목 추가

```bash
curl -X POST http://localhost:8080/api/v1/admin/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "066570",
    "name": "LG전자",
    "market": "KR",
    "exchange": "KOSPI",
    "sector": "Technology"
  }'
```

### 커스텀 데이터 수집

```go
// 특정 종목 30일 일봉 데이터 수집
err := dataCollector.CollectDailyData("005930", 30)
```

### API 클라이언트 직접 사용

```go
// DB증권 API 클라이언트 생성
apiClient := client.NewDBSecClient(cfg)

// 현재가 조회
price, err := apiClient.GetDomesticStockPrice("005930")

// 호가 정보 조회
asking, err := apiClient.GetDomesticStockAskingPrice("005930")
```

## 📊 모니터링

### 로그 확인

```bash
# API 호출 로그
docker-compose logs data-collector

# 데이터베이스 연결 상태
docker-compose logs postgres

# 전체 시스템 상태
docker-compose ps
```

### 성능 메트릭

```bash
# 데이터베이스 통계
curl http://localhost:8080/api/v1/admin/database/stats

# API 응답 시간 확인
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/api/v1/stocks/005930/price
```

## 🔒 보안 고려사항

1. **API 키 보안**
   - `.env` 파일을 Git에 커밋하지 마세요
   - 프로덕션에서는 환경변수 또는 시크릿 관리 도구 사용

2. **네트워크 보안**
   - HTTPS 사용 (프로덕션)
   - API 키 로테이션 정기 실행

3. **데이터 무결성**
   - 중복 데이터 체크
   - 데이터 검증 및 sanitization

## 🚨 문제 해결

### 일반적인 문제

#### 1. 인증 실패
```bash
# API 키 확인
echo $DBSEC_APP_KEY
echo $DBSEC_APP_SECRET

# 설정 재확인
curl http://localhost:8080/api/v1/admin/api-status
```

#### 2. 데이터 수집 실패
```bash
# 로그 확인
docker-compose logs data-collector | grep ERROR

# 수동 수집 시도
curl -X POST http://localhost:8080/api/v1/admin/collect/005930
```

#### 3. Rate Limit 초과
- 수집 주기 조정 (5분 → 10분)
- 종목 수 제한
- 병렬 처리 제한

## 📈 성능 최적화

### 데이터베이스 최적화
- 월별 파티셔닝 적용
- 인덱스 최적화
- 연결 풀 설정

### 캐싱 전략
- Redis 캐싱 (TTL: 1분)
- CDN 적용 (정적 데이터)
- 응답 압축

### API 호출 최적화
- 배치 처리
- 비동기 호출
- 중요도별 우선순위

이 통합을 통해 실제 시장 데이터를 기반으로 한 정확한 투자 추천 서비스를 제공할 수 있습니다.