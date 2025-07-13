# 수정된 개발 계획

## AI/ML 범위 조정 및 우선순위 재정의

### AI/ML 역할 재정의
- **기술지표 계산**: Go로 구현 (RSI, MACD, Bollinger Bands 등)
- **AI/ML 역할**: 기술지표 값들을 종합하여 최종 매매 판단
- **구현 방식**: Interface 중심 설계, 추후 실제 모델 구현

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Stock Data    │───►│  Technical      │───►│   AI Decision   │
│   (Raw OHLCV)   │    │  Indicators     │    │   Engine        │
│                 │    │  (Go)           │    │   (Python)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                        │
                                                        ▼
                                               ┌─────────────────┐
                                               │ Trading Signal  │
                                               │ (BUY/SELL/HOLD) │
                                               └─────────────────┘
```

### 개발 단계별 우선순위

#### 1단계: 핵심 백엔드 인프라 (최우선 - High)
1. **개발 환경 초기 설정**
   - Docker Compose 구성
   - Go modules 초기화
   - 프로젝트 디렉토리 구조

2. **PostgreSQL 스키마 설계 및 구현**
   - 시계열 데이터 파티셔닝
   - 기술지표 저장 구조
   - 인덱스 최적화

3. **Redis 캐싱 전략 구현**
   - 실시간 데이터 캐싱
   - 지표 결과 캐싱
   - TTL 전략

4. **DB증권 Open API 연동 모듈**
   - API 클라이언트 구현
   - Rate limiting
   - 데이터 수집 스케줄러

5. **주식 기술 지표 분석 모듈 (Go)**
   - 10가지 이상 기술지표 구현
   - 실시간 계산 엔진
   - 결과 저장 로직

#### 2단계: 메시징 및 AI 인터페이스 (Medium)
6. **RabbitMQ 메시지 큐 시스템**
   - Exchange/Queue 설계
   - 워커 패턴 구현
   - 메시지 신뢰성 보장

7. **AI/ML Interface 설계 및 구현**
   - Python 서비스 인터페이스
   - 기술지표 데이터 전달 구조
   - 판단 결과 반환 형식
   - **TODO**: 실제 ML 모델 구현

8. **매매 신호 생성 알고리즘**
   - AI 판단 결과 해석
   - 리스크 관리 로직
   - 신호 강도 계산

#### 3단계: 부가 기능 (Medium)
9. **웹 크롤링 모듈**
   - 뉴스 수집
   - 감성 분석 (기본)
   - **TODO**: 고도화된 NLP 분석

#### 4단계: 사용자 기능 (Low)
10. **알림 시스템**
    - 실시간 알림
    - 다중 채널 지원

11. **사용자 인터페이스 (React)**
    - 대시보드
    - 차트 컴포넌트
    - 실시간 업데이트

12. **백테스팅 시스템**
    - 과거 데이터 검증
    - 성과 지표 계산

#### 5단계: 최종 검증 (Medium)
13. **통합 테스트 및 성능 최적화**
    - 부하 테스트
    - 성능 튜닝

## AI/ML 인터페이스 설계

### 데이터 플로우
```go
// Go에서 계산한 기술지표
type TechnicalIndicators struct {
    Symbol    string                 `json:"symbol"`
    Timestamp time.Time              `json:"timestamp"`
    Price     StockPrice            `json:"price"`
    Indicators map[string]float64   `json:"indicators"`
}

// AI 서비스로 전송할 데이터
type AIDecisionRequest struct {
    Symbol      string              `json:"symbol"`
    Market      string              `json:"market"`
    Indicators  TechnicalIndicators `json:"indicators"`
    NewsScore   float64             `json:"news_score,omitempty"`
}

// AI 서비스에서 받을 응답
type AIDecisionResponse struct {
    Symbol      string    `json:"symbol"`
    Decision    string    `json:"decision"`    // BUY/SELL/HOLD
    Confidence  float64   `json:"confidence"`  // 0.0 ~ 1.0
    Reasoning   []string  `json:"reasoning"`
    Timestamp   time.Time `json:"timestamp"`
}
```

### Python AI Service Interface
```python
# ai/interface/decision_engine.py
from abc import ABC, abstractmethod
from typing import Dict, List
from dataclasses import dataclass

@dataclass
class TechnicalIndicators:
    symbol: str
    timestamp: str
    price: Dict[str, float]
    indicators: Dict[str, float]

@dataclass
class AIDecision:
    symbol: str
    decision: str  # BUY/SELL/HOLD
    confidence: float  # 0.0 ~ 1.0
    reasoning: List[str]
    timestamp: str

class DecisionEngine(ABC):
    """AI 결정 엔진 인터페이스"""
    
    @abstractmethod
    def make_decision(self, indicators: TechnicalIndicators) -> AIDecision:
        """기술지표를 바탕으로 매매 결정"""
        pass
    
    @abstractmethod
    def update_model(self, training_data: List[Dict]) -> bool:
        """모델 업데이트 (TODO: 실제 구현)"""
        pass

