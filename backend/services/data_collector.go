package services

import (
	"fmt"
	"log"
	"time"

	"stock-recommender/backend/config"
	"stock-recommender/backend/models"
	"stock-recommender/backend/openapi/client"
	apimodels "stock-recommender/backend/openapi/models"
	"gorm.io/gorm"
)

type DataCollectorService struct {
	db        *gorm.DB
	apiClient *client.DBSecClient
	config    *config.Config
}

func NewDataCollectorService(db *gorm.DB, cfg *config.Config) *DataCollectorService {
	return &DataCollectorService{
		db:        db,
		apiClient: client.NewDBSecClient(cfg),
		config:    cfg,
	}
}

// 전체 종목 데이터 수집
func (s *DataCollectorService) CollectAllStocks() error {
	log.Println("Starting data collection for all stocks...")

	// 등록된 종목 목록 조회
	var stocks []models.Stock
	if err := s.db.Find(&stocks).Error; err != nil {
		return fmt.Errorf("failed to get stocks: %w", err)
	}

	successCount := 0
	errorCount := 0

	for _, stock := range stocks {
		err := s.CollectStockData(stock.Symbol, stock.Market)
		if err != nil {
			log.Printf("Failed to collect data for %s (%s): %v", stock.Symbol, stock.Name, err)
			errorCount++
		} else {
			successCount++
		}

		// API 호출 제한을 위한 지연
		time.Sleep(100 * time.Millisecond)
	}

	log.Printf("Data collection completed: %d success, %d errors", successCount, errorCount)
	return nil
}

// 특정 종목 데이터 수집
func (s *DataCollectorService) CollectStockData(symbol, market string) error {
	// API에서 데이터 수집
	priceData, askingData, err := s.apiClient.CollectStockData(symbol, market)
	if err != nil {
		// API 실패시 Mock 데이터 사용 (개발용)
		if !s.apiClient.HasValidCredentials() {
			log.Printf("API credentials not available, using mock data for %s", symbol)
			return s.generateMockData(symbol, market)
		}
		return fmt.Errorf("failed to collect data from API: %w", err)
	}

	// 주가 데이터 저장
	if priceData != nil {
		if err := s.saveStockPrice(priceData); err != nil {
			return fmt.Errorf("failed to save price data: %w", err)
		}
	}

	// 호가 데이터 저장 (있는 경우)
	if askingData != nil {
		if err := s.saveAskingPrice(askingData); err != nil {
			return fmt.Errorf("failed to save asking price data: %w", err)
		}
	}

	log.Printf("Successfully collected data for %s", symbol)
	return nil
}

// 주가 데이터 저장
func (s *DataCollectorService) saveStockPrice(priceData *apimodels.ParsedStockPrice) error {
	stockPrice := models.StockPrice{
		Symbol:         priceData.Symbol,
		OpenPrice:      priceData.OpenPrice,
		HighPrice:      priceData.HighPrice,
		LowPrice:       priceData.LowPrice,
		ClosePrice:     priceData.CurrentPrice,
		Volume:         priceData.Volume,
		TradeAmount:    priceData.TradeAmount,
		PrevClosePrice: priceData.PrevClosePrice,
		Change:         priceData.Change,
		ChangeRate:     priceData.ChangeRate,
		Timestamp:      priceData.Timestamp,
		Market:         priceData.Market,
	}

	// 중복 데이터 체크 (같은 시각, 같은 종목)
	var existing models.StockPrice
	result := s.db.Where("symbol = ? AND timestamp = ?", stockPrice.Symbol, stockPrice.Timestamp).First(&existing)
	
	if result.Error == gorm.ErrRecordNotFound {
		// 새 데이터 삽입
		return s.db.Create(&stockPrice).Error
	} else if result.Error != nil {
		return result.Error
	}

	// 기존 데이터 업데이트
	return s.db.Model(&existing).Updates(stockPrice).Error
}

// 호가 데이터 저장
func (s *DataCollectorService) saveAskingPrice(askingData *apimodels.ParsedAskingPrice) error {
	askingPrice := models.AskingPrice{
		Symbol:      askingData.Symbol,
		AskPrice1:   askingData.AskPrices[0],
		AskPrice2:   askingData.AskPrices[1],
		AskPrice3:   askingData.AskPrices[2],
		AskPrice4:   askingData.AskPrices[3],
		AskPrice5:   askingData.AskPrices[4],
		BidPrice1:   askingData.BidPrices[0],
		BidPrice2:   askingData.BidPrices[1],
		BidPrice3:   askingData.BidPrices[2],
		BidPrice4:   askingData.BidPrices[3],
		BidPrice5:   askingData.BidPrices[4],
		AskVolume1:  askingData.AskVolumes[0],
		AskVolume2:  askingData.AskVolumes[1],
		AskVolume3:  askingData.AskVolumes[2],
		AskVolume4:  askingData.AskVolumes[3],
		AskVolume5:  askingData.AskVolumes[4],
		BidVolume1:  askingData.BidVolumes[0],
		BidVolume2:  askingData.BidVolumes[1],
		BidVolume3:  askingData.BidVolumes[2],
		BidVolume4:  askingData.BidVolumes[3],
		BidVolume5:  askingData.BidVolumes[4],
		TotalAskVol: askingData.TotalAskVol,
		TotalBidVol: askingData.TotalBidVol,
		Timestamp:   askingData.Timestamp,
	}

	// 최신 호가 정보만 유지 (이전 데이터 삭제)
	if err := s.db.Where("symbol = ?", askingPrice.Symbol).Delete(&models.AskingPrice{}).Error; err != nil {
		log.Printf("Warning: failed to delete old asking price data for %s: %v", askingPrice.Symbol, err)
	}

	return s.db.Create(&askingPrice).Error
}

