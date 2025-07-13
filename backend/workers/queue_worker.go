package workers

import (
	"encoding/json"
	"log"
	"stock-recommender/backend/models"
	"stock-recommender/backend/services"
	"time"

	"gorm.io/gorm"
)

type QueueWorker struct {
	db               *gorm.DB
	queueService     *services.QueueService
	indicatorService *services.IndicatorService
	signalGenerator  *services.SignalGeneratorService
	aiClient         *services.AIClient
	cacheService     *services.CacheService
}

func NewQueueWorker(
	db *gorm.DB,
	queueService *services.QueueService,
	indicatorService *services.IndicatorService,
	signalGenerator *services.SignalGeneratorService,
	aiClient *services.AIClient,
	cacheService *services.CacheService,
) *QueueWorker {
	return &QueueWorker{
		db:               db,
		queueService:     queueService,
		indicatorService: indicatorService,
		signalGenerator:  signalGenerator,
		aiClient:         aiClient,
		cacheService:     cacheService,
	}
}

func (w *QueueWorker) StartWorkers() error {
	log.Println("Starting queue workers...")

	// Price update worker
	err := w.queueService.Subscribe("price.updates", w.handlePriceUpdate)
	if err != nil {
		return err
	}

	// Indicator calculation worker
	err = w.queueService.Subscribe("indicator.calculation", w.handleIndicatorCalculation)
	if err != nil {
		return err
	}

	// AI request worker
	err = w.queueService.Subscribe("ai.requests", w.handleAIRequest)
	if err != nil {
		return err
	}

	// Signal generation worker
	err = w.queueService.Subscribe("signal.generation", w.handleSignalGeneration)
	if err != nil {
		return err
	}

	log.Println("All queue workers started successfully")
	return nil
}

func (w *QueueWorker) handlePriceUpdate(message services.Message) error {
	log.Printf("Processing price update for %s", message.Symbol)

	// Trigger indicator calculation
	err := w.queueService.PublishIndicatorRequest(message.Symbol, message.Market)
	if err != nil {
		log.Printf("Failed to publish indicator request: %v", err)
		return err
	}

	// Invalidate cache
	w.cacheService.InvalidateStock(message.Symbol)

	return nil
}

func (w *QueueWorker) handleIndicatorCalculation(message services.Message) error {
	log.Printf("Calculating indicators for %s", message.Symbol)

	// Fetch recent price data
	var prices []models.StockPrice
	err := w.db.Where("symbol = ? AND market = ?", message.Symbol, message.Market).
		Order("timestamp desc").
		Limit(50).
		Find(&prices).Error
	if err != nil {
		log.Printf("Failed to fetch prices for %s: %v", message.Symbol, err)
		return err
	}

	if len(prices) < 20 {
		log.Printf("Insufficient price data for %s", message.Symbol)
		return nil
	}

	// Calculate indicators
	indicators := w.indicatorService.CalculateAll(prices)
	if indicators == nil {
		log.Printf("Failed to calculate indicators for %s", message.Symbol)
		return nil
	}

	// Save indicators to database
	err = w.saveIndicators(message.Symbol, indicators)
	if err != nil {
		log.Printf("Failed to save indicators for %s: %v", message.Symbol, err)
		return err
	}

	// Update cache
	indicatorMap := w.convertIndicatorsToMap(indicators)
	w.cacheService.SetIndicators(message.Symbol, indicatorMap)

	// Trigger AI analysis
	err = w.queueService.PublishAIRequest(message.Symbol, message.Market, indicatorMap)
	if err != nil {
		log.Printf("Failed to publish AI request: %v", err)
	}

	return nil
}

func (w *QueueWorker) handleAIRequest(message services.Message) error {
	log.Printf("Processing AI request for %s", message.Symbol)

	// Get latest price
	var latestPrice models.StockPrice
	err := w.db.Where("symbol = ?", message.Symbol).
		Order("timestamp desc").
		First(&latestPrice).Error
	if err != nil {
		log.Printf("Failed to get latest price for %s: %v", message.Symbol, err)
		return err
	}

	// Prepare indicators data
	indicatorData, ok := message.Data.(map[string]interface{})
	if !ok {
		log.Printf("Invalid indicator data format for %s", message.Symbol)
		return nil
	}

	indicators := make(map[string]float64)
	for k, v := range indicatorData {
		if val, ok := v.(float64); ok {
			indicators[k] = val
		}
	}

	// Create AI request
	aiRequest := models.AIDecisionRequest{
		Symbol:     message.Symbol,
		Market:     message.Market,
		Price:      latestPrice,
		Indicators: indicators,
		Metadata: map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"worker":    "queue_worker",
		},
	}

	// Get AI decision
	decision, err := w.aiClient.GetDecision(aiRequest)
	if err != nil {
		log.Printf("AI request failed for %s: %v", message.Symbol, err)
		return err
	}

	// Create trading signal
	signal := &models.TradingSignal{
		Symbol:     message.Symbol,
		SignalType: decision.Decision,
		Strength:   w.calculateSignalStrength(decision.Confidence),
		Confidence: decision.Confidence,
		Reasons:    w.reasonsToJSON(decision.Reasoning),
		Source:     "AI",
		CreatedAt:  time.Now(),
	}

	// Save signal
	err = w.db.Create(signal).Error
	if err != nil {
		log.Printf("Failed to save signal for %s: %v", message.Symbol, err)
		return err
	}

	// Publish signal
	err = w.queueService.PublishSignal(message.Symbol, message.Market, signal)
	if err != nil {
		log.Printf("Failed to publish signal: %v", err)
	}

	return nil
}

func (w *QueueWorker) handleSignalGeneration(message services.Message) error {
	log.Printf("Processing signal generation for %s", message.Symbol)

	// This can trigger notifications, alerts, etc.
	// For now, just log the signal
	if signalData, ok := message.Data.(*models.TradingSignal); ok {
		log.Printf("New trading signal: %s - %s (confidence: %.2f)", 
			signalData.Symbol, signalData.SignalType, signalData.Confidence)
	}

	return nil
}

// Helper functions
func (w *QueueWorker) saveIndicators(symbol string, indicators *services.IndicatorResult) error {
	indicatorMap := w.convertIndicatorsToMap(indicators)
	
	for name, value := range indicatorMap {
		indicatorRecord := &models.TechnicalIndicator{
			Symbol:         symbol,
			IndicatorName:  name,
			IndicatorValue: w.valueToJSON(value),
			CalculatedAt:   time.Now(),
			CreatedAt:      time.Now(),
		}
		
		err := w.db.Create(indicatorRecord).Error
		if err != nil {
			return err
		}
	}
	
	return nil
}

func (w *QueueWorker) convertIndicatorsToMap(indicators *services.IndicatorResult) map[string]float64 {
	return map[string]float64{
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
}

func (w *QueueWorker) valueToJSON(value float64) string {
	data, _ := json.Marshal(map[string]float64{"value": value})
	return string(data)
}

func (w *QueueWorker) reasonsToJSON(reasons []string) string {
	data, _ := json.Marshal(reasons)
	return string(data)
}

func (w *QueueWorker) calculateSignalStrength(confidence float64) float64 {
	// Simple mapping of confidence to strength
	if confidence > 0.8 {
		return 0.9
	} else if confidence > 0.6 {
		return 0.7
	} else if confidence > 0.4 {
		return 0.5
	}
	return 0.3
}