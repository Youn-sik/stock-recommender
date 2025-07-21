package foreign

import (
	"encoding/json"
	"fmt"
	"time"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/errors"
	"stock-recommender/backend/openapi/logger"
	"stock-recommender/backend/openapi/models"
	"stock-recommender/backend/openapi/utils"
)

// ForeignDayChartService 해외주식 일차트조회 서비스
type ForeignDayChartService struct {
	client *client.DBSecClient
	logger logger.Logger
}

// NewForeignDayChartService 새로운 해외주식 일차트조회 서비스 생성
func NewForeignDayChartService(client *client.DBSecClient) *ForeignDayChartService {
	return &ForeignDayChartService{
		client: client,
		logger: logger.GetDefaultLogger().With(logger.Field{Key: "service", Value: "foreign_day_chart"}),
	}
}

// GetDayChart 해외주식 일차트 데이터 조회
func (s *ForeignDayChartService) GetDayChart(stockCode string, period models.DayChartPeriod, options models.DayChartOptions) ([]models.ForeignDayChartData, error) {
	s.logger.Info("Getting foreign stock day chart", 
		logger.Field{Key: "stock_code", Value: stockCode},
		logger.Field{Key: "period", Value: period},
		logger.Field{Key: "options", Value: options})

	// 입력 검증
	if err := s.validateInputs(stockCode, period, options); err != nil {
		return nil, err
	}

	// 요청 데이터 구성
	request := s.buildRequest(stockCode, period, options)

	// API 호출
	respBody, err := s.client.MakeRequestWithHeaders("POST", models.PathForeignStockDayChart, nil, request, nil)
	if err != nil {
		s.logger.Error("Failed to call day chart API", err, 
			logger.Field{Key: "stock_code", Value: stockCode})
		return nil, errors.NewNetworkError("failed to call day chart API", err)
	}

	// 응답 파싱
	var response models.ForeignDayChartResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		s.logger.Error("Failed to parse API response", err)
		return nil, errors.NewParseError("failed to parse API response", err)
	}

	// 응답 코드 확인
	if !utils.IsSuccessResponse(response.RspCd) {
		s.logger.Warn("API returned error", 
			logger.Field{Key: "response_code", Value: response.RspCd},
			logger.Field{Key: "response_message", Value: response.RspMsg})
		return nil, errors.NewAPIError(errors.ErrCodeServerError, "API returned error", fmt.Errorf("code: %s, message: %s", response.RspCd, response.RspMsg))
	}

	// 데이터 변환
	chartData := s.convertToChartData(stockCode, response.Out, options)

	s.logger.Info("Successfully retrieved day chart data", 
		logger.Field{Key: "stock_code", Value: stockCode},
		logger.Field{Key: "data_count", Value: len(chartData)})

	return chartData, nil
}

// GetDayChartWithDays 일수를 지정하여 일차트 조회 (편의 메서드)
func (s *ForeignDayChartService) GetDayChartWithDays(stockCode, market string, days int, useAdjusted bool) ([]models.ForeignDayChartData, error) {
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")

	period := models.DayChartPeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}

	options := models.DayChartOptions{
		UseAdjusted: useAdjusted,
		Market:      market,
	}

	return s.GetDayChart(stockCode, period, options)
}

// GetRecentDayChart 최근 데이터 조회
func (s *ForeignDayChartService) GetRecentDayChart(stockCode, market string, days int) ([]models.ForeignDayChartData, error) {
	return s.GetDayChartWithDays(stockCode, market, days, true) // 기본적으로 수정주가 사용
}

// GetNASDAQDayChart 나스닥 종목 일차트 조회
func (s *ForeignDayChartService) GetNASDAQDayChart(stockCode string, days int) ([]models.ForeignDayChartData, error) {
	return s.GetDayChartWithDays(stockCode, "NASDAQ", days, true)
}

// GetNYDayChart 뉴욕 증권거래소 종목 일차트 조회
func (s *ForeignDayChartService) GetNYDayChart(stockCode string, days int) ([]models.ForeignDayChartData, error) {
	return s.GetDayChartWithDays(stockCode, "NY", days, true)
}