class MockDecisionEngine(DecisionEngine):
    """개발용 Mock 엔진"""
    
    def make_decision(self, indicators: TechnicalIndicators) -> AIDecision:
        # 간단한 규칙 기반 로직 (임시)
        rsi = indicators.indicators.get('rsi', 50)
        macd = indicators.indicators.get('macd', 0)
        
        if rsi < 30 and macd > 0:
            decision = "BUY"
            confidence = 0.7
            reasoning = ["RSI oversold", "MACD positive"]
        elif rsi > 70 and macd < 0:
            decision = "SELL"
            confidence = 0.7
            reasoning = ["RSI overbought", "MACD negative"]
        else:
            decision = "HOLD"
            confidence = 0.5
            reasoning = ["No clear signal"]
        
        return AIDecision(
            symbol=indicators.symbol,
            decision=decision,
            confidence=confidence,
            reasoning=reasoning,
            timestamp=indicators.timestamp
        )
    
    def update_model(self, training_data: List[Dict]) -> bool:
        # TODO: 실제 ML 모델 학습 구현
        print("Model update scheduled for future implementation")
        return True

# TODO: 실제 ML 모델 구현
class MLDecisionEngine(DecisionEngine):
    """실제 ML 모델 기반 엔진 (미구현)"""
    
    def __init__(self):
        # TODO: 모델 로드
        pass
    
    def make_decision(self, indicators: TechnicalIndicators) -> AIDecision:
        # TODO: LSTM, Random Forest 등 실제 ML 모델 적용
        raise NotImplementedError("ML model implementation pending")
    
    def update_model(self, training_data: List[Dict]) -> bool:
        # TODO: 실제 모델 재학습
        raise NotImplementedError("Model training implementation pending")
```

### Go-Python 통신 구조
```go
// services/ai_client.go
type AIClient struct {
    baseURL string
    client  *http.Client
}

func (c *AIClient) GetDecision(indicators TechnicalIndicators) (*AIDecisionResponse, error) {
    request := AIDecisionRequest{
        Symbol:     indicators.Symbol,
        Market:     "KR", // or "US"
        Indicators: indicators,
    }
    
    // HTTP POST to Python AI service
    resp, err := c.client.Post(
        c.baseURL+"/api/v1/decision",
        "application/json",
        bytes.NewBuffer(jsonBytes),
    )
    
    var decision AIDecisionResponse
    err = json.NewDecoder(resp.Body).Decode(&decision)
    return &decision, err
}
```

## 구현 TODO 목록

### AI/ML 관련 TODO
```markdown
### Phase 1: Interface Implementation (현재)
- [x] AI service interface 설계
- [x] Mock decision engine 구현
- [ ] Go-Python HTTP 통신 구현
- [ ] 기본 규칙 기반 로직 구현

### Phase 2: ML Model Development (미래)
- [ ] 데이터 수집 및 전처리 파이프라인
- [ ] Feature engineering (기술지표 조합)
- [ ] LSTM 모델 구현 (시계열 예측)
- [ ] Random Forest 모델 구현 (분류)
- [ ] Ensemble 모델 구현
- [ ] 백테스팅 및 모델 검증
- [ ] 실시간 모델 업데이트 로직
- [ ] A/B 테스트 프레임워크

### Phase 3: Advanced Features (미래)
- [ ] 감성 분석 통합
- [ ] 뉴스 임베딩 활용
- [ ] 리인포스먼트 러닝 적용
- [ ] 실시간 모델 성능 모니터링
```

## 문서화 전략

### 개발 진행 상황 추적
1. **일일 개발 로그**: `docs/dev-logs/YYYY-MM-DD.md`
2. **API 문서**: `docs/api/` 디렉토리에 OpenAPI 스펙
3. **아키텍처 결정 기록**: `docs/adr/` (Architecture Decision Records)
4. **성능 테스트 결과**: `docs/performance/`

### 문서 템플릿
```markdown
# 개발 로그 - YYYY-MM-DD

## 완료된 작업
- [ ] 작업 항목 1
- [ ] 작업 항목 2

## 진행 중인 작업
- 현재 작업 내용

## 발견된 이슈
- 이슈 설명 및 해결 방안

## 다음 단계
- 예정된 작업

## 코드 변경사항
- 주요 파일 수정 내역
- 새로 추가된 기능

## 테스트 결과
- 단위 테스트 결과
- 통합 테스트 결과
```

## 즉시 시작할 작업

1. **개발 환경 설정**
   - Docker Compose 파일 작성
   - Go modules 초기화
   - 프로젝트 구조 생성

2. **데이터베이스 스키마 구현**
   - PostgreSQL 테이블 설계
   - 마이그레이션 스크립트

3. **기본 API 서버 구축**
   - Gin 웹 서버 설정
   - 기본 라우팅

이제 개발 환경 설정부터 시작하겠습니다.