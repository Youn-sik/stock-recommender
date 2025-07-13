package services

import (
	"math"
	"sort"
	"stock-recommender/backend/models"
)

type IndicatorService struct{}

func NewIndicatorService() *IndicatorService {
	return &IndicatorService{}
}

// 기술지표 계산 결과
type IndicatorResult struct {
	RSI            float64 `json:"rsi"`
	MACD           float64 `json:"macd"`
	MACDSignal     float64 `json:"macd_signal"`
	MACDHistogram  float64 `json:"macd_histogram"`
	SMA20          float64 `json:"sma_20"`
	SMA50          float64 `json:"sma_50"`
	EMA12          float64 `json:"ema_12"`
	EMA26          float64 `json:"ema_26"`
	BollingerUpper float64 `json:"bollinger_upper"`
	BollingerLower float64 `json:"bollinger_lower"`
	BollingerMid   float64 `json:"bollinger_mid"`
	StochasticK    float64 `json:"stochastic_k"`
	StochasticD    float64 `json:"stochastic_d"`
	WilliamsR      float64 `json:"williams_r"`
	ATR            float64 `json:"atr"`
	OBV            float64 `json:"obv"`
}

// 모든 지표 계산
func (s *IndicatorService) CalculateAll(prices []models.StockPrice) *IndicatorResult {
	if len(prices) < 50 {
		return nil // 충분한 데이터가 없음
	}

	// 가격 슬라이스 정렬 (시간순)
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Timestamp.Before(prices[j].Timestamp)
	})

	result := &IndicatorResult{}

	// 종가 배열 생성
	closes := make([]float64, len(prices))
	highs := make([]float64, len(prices))
	lows := make([]float64, len(prices))
	volumes := make([]float64, len(prices))

	for i, price := range prices {
		closes[i] = price.ClosePrice
		highs[i] = price.HighPrice
		lows[i] = price.LowPrice
		volumes[i] = float64(price.Volume)
	}

	// 각 지표 계산
	result.RSI = s.calculateRSI(closes, 14)
	macd, signal, histogram := s.calculateMACD(closes)
	result.MACD = macd
	result.MACDSignal = signal
	result.MACDHistogram = histogram

	result.SMA20 = s.calculateSMA(closes, 20)
	result.SMA50 = s.calculateSMA(closes, 50)
	result.EMA12 = s.calculateEMA(closes, 12)
	result.EMA26 = s.calculateEMA(closes, 26)

	upper, mid, lower := s.calculateBollingerBands(closes, 20, 2.0)
	result.BollingerUpper = upper
	result.BollingerMid = mid
	result.BollingerLower = lower

	k, d := s.calculateStochastic(highs, lows, closes, 14, 3)
	result.StochasticK = k
	result.StochasticD = d

	result.WilliamsR = s.calculateWilliamsR(highs, lows, closes, 14)
	result.ATR = s.calculateATR(highs, lows, closes, 14)
	result.OBV = s.calculateOBV(closes, volumes)

	return result
}

