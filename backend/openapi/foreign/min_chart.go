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

// ForeignMinChartService 해외주식 분차트조회 서비스
type ForeignMinChartService struct {
	client *client.DBSecClient
	logger logger.Logger
}

// NewForeignMinChartService 새로운 해외주식 분차트조회 서비스 생성
func NewForeignMinChartService(client *client.DBSecClient) *ForeignMinChartService {
	return &ForeignMinChartService{
		client: client,
		logger: logger.GetDefaultLogger().With(logger.Field{Key: "service", Value: "foreign_min_chart"}),
	}
}

// GetMinChart 해외주식 분차트 데이터 조회
func (s *ForeignMinChartService) GetMinChart(stockCode string, period models.ChartPeriod, options models.ChartOptions) ([]models.ForeignMinChartData, error) {
	s.logger.Info("Getting foreign stock min chart", 
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
	respBody, err := s.client.MakeRequestWithHeaders("POST", models.PathForeignStockMinChart, nil, request, nil)
	if err != nil {
		s.logger.Error("Failed to call min chart API", err, 
			logger.Field{Key: "stock_code", Value: stockCode})
		return nil, errors.NewNetworkError("failed to call min chart API", err)
	}

	// 응답 파싱
	var response models.ForeignMinChartResponse
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

	s.logger.Info("Successfully retrieved min chart data", 
		logger.Field{Key: "stock_code", Value: stockCode},
		logger.Field{Key: "data_count", Value: len(chartData)})

	return chartData, nil
}

// GetMinChartWithOptions 옵션을 사용한 분차트 조회 (편의 메서드)
func (s *ForeignMinChartService) GetMinChartWithOptions(stockCode, market, interval string, days int, useAdjusted bool) ([]models.ForeignMinChartData, error) {
	endDate := time.Now().Format("20060102")
	startDate := time.Now().AddDate(0, 0, -days).Format("20060102")

	period := models.ChartPeriod{
		StartDate: startDate,
		EndDate:   endDate,
		IsRange:   true,
	}

	options := models.ChartOptions{
		Interval:    interval,
		UseAdjusted: useAdjusted,
		Market:      market,
		DataCount:   0, // 기본값 사용
	}

	return s.GetMinChart(stockCode, period, options)
}

// GetLatestMinChart 최근 데이터 조회 (기간 미지정)
func (s *ForeignMinChartService) GetLatestMinChart(stockCode, market, interval string, dataCount int) ([]models.ForeignMinChartData, error) {
	endDate := time.Now().Format("20060102")

	period := models.ChartPeriod{
		EndDate: endDate,
		IsRange: false, // 기간 미지정
	}

	options := models.ChartOptions{
		Interval:    interval,
		UseAdjusted: true, // 기본적으로 수정주가 사용
		Market:      market,
		DataCount:   dataCount,
	}

	return s.GetMinChart(stockCode, period, options)
}

// GetNASDAQMinChart 나스닥 종목 분차트 조회
func (s *ForeignMinChartService) GetNASDAQMinChart(stockCode, interval string, days int) ([]models.ForeignMinChartData, error) {
	return s.GetMinChartWithOptions(stockCode, "NASDAQ", interval, days, true)
}

// GetNYMinChart 뉴욕 증권거래소 종목 분차트 조회
func (s *ForeignMinChartService) GetNYMinChart(stockCode, interval string, days int) ([]models.ForeignMinChartData, error) {
	return s.GetMinChartWithOptions(stockCode, "NY", interval, days, true)
}

// GetAMEXMinChart 아멕스 종목 분차트 조회
func (s *ForeignMinChartService) GetAMEXMinChart(stockCode, interval string, days int) ([]models.ForeignMinChartData, error) {
	return s.GetMinChartWithOptions(stockCode, "AMEX", interval, days, true)
}

// GetPopularStocksMinChart 인기 종목들의 분차트 조회
func (s *ForeignMinChartService) GetPopularStocksMinChart(interval string, days int) (map[string][]models.ForeignMinChartData, error) {
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
	}

	results := make(map[string][]models.ForeignMinChartData)

	for _, stock := range popularStocks {
		data, err := s.GetMinChartWithOptions(stock.code, stock.market, interval, days, true)
		if err != nil {
			s.logger.Warn("Failed to get chart data for stock", 
				logger.Field{Key: "stock_code", Value: stock.code},
				logger.Field{Key: "error", Value: err.Error()})
			continue
		}
		results[stock.code] = data
	}

	return results, nil
}

