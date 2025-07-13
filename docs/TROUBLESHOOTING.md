# 🛠️ 문제 해결 가이드

## 개요

주식 투자 추천 시스템 운영 중 발생할 수 있는 일반적인 문제들과 해결 방법을 정리한 가이드입니다.

## 🚨 시스템 시작 관련 문제

### 1. Docker 컨테이너가 시작되지 않는 경우

#### 증상
```bash
docker-compose up -d
# 일부 서비스가 실행되지 않음
```

#### 진단
```bash
# 컨테이너 상태 확인
docker-compose ps

# 로그 확인
docker-compose logs [service-name]

# 시스템 리소스 확인
docker stats
```

#### 해결 방법

**메모리 부족**
```bash
# 불필요한 컨테이너 정리
docker system prune -f

# 메모리 사용량 확인
free -h

# 스왑 메모리 확인
swapon -s
```

**포트 충돌**
```bash
# 포트 사용 확인
netstat -tlnp | grep :5432
netstat -tlnp | grep :6379
netstat -tlnp | grep :5672

# 충돌하는 프로세스 종료
sudo kill -9 [PID]
```

**Docker 권한 문제**
```bash
# Docker 그룹에 사용자 추가
sudo usermod -aG docker $USER
newgrp docker

# Docker 서비스 재시작
sudo systemctl restart docker
```

### 2. 환경 변수 설정 오류

#### 증상
```
Error: failed to authenticate with DBSec API
Database connection failed
```

#### 진단
```bash
# .env 파일 존재 확인
ls -la .env

# 환경 변수 확인
docker-compose exec backend env | grep DB
docker-compose exec backend env | grep DBSEC
```

#### 해결 방법

**.env 파일 생성**
```bash
# 예시 파일 복사
cp .env.example .env

# 필수 변수 설정
vim .env
```

**필수 환경 변수 확인**
```bash
# .env 파일 내용 예시
DB_HOST=postgres
DB_PORT=5432
DB_USER=stockuser
DB_PASSWORD=stockpass
DB_NAME=stockdb

DBSEC_APP_KEY=your_actual_key
DBSEC_APP_SECRET=your_actual_secret
```

## 🗄️ 데이터베이스 관련 문제

### 1. PostgreSQL 연결 실패

#### 증상
```
failed to connect to database: connection refused
```

#### 진단
```bash
# PostgreSQL 컨테이너 상태 확인
docker-compose logs postgres

# 네트워크 연결 확인
docker-compose exec backend ping postgres

# 포트 확인
docker-compose exec postgres netstat -tlnp | grep 5432
```

#### 해결 방법

**컨테이너 재시작**
```bash
docker-compose restart postgres
```

**데이터베이스 초기화**
```bash
# 주의: 모든 데이터가 삭제됩니다
docker-compose down -v
docker-compose up -d postgres
```

**수동 연결 테스트**
```bash
# PostgreSQL 컨테이너 내부 접속
docker-compose exec postgres psql -U stockuser -d stockdb

# 테이블 확인
\dt

# 연결 종료
\q
```

### 2. 마이그레이션 실패

#### 증상
```
failed to migrate database: table already exists
failed to migrate database: column does not exist
```

#### 해결 방법

**수동 마이그레이션**
```bash
# 백엔드 컨테이너에서 실행
docker-compose exec backend go run main.go migrate

# 또는 직접 SQL 실행
docker-compose exec postgres psql -U stockuser -d stockdb -f /sql/schema.sql
```

**스키마 리셋 (개발 환경)**
```bash
# 데이터베이스 삭제 후 재생성
docker-compose exec postgres dropdb -U stockuser stockdb
docker-compose exec postgres createdb -U stockuser stockdb
docker-compose restart backend
```

### 3. 파티션 관련 오류

#### 증상
```
failed to create partitions
partition does not exist
```

#### 해결 방법

**파티션 수동 생성**
```sql
-- PostgreSQL에 접속하여 실행
docker-compose exec postgres psql -U stockuser -d stockdb

-- 현재 파티션 확인
SELECT schemaname, tablename 
FROM pg_tables 
WHERE tablename LIKE 'stock_prices_%';

-- 파티션 수동 생성 (예: 2024년 7월)
CREATE TABLE stock_prices_202407 PARTITION OF stock_prices
FOR VALUES FROM ('2024-07-01') TO ('2024-08-01');
```

## 🔌 API 연동 관련 문제

### 1. DB증권 API 인증 실패

#### 증상
```
authentication failed with status 401
API credentials not configured
```

#### 진단
```bash
# API 상태 확인
curl http://localhost:8080/api/v1/admin/api-status

# 환경 변수 확인
docker-compose exec backend env | grep DBSEC
```

#### 해결 방법