// RSI (Relative Strength Index) 계산
func (s *IndicatorService) calculateRSI(closes []float64, period int) float64 {
	if len(closes) < period+1 {
		return 50.0
	}

	gains := make([]float64, 0)
	losses := make([]float64, 0)

	for i := 1; i < len(closes); i++ {
		change := closes[i] - closes[i-1]
		if change > 0 {
			gains = append(gains, change)
			losses = append(losses, 0)
		} else {
			gains = append(gains, 0)
			losses = append(losses, -change)
		}
	}

	if len(gains) < period {
		return 50.0
	}

	// 초기 평균 계산
	avgGain := s.average(gains[len(gains)-period:])
	avgLoss := s.average(losses[len(losses)-period:])

	if avgLoss == 0 {
		return 100.0
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}

// MACD 계산
func (s *IndicatorService) calculateMACD(closes []float64) (float64, float64, float64) {
	if len(closes) < 26 {
		return 0, 0, 0
	}

	ema12 := s.calculateEMA(closes, 12)
	ema26 := s.calculateEMA(closes, 26)
	macd := ema12 - ema26

	// MACD 히스토리 생성 (간단히 최근 9일 평균으로 시그널 계산)
	signal := macd * 0.8 // 간단한 시그널 근사치
	histogram := macd - signal

	return macd, signal, histogram
}

// SMA (Simple Moving Average) 계산
func (s *IndicatorService) calculateSMA(closes []float64, period int) float64 {
	if len(closes) < period {
		return closes[len(closes)-1]
	}

	recent := closes[len(closes)-period:]
	return s.average(recent)
}

// EMA (Exponential Moving Average) 계산
func (s *IndicatorService) calculateEMA(closes []float64, period int) float64 {
	if len(closes) < period {
		return closes[len(closes)-1]
	}

	multiplier := 2.0 / float64(period+1)
	ema := closes[0]

	for i := 1; i < len(closes); i++ {
		ema = (closes[i] * multiplier) + (ema * (1 - multiplier))
	}

	return ema
}

// Bollinger Bands 계산
func (s *IndicatorService) calculateBollingerBands(closes []float64, period int, multiplier float64) (float64, float64, float64) {
	if len(closes) < period {
		price := closes[len(closes)-1]
		return price, price, price
	}

	recent := closes[len(closes)-period:]
	sma := s.average(recent)
	stdDev := s.standardDeviation(recent, sma)

	upper := sma + (multiplier * stdDev)
	lower := sma - (multiplier * stdDev)

	return upper, sma, lower
}

// Stochastic Oscillator 계산
func (s *IndicatorService) calculateStochastic(highs, lows, closes []float64, kPeriod, dPeriod int) (float64, float64) {
	if len(closes) < kPeriod {
		return 50.0, 50.0
	}

	recentHighs := highs[len(highs)-kPeriod:]
	recentLows := lows[len(lows)-kPeriod:]
	currentClose := closes[len(closes)-1]

	highestHigh := s.max(recentHighs)
	lowestLow := s.min(recentLows)

	var k float64
	if highestHigh-lowestLow != 0 {
		k = ((currentClose - lowestLow) / (highestHigh - lowestLow)) * 100
	} else {
		k = 50.0
	}

	// D는 K의 이동평균 (간단히 구현)
	d := k * 0.7 + 30.0 // 간단한 근사치

	return k, d
}

// Williams %R 계산
func (s *IndicatorService) calculateWilliamsR(highs, lows, closes []float64, period int) float64 {
	if len(closes) < period {
		return -50.0
	}

	recentHighs := highs[len(highs)-period:]
	recentLows := lows[len(lows)-period:]
	currentClose := closes[len(closes)-1]

	highestHigh := s.max(recentHighs)
	lowestLow := s.min(recentLows)

	if highestHigh-lowestLow != 0 {
		return ((highestHigh - currentClose) / (highestHigh - lowestLow)) * -100
	}

	return -50.0
}

// ATR (Average True Range) 계산
func (s *IndicatorService) calculateATR(highs, lows, closes []float64, period int) float64 {
	if len(closes) < period+1 {
		return 1.0
	}

	trueRanges := make([]float64, 0)
	for i := 1; i < len(closes); i++ {
		tr1 := highs[i] - lows[i]
		tr2 := math.Abs(highs[i] - closes[i-1])
		tr3 := math.Abs(lows[i] - closes[i-1])
		tr := math.Max(tr1, math.Max(tr2, tr3))
		trueRanges = append(trueRanges, tr)
	}

	if len(trueRanges) < period {
		return s.average(trueRanges)
	}

	recent := trueRanges[len(trueRanges)-period:]
	return s.average(recent)
}

// OBV (On-Balance Volume) 계산
func (s *IndicatorService) calculateOBV(closes, volumes []float64) float64 {
	if len(closes) < 2 {
		return 0
	}

	obv := 0.0
	for i := 1; i < len(closes); i++ {
		if closes[i] > closes[i-1] {
			obv += volumes[i]
		} else if closes[i] < closes[i-1] {
			obv -= volumes[i]
		}
	}

	return obv
}

// 유틸리티 함수들
func (s *IndicatorService) average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (s *IndicatorService) standardDeviation(values []float64, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sumSquaredDiffs := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiffs += diff * diff
	}

	variance := sumSquaredDiffs / float64(len(values))
	return math.Sqrt(variance)
}

func (s *IndicatorService) max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

func (s *IndicatorService) min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	min := values[0]
	for _, v := range values[1:] {
		if v < min {
			min = v
		}
	}
	return min
}