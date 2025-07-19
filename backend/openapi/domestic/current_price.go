package domestic

import (
	"fmt"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/errors"
	"stock-recommender/backend/openapi/models"
	"stock-recommender/backend/openapi/utils"
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
		return nil, errors.NewNetworkError("failed to call current price API", err)
	}

	// 응답 파싱 및 검증
	var response models.CurrentPriceResponse
	if err := utils.ParseAPIResponse(respBody, &response); err != nil {
		return nil, errors.NewParseError("failed to parse current price response", err)
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
		BasePrice:        utils.ParseFloat(output.Sdpr),
		CurrentPrice:     utils.ParseFloat(output.Prpr),
		UpperLimit:       utils.ParseFloat(output.Mxpr),
		LowerLimit:       utils.ParseFloat(output.Llam),
		OpenPrice:        utils.ParseFloat(output.Oprc),
		HighPrice:        utils.ParseFloat(output.Hprc),
		LowPrice:         utils.ParseFloat(output.Lprc),
		PriceChange:      utils.ParseFloat(output.PrdyVrss),
		PriceChangeRate:  utils.ParseFloat(output.PrdyCtrt),
		PER:              utils.ParseFloat(output.Per),
		PBR:              utils.ParseFloat(output.Pbr),
		TradingValue:     utils.ParseFloat(output.AcmlTrPbmn),
		TradingVolume:    utils.ParseInt(output.AcmlVol),
		YesterdayVolume:  utils.ParseInt(output.PrdyVol),
		BidPrice:         utils.ParseFloat(output.Bidp1),
		AskPrice:         utils.ParseFloat(output.Askp1),
		MarketOpenRate:   utils.ParseFloat(output.SdprVrssMrktRate),
		CurrentOpenRate:  utils.ParseFloat(output.PrprVrssOprcRate),
		MarketHighRate:   utils.ParseFloat(output.SdprVrssHgprRate),
		CurrentHighRate:  utils.ParseFloat(output.PrprVrssHgprRate),
		MarketLowRate:    utils.ParseFloat(output.SdprVrssLwprRate),
		CurrentLowRate:   utils.ParseFloat(output.PrprVrssLwprRate),
	}
}