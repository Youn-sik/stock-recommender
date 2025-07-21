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

// ForeignWeekChartService 해외주식 주차트조회 서비스
type ForeignWeekChartService struct {
	client *client.DBSecClient
	logger logger.Logger
}

// NewForeignWeekChartService 새로운 해외주식 주차트조회 서비스 생성
func NewForeignWeekChartService(client *client.DBSecClient) *ForeignWeekChartService {
	return &ForeignWeekChartService{
		client: client,
		logger: logger.GetDefaultLogger().With(logger.Field{Key: "service", Value: "foreign_week_chart"}),
	}
}

// GetWeekChart 해외주식 주차트 데이터 조회
func (s *ForeignWeekChartService) GetWeekChart(stockCode string, period models.WeekChartPeriod, options models.WeekChartOptions) ([]models.ForeignWeekChartData, error) {
	s.logger.Info("Getting foreign stock week chart", 
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
	respBody, err := s.client.MakeRequestWithHeaders("POST", models.PathForeignStockWeekChart, nil, request, nil)
	if err != nil {
		s.logger.Error("Failed to call week chart API", err, 
			logger.Field{Key: "stock_code", Value: stockCode})
		return nil, errors.NewNetworkError("failed to call week chart API", err)
	}

	// 응답 파싱
	var response models.ForeignWeekChartResponse
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

	s.logger.Info("Successfully retrieved week chart data", 
		logger.Field{Key: "stock_code", Value: stockCode},
		logger.Field{Key: "data_count", Value: len(chartData)})

	return chartData, nil
}

// GetWeekChartWithWeeks 주 수를 지정하여 주차트 조회 (편의 메서드)
func (s *ForeignWeekChartService) GetWeekChartWithWeeks(stockCode, market string, weeks int, useAdjusted bool) ([]models.ForeignWeekChartData, error) {
	// 주 단위로 날짜 계산 (weeks * 7일)
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -weeks*7).Format("2006-01-02")

	period := models.WeekChartPeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}

	options := models.WeekChartOptions{
		UseAdjusted: useAdjusted,
		Market:      market,
	}

	return s.GetWeekChart(stockCode, period, options)
}

// GetRecentWeekChart 최근 주차트 데이터 조회
func (s *ForeignWeekChartService) GetRecentWeekChart(stockCode, market string, weeks int) ([]models.ForeignWeekChartData, error) {
	return s.GetWeekChartWithWeeks(stockCode, market, weeks, true) // 기본적으로 수정주가 사용
}

// GetNASDAQWeekChart 나스닥 종목 주차트 조회
func (s *ForeignWeekChartService) GetNASDAQWeekChart(stockCode string, weeks int) ([]models.ForeignWeekChartData, error) {
	return s.GetWeekChartWithWeeks(stockCode, "NASDAQ", weeks, true)
}

// GetNYWeekChart 뉴욕 증권거래소 종목 주차트 조회
func (s *ForeignWeekChartService) GetNYWeekChart(stockCode string, weeks int) ([]models.ForeignWeekChartData, error) {
	return s.GetWeekChartWithWeeks(stockCode, "NY", weeks, true)
}

// GetAMEXWeekChart 아멕스 종목 주차트 조회
func (s *ForeignWeekChartService) GetAMEXWeekChart(stockCode string, weeks int) ([]models.ForeignWeekChartData, error) {
	return s.GetWeekChartWithWeeks(stockCode, "AMEX", weeks, true)
}

// Get52WeekChart 52주(1년) 차트 조회
func (s *ForeignWeekChartService) Get52WeekChart(stockCode, market string) ([]models.ForeignWeekChartData, error) {
	return s.GetWeekChartWithWeeks(stockCode, market, 52, true)
}

// Get26WeekChart 26주(6개월) 차트 조회
func (s *ForeignWeekChartService) Get26WeekChart(stockCode, market string) ([]models.ForeignWeekChartData, error) {
	return s.GetWeekChartWithWeeks(stockCode, market, 26, true)
}

