# 주식 투자 추천 프로그램

AI 기반 주식 투자 추천 시스템으로, 한국 및 미국 지수 종목을 대상으로 기술지표 분석과 뉴스 감성 분석을 통해 매매 신호를 제공합니다.

## 🚀 주요 기능

- **실시간 주가 데이터 수집** (DB증권 Open API)
- **10가지 이상 기술지표 분석** (RSI, MACD, Bollinger Bands 등)
- **AI 기반 매매 신호 생성** (Interface 설계 완료, ML 모델 확장 가능)
- **뉴스 크롤링 및 감성 분석** (네이버/다음 금융뉴스)
- **실시간 알림 시스템** (RabbitMQ 기반)
- **고성능 캐싱** (Redis)

## 🏗 시스템 아키텍처

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   Backend       │
│   (React)       │◄──►│   (Go/Gin)      │◄──►│   Services      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                       ┌─────────────────┐            │
                       │   Message Queue │◄───────────┘
                       │   (RabbitMQ)    │
                       └─────────────────┘
                                │
        ┌─────────────────┬─────────────────┬─────────────────┐
        │   Data Collector│   AI Service    │   News Crawler  │
        │   (DB증권 API)   │   (Python)      │   (뉴스/감성)    │
        └─────────────────┴─────────────────┴─────────────────┘
                                │
                    ┌─────────────────────────────────┐
                    │          Database               │
                    │  ┌─────────────┬─────────────┐  │
                    │  │ PostgreSQL  │   Redis     │  │
                    │  │ (시계열)     │  (캐시)      │  │
                    │  └─────────────┴─────────────┘  │
                    └─────────────────────────────────┘
```

## 🛠 기술 스택

### Backend
- **Go** - 고성능 백엔드 서버
- **Gin** - 웹 프레임워크
- **GORM** - ORM
- **PostgreSQL** - 시계열 데이터베이스 (파티셔닝)
- **Redis** - 캐싱
- **RabbitMQ** - 메시지 큐

### AI/ML
- **Python** - AI 서비스
- **FastAPI** - API 프레임워크
- **Mock Decision Engine** - 규칙 기반 매매 결정 (ML 확장 가능)

### 크롤링
- **Python** - 웹 크롤러
- **aiohttp** - 비동기 HTTP 클라이언트
- **BeautifulSoup** - HTML 파싱
- **감성 분석** - 뉴스 텍스트 분석

### Frontend (예정)
- **React** - 사용자 인터페이스
- **TypeScript** - 타입 안전성
- **Material-UI** - UI 컴포넌트

## 📊 구현된 기술지표

1. **RSI** (Relative Strength Index)
2. **MACD** (Moving Average Convergence Divergence)
3. **SMA 20/50** (Simple Moving Average)
4. **EMA 12/26** (Exponential Moving Average)
5. **Bollinger Bands** (상단/중간/하단)
6. **Stochastic Oscillator** (K, D)
7. **Williams %R**
8. **ATR** (Average True Range)
9. **OBV** (On-Balance Volume)

## 🚀 빠른 시작

### 사전 요구사항
- Docker & Docker Compose
- Go 1.21+
- Python 3.11+
- DB증권 Open API 키

### 1. 환경 설정
```bash
# 환경 변수 파일 생성
cp .env.example .env

# API 키 설정
vim .env  # DBSEC_API_KEY 설정
```

### 2. Docker 컨테이너 실행
```bash
# 데이터베이스 서비스 시작
docker-compose up -d postgres redis rabbitmq

# 백엔드 서비스 빌드 및 실행
docker-compose up -d backend data-collector

# AI 서비스 실행
docker-compose up -d ai-service

# 크롤러 실행
docker-compose up -d crawler
```

### 3. 개발 모드 실행 (로컬)
```bash
# 백엔드 실행
go run main.go

# AI 서비스 실행
cd ai && python main.py

# 크롤러 실행
cd crawler && python main.py
```

### 4. DB증권 API 연동 설정
```bash
# .env 파일에 API 키 설정
DBSEC_APP_KEY=your_dbsec_app_key
DBSEC_APP_SECRET=your_dbsec_app_secret

# 주요 종목 초기화
curl -X POST http://localhost:8080/api/v1/admin/initialize/major-stocks

# API 연결 상태 확인
curl http://localhost:8080/api/v1/admin/api-status
```

## 📚 상세 문서

- **📊 [프로젝트 진행 상황](docs/PROJECT_STATUS.md)** - 전체 작업 현황 및 완료 상태
- **🏗️ [기술 아키텍처](docs/TECHNICAL_ARCHITECTURE.md)** - 시스템 설계 및 구조 상세 설명
- **🚀 [배포 가이드](docs/DEPLOYMENT_GUIDE.md)** - Docker 기반 배포 및 운영 가이드
- **🔌 [DB증권 API 연동](docs/DB_SECURITIES_API_INTEGRATION.md)** - 실시간 데이터 수집 API 가이드
- **📡 [API 문서](docs/API_DOCUMENTATION.md)** - REST API 엔드포인트 상세 문서
- **🛠️ [문제 해결 가이드](docs/TROUBLESHOOTING.md)** - 일반적인 문제 및 해결 방법

## 📡 주요 API 엔드포인트

### 🏥 시스템 상태
- `GET /health` - 헬스 체크

### 📈 주식 정보
- `GET /api/v1/stocks` - 종목 목록
- `GET /api/v1/stocks/{symbol}` - 종목 상세 정보
- `GET /api/v1/stocks/{symbol}/price` - 실시간 주가
- `GET /api/v1/stocks/{symbol}/indicators` - 기술지표

### 🎯 매매 신호
- `GET /api/v1/signals` - 전체 매매 신호
- `GET /api/v1/signals/{symbol}` - 종목별 신호

### 🔧 관리자 API
- `POST /api/v1/admin/stocks` - 종목 등록
- `POST /api/v1/admin/collect/{symbol}` - 데이터 수집 트리거
- `GET /api/v1/admin/api-status` - API 연결 상태
- `GET /api/v1/admin/database/stats` - 시스템 통계

> **📋 상세한 API 문서는 [API Documentation](docs/API_DOCUMENTATION.md)을 참조하세요.**

## 🤖 AI 서비스 API

- `POST http://localhost:8001/api/v1/decision` - AI 매매 결정 요청
- `GET http://localhost:8001/api/v1/models/status` - 모델 상태 확인
- `GET http://localhost:8001/health` - AI 서비스 헬스 체크

