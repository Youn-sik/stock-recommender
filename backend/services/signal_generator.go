package services

import (
	"encoding/json"
	"fmt"
	"log"
	"stock-recommender/backend/models"
	"time"

	"gorm.io/gorm"
)

type SignalGeneratorService struct {
	db               *gorm.DB
	indicatorService *IndicatorService
	aiClient         *AIClient
	cacheService     *CacheService
	queueService     *QueueService
}

func NewSignalGeneratorService(
	db *gorm.DB,
	indicatorService *IndicatorService,
	aiClient *AIClient,
	cacheService *CacheService,
	queueService *QueueService,
) *SignalGeneratorService {
	return &SignalGeneratorService{
		db:               db,
		indicatorService: indicatorService,
		aiClient:         aiClient,
		cacheService:     cacheService,
		queueService:     queueService,
	}
}

// 특정 종목에 대한 매매 신호 생성
func (s *SignalGeneratorService) GenerateSignal(symbol, market string) (*models.TradingSignal, error) {
	log.Printf("Generating signal for %s (%s)", symbol, market)

	// 1. 최근 주가 데이터 조회 (50일치)
	var prices []models.StockPrice
	err := s.db.Where("symbol = ? AND market = ?", symbol, market).
		Order("timestamp desc").
		Limit(50).
		Find(&prices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch price data: %w", err)
	}

	if len(prices) < 20 {
		return nil, fmt.Errorf("insufficient price data for %s", symbol)
	}

	// 2. 기술지표 계산
	indicators := s.indicatorService.CalculateAll(prices)
	if indicators == nil {
		return nil, fmt.Errorf("failed to calculate indicators for %s", symbol)
	}

	// 3. 기술지표를 map으로 변환
	indicatorMap := map[string]float64{
		"rsi":             indicators.RSI,
		"macd":            indicators.MACD,
		"macd_signal":     indicators.MACDSignal,
		"macd_histogram":  indicators.MACDHistogram,
		"sma_20":          indicators.SMA20,
		"sma_50":          indicators.SMA50,
		"ema_12":          indicators.EMA12,
		"ema_26":          indicators.EMA26,
		"bollinger_upper": indicators.BollingerUpper,
		"bollinger_lower": indicators.BollingerLower,
		"bollinger_mid":   indicators.BollingerMid,
		"stochastic_k":    indicators.StochasticK,
		"stochastic_d":    indicators.StochasticD,
		"williams_r":      indicators.WilliamsR,
		"atr":             indicators.ATR,
		"obv":             indicators.OBV,
	}

	// 4. 최신 주가 정보
	latestPrice := prices[0]

	// 5. AI 서비스에 의사결정 요청
	aiRequest := models.AIDecisionRequest{
		Symbol:     symbol,
		Market:     market,
		Price:      latestPrice,
		Indicators: indicatorMap,
		Metadata: map[string]interface{}{
			"data_points": len(prices),
			"timestamp":   time.Now().Unix(),
		},
	}

	aiResponse, err := s.aiClient.GetDecision(aiRequest)
	if err != nil {
		log.Printf("AI service error for %s: %v", symbol, err)
		// AI 서비스 실패 시 규칙 기반 fallback
		return s.generateRuleBasedSignal(symbol, market, indicatorMap, latestPrice)
	}

	// 6. AI 응답을 TradingSignal로 변환
	signal := &models.TradingSignal{
		Symbol:     symbol,
		SignalType: aiResponse.Decision,
		Strength:   s.calculateStrength(aiResponse.Confidence, indicatorMap),
		Confidence: aiResponse.Confidence,
		Reasons:    s.reasonsToJSON(aiResponse.Reasoning),
		Source:     "AI",
		CreatedAt:  time.Now(),
	}

	// 7. 데이터베이스에 저장
	if err := s.db.Create(signal).Error; err != nil {
		return nil, fmt.Errorf("failed to save signal: %w", err)
	}

	// 8. 캐시 무효화
	s.cacheService.InvalidateStock(symbol)

	// 9. 메시지 큐에 신호 발행
	if s.queueService != nil {
		s.queueService.PublishSignal(symbol, market, signal)
	}

	log.Printf("Generated signal for %s: %s (confidence: %.2f)", symbol, signal.SignalType, signal.Confidence)
	return signal, nil
}