**API 키 확인**
```bash
# .env 파일에서 키 확인
grep DBSEC .env

# 키 유효성 확인 (DB증권 포털에서)
# https://openapi.dbsec.co.kr/
```

**수동 인증 테스트**
```bash
# 직접 인증 API 호출
curl -X POST https://openapi.dbsec.co.kr/oauth2/tokenP \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials&appkey=${DBSEC_APP_KEY}&appsecret=${DBSEC_APP_SECRET}"
```

### 2. Rate Limit 초과

#### 증상
```
API request failed with status 429
rate limit exceeded
```

#### 해결 방법

**수집 주기 조정**
```go
// data_collector.go에서 수정
ticker := time.NewTicker(10 * time.Minute) // 5분 → 10분으로 변경
```

**요청 우선순위 조정**
```bash
# 중요한 종목만 수집
curl -X PUT http://localhost:8080/api/v1/admin/stocks/005930/status \
  -H "Content-Type: application/json" \
  -d '{"is_active": true}'

# 불필요한 종목 비활성화
curl -X PUT http://localhost:8080/api/v1/admin/stocks/OTHER/status \
  -H "Content-Type: application/json" \
  -d '{"is_active": false}'
```

### 3. Mock 데이터 사용 확인

#### Mock 모드 여부 확인
```bash
# 로그에서 Mock 데이터 사용 확인
docker-compose logs data-collector | grep "mock"
docker-compose logs data-collector | grep "API credentials not available"
```

#### 실제 API로 전환
```bash
# .env 파일에 올바른 API 키 설정
DBSEC_APP_KEY=your_real_key
DBSEC_APP_SECRET=your_real_secret

# 서비스 재시작
docker-compose restart backend data-collector
```

## 🧮 서비스 통신 문제

### 1. AI 서비스 연결 실패

#### 증상
```
failed to connect to AI service
AI service unavailable, using rule-based fallback
```

#### 진단
```bash
# AI 서비스 상태 확인
docker-compose logs ai-service

# 네트워크 연결 확인
docker-compose exec backend ping ai-service

# AI 서비스 헬스 체크
curl http://localhost:8001/health
```

#### 해결 방법

**AI 서비스 재시작**
```bash
docker-compose restart ai-service
```

**수동 연결 테스트**
```bash
# AI 의사결정 API 직접 호출
curl -X POST http://localhost:8001/api/v1/decision \
  -H "Content-Type: application/json" \
  -d '{
    "symbol": "005930",
    "market": "KR",
    "indicators": {"rsi": 65.4, "macd": 1250.5}
  }'
```

### 2. RabbitMQ 연결 문제

#### 증상
```
failed to connect to RabbitMQ
queue worker failed to start
```

#### 진단
```bash
# RabbitMQ 상태 확인
docker-compose logs rabbitmq

# 관리 UI 접속 (웹 브라우저)
http://localhost:15672
# 기본 계정: guest/guest
```

#### 해결 방법

**RabbitMQ 재시작**
```bash
docker-compose restart rabbitmq
```

**큐 상태 확인**
```bash
# RabbitMQ 컨테이너 내부에서 큐 확인
docker-compose exec rabbitmq rabbitmqctl list_queues
```

### 3. Redis 연결 문제

#### 증상
```
failed to connect to redis
cache service unavailable
```

#### 진단
```bash
# Redis 상태 확인
docker-compose logs redis

# 연결 테스트
docker-compose exec redis redis-cli ping
```

#### 해결 방법

**Redis 재시작**
```bash
docker-compose restart redis
```

**캐시 데이터 확인**
```bash
# Redis 데이터 확인
docker-compose exec redis redis-cli
> KEYS *
> GET stock:005930:price
> quit
```

## 📊 성능 관련 문제

### 1. 응답 속도 저하

#### 진단
```bash
# API 응답 시간 측정
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/api/v1/stocks/005930/price

# 시스템 리소스 확인
docker stats

# 데이터베이스 성능 확인
docker-compose exec postgres psql -U stockuser -d stockdb -c "
SELECT query, mean_time, calls 
FROM pg_stat_statements 
ORDER BY mean_time DESC LIMIT 5;"
```

#### 해결 방법

**인덱스 최적화**
```sql
-- 자주 사용되는 쿼리의 인덱스 확인
EXPLAIN ANALYZE SELECT * FROM stock_prices WHERE symbol = '005930' ORDER BY timestamp DESC LIMIT 1;

-- 필요시 인덱스 추가
CREATE INDEX CONCURRENTLY idx_stock_prices_symbol_timestamp ON stock_prices(symbol, timestamp DESC);
```

**캐시 설정 조정**
```go
// cache.go에서 TTL 조정
const (
    StockPriceTTL = 1 * time.Minute  // 캐시 유지 시간 조정
    IndicatorTTL  = 5 * time.Minute
)
```