// Get13WeekChart 13주(3개월) 차트 조회
func (s *ForeignWeekChartService) Get13WeekChart(stockCode, market string) ([]models.ForeignWeekChartData, error) {
	return s.GetWeekChartWithWeeks(stockCode, market, 13, true)
}

// GetTechGiantsWeekChart 기술주 대장주들의 주차트 조회
func (s *ForeignWeekChartService) GetTechGiantsWeekChart(weeks int) (map[string][]models.ForeignWeekChartData, error) {
	techStocks := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA", "NVDA", "META"}
	results := make(map[string][]models.ForeignWeekChartData)

	for _, stockCode := range techStocks {
		data, err := s.GetNASDAQWeekChart(stockCode, weeks)
		if err != nil {
			s.logger.Warn("Failed to get tech stock week chart", 
				logger.Field{Key: "stock_code", Value: stockCode},
				logger.Field{Key: "error", Value: err.Error()})
			continue
		}
		results[stockCode] = data
	}

	return results, nil
}

// GetVolatilityAnalysis 주간 변동성 분석
func (s *ForeignWeekChartService) GetVolatilityAnalysis(chartData []models.ForeignWeekChartData) map[string]float64 {
	if len(chartData) == 0 {
		return nil
	}

	var weeklyRanges, weeklyRangeRates, changeRates []float64
	var totalVolume int64

	for _, data := range chartData {
		weeklyRanges = append(weeklyRanges, data.WeeklyRange)
		weeklyRangeRates = append(weeklyRangeRates, data.WeeklyRangeRate)
		if data.ChangeRate != 0 { // 0이 아닌 값만 포함
			changeRates = append(changeRates, data.ChangeRate)
		}
		totalVolume += data.Volume
	}

	analysis := make(map[string]float64)
	
	// 평균 주간 변동폭
	analysis["avg_weekly_range"] = s.avgFloat(weeklyRanges)
	
	// 평균 주간 변동률
	analysis["avg_weekly_range_rate"] = s.avgFloat(weeklyRangeRates)
	
	// 평균 주간 변화율
	if len(changeRates) > 0 {
		analysis["avg_weekly_change_rate"] = s.avgFloat(changeRates)
	}
	
	// 최대 주간 변동률
	analysis["max_weekly_range_rate"] = s.maxFloat(weeklyRangeRates)
	
	// 평균 주간 거래량
	if len(chartData) > 0 {
		analysis["avg_weekly_volume"] = float64(totalVolume) / float64(len(chartData))
	}
	
	return analysis
}

