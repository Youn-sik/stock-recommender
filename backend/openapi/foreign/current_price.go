package foreign

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
)

// ForeignCurrentPriceService 해외주식현재가조회 서비스
type ForeignCurrentPriceService struct {
	client *client.DBSecClient
}

// NewForeignCurrentPriceService 새로운 해외주식현재가조회 서비스 생성
func NewForeignCurrentPriceService(client *client.DBSecClient) *ForeignCurrentPriceService {
	return &ForeignCurrentPriceService{
		client: client,
	}
}

// GetForeignCurrentPrice 해외주식 현재가 조회
// stockCode: 해외주식종목코드 (예: TSLA, AAPL)
// marketDiv: 시장분류코드 (FY: 뉴욕, FN: 나스닥, FA: 아멕스)
func (s *ForeignCurrentPriceService) GetForeignCurrentPrice(stockCode string, marketDiv string) (*models.ForeignCurrentPriceData, error) {
	// 요청 데이터 구성
	reqBody := models.ForeignCurrentPriceRequest{
		In: models.ForeignCurrentPriceInput{
			InputCondMrktDivCode: marketDiv,
			InputIscd1:           stockCode,
		},
	}

	// API 호출
	respBody, err := s.client.MakeRequestWithHeaders("POST", models.PathForeignStockCurrentPrice, nil, reqBody, map[string]string{
		"cont_yn": "N",
		"tr_id":   models.TrIdForeignStockCurrentPrice,
	})
	if err != nil {
		return nil, err
	}

	// 응답 파싱
	var response models.ForeignCurrentPriceResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 응답 코드 확인
	if response.RspCd != "00000" {
		return nil, fmt.Errorf("API error %s: %s", response.RspCd, response.RspMsg)
	}

	// 데이터 변환
	return s.convertToForeignCurrentPriceData(stockCode, marketDiv, &response.Out), nil
}

// GetNYStockPrice 뉴욕 거래소 주식 현재가 조회
func (s *ForeignCurrentPriceService) GetNYStockPrice(stockCode string) (*models.ForeignCurrentPriceData, error) {
	return s.GetForeignCurrentPrice(stockCode, models.ForeignMarketNY)
}

// GetNASDAQStockPrice 나스닥 주식 현재가 조회
func (s *ForeignCurrentPriceService) GetNASDAQStockPrice(stockCode string) (*models.ForeignCurrentPriceData, error) {
	return s.GetForeignCurrentPrice(stockCode, models.ForeignMarketNASDAQ)
}

// GetAMEXStockPrice 아멕스 주식 현재가 조회
func (s *ForeignCurrentPriceService) GetAMEXStockPrice(stockCode string) (*models.ForeignCurrentPriceData, error) {
	return s.GetForeignCurrentPrice(stockCode, models.ForeignMarketAMEX)
}

// GetUSStockPrice 미국 주식 현재가 조회 (자동 거래소 감지)
func (s *ForeignCurrentPriceService) GetUSStockPrice(stockCode string) (*models.ForeignCurrentPriceData, error) {
	// 먼저 나스닥에서 시도
	if data, err := s.GetNASDAQStockPrice(stockCode); err == nil {
		return data, nil
	}

	// 나스닥에서 실패하면 뉴욕 거래소에서 시도
	if data, err := s.GetNYStockPrice(stockCode); err == nil {
		return data, nil
	}

	// 뉴욕에서도 실패하면 아멕스에서 시도
	return s.GetAMEXStockPrice(stockCode)
}

// GetMultipleForeignStockPrices 여러 해외 주식의 현재가 일괄 조회
func (s *ForeignCurrentPriceService) GetMultipleForeignStockPrices(stockCodes []string, marketDiv string) (map[string]*models.ForeignCurrentPriceData, error) {
	result := make(map[string]*models.ForeignCurrentPriceData)
	
	for _, code := range stockCodes {
		data, err := s.GetForeignCurrentPrice(code, marketDiv)
		if err != nil {
			// 개별 오류는 로그하고 계속 진행
			fmt.Printf("Failed to get price for %s: %v\n", code, err)
			continue
		}
		result[code] = data
	}
	
	return result, nil
}