### 2. 메모리 사용량 증가

#### 진단
```bash
# 컨테이너별 메모리 사용량
docker stats --format "table {{.Container}}\t{{.MemUsage}}\t{{.MemPerc}}"

# 호스트 메모리 상태
free -h
```

#### 해결 방법

**메모리 제한 설정**
```yaml
# docker-compose.yml에 추가
services:
  backend:
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M
```

**가비지 컬렉션 조정**
```bash
# Go 런타임 환경 변수 설정
export GOGC=100  # 기본값 100, 낮추면 더 자주 GC 실행
```

## 🔍 로그 및 디버깅

### 1. 로그 레벨 조정

```bash
# 자세한 로그 출력
docker-compose exec backend go run main.go --log-level=debug

# 특정 서비스 로그만 확인
docker-compose logs -f backend
docker-compose logs -f ai-service | grep ERROR
```

### 2. 로그 파일 저장

```bash
# 로그를 파일로 저장
docker-compose logs > system.log 2>&1

# 특정 시간 범위의 로그
docker-compose logs --since="2024-07-13T15:00:00" --until="2024-07-13T16:00:00"
```

### 3. 디버그 모드 실행

```bash
# 개발 모드로 실행 (더 자세한 로그)
export GIN_MODE=debug
docker-compose up backend
```

## 🔄 데이터 복구

### 1. 백업에서 복구

```bash
# 데이터베이스 백업 복구
docker-compose exec -T postgres psql -U stockuser stockdb < backup_20240713.sql

# 특정 테이블만 복구
docker-compose exec -T postgres psql -U stockuser stockdb < stock_prices_backup.sql
```

### 2. 데이터 재수집

```bash
# 전체 데이터 재수집
curl -X POST http://localhost:8080/api/v1/admin/collect/all

# 특정 기간 일봉 데이터 재수집
# (현재는 API를 통해 지원하지 않음, 직접 코드 수정 필요)
```

## ⚡ 긴급 상황 대응

### 1. 전체 시스템 재시작

```bash
# 모든 서비스 중지
docker-compose down

# 볼륨 제외하고 재시작 (데이터 보존)
docker-compose up -d

# 필요시 볼륨까지 초기화 (주의: 모든 데이터 삭제)
docker-compose down -v && docker-compose up -d
```

### 2. 서비스별 독립 실행

```bash
# 데이터베이스만 실행
docker-compose up -d postgres redis

# 백엔드만 실행 (개발 모드)
go run main.go

# AI 서비스만 실행
cd ai && python main.py
```

### 3. 응급 헬스 체크

```bash
#!/bin/bash
# health_check.sh

echo "=== 시스템 헬스 체크 ==="

# 기본 서비스 상태
echo "1. 서비스 상태:"
docker-compose ps

# API 응답 확인
echo "2. API 응답:"
curl -s http://localhost:8080/health | jq .

# 데이터베이스 연결
echo "3. 데이터베이스:"
docker-compose exec postgres psql -U stockuser -d stockdb -c "SELECT 1;" >/dev/null 2>&1 && echo "OK" || echo "FAIL"

# 최신 데이터 확인
echo "4. 최신 데이터:"
curl -s http://localhost:8080/api/v1/stocks/005930/price | jq .price.timestamp

echo "=== 헬스 체크 완료 ==="
```

## 📞 추가 지원

### 로그 수집 스크립트

```bash
#!/bin/bash
# collect_logs.sh

echo "시스템 상태 수집 중..."

# 시스템 정보
echo "=== 시스템 정보 ===" > system_info.txt
docker --version >> system_info.txt
docker-compose --version >> system_info.txt
free -h >> system_info.txt
df -h >> system_info.txt

# 컨테이너 상태
echo "=== 컨테이너 상태 ===" > container_status.txt
docker-compose ps >> container_status.txt
docker stats --no-stream >> container_status.txt

# 로그 수집
docker-compose logs > application_logs.txt

echo "로그 수집 완료: system_info.txt, container_status.txt, application_logs.txt"
```

### 자주 사용하는 명령어 모음

```bash
# 빠른 재시작
alias restart-stock="docker-compose restart backend data-collector"

# 로그 실시간 모니터링
alias watch-logs="docker-compose logs -f backend ai-service"

# API 상태 확인
alias check-api="curl -s http://localhost:8080/health | jq ."

# 데이터베이스 접속
alias db-connect="docker-compose exec postgres psql -U stockuser -d stockdb"
```

---

**📅 문서 업데이트**: 2024년 7월 13일  
**🔗 관련 문서**: [배포 가이드](DEPLOYMENT_GUIDE.md), [API 문서](API_DOCUMENTATION.md)