// 종목별 일봉 데이터 수집
func (s *DataCollectorService) CollectDailyData(symbol string, days int) error {
	endDate := time.Now().Format("20060102")
	startDate := time.Now().AddDate(0, 0, -days).Format("20060102")

	dailyData, err := s.apiClient.GetDomesticStockDaily(symbol, startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to get daily data: %w", err)
	}

	for _, data := range dailyData {
		stockPrice := models.StockPrice{
			Symbol:      data.Symbol,
			OpenPrice:   data.OpenPrice,
			HighPrice:   data.HighPrice,
			LowPrice:    data.LowPrice,
			ClosePrice:  data.ClosePrice,
			Volume:      data.Volume,
			TradeAmount: data.TradeAmount,
			Timestamp:   data.Date,
			Market:      "KR",
		}

		// 중복 체크 후 저장
		var existing models.StockPrice
		result := s.db.Where("symbol = ? AND DATE(timestamp) = DATE(?)", 
			stockPrice.Symbol, stockPrice.Timestamp).First(&existing)
		
		if result.Error == gorm.ErrRecordNotFound {
			if err := s.db.Create(&stockPrice).Error; err != nil {
				log.Printf("Failed to save daily data for %s on %s: %v", symbol, data.Date.Format("2006-01-02"), err)
			}
		}
	}

	return nil
}

// Mock 데이터 생성 (개발 및 테스트용)
func (s *DataCollectorService) generateMockData(symbol, market string) error {
	log.Printf("Generating mock data for %s (%s)", symbol, market)

	// 기본 주가 설정 (종목별로 다르게)
	basePrice := 50000.0
	switch symbol {
	case "005930": // 삼성전자
		basePrice = 70000.0
	case "000660": // SK하이닉스
		basePrice = 120000.0
	case "AAPL":
		basePrice = 180.0
	case "TSLA":
		basePrice = 250.0
	}

	// 랜덤 변동 (-2% ~ +2%)
	variation := float64(time.Now().Unix()%400 - 200) / 10000.0 // -0.02 ~ 0.02
	currentPrice := basePrice * (1 + variation)
	
	mockPrice := models.StockPrice{
		Symbol:         symbol,
		OpenPrice:      currentPrice * 0.998,
		HighPrice:      currentPrice * 1.003,
		LowPrice:       currentPrice * 0.995,
		ClosePrice:     currentPrice,
		Volume:         int64(100000 + time.Now().Unix()%500000),
		TradeAmount:    int64(currentPrice * float64(100000+time.Now().Unix()%500000)),
		PrevClosePrice: basePrice,
		Change:         currentPrice - basePrice,
		ChangeRate:     ((currentPrice - basePrice) / basePrice) * 100,
		Timestamp:      time.Now(),
		Market:         market,
	}

	return s.db.Create(&mockPrice).Error
}

// 시장별 주요 종목 초기화
func (s *DataCollectorService) InitializeMajorStocks() error {
	majorStocks := s.apiClient.GetMajorStocks()

	for market, symbols := range majorStocks {
		for _, symbol := range symbols {
			// 종목이 이미 등록되어 있는지 확인
			var existing models.Stock
			result := s.db.Where("symbol = ? AND market = ?", symbol, market).First(&existing)
			
			if result.Error == gorm.ErrRecordNotFound {
				// 새 종목 등록
				stock := models.Stock{
					Symbol:   symbol,
					Market:   market,
					IsActive: true,
				}

				// 종목명 설정
				switch symbol {
				case "005930":
					stock.Name = "삼성전자"
					stock.Exchange = "KOSPI"
					stock.Sector = "Technology"
				case "000660":
					stock.Name = "SK하이닉스"
					stock.Exchange = "KOSPI"
					stock.Sector = "Technology"
				case "AAPL":
					stock.Name = "Apple Inc."
					stock.Exchange = "NASDAQ"
					stock.Sector = "Technology"
				case "TSLA":
					stock.Name = "Tesla Inc."
					stock.Exchange = "NASDAQ"
					stock.Sector = "Automotive"
				default:
					stock.Name = symbol // 기본값
				}

				if err := s.db.Create(&stock).Error; err != nil {
					log.Printf("Failed to create stock %s: %v", symbol, err)
				} else {
					log.Printf("Added new stock: %s (%s)", stock.Name, symbol)
				}
			}
		}
	}

	return nil
}

// 정기 수집 작업 시작
func (s *DataCollectorService) StartScheduledCollection() {
	log.Println("Starting scheduled data collection...")

	// 주요 종목 초기화
	if err := s.InitializeMajorStocks(); err != nil {
		log.Printf("Failed to initialize major stocks: %v", err)
	}

	// 즉시 한 번 수집
	if err := s.CollectAllStocks(); err != nil {
		log.Printf("Initial data collection failed: %v", err)
	}

	// 정기 수집 (5분마다)
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.CollectAllStocks(); err != nil {
					log.Printf("Scheduled data collection failed: %v", err)
				}
			}
		}
	}()
}

// API 상태 확인
func (s *DataCollectorService) GetAPIStatus() map[string]interface{} {
	return s.apiClient.GetAPIStatus()
}