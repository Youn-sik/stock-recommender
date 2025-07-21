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

// ForeignMonthChartService 해외주식 월차트조회 서비스
type ForeignMonthChartService struct {
	client *client.DBSecClient
	logger logger.Logger
}

// NewForeignMonthChartService 새로운 해외주식 월차트조회 서비스 생성
func NewForeignMonthChartService(client *client.DBSecClient) *ForeignMonthChartService {
	return &ForeignMonthChartService{
		client: client,
		logger: logger.GetDefaultLogger().With(logger.Field{Key: "service", Value: "foreign_month_chart"}),
	}
}

// GetMonthChart 해외주식 월차트 데이터 조회
func (s *ForeignMonthChartService) GetMonthChart(stockCode string, period models.MonthChartPeriod, options models.MonthChartOptions) ([]models.ForeignMonthChartData, error) {
	s.logger.Info("Getting foreign stock month chart", 
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
	respBody, err := s.client.MakeRequestWithHeaders("POST", models.PathForeignStockMonthChart, nil, request, nil)
	if err != nil {
		s.logger.Error("Failed to call month chart API", err, 
			logger.Field{Key: "stock_code", Value: stockCode})
		return nil, errors.NewNetworkError("failed to call month chart API", err)
	}

	// 응답 파싱
	var response models.ForeignMonthChartResponse
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

	s.logger.Info("Successfully retrieved month chart data", 
		logger.Field{Key: "stock_code", Value: stockCode},
		logger.Field{Key: "data_count", Value: len(chartData)})

	return chartData, nil
}

// GetMonthChartWithMonths 월 수를 지정하여 월차트 조회 (편의 메서드)
func (s *ForeignMonthChartService) GetMonthChartWithMonths(stockCode, market string, months int, useAdjusted bool) ([]models.ForeignMonthChartData, error) {
	// 월 단위로 날짜 계산
	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, -months, 0).Format("2006-01-02")

	period := models.MonthChartPeriod{
		StartDate: startDate,
		EndDate:   endDate,
	}

	options := models.MonthChartOptions{
		UseAdjusted: useAdjusted,
		Market:      market,
	}

	return s.GetMonthChart(stockCode, period, options)
}

// GetRecentMonthChart 최근 월차트 데이터 조회
func (s *ForeignMonthChartService) GetRecentMonthChart(stockCode, market string, months int) ([]models.ForeignMonthChartData, error) {
	return s.GetMonthChartWithMonths(stockCode, market, months, true) // 기본적으로 수정주가 사용
}

// GetNASDAQMonthChart 나스닥 종목 월차트 조회
func (s *ForeignMonthChartService) GetNASDAQMonthChart(stockCode string, months int) ([]models.ForeignMonthChartData, error) {
	return s.GetMonthChartWithMonths(stockCode, "NASDAQ", months, true)
}

// GetNYMonthChart 뉴욕 증권거래소 종목 월차트 조회
func (s *ForeignMonthChartService) GetNYMonthChart(stockCode string, months int) ([]models.ForeignMonthChartData, error) {
	return s.GetMonthChartWithMonths(stockCode, "NY", months, true)
}

// GetAMEXMonthChart 아멕스 종목 월차트 조회
func (s *ForeignMonthChartService) GetAMEXMonthChart(stockCode string, months int) ([]models.ForeignMonthChartData, error) {
	return s.GetMonthChartWithMonths(stockCode, "AMEX", months, true)
}

// Get12MonthChart 12개월(1년) 차트 조회
func (s *ForeignMonthChartService) Get12MonthChart(stockCode, market string) ([]models.ForeignMonthChartData, error) {
	return s.GetMonthChartWithMonths(stockCode, market, 12, true)
}

// Get24MonthChart 24개월(2년) 차트 조회
func (s *ForeignMonthChartService) Get24MonthChart(stockCode, market string) ([]models.ForeignMonthChartData, error) {
	return s.GetMonthChartWithMonths(stockCode, market, 24, true)
}

// Get36MonthChart 36개월(3년) 차트 조회
func (s *ForeignMonthChartService) Get36MonthChart(stockCode, market string) ([]models.ForeignMonthChartData, error) {
	return s.GetMonthChartWithMonths(stockCode, market, 36, true)
}