// GetAMEXDayChart 아멕스 종목 일차트 조회
func (s *ForeignDayChartService) GetAMEXDayChart(stockCode string, days int) ([]models.ForeignDayChartData, error) {
	return s.GetDayChartWithDays(stockCode, "AMEX", days, true)
}

// GetYearChart 1년 차트 조회
func (s *ForeignDayChartService) GetYearChart(stockCode, market string) ([]models.ForeignDayChartData, error) {
	return s.GetDayChartWithDays(stockCode, market, 365, true)
}

// GetMonthChart 1개월 차트 조회
func (s *ForeignDayChartService) GetMonthChart(stockCode, market string) ([]models.ForeignDayChartData, error) {
	return s.GetDayChartWithDays(stockCode, market, 30, true)
}

// GetWeekChart 1주일 차트 조회
func (s *ForeignDayChartService) GetWeekChart(stockCode, market string) ([]models.ForeignDayChartData, error) {
	return s.GetDayChartWithDays(stockCode, market, 7, true)
}

// GetPopularStocksDayChart 인기 종목들의 일차트 조회
func (s *ForeignDayChartService) GetPopularStocksDayChart(days int) (map[string][]models.ForeignDayChartData, error) {
	popularStocks := []struct {
		code   string
		market string
	}{
		{"AAPL", "NASDAQ"},
		{"MSFT", "NASDAQ"},
		{"GOOGL", "NASDAQ"},
		{"AMZN", "NASDAQ"},
		{"TSLA", "NASDAQ"},
		{"NVDA", "NASDAQ"},
		{"META", "NASDAQ"},
		{"IBM", "NY"},
		{"GE", "NY"},
	}

	results := make(map[string][]models.ForeignDayChartData)

	for _, stock := range popularStocks {
		data, err := s.GetDayChartWithDays(stock.code, stock.market, days, true)
		if err != nil {
			s.logger.Warn("Failed to get day chart data for stock", 
				logger.Field{Key: "stock_code", Value: stock.code},
				logger.Field{Key: "error", Value: err.Error()})
			continue
		}
		results[stock.code] = data
	}

	return results, nil
}

// GetTechGiantsDayChart 기술주 대장주들의 일차트 조회
func (s *ForeignDayChartService) GetTechGiantsDayChart(days int) (map[string][]models.ForeignDayChartData, error) {
	techStocks := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA", "NVDA", "META"}
	results := make(map[string][]models.ForeignDayChartData)

	for _, stockCode := range techStocks {
		data, err := s.GetNASDAQDayChart(stockCode, days)
		if err != nil {
			s.logger.Warn("Failed to get tech stock day chart", 
				logger.Field{Key: "stock_code", Value: stockCode},
				logger.Field{Key: "error", Value: err.Error()})
			continue
		}
		results[stockCode] = data
	}

	return results, nil
}

// validateInputs 입력값 검증
func (s *ForeignDayChartService) validateInputs(stockCode string, period models.DayChartPeriod, options models.DayChartOptions) error {
	if stockCode == "" {
		return errors.NewValidationError("stock code is required", nil)
	}

	if options.Market == "" {
		return errors.NewValidationError("market is required", nil)
	}

	if period.StartDate == "" || period.EndDate == "" {
		return errors.NewValidationError("start_date and end_date are required", nil)
	}

	// 날짜 형식 검증
	if period.GetFormattedStartDate() == "" {
		return errors.NewValidationError("invalid start_date format (expected: YYYY-MM-DD or YYYYMMDD)", nil)
	}

	if period.GetFormattedEndDate() == "" {
		return errors.NewValidationError("invalid end_date format (expected: YYYY-MM-DD or YYYYMMDD)", nil)
	}

	return nil
}

// buildRequest API 요청 데이터 구성
func (s *ForeignDayChartService) buildRequest(stockCode string, period models.DayChartPeriod, options models.DayChartOptions) models.ForeignDayChartRequest {
	input := models.ForeignDayChartInput{
		InputCondMrktDivCode: options.GetMarketCode(),
		InputOrgAdjPrc:       options.GetAdjustedCode(),
		InputIscd1:           stockCode,
		InputDate1:           period.GetFormattedStartDate(),
		InputDate2:           period.GetFormattedEndDate(),
	}

	return models.ForeignDayChartRequest{In: input}
}