// 규칙 기반 fallback 신호 생성
func (s *SignalGeneratorService) generateRuleBasedSignal(symbol, market string, indicators map[string]float64, price models.StockPrice) (*models.TradingSignal, error) {
	log.Printf("Using rule-based fallback for %s", symbol)

	decision := "HOLD"
	confidence := 0.5
	reasons := []string{"AI service unavailable, using rule-based analysis"}

	// 간단한 규칙 기반 로직
	rsi := indicators["rsi"]
	macd := indicators["macd"]
	sma20 := indicators["sma_20"]
	sma50 := indicators["sma_50"]

	buySignals := 0
	sellSignals := 0

	if rsi < 30 {
		buySignals++
		reasons = append(reasons, "RSI oversold")
	} else if rsi > 70 {
		sellSignals++
		reasons = append(reasons, "RSI overbought")
	}

	if macd > 0 {
		buySignals++
		reasons = append(reasons, "MACD positive")
	} else {
		sellSignals++
		reasons = append(reasons, "MACD negative")
	}

	if sma20 > sma50 {
		buySignals++
		reasons = append(reasons, "SMA20 > SMA50")
	} else {
		sellSignals++
		reasons = append(reasons, "SMA20 < SMA50")
	}

	if buySignals > sellSignals {
		decision = "BUY"
		confidence = 0.6
	} else if sellSignals > buySignals {
		decision = "SELL"
		confidence = 0.6
	}

	signal := &models.TradingSignal{
		Symbol:     symbol,
		SignalType: decision,
		Strength:   confidence * 0.8, // Rule-based는 약간 낮은 강도
		Confidence: confidence,
		Reasons:    s.reasonsToJSON(reasons),
		Source:     "RULE",
		CreatedAt:  time.Now(),
	}

	if err := s.db.Create(signal).Error; err != nil {
		return nil, fmt.Errorf("failed to save rule-based signal: %w", err)
	}

	return signal, nil
}

// 모든 활성 종목에 대한 신호 생성
func (s *SignalGeneratorService) GenerateSignalsForAllStocks() error {
	log.Println("Generating signals for all active stocks")

	var stocks []models.Stock
	err := s.db.Where("is_active = ?", true).Find(&stocks).Error
	if err != nil {
		return fmt.Errorf("failed to fetch active stocks: %w", err)
	}

	successCount := 0
	errorCount := 0

	for _, stock := range stocks {
		_, err := s.GenerateSignal(stock.Symbol, stock.Market)
		if err != nil {
			log.Printf("Failed to generate signal for %s: %v", stock.Symbol, err)
			errorCount++
		} else {
			successCount++
		}

		// API 호출 제한을 위한 지연
		time.Sleep(500 * time.Millisecond)
	}

	log.Printf("Signal generation completed: %d success, %d errors", successCount, errorCount)
	return nil
}

// 유틸리티 함수들
func (s *SignalGeneratorService) calculateStrength(confidence float64, indicators map[string]float64) float64 {
	// 신뢰도와 지표 강도를 결합하여 강도 계산
	strength := confidence

	// 볼린저 밴드 기반 조정
	if upper, ok := indicators["bollinger_upper"]; ok {
		if lower, ok := indicators["bollinger_lower"]; ok {
			// 밴드 위치에 따른 강도 조정
			bandWidth := upper - lower
			if bandWidth > 0 {
				strength += 0.1 // 변동성이 클 때 강도 증가
			}
		}
	}

	// ATR 기반 조정
	if atr, ok := indicators["atr"]; ok {
		if atr > 1000 { // 높은 변동성
			strength += 0.05
		}
	}

	// 최대 1.0으로 제한
	if strength > 1.0 {
		strength = 1.0
	}

	return strength
}

func (s *SignalGeneratorService) reasonsToJSON(reasons []string) string {
	data, err := json.Marshal(reasons)
	if err != nil {
		return `["Failed to encode reasons"]`
	}
	return string(data)
}