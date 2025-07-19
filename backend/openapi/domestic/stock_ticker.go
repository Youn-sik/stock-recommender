package domestic

import (
	"encoding/json"
	"fmt"
	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
)

// StockTickerService 주식종목 조회 서비스
type StockTickerService struct {
	client *client.DBSecClient
}

// NewStockTickerService 새로운 주식종목 조회 서비스 생성
func NewStockTickerService(client *client.DBSecClient) *StockTickerService {
	return &StockTickerService{
		client: client,
	}
}

// GetStockTickers 주식종목 조회
// marketDiv: 시장분류코드 (J: 주식, E: ETF, EN: ETN)
// contKey: 연속키 (optional, 추가 데이터 조회시 사용)
func (s *StockTickerService) GetStockTickers(marketDiv string, contKey string) (*models.StockTickerResponse, string, error) {
	// 요청 데이터 구성
	reqBody := models.StockTickerRequest{
		In: models.StockTickerInput{
			InputCondMrktDivCode: marketDiv,
		},
	}

	// 헤더 설정을 위한 커스텀 요청 함수
	respBody, nextContKey, err := s.makeRequestWithContKey(reqBody, contKey)
	if err != nil {
		return nil, "", err
	}

	// 응답 파싱
	var response models.StockTickerResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, "", fmt.Errorf("failed to parse response: %w", err)
	}

	// 응답 코드 확인
	if response.RspCd != "00000" {
		return nil, "", fmt.Errorf("API error %s: %s", response.RspCd, response.RspMsg)
	}

	return &response, nextContKey, nil
}

// GetAllStockTickers 모든 주식종목 조회 (페이지네이션 포함)
func (s *StockTickerService) GetAllStockTickers(marketDiv string) ([]models.StockTickerOutput, error) {
	var allStocks []models.StockTickerOutput
	contKey := ""

	for {
		response, nextContKey, err := s.GetStockTickers(marketDiv, contKey)
		if err != nil {
			return nil, err
		}

		allStocks = append(allStocks, response.Out...)

		// 더 이상 데이터가 없으면 종료
		if nextContKey == "" || nextContKey == "N" {
			break
		}

		contKey = nextContKey
	}

	return allStocks, nil
}

// GetStocks 주식 종목만 조회
func (s *StockTickerService) GetStocks() ([]models.StockTickerOutput, error) {
	return s.GetAllStockTickers(models.MarketDivStock)
}

// GetETFs ETF 종목만 조회
func (s *StockTickerService) GetETFs() ([]models.StockTickerOutput, error) {
	return s.GetAllStockTickers(models.MarketDivETF)
}

// GetETNs ETN 종목만 조회
func (s *StockTickerService) GetETNs() ([]models.StockTickerOutput, error) {
	return s.GetAllStockTickers(models.MarketDivETN)
}

// makeRequestWithContKey 연속키를 포함한 요청 처리
func (s *StockTickerService) makeRequestWithContKey(reqBody interface{}, contKey string) ([]byte, string, error) {
	// 추가 헤더 설정
	headers := map[string]string{
		"cont_yn":  s.getContYn(contKey),
		"cont_key": contKey,
		"tr_id":    models.TrIdStockTicker,
	}

	// API 호출
	response, err := s.client.MakeRequestWithFullResponse("POST", models.PathDomesticStockTicker, nil, reqBody, headers)
	if err != nil {
		return nil, "", err
	}

	// 응답 헤더에서 연속키 추출
	nextContKey := response.Headers.Get("cont_key")
	if nextContKey == "" {
		nextContKey = response.Headers.Get("CONT_KEY")
	}

	// 연속 여부 확인
	contYn := response.Headers.Get("cont_yn")
	if contYn == "" {
		contYn = response.Headers.Get("CONT_YN")
	}

	// 연속 거래가 없으면 빈 문자열 반환
	if contYn != "Y" {
		nextContKey = ""
	}

	return response.Body, nextContKey, nil
}

// getContYn 연속거래 여부 반환
func (s *StockTickerService) getContYn(contKey string) string {
	if contKey == "" {
		return "N"
	}
	return "Y"
}