// Get60MonthChart 60개월(5년) 차트 조회
func (s *ForeignMonthChartService) Get60MonthChart(stockCode, market string) ([]models.ForeignMonthChartData, error) {
	return s.GetMonthChartWithMonths(stockCode, market, 60, true)
}

// GetTechGiantsMonthChart 기술주 대장주들의 월차트 조회
func (s *ForeignMonthChartService) GetTechGiantsMonthChart(months int) (map[string][]models.ForeignMonthChartData, error) {
	techStocks := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA", "NVDA", "META"}
	results := make(map[string][]models.ForeignMonthChartData)

	for _, stockCode := range techStocks {
		data, err := s.GetNASDAQMonthChart(stockCode, months)
		if err != nil {
			s.logger.Warn("Failed to get tech stock month chart", 
				logger.Field{Key: "stock_code", Value: stockCode},
				logger.Field{Key: "error", Value: err.Error()})
			continue
		}
		results[stockCode] = data
	}

	return results, nil
}

// GetVolatilityAnalysis 월간 변동성 분석
func (s *ForeignMonthChartService) GetVolatilityAnalysis(chartData []models.ForeignMonthChartData) map[string]float64 {
	if len(chartData) == 0 {
		return nil
	}

	var monthlyRanges, monthlyRangeRates, changeRates []float64
	var totalVolume int64

	for _, data := range chartData {
		monthlyRanges = append(monthlyRanges, data.MonthlyRange)
		monthlyRangeRates = append(monthlyRangeRates, data.MonthlyRangeRate)
		if data.ChangeRate != 0 { // 0이 아닌 값만 포함
			changeRates = append(changeRates, data.ChangeRate)
		}
		totalVolume += data.Volume
	}

	analysis := make(map[string]float64)
	
	// 평균 월간 변동폭
	analysis["avg_monthly_range"] = s.avgFloat(monthlyRanges)
	
	// 평균 월간 변동률
	analysis["avg_monthly_range_rate"] = s.avgFloat(monthlyRangeRates)
	
	// 평균 월간 변화율
	if len(changeRates) > 0 {
		analysis["avg_monthly_change_rate"] = s.avgFloat(changeRates)
	}
	
	// 최대 월간 변동률
	analysis["max_monthly_range_rate"] = s.maxFloat(monthlyRangeRates)
	
	// 평균 월간 거래량
	if len(chartData) > 0 {
		analysis["avg_monthly_volume"] = float64(totalVolume) / float64(len(chartData))
	}
	
	return analysis
}

