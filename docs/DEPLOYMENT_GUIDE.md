# 🚀 주식 투자 추천 프로그램 배포 가이드

## 📋 시스템 요구사항

### 최소 하드웨어 요구사항
- **CPU**: 2 코어 이상
- **메모리**: 4GB RAM 이상
- **저장공간**: 20GB 이상 (데이터베이스 포함)
- **네트워크**: 인터넷 연결 (API 호출용)

### 소프트웨어 요구사항
- **Docker**: 20.10 이상
- **Docker Compose**: 2.0 이상
- **Git**: 버전 관리용

### 필수 API 키
- **DB증권 Open API 키**: [https://openapi.dbsec.co.kr/](https://openapi.dbsec.co.kr/)에서 신청

## 🔧 설치 및 설정

### 1. 소스코드 다운로드
```bash
# Git 클론
git clone <repository-url> stock-recommender
cd stock-recommender

# 또는 압축파일 다운로드 후 압축 해제
unzip stock-recommender.zip
cd stock-recommender
```

### 2. 환경 설정
```bash
# 환경 변수 파일 생성
cp .env.example .env

# 환경 변수 설정 (중요!)
vim .env
```

#### 필수 환경 변수 설정
```bash
# .env 파일 내용

# 데이터베이스 설정
DB_HOST=postgres
DB_PORT=5432
DB_USER=stockuser
DB_PASSWORD=stockpass
DB_NAME=stockdb

# Redis 설정
REDIS_HOST=redis
REDIS_PORT=6379

# RabbitMQ 설정
RABBITMQ_HOST=rabbitmq
RABBITMQ_PORT=5672
RABBITMQ_USER=stockmq
RABBITMQ_PASS=stockmqpass

# 🔑 중요: DB증권 API 키 설정
DBSEC_APP_KEY=your_actual_api_key_here

# AI 서비스 설정
AI_SERVICE_URL=http://ai-service:8001

# 애플리케이션 설정
PORT=8080
GIN_MODE=release
```

### 3. 시스템 검증
```bash
# 테스트 스크립트 실행
./scripts/run-tests.sh

# Docker Compose 설정 검증
docker-compose config
```

## 🚀 서비스 실행

### 전체 시스템 실행 (권장)
```bash
# 모든 서비스 실행
docker-compose up -d

# 로그 확인
docker-compose logs -f
```

### 단계별 실행
```bash
# 1단계: 데이터베이스 서비스
docker-compose up -d postgres redis rabbitmq

# 2단계: 데이터베이스 초기화 확인
docker-compose logs postgres

# 3단계: 백엔드 서비스
docker-compose up -d backend data-collector

# 4단계: AI 및 크롤러 서비스  
docker-compose up -d ai-service crawler
```

## ✅ 서비스 상태 확인

### 헬스 체크
```bash
# 백엔드 API 상태 확인
curl http://localhost:8080/health

# AI 서비스 상태 확인
curl http://localhost:8001/health

# 모든 컨테이너 상태 확인
docker-compose ps
```

### 예상 응답
```json
// http://localhost:8080/health
{
  "status": "ok",
  "timestamp": "2024-07-13T15:30:00Z",
  "database": "connected",
  "version": "1.0.0"
}

// http://localhost:8001/health  
{
  "status": "healthy",
  "timestamp": "2024-07-13T15:30:00Z",
  "service": "ai-decision-service",
  "version": "1.0.0"
}
```

## 📊 API 사용 방법

### 기본 API 엔드포인트

#### 1. 종목 정보 조회
```bash
# 전체 종목 목록
curl http://localhost:8080/api/v1/stocks

# 특정 종목 정보
curl http://localhost:8080/api/v1/stocks/005930

# 최신 주가 정보
curl http://localhost:8080/api/v1/stocks/005930/price

# 기술지표 정보
curl http://localhost:8080/api/v1/stocks/005930/indicators
```

#### 2. 매매 신호 조회
```bash
# 전체 매매 신호
curl http://localhost:8080/api/v1/signals

# 특정 종목 매매 신호
curl http://localhost:8080/api/v1/signals/005930

# 신호 타입별 필터링
curl "http://localhost:8080/api/v1/signals?signal_type=BUY"

# 시장별 필터링
curl "http://localhost:8080/api/v1/signals?market=KR"
```

#### 3. AI 분석 요청 (직접)
```bash
curl -X POST http://localhost:8001/api/v1/decision \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "005930",
    "market": "KR",
    "indicators": {
      "rsi": 65.4,
      "macd": 1250.5,
      "sma_20": 70800.0
    }
  }'
```

## 🔧 관리자 기능

### 종목 등록
```bash
curl -X POST http://localhost:8080/api/v1/admin/stocks \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "TEST001",
    "name": "테스트 종목",
    "market": "KR",
    "exchange": "KOSPI",
    "sector": "Technology",
    "industry": "Software"
  }'
```

### 데이터 수집 트리거
```bash
# 특정 종목 데이터 수집
curl -X POST http://localhost:8080/api/v1/admin/collect/005930
```

## 📈 데이터 플로우 확인

### 1. 데이터 수집 확인
```bash
# 데이터 컬렉터 로그 확인
docker-compose logs data-collector

# 최신 주가 데이터 확인
curl http://localhost:8080/api/v1/stocks/005930/price
```

### 2. 지표 계산 확인
```bash
# 기술지표 확인
curl http://localhost:8080/api/v1/stocks/005930/indicators
```

### 3. AI 분석 결과 확인
```bash
# 매매 신호 확인
curl http://localhost:8080/api/v1/signals/005930
```

## 🗄️ 데이터베이스 관리

### 데이터베이스 접속
```bash
# PostgreSQL 컨테이너 접속
docker-compose exec postgres psql -U stockuser -d stockdb

# 주요 테이블 확인
\dt

# 데이터 확인 예시
SELECT * FROM stocks LIMIT 5;
SELECT * FROM stock_prices ORDER BY timestamp DESC LIMIT 10;
SELECT * FROM trading_signals ORDER BY created_at DESC LIMIT 5;
```

### 파티션 관리
```sql
-- 파티션 목록 확인
SELECT schemaname, tablename, partitionname 
FROM pg_partitions 
WHERE tablename = 'stock_prices';

-- 파티션 크기 확인  
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables 
WHERE tablename LIKE 'stock_prices_%' 
ORDER BY tablename;
```

## 🔍 문제 해결

### 일반적인 문제

#### 1. 컨테이너가 시작되지 않는 경우
```bash
# 컨테이너 상태 확인
docker-compose ps

# 로그 확인
docker-compose logs [service-name]

# 컨테이너 재시작
docker-compose restart [service-name]
```

#### 2. 데이터베이스 연결 실패
```bash
# PostgreSQL 컨테이너 상태 확인
docker-compose logs postgres

# 네트워크 확인
docker-compose exec backend ping postgres

# 포트 확인
netstat -tlnp | grep 5432
```

#### 3. API 키 관련 오류
```bash
# 환경 변수 확인
docker-compose exec backend env | grep DBSEC

# API 키가 없는 경우 (Mock 데이터 사용)
docker-compose logs data-collector | grep "mock"
```

#### 4. 메모리 부족
```bash
# 메모리 사용량 확인
docker stats

# 불필요한 컨테이너 정리
docker system prune -f
```

### 로그 파일 위치
```bash
# 컨테이너 로그
docker-compose logs [service-name]

# 영구 로그 (볼륨 마운트된 경우)
./logs/
```

## 🚀 성능 최적화

### 프로덕션 설정
```yaml
# docker-compose.prod.yml 예시
version: '3.8'
services:
  backend:
    deploy:
      replicas: 2
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
    environment:
      - GIN_MODE=release
      
  postgres:
    command: postgres -c max_connections=100 -c shared_buffers=256MB
    deploy:
      resources:
        limits:
          memory: 1G
```

### 성능 모니터링
```bash
# 리소스 사용량 모니터링
docker stats --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"

# 데이터베이스 성능 확인
docker-compose exec postgres psql -U stockuser -d stockdb -c "
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC LIMIT 5;"
```

## 📊 백업 및 복구

### 데이터베이스 백업
```bash
# 전체 백업
docker-compose exec postgres pg_dump -U stockuser stockdb > backup_$(date +%Y%m%d).sql

# 특정 테이블 백업
docker-compose exec postgres pg_dump -U stockuser -t stock_prices stockdb > prices_backup.sql
```

### 복구
```bash
# 백업 복구
docker-compose exec -T postgres psql -U stockuser stockdb < backup_20240713.sql
```

## 🔐 보안 고려사항

### 프로덕션 환경 보안
1. **API 키 보안**
   - `.env` 파일 권한 설정: `chmod 600 .env`
   - Git에 API 키 커밋 금지

2. **네트워크 보안**
   - 외부 접근 차단: PostgreSQL, Redis, RabbitMQ 포트
   - API 서버만 외부 노출

3. **컨테이너 보안**
   - 정기적인 이미지 업데이트
   - 불필요한 권한 제거

### 방화벽 설정 (선택사항)
```bash
# API 서버만 외부 접근 허용
sudo ufw allow 8080
sudo ufw deny 5432,6379,5672,15672
```

## 📞 지원 및 문의

### 로그 수집 (문의시 첨부)
```bash
# 전체 시스템 상태 수집
./scripts/collect-logs.sh > system-status.txt
```

### 유용한 명령어 모음
```bash
# 전체 재시작
docker-compose down && docker-compose up -d

# 특정 서비스 재시작  
docker-compose restart backend

# 로그 실시간 모니터링
docker-compose logs -f backend ai-service

# 데이터베이스 초기화 (주의!)
docker-compose down -v && docker-compose up -d
```

이 가이드를 따라하면 주식 투자 추천 시스템을 성공적으로 배포하고 운영할 수 있습니다.