// GetMultipleUSStockPrices 여러 미국 주식의 현재가 일괄 조회 (자동 거래소 감지)
func (s *ForeignCurrentPriceService) GetMultipleUSStockPrices(stockCodes []string) (map[string]*models.ForeignCurrentPriceData, error) {
	result := make(map[string]*models.ForeignCurrentPriceData)
	
	for _, code := range stockCodes {
		data, err := s.GetUSStockPrice(code)
		if err != nil {
			// 개별 오류는 로그하고 계속 진행
			fmt.Printf("Failed to get price for %s: %v\n", code, err)
			continue
		}
		result[code] = data
	}
	
	return result, nil
}

// GetPopularStockPrices 인기 주식들의 현재가 조회
func (s *ForeignCurrentPriceService) GetPopularStockPrices() (map[string]*models.ForeignCurrentPriceData, error) {
	popularStocks := []string{
		"AAPL", // 애플
		"MSFT", // 마이크로소프트  
		"GOOGL", // 알파벳
		"AMZN", // 아마존
		"TSLA", // 테슬라
		"META", // 메타
		"NVDA", // 엔비디아
		"NFLX", // 넷플릭스
	}
	
	return s.GetMultipleUSStockPrices(popularStocks)
}

// GetTechGiantsPrices 빅테크 기업들의 현재가 조회
func (s *ForeignCurrentPriceService) GetTechGiantsPrices() (map[string]*models.ForeignCurrentPriceData, error) {
	techGiants := []string{
		"AAPL",  // 애플
		"MSFT",  // 마이크로소프트
		"GOOGL", // 알파벳
		"AMZN",  // 아마존
		"META",  // 메타
		"NVDA",  // 엔비디아
	}
	
	return s.GetMultipleUSStockPrices(techGiants)
}

// convertToForeignCurrentPriceData 응답 데이터를 구조화된 형식으로 변환
func (s *ForeignCurrentPriceService) convertToForeignCurrentPriceData(stockCode string, marketDiv string, output *models.ForeignCurrentPriceOutput) *models.ForeignCurrentPriceData {
	return &models.ForeignCurrentPriceData{
		StockCode:        stockCode,
		Market:           s.getMarketName(marketDiv),
		BasePrice:        s.parseFloat(output.Sdpr),
		CurrentPrice:     s.parseFloat(output.Prpr),
		UpperLimit:       s.parseFloat(output.Mxpr),
		LowerLimit:       s.parseFloat(output.Llam),
		OpenPrice:        s.parseFloat(output.Oprc),
		HighPrice:        s.parseFloat(output.Hprc),
		LowPrice:         s.parseFloat(output.Lprc),
		PriceChange:      s.parseFloat(output.PrdyVrss),
		PriceChangeRate:  s.parseFloat(output.PrdyCtrt),
		PER:              s.parseFloat(output.Per),
		TradingValue:     s.parseFloat(output.AcmlTrPbmn),
		TradingVolume:    s.parseInt(output.AcmlVol),
		YesterdayVolume:  s.parseInt(output.PrdyVol),
		BidPrice:         s.parseFloat(output.Bidp1),
		AskPrice:         s.parseFloat(output.Askp1),
		MarketOpenRate:   s.parseFloat(output.SdprVrssMrktRate),
		CurrentOpenRate:  s.parseFloat(output.PrprVrssOprcRate),
		MarketHighRate:   s.parseFloat(output.SdprVrssHgprRate),
		CurrentHighRate:  s.parseFloat(output.PrprVrssHgprRate),
		MarketLowRate:    s.parseFloat(output.SdprVrssLwprRate),
		CurrentLowRate:   s.parseFloat(output.PrprVrssLwprRate),
		Currency:         "USD",
	}
}

// getMarketName 시장분류코드를 시장명으로 변환
func (s *ForeignCurrentPriceService) getMarketName(marketDiv string) string {
	switch marketDiv {
	case models.ForeignMarketNY:
		return "뉴욕"
	case models.ForeignMarketNASDAQ:
		return "나스닥"
	case models.ForeignMarketAMEX:
		return "아멕스"
	default:
		return marketDiv
	}
}

// 유틸리티 함수들
func (s *ForeignCurrentPriceService) parseFloat(str string) float64 {
	if str == "" {
		return 0
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		return 0
	}
	return val
}

func (s *ForeignCurrentPriceService) parseInt(str string) int64 {
	if str == "" {
		return 0
	}
	val, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		return 0
	}
	return val
}