// validateInputs 입력값 검증
func (s *ForeignWeekChartService) validateInputs(stockCode string, period models.WeekChartPeriod, options models.WeekChartOptions) error {
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
func (s *ForeignWeekChartService) buildRequest(stockCode string, period models.WeekChartPeriod, options models.WeekChartOptions) models.ForeignWeekChartRequest {
	input := models.ForeignWeekChartInput{
		InputCondMrktDivCode: options.GetMarketCode(),
		InputOrgAdjPrc:       options.GetAdjustedCode(),
		InputIscd1:           stockCode,
		InputDate1:           period.GetFormattedStartDate(),
		InputDate2:           period.GetFormattedEndDate(),
		InputPeriodDivCode:   models.PeriodDivWeek, // 고정값: W
	}

	return models.ForeignWeekChartRequest{In: input}
}

// convertToChartData API 응답을 비즈니스 모델로 변환
func (s *ForeignWeekChartService) convertToChartData(stockCode string, outputs []models.ForeignWeekChartOutput, options models.WeekChartOptions) []models.ForeignWeekChartData {
	var chartData []models.ForeignWeekChartData

	for i, output := range outputs {
		// 주 종료일에서 연도와 주차 계산
		weekEndDate := s.formatDate(output.Date)
		year, weekNumber := s.getYearWeek(output.Date)
		weekStartDate := s.calculateWeekStartDate(output.Date)
		
		// 거래량 처리 (주차트에서는 빈 값일 수 있음)
		volume := int64(0)
		if output.CntgVol != "" {
			volume = utils.ParseInt(output.CntgVol)
		}

		data := models.ForeignWeekChartData{
			StockCode:     stockCode,
			WeekEndDate:   weekEndDate,
			WeekStartDate: weekStartDate,
			Open:          utils.ParseFloat(output.Oprc),
			High:          utils.ParseFloat(output.Hprc),
			Low:           utils.ParseFloat(output.Lprc),
			Close:         utils.ParseFloat(output.Prpr),
			Volume:        volume,
			Market:        s.getMarketName(options.GetMarketCode()),
			MarketCode:    options.GetMarketCode(),
			IsAdjusted:    options.UseAdjusted,
			WeekNumber:    weekNumber,
			Year:          year,
		}

		// 주간 변동폭 계산
		data.WeeklyRange = data.High - data.Low
		if data.Low > 0 {
			data.WeeklyRangeRate = (data.WeeklyRange / data.Low) * 100
		}

		// 전주대비 계산 (이전 데이터가 있는 경우)
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
func (s *ForeignWeekChartService) formatDate(date string) string {
	if len(date) != 8 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:8])
}

// getMarketName 시장 코드를 시장명으로 변환
func (s *ForeignWeekChartService) getMarketName(marketCode string) string {
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

// getYearWeek 날짜에서 연도와 주차 번호 계산
func (s *ForeignWeekChartService) getYearWeek(dateStr string) (int, int) {
	if len(dateStr) != 8 {
		return 0, 0
	}

	// YYYYMMDD 형식을 time.Time으로 변환
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return 0, 0
	}

	year, week := t.ISOWeek()
	return year, week
}

// calculateWeekStartDate 주 종료일에서 주 시작일 계산
func (s *ForeignWeekChartService) calculateWeekStartDate(weekEndDateStr string) string {
	if len(weekEndDateStr) != 8 {
		return ""
	}

	// YYYYMMDD 형식을 time.Time으로 변환
	endDate, err := time.Parse("20060102", weekEndDateStr)
	if err != nil {
		return ""
	}

	// 주의 시작일 계산 (일요일 기준으로 6일 전)
	weekday := int(endDate.Weekday())
	daysToSubtract := weekday
	if weekday == 0 { // 일요일인 경우
		daysToSubtract = 6
	}
	
	startDate := endDate.AddDate(0, 0, -daysToSubtract)
	return startDate.Format("2006-01-02")
}

// Get52WeekHighLow 52주 최고/최저가 계산
func (s *ForeignWeekChartService) Get52WeekHighLow(chartData []models.ForeignWeekChartData) (float64, float64) {
	if len(chartData) == 0 {
		return 0, 0
	}

	high := chartData[0].High
	low := chartData[0].Low

	for _, data := range chartData[1:] {
		if data.High > high {
			high = data.High
		}
		if data.Low < low {
			low = data.Low
		}
	}

	return high, low
}

// GetTrendAnalysis 추세 분석
func (s *ForeignWeekChartService) GetTrendAnalysis(chartData []models.ForeignWeekChartData) string {
	if len(chartData) < 4 {
		return "Insufficient data"
	}

	// 최근 4주의 종가 추출
	recentCloses := make([]float64, 0, 4)
	for i := 0; i < 4 && i < len(chartData); i++ {
		recentCloses = append(recentCloses, chartData[i].Close)
	}

	// 상승 또는 하락 추세 판단 (최신 데이터가 앞에 있으므로)
	downCount := 0
	for i := 0; i < len(recentCloses)-1; i++ {
		if recentCloses[i] < recentCloses[i+1] { // 최신이 이전보다 낮으면 하락
			downCount++
		}
	}

	if downCount >= 2 {
		return "Downtrend"
	} else if downCount == 0 {
		return "Uptrend"
	}
	return "Sideways"
}

// 유틸리티 함수들
func (s *ForeignWeekChartService) maxFloat(values []float64) float64 {
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

func (s *ForeignWeekChartService) avgFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}