// validateInputs 입력값 검증
func (s *ForeignMonthChartService) validateInputs(stockCode string, period models.MonthChartPeriod, options models.MonthChartOptions) error {
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
func (s *ForeignMonthChartService) buildRequest(stockCode string, period models.MonthChartPeriod, options models.MonthChartOptions) models.ForeignMonthChartRequest {
	input := models.ForeignMonthChartInput{
		InputOrgAdjPrc:       options.GetAdjustedCode(),
		InputCondMrktDivCode: options.GetMarketCode(),
		InputIscd1:           stockCode,
		InputDate1:           period.GetFormattedStartDate(),
		InputDate2:           period.GetFormattedEndDate(),
	}

	return models.ForeignMonthChartRequest{In: input}
}

// convertToChartData API 응답을 비즈니스 모델로 변환
func (s *ForeignMonthChartService) convertToChartData(stockCode string, outputs []models.ForeignMonthChartOutput, options models.MonthChartOptions) []models.ForeignMonthChartData {
	var chartData []models.ForeignMonthChartData

	for i, output := range outputs {
		// 월 종료일에서 연도와 월 계산
		monthEndDate := s.formatDate(output.Date)
		year, month := s.getYearMonth(output.Date)
		monthStartDate := s.calculateMonthStartDate(output.Date)
		
		// 거래량 처리
		volume := int64(0)
		if output.AcmlVol != "" {
			volume = utils.ParseInt(output.AcmlVol)
		}

		data := models.ForeignMonthChartData{
			StockCode:      stockCode,
			MonthEndDate:   monthEndDate,
			MonthStartDate: monthStartDate,
			Open:           utils.ParseFloat(output.Oprc),
			High:           utils.ParseFloat(output.Hprc),
			Low:            utils.ParseFloat(output.Lprc),
			Close:          utils.ParseFloat(output.Prpr),
			Volume:         volume,
			Market:         s.getMarketName(options.GetMarketCode()),
			MarketCode:     options.GetMarketCode(),
			IsAdjusted:     options.UseAdjusted,
			Year:           year,
			Month:          month,
		}

		// 월간 변동폭 계산
		data.MonthlyRange = data.High - data.Low
		if data.Low > 0 {
			data.MonthlyRangeRate = (data.MonthlyRange / data.Low) * 100
		}

		// 전월대비 계산 (이전 데이터가 있는 경우)
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
func (s *ForeignMonthChartService) formatDate(date string) string {
	if len(date) != 8 {
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", date[:4], date[4:6], date[6:8])
}

// getMarketName 시장 코드를 시장명으로 변환
func (s *ForeignMonthChartService) getMarketName(marketCode string) string {
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

// getYearMonth 날짜에서 연도와 월 추출
func (s *ForeignMonthChartService) getYearMonth(dateStr string) (int, int) {
	if len(dateStr) != 8 {
		return 0, 0
	}

	// YYYYMMDD 형식을 time.Time으로 변환
	t, err := time.Parse("20060102", dateStr)
	if err != nil {
		return 0, 0
	}

	return t.Year(), int(t.Month())
}

// calculateMonthStartDate 월 종료일에서 월 시작일 계산
func (s *ForeignMonthChartService) calculateMonthStartDate(monthEndDateStr string) string {
	if len(monthEndDateStr) != 8 {
		return ""
	}

	// YYYYMMDD 형식을 time.Time으로 변환
	endDate, err := time.Parse("20060102", monthEndDateStr)
	if err != nil {
		return ""
	}

	// 해당 월의 첫 번째 날
	startDate := time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())
	return startDate.Format("2006-01-02")
}

// GetLongTermTrend 장기 추세 분석 (12개월 기준)
func (s *ForeignMonthChartService) GetLongTermTrend(chartData []models.ForeignMonthChartData) string {
	if len(chartData) < 6 {
		return "Insufficient data"
	}

	// 최근 6개월의 종가 추출
	recentCloses := make([]float64, 0, 6)
	for i := 0; i < 6 && i < len(chartData); i++ {
		recentCloses = append(recentCloses, chartData[i].Close)
	}

	// 상승 또는 하락 추세 판단 (최신 데이터가 앞에 있으므로)
	downCount := 0
	for i := 0; i < len(recentCloses)-1; i++ {
		if recentCloses[i] < recentCloses[i+1] { // 최신이 이전보다 낮으면 하락
			downCount++
		}
	}

	if downCount >= 4 {
		return "Long-term Downtrend"
	} else if downCount <= 1 {
		return "Long-term Uptrend"
	}
	return "Long-term Sideways"
}

// Get12MonthHighLow 12개월 최고/최저가 계산
func (s *ForeignMonthChartService) Get12MonthHighLow(chartData []models.ForeignMonthChartData) (float64, float64) {
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

// GetSeasonalAnalysis 계절성 분석 (월별 평균 수익률)
func (s *ForeignMonthChartService) GetSeasonalAnalysis(chartData []models.ForeignMonthChartData) map[int]float64 {
	monthlyReturns := make(map[int][]float64)
	
	// 월별 수익률 분류
	for _, data := range chartData {
		if data.ChangeRate != 0 {
			monthlyReturns[data.Month] = append(monthlyReturns[data.Month], data.ChangeRate)
		}
	}
	
	// 월별 평균 수익률 계산
	seasonalData := make(map[int]float64)
	for month, returns := range monthlyReturns {
		if len(returns) > 0 {
			seasonalData[month] = s.avgFloat(returns)
		}
	}
	
	return seasonalData
}

// 유틸리티 함수들
func (s *ForeignMonthChartService) maxFloat(values []float64) float64 {
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

func (s *ForeignMonthChartService) avgFloat(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}