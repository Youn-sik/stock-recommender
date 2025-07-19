package foreign

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
)

// ForeignStockTickerService 해외주식종목 조회 서비스
type ForeignStockTickerService struct {
	client *client.DBSecClient
}

// NewForeignStockTickerService 새로운 해외주식종목 조회 서비스 생성
func NewForeignStockTickerService(client *client.DBSecClient) *ForeignStockTickerService {
	return &ForeignStockTickerService{
		client: client,
	}
}

// GetForeignStockTickers 해외주식종목 조회
// exchangeCode: 해외증시구분코드 (NY: 뉴욕, NA: 나스닥, AM: 아멕스)
// contKey: 연속키 (optional, 추가 데이터 조회시 사용)
func (s *ForeignStockTickerService) GetForeignStockTickers(exchangeCode string, contKey string) (*models.ForeignStockTickerResponse, string, error) {
	// 요청 데이터 구성
	reqBody := models.ForeignStockTickerRequest{
		In: models.ForeignStockTickerInput{
			InputDataCode: exchangeCode,
		},
	}

	// 헤더 설정을 위한 커스텀 요청 함수
	respBody, nextContKey, err := s.makeRequestWithContKey(reqBody, contKey)
	if err != nil {
		return nil, "", err
	}

	// 응답 파싱
	var response models.ForeignStockTickerResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, "", fmt.Errorf("failed to parse response: %w", err)
	}

	// 응답 코드 확인
	if response.RspCd != "00000" {
		return nil, "", fmt.Errorf("API error %s: %s", response.RspCd, response.RspMsg)
	}

	return &response, nextContKey, nil
}

// GetAllForeignStockTickers 모든 해외주식종목 조회 (페이지네이션 포함)
func (s *ForeignStockTickerService) GetAllForeignStockTickers(exchangeCode string) ([]models.ForeignStockData, error) {
	var allStocks []models.ForeignStockData
	contKey := ""

	for {
		response, nextContKey, err := s.GetForeignStockTickers(exchangeCode, contKey)
		if err != nil {
			return nil, err
		}

		// 응답 데이터를 변환된 형식으로 변환
		for _, stock := range response.Out {
			data := s.convertToForeignStockData(exchangeCode, &stock)
			allStocks = append(allStocks, *data)
		}

		// 더 이상 데이터가 없으면 종료
		if nextContKey == "" || nextContKey == "N" {
			break
		}

		contKey = nextContKey
	}

	return allStocks, nil
}

// GetNYStocks 뉴욕 거래소 종목 조회
func (s *ForeignStockTickerService) GetNYStocks() ([]models.ForeignStockData, error) {
	return s.GetAllForeignStockTickers(models.ExchangeNY)
}

// GetNASDAQStocks 나스닥 거래소 종목 조회
func (s *ForeignStockTickerService) GetNASDAQStocks() ([]models.ForeignStockData, error) {
	return s.GetAllForeignStockTickers(models.ExchangeNASDAQ)
}

// GetAMEXStocks 아멕스 거래소 종목 조회
func (s *ForeignStockTickerService) GetAMEXStocks() ([]models.ForeignStockData, error) {
	return s.GetAllForeignStockTickers(models.ExchangeAMEX)
}

// GetAllUSStocks 모든 미국 주식 종목 조회 (NY + NASDAQ + AMEX)
func (s *ForeignStockTickerService) GetAllUSStocks() (map[string][]models.ForeignStockData, error) {
	result := make(map[string][]models.ForeignStockData)
	
	// 뉴욕 거래소
	nyStocks, err := s.GetNYStocks()
	if err != nil {
		return nil, fmt.Errorf("failed to get NY stocks: %w", err)
	}
	result["NY"] = nyStocks

	// 나스닥
	nasdaqStocks, err := s.GetNASDAQStocks()
	if err != nil {
		return nil, fmt.Errorf("failed to get NASDAQ stocks: %w", err)
	}
	result["NASDAQ"] = nasdaqStocks

	// 아멕스
	amexStocks, err := s.GetAMEXStocks()
	if err != nil {
		return nil, fmt.Errorf("failed to get AMEX stocks: %w", err)
	}
	result["AMEX"] = amexStocks

	return result, nil
}

// GetStocksBySector 업종별 종목 조회
func (s *ForeignStockTickerService) GetStocksBySector(exchangeCode string, sectorName string) ([]models.ForeignStockData, error) {
	allStocks, err := s.GetAllForeignStockTickers(exchangeCode)
	if err != nil {
		return nil, err
	}

	var filteredStocks []models.ForeignStockData
	for _, stock := range allStocks {
		if strings.Contains(stock.SectorName, sectorName) {
			filteredStocks = append(filteredStocks, stock)
		}
	}

	return filteredStocks, nil
}

// GetTechStocks IT 업종 종목 조회 (나스닥 기준)
func (s *ForeignStockTickerService) GetTechStocks() ([]models.ForeignStockData, error) {
	return s.GetStocksBySector(models.ExchangeNASDAQ, "IT")
}

// makeRequestWithContKey 연속키를 포함한 요청 처리
func (s *ForeignStockTickerService) makeRequestWithContKey(reqBody interface{}, contKey string) ([]byte, string, error) {
	// 추가 헤더 설정
	headers := map[string]string{
		"cont_yn":  s.getContYn(contKey),
		"cont_key": contKey,
		"tr_id":    models.TrIdForeignStockTicker,
	}

	// API 호출
	response, err := s.client.MakeRequestWithFullResponse("POST", models.PathForeignStockTicker, nil, reqBody, headers)
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
func (s *ForeignStockTickerService) getContYn(contKey string) string {
	if contKey == "" {
		return "N"
	}
	return "Y"
}

// convertToForeignStockData 응답 데이터를 구조화된 형식으로 변환
func (s *ForeignStockTickerService) convertToForeignStockData(exchangeCode string, output *models.ForeignStockTickerOutput) *models.ForeignStockData {
	return &models.ForeignStockData{
		StockCode:    output.Iscd,
		KoreanName:   output.KorIsnm,
		SectorName:   output.BstpLargName,
		ExchangeCode: output.ExchClsCode2,
		Exchange:     s.getExchangeName(exchangeCode),
		SellUnit:     s.parseInt(output.SelnVolUnit),
		BuyUnit:      s.parseInt(output.ShnuVolUnit),
	}
}

// getExchangeName 거래소 코드를 거래소명으로 변환
func (s *ForeignStockTickerService) getExchangeName(exchangeCode string) string {
	switch exchangeCode {
	case models.ExchangeNY:
		return "뉴욕"
	case models.ExchangeNASDAQ:
		return "나스닥"
	case models.ExchangeAMEX:
		return "아멕스"
	default:
		return exchangeCode
	}
}

// parseInt 문자열을 정수로 변환
func (s *ForeignStockTickerService) parseInt(str string) int64 {
	if str == "" {
		return 0
	}
	val, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		return 0
	}
	return val
}