// validateInputs 입력값 검증
func (s *ForeignMinChartService) validateInputs(stockCode string, period models.ChartPeriod, options models.ChartOptions) error {
	if stockCode == "" {
		return errors.NewValidationError("stock code is required", nil)
	}

	if options.Market == "" {
		return errors.NewValidationError("market is required", nil)
	}

	if options.Interval == "" {
		return errors.NewValidationError("interval is required", nil)
	}

	// 날짜 형식 검증 (기간 지정시)
	if period.IsRange {
		if period.StartDate == "" || period.EndDate == "" {
			return errors.NewValidationError("start_date and end_date are required when period is specified", nil)
		}
		
		if len(period.StartDate) != 8 || len(period.EndDate) != 8 {
			return errors.NewValidationError("date format should be YYYYMMDD", nil)
		}
	} else {
		if period.EndDate == "" {
			return errors.NewValidationError("end_date is required", nil)
		}
	}

	return nil
}

// buildRequest API 요청 데이터 구성
func (s *ForeignMinChartService) buildRequest(stockCode string, period models.ChartPeriod, options models.ChartOptions) models.ForeignMinChartRequest {
	input := models.ForeignMinChartInput{
		InputCondMrktDivCode: options.GetMarketCode(),
		InputIscd1:           stockCode,
		InputHourClsCode:     models.HourClassCode,
		InputDivXtick:        options.GetIntervalCode(),
		InputOrgAdjPrc:       options.GetAdjustedCode(),
		DataCnt:              options.GetDataCountString(),
	}

	if period.IsRange {
		input.InputDate1 = period.StartDate
		input.InputDate2 = period.EndDate
		input.InputPwDataIncuYn = models.PeriodSpecified
	} else {
		input.InputDate1 = ""
		input.InputDate2 = period.EndDate
		input.InputPwDataIncuYn = models.PeriodNotSpecified
	}

	return models.ForeignMinChartRequest{In: input}
}

// convertToChartData API 응답을 비즈니스 모델로 변환
func (s *ForeignMinChartService) convertToChartData(stockCode string, outputs []models.ForeignMinChartOutput, options models.ChartOptions) []models.ForeignMinChartData {
	var chartData []models.ForeignMinChartData

	for _, output := range outputs {
		data := models.ForeignMinChartData{
			StockCode:    stockCode,
			DateTime:     s.formatDateTime(output.Date, output.Hour),
			Date:         s.formatDate(output.Date),
			Time:         s.formatTime(output.Hour),
			Open:         utils.ParseFloat(output.Oprc),
			High:         utils.ParseFloat(output.Hprc),
			Low:          utils.ParseFloat(output.Lprc),
			Close:        utils.ParseFloat(output.Prpr),
			Volume:       utils.ParseInt(output.CntgVol),
			Market:       s.getMarketName(options.GetMarketCode()),
			MarketCode:   options.GetMarketCode(),
			Interval:     options.Interval,
			IntervalCode: options.GetIntervalCode(),
			IsAdjusted:   options.UseAdjusted,
		}
		chartData = append(chartData, data)
	}

	return chartData
}

// formatDateTime 날짜와 시간을 ISO 형식으로 변환
func (s *ForeignMinChartService) formatDateTime(date, hour string) string {
	if len(date) != 8 || len(hour) != 6 {
		return ""
	}
	
	dateFormatted := fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:8])
	timeFormatted := fmt.Sprintf("%s:%s:%s", hour[:2], hour[2:4], hour[4:6])
	
	return fmt.Sprintf("%s %s", dateFormatted, timeFormatted)
}

// formatDate 날짜를 YYYY-MM-DD 형식으로 변환
func (s *ForeignMinChartService) formatDate(date string) string {
	if len(date) != 8 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:8])
}

// formatTime 시간을 HH:MM:SS 형식으로 변환
func (s *ForeignMinChartService) formatTime(hour string) string {
	if len(hour) != 6 {
		return ""
	}
	return fmt.Sprintf("%s:%s:%s", hour[:2], hour[2:4], hour[4:6])
}

// getMarketName 시장 코드를 시장명으로 변환
func (s *ForeignMinChartService) getMarketName(marketCode string) string {
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

// GetIntervalDescription 시간간격 코드를 설명으로 변환
func (s *ForeignMinChartService) GetIntervalDescription(intervalCode string) string {
	switch intervalCode {
	case models.ChartInterval30Sec:
		return "30초"
	case models.ChartInterval1Min:
		return "1분"
	case models.ChartInterval2Min:
		return "2분"
	case models.ChartInterval5Min:
		return "5분"
	case models.ChartInterval10Min:
		return "10분"
	case models.ChartInterval60Min:
		return "60분"
	default:
		return "알 수 없음"
	}
}