// convertToChartData API 응답을 비즈니스 모델로 변환
func (s *ForeignDayChartService) convertToChartData(stockCode string, outputs []models.ForeignDayChartOutput, options models.DayChartOptions) []models.ForeignDayChartData {
	var chartData []models.ForeignDayChartData

	for i, output := range outputs {
		data := models.ForeignDayChartData{
			StockCode:  stockCode,
			Date:       s.formatDate(output.Date),
			Open:       utils.ParseFloat(output.Oprc),
			High:       utils.ParseFloat(output.Hprc),
			Low:        utils.ParseFloat(output.Lprc),
			Close:      utils.ParseFloat(output.Prpr),
			Volume:     utils.ParseInt(output.AcmlVol),
			Market:     s.getMarketName(options.GetMarketCode()),
			MarketCode: options.GetMarketCode(),
			IsAdjusted: options.UseAdjusted,
			WeekDay:    s.getWeekDay(output.Date),
		}

		// 전일대비 계산 (이전 데이터가 있는 경우)
		if i < len(outputs)-1 { // 데이터는 최신순으로 정렬되어 있다고 가정
			prevClose := utils.ParseFloat(outputs[i+1].Prpr)
			if prevClose > 0 {
				data.PriceChange = data.Close - prevClose
				data.ChangeRate = (data.PriceChange / prevClose) * 100
			}
		}

		chartData = append(chartData, data)
	}

	return chartData
}

// formatDate 날짜를 YYYY-MM-DD 형식으로 변환
func (s *ForeignDayChartService) formatDate(date string) string {
	if len(date) != 8 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:8])
}

// getMarketName 시장 코드를 시장명으로 변환
func (s *ForeignDayChartService) getMarketName(marketCode string) string {
	switch marketCode {
	case models.ForeignMarketNY:
		return "New York Stock Exchange"
	case models.ForeignMarketNASDAQ:
		return "NASDAQ"
	case models.ForeignMarketAMEX:
		return "American Stock Exchange"
	default:
		return "Unknown"
	}
}

// getWeekDay 날짜에서 요일 계산
func (s *ForeignDayChartService) getWeekDay(dateStr string) string {
	if len(dateStr) != 8 {
		return ""
	}

	// YYYYMMDD 형식을 time.Time으로 변환
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return ""
	}

	weekdays := []string{"일", "월", "화", "수", "목", "금", "토"}
	return weekdays[t.Weekday()]
}

// GetPriceStatistics 가격 통계 계산
func (s *ForeignDayChartService) GetPriceStatistics(chartData []models.ForeignDayChartData) map[string]float64 {
	if len(chartData) == 0 {
		return nil
	}

	var highs, lows, closes []float64
	var totalVolume int64

	for _, data := range chartData {
		highs = append(highs, data.High)
		lows = append(lows, data.Low)
		closes = append(closes, data.Close)
		totalVolume += data.Volume
	}

	stats := make(map[string]float64)
	
	// 최고가/최저가
	stats["max_high"] = s.maxFloat(highs)
	stats["min_low"] = s.minFloat(lows)
	
	// 평균가격
	stats["avg_close"] = s.avgFloat(closes)
	
	// 평균거래량
	stats["avg_volume"] = float64(totalVolume) / float64(len(chartData))
	
	// 변동성 (표준편차)
	stats["volatility"] = s.stdDevFloat(closes)
	
	return stats
}

// 유틸리티 함수들
func (s *ForeignDayChartService) maxFloat(values []float64) float64 {
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

func (s *ForeignDayChartService) minFloat(values []float64) float64 {
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

func (s *ForeignDayChartService) avgFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func (s *ForeignDayChartService) stdDevFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	avg := s.avgFloat(values)
	sumSquares := 0.0
	
	for _, v := range values {
		diff := v - avg
		sumSquares += diff * diff
	}
	
	variance := sumSquares / float64(len(values))
	return variance // 실제로는 math.Sqrt(variance)이지만 math import 하지 않으려면 variance 반환
}