package domestic

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
)

// CurrentPriceService 현재가조회 서비스
type CurrentPriceService struct {
	client *client.DBSecClient
}

// NewCurrentPriceService 새로운 현재가조회 서비스 생성
func NewCurrentPriceService(client *client.DBSecClient) *CurrentPriceService {
	return &CurrentPriceService{
		client: client,
	}
}

// GetCurrentPrice 현재가 조회
// stockCode: 종목코드 (6자리) 또는 지수코드
// marketDiv: 시장분류코드 (J: 주식, E: ETF, EN: ETN, W: ELW, U: 업종&지수)
func (s *CurrentPriceService) GetCurrentPrice(stockCode string, marketDiv string) (*models.CurrentPriceData, error) {
	// 요청 데이터 구성
	reqBody := models.CurrentPriceRequest{
		In: models.CurrentPriceInput{
			InputCondMrktDivCode: marketDiv,
			InputIscd1:           stockCode,
		},
	}

	// API 호출
	respBody, err := s.client.MakeRequestWithHeaders("POST", models.PathDomesticStockCurrentPrice, nil, reqBody, map[string]string{
		"cont_yn": "N",
		"tr_id":   models.TrIdStockCurrentPrice,
	})
	if err != nil {
		return nil, err
	}

	// 응답 파싱
	var response models.CurrentPriceResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 응답 코드 확인
	if response.RspCd != "00000" {
		return nil, fmt.Errorf("API error %s: %s", response.RspCd, response.RspMsg)
	}

	// 데이터 변환
	return s.convertToCurrentPriceData(stockCode, &response.Out), nil
}

// GetStockPrice 주식 현재가 조회
func (s *CurrentPriceService) GetStockPrice(stockCode string) (*models.CurrentPriceData, error) {
	return s.GetCurrentPrice(stockCode, models.MarketDivStock)
}

// GetETFPrice ETF 현재가 조회
func (s *CurrentPriceService) GetETFPrice(stockCode string) (*models.CurrentPriceData, error) {
	return s.GetCurrentPrice(stockCode, models.MarketDivETF)
}

// GetETNPrice ETN 현재가 조회
func (s *CurrentPriceService) GetETNPrice(stockCode string) (*models.CurrentPriceData, error) {
	return s.GetCurrentPrice(stockCode, models.MarketDivETN)
}

// GetELWPrice ELW 현재가 조회
func (s *CurrentPriceService) GetELWPrice(stockCode string) (*models.CurrentPriceData, error) {
	return s.GetCurrentPrice(stockCode, models.MarketDivELW)
}

// GetIndexPrice 지수 현재가 조회
func (s *CurrentPriceService) GetIndexPrice(indexCode string) (*models.CurrentPriceData, error) {
	return s.GetCurrentPrice(indexCode, models.MarketDivIndex)
}

// GetKOSPIPrice KOSPI 지수 조회
func (s *CurrentPriceService) GetKOSPIPrice() (*models.CurrentPriceData, error) {
	return s.GetIndexPrice(models.IndexKOSPI)
}

// GetKOSDAQPrice KOSDAQ 지수 조회
func (s *CurrentPriceService) GetKOSDAQPrice() (*models.CurrentPriceData, error) {
	return s.GetIndexPrice(models.IndexKOSDAQ)
}

// GetKOSPI200Price KOSPI200 지수 조회
func (s *CurrentPriceService) GetKOSPI200Price() (*models.CurrentPriceData, error) {
	return s.GetIndexPrice(models.IndexKOSPI200)
}

// GetMultipleStockPrices 여러 종목의 현재가 일괄 조회
func (s *CurrentPriceService) GetMultipleStockPrices(stockCodes []string) (map[string]*models.CurrentPriceData, error) {
	result := make(map[string]*models.CurrentPriceData)
	
	for _, code := range stockCodes {
		data, err := s.GetStockPrice(code)
		if err != nil {
			// 개별 오류는 로그하고 계속 진행
			fmt.Printf("Failed to get price for %s: %v\n", code, err)
			continue
		}
		result[code] = data
	}
	
	return result, nil
}

// convertToCurrentPriceData 응답 데이터를 구조화된 형식으로 변환
func (s *CurrentPriceService) convertToCurrentPriceData(stockCode string, output *models.CurrentPriceOutput) *models.CurrentPriceData {
	return &models.CurrentPriceData{
		StockCode:        stockCode,
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
		PBR:              s.parseFloat(output.Pbr),
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
	}
}

// 유틸리티 함수들
func (s *CurrentPriceService) parseFloat(str string) float64 {
	if str == "" {
		return 0
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		return 0
	}
	return val
}

func (s *CurrentPriceService) parseInt(str string) int64 {
	if str == "" {
		return 0
	}
	val, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		return 0
	}
	return val
}