## 📈 데이터 플로우

1. **데이터 수집**: DB증권 API → 주가 데이터
2. **지표 계산**: Go 엔진 → 10가지 기술지표
3. **AI 분석**: Python 서비스 → 매매 결정
4. **신호 생성**: 종합 분석 → BUY/SELL/HOLD
5. **캐싱**: Redis → 성능 최적화
6. **알림**: RabbitMQ → 실시간 알림

## 📁 프로젝트 구조

```
stock-recommender/
├── backend/                 # Go 백엔드
│   ├── config/             # 설정 관리
│   ├── database/           # DB 연결
│   ├── handlers/           # HTTP 핸들러
│   ├── models/             # 데이터 모델
│   └── services/           # 비즈니스 로직
├── ai/                     # Python AI 서비스
│   ├── interface/          # API 인터페이스
│   ├── models/             # 데이터 모델
│   └── main.py
├── crawler/                # 뉴스 크롤러
│   ├── crawlers/           # 크롤링 모듈
│   ├── sentiment/          # 감성 분석
│   └── database/           # DB 연결
├── sql/                    # 데이터베이스 스키마
├── docs/                   # 문서
└── docker-compose.yml      # 컨테이너 설정
```

## 🔧 개발 정보

### 코드 스타일
- **Go**: gofmt, golint 준수
- **Python**: Black, flake8 준수
- **SQL**: PostgreSQL 표준 준수

### 테스트 실행
```bash
# Go 백엔드 테스트
go test ./...

# Python AI 서비스 테스트
cd ai && pytest

# Python 크롤러 테스트
cd crawler && pytest
```

> **🔧 상세한 개발 환경 설정은 [기술 아키텍처](docs/TECHNICAL_ARCHITECTURE.md)를 참조하세요.**

## ✅ 구현 완료된 기능

### 백엔드 시스템 (완료)
- [x] **DB증권 Open API 통합** - 실시간 주가 데이터 수집
- [x] **10가지 기술지표 분석** - RSI, MACD, Bollinger Bands 등
- [x] **PostgreSQL 시계열 DB** - 월별 파티셔닝 적용
- [x] **Redis 캐싱 시스템** - 성능 최적화
- [x] **RabbitMQ 메시지 큐** - 비동기 처리
- [x] **AI 의사결정 인터페이스** - Python FastAPI 서비스
- [x] **매매 신호 생성 엔진** - 규칙 기반 + AI 결합
- [x] **뉴스 크롤링 & 감성 분석** - 한국어 뉴스 처리
- [x] **Docker 컨테이너 환경** - 완전 자동화 배포

### 데이터 수집 (완료)
- [x] **국내주식**: KOSPI, KOSDAQ 주요 종목
- [x] **해외주식**: NASDAQ, NYSE 주요 종목  
- [x] **주요 지수**: KOSPI, KOSDAQ, KOSPI200
- [x] **실시간 호가**: 매수/매도 5단계 호가 정보
- [x] **일봉 차트**: 과거 데이터 수집 및 저장

## 📋 TODO 및 확장 계획

### AI/ML 고도화 (우선순위 높음)
- [ ] LSTM 주가 예측 모델
- [ ] Random Forest 분류 모델
- [ ] Ensemble 의사결정
- [ ] 실시간 모델 업데이트
- [ ] 백테스팅 시스템

### 추가 기능 (우선순위 중간)
- [ ] 웹 프론트엔드 (React)
- [ ] 실시간 알림 (이메일, 카카오톡)
- [ ] 포트폴리오 관리
- [ ] 리스크 관리 도구

### 성능 최적화 (우선순위 낮음)
- [ ] 분산 처리 시스템
- [ ] 더 많은 데이터 소스
- [ ] 고주파 거래 지원

---

## 📖 문서 목차

### 🚀 시작하기
- **[README.md](README.md)** - 프로젝트 개요 및 빠른 시작
- **[배포 가이드](docs/DEPLOYMENT_GUIDE.md)** - Docker 기반 설치 및 배포

### 📊 프로젝트 정보
- **[프로젝트 진행 상황](docs/PROJECT_STATUS.md)** - 작업 현황 및 완료도
- **[기술 아키텍처](docs/TECHNICAL_ARCHITECTURE.md)** - 시스템 설계 및 구조

### 🔧 개발 및 운영
- **[DB증권 API 연동](docs/DB_SECURITIES_API_INTEGRATION.md)** - 실시간 데이터 수집
- **[API 문서](docs/API_DOCUMENTATION.md)** - REST API 상세 명세
- **[문제 해결 가이드](docs/TROUBLESHOOTING.md)** - 오류 진단 및 해결

### 💡 추가 정보
- **환경 설정**: `.env.example` 파일 참조
- **Docker 설정**: `docker-compose.yml` 파일 참조
- **데이터베이스 스키마**: `sql/` 디렉토리 참조

---

**🏗️ 개발**: AI 기반 주식 투자 추천 시스템 | **📅 업데이트**: 2024년 7월 13일
