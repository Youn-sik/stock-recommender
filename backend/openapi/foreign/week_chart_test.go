package foreign

import (
	"testing"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
	"stock-recommender/backend/openapi/utils"
)

func TestForeignWeekChartService_GetWeekChart(t *testing.T) {
	// 모의 데이터 생성
	mockData := []models.ForeignWeekChartOutput{
		{
			Hour:    "",
			Date:    "20240129",
			Prpr:    "187.9100",
			Oprc:    "185.6300",
			Hprc:    "196.3593",
			Lprc:    "182.0000",
			CntgVol: "",
		},
		{
			Hour:    "",
			Date:    "20240122",
			Prpr:    "183.2500",
			Oprc:    "212.2600",
			Hprc:    "217.8000",
			Lprc:    "180.0600",
			CntgVol: "",
		},
	}

	// 모의 서버 생성
	handler := utils.CreateForeignWeekChartMockHandler(t, models.PathForeignStockWeekChart, "TSLA", mockData)
	mockServer := utils.NewMockServer(t, handler)
	defer mockServer.Close()

	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignWeekChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	// 테스트 실행
	t.Run("GetWeekChart", func(t *testing.T) {
		period := models.WeekChartPeriod{
			StartDate: "20230101",
			EndDate:   "20240201",
		}

		options := models.WeekChartOptions{
			UseAdjusted: true,
			Market:      "NASDAQ",
		}

		data, err := service.GetWeekChart("TSLA", period, options)
		if err != nil {
			t.Fatalf("Failed to get week chart: %v", err)
		}

		// 데이터 검증
		if len(data) == 0 {
			t.Error("Expected week chart data, but got empty result")
		}

		// 첫 번째 데이터 검증
		if len(data) > 0 {
			firstData := data[0]
			utils.AssertStringEqual(t, "TSLA", firstData.StockCode, "Stock code")
			utils.AssertFloatEqual(t, 187.91, firstData.Close, "Close price")
			utils.AssertStringEqual(t, "2024-01-29", firstData.WeekEndDate, "Week end date format")
			
			// 주간 변동폭 계산 확인
			expectedRange := 196.3593 - 182.0000
			if firstData.WeeklyRange < expectedRange-0.01 || firstData.WeeklyRange > expectedRange+0.01 {
				t.Errorf("Expected weekly range %.4f, got %.4f", expectedRange, firstData.WeeklyRange)
			}
		}
	})

	t.Run("GetNASDAQWeekChart", func(t *testing.T) {
		data, err := service.GetNASDAQWeekChart("AAPL", 13) // 13주
		if err != nil {
			t.Fatalf("Failed to get NASDAQ week chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected week chart data for NASDAQ")
		}
	})

	t.Run("GetNYWeekChart", func(t *testing.T) {
		data, err := service.GetNYWeekChart("IBM", 13)
		if err != nil {
			t.Fatalf("Failed to get NY week chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected week chart data for NY")
		}
	})

	t.Run("GetRecentWeekChart", func(t *testing.T) {
		data, err := service.GetRecentWeekChart("MSFT", "NASDAQ", 4) // 최근 4주
		if err != nil {
			t.Fatalf("Failed to get recent week chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected recent week chart data")
		}
	})
}

func TestForeignWeekChartService_PeriodMethods(t *testing.T) {
	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignWeekChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	t.Run("Get52WeekChart", func(t *testing.T) {
		data, err := service.Get52WeekChart("AAPL", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get 52 week chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected 52 week chart data")
		}
	})

	t.Run("Get26WeekChart", func(t *testing.T) {
		data, err := service.Get26WeekChart("TSLA", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get 26 week chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected 26 week chart data")
		}
	})

	t.Run("Get13WeekChart", func(t *testing.T) {
		data, err := service.Get13WeekChart("NVDA", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get 13 week chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected 13 week chart data")
		}
	})
}

func TestForeignWeekChartService_GetTechGiantsWeekChart(t *testing.T) {
	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignWeekChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	results, err := service.GetTechGiantsWeekChart(13) // 13주
	if err != nil {
		t.Fatalf("Failed to get tech giants week chart: %v", err)
	}

	// 최소 하나 이상의 결과가 있어야 함
	if len(results) == 0 {
		t.Error("Expected at least one tech stock chart result")
	}

	// 각 종목 데이터 확인
	for stockCode, data := range results {
		if len(data) == 0 {
			t.Errorf("Expected week chart data for stock %s", stockCode)
		}

		// 첫 번째 데이터의 종목 코드 확인
		if len(data) > 0 && data[0].StockCode != stockCode {
			t.Errorf("Expected stock code %s, got %s", stockCode, data[0].StockCode)
		}
	}
}

func TestForeignWeekChartService_DataConversion(t *testing.T) {
	service := &ForeignWeekChartService{}

	// 테스트 데이터
	outputs := []models.ForeignWeekChartOutput{
		{
			Hour:    "",
			Date:    "20240129",
			Prpr:    "187.9100",
			Oprc:    "185.6300",
			Hprc:    "196.3593",
			Lprc:    "182.0000",
			CntgVol: "",
		},
		{
			Hour:    "",
			Date:    "20240122",
			Prpr:    "183.2500",
			Oprc:    "212.2600",
			Hprc:    "217.8000",
			Lprc:    "180.0600",
			CntgVol: "1000000",
		},
	}

	options := models.WeekChartOptions{
		UseAdjusted: true,
		Market:      "NASDAQ",
	}

	// 변환 테스트
	data := service.convertToChartData("TSLA", outputs, options)

	// 검증
	if len(data) != 2 {
		t.Errorf("Expected 2 week chart data, got %d", len(data))
	}

	if len(data) > 0 {
		chartData := data[0]
		if chartData.StockCode != "TSLA" {
			t.Errorf("Expected stock code TSLA, got %s", chartData.StockCode)
		}
		if chartData.Open != 185.63 {
			t.Errorf("Expected open price 185.63, got %.2f", chartData.Open)
		}
		if chartData.High != 196.3593 {
			t.Errorf("Expected high price 196.3593, got %.4f", chartData.High)
		}
		if chartData.Low != 182.00 {
			t.Errorf("Expected low price 182.00, got %.2f", chartData.Low)
		}
		if chartData.Close != 187.91 {
			t.Errorf("Expected close price 187.91, got %.2f", chartData.Close)
		}
		if chartData.Volume != 0 {
			t.Errorf("Expected volume 0 for empty string, got %d", chartData.Volume)
		}
		if chartData.Market != "NASDAQ" {
			t.Errorf("Expected market NASDAQ, got %s", chartData.Market)
		}
		if chartData.WeekEndDate != "2024-01-29" {
			t.Errorf("Expected week end date 2024-01-29, got %s", chartData.WeekEndDate)
		}
		if chartData.WeekStartDate == "" {
			t.Error("Expected week start date to be calculated")
		}
		if chartData.Year == 0 || chartData.WeekNumber == 0 {
			t.Error("Expected year and week number to be set")
		}
		
		// 주간 변동폭 검증
		expectedRange := 196.3593 - 182.0000
		if chartData.WeeklyRange < expectedRange-0.01 || chartData.WeeklyRange > expectedRange+0.01 {
			t.Errorf("Expected weekly range %.4f, got %.4f", expectedRange, chartData.WeeklyRange)
		}
	}
	
	// 두 번째 데이터의 거래량 확인
	if len(data) > 1 && data[1].Volume != 1000000 {
		t.Errorf("Expected volume 1000000, got %d", data[1].Volume)
	}
}

func TestForeignWeekChartService_UtilityFunctions(t *testing.T) {
	service := &ForeignWeekChartService{}

	t.Run("formatDate", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"20240129", "2024-01-29"},
			{"20231225", "2023-12-25"},
			{"", ""},
			{"2024012", ""},
		}

		for _, test := range tests {
			result := service.formatDate(test.input)
			if result != test.expected {
				t.Errorf("formatDate(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	})

	t.Run("getMarketName", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{models.ForeignMarketNY, "New York Stock Exchange"},
			{models.ForeignMarketNASDAQ, "NASDAQ"},
			{models.ForeignMarketAMEX, "American Stock Exchange"},
			{"UNKNOWN", "Unknown"},
		}

		for _, test := range tests {
			result := service.getMarketName(test.input)
			if result != test.expected {
				t.Errorf("getMarketName(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	})

	t.Run("getYearWeek", func(t *testing.T) {
		tests := []struct {
			input        string
			expectedYear int
			expectedWeek int
		}{
			{"20240129", 2024, 5},  // 2024년 1월 29일은 5주차
			{"20231225", 2023, 52}, // 2023년 12월 25일은 52주차
			{"", 0, 0},
			{"2024012", 0, 0},
		}

		for _, test := range tests {
			year, week := service.getYearWeek(test.input)
			if year != test.expectedYear || week != test.expectedWeek {
				t.Errorf("getYearWeek(%s) = (%d, %d), expected (%d, %d)", 
					test.input, year, week, test.expectedYear, test.expectedWeek)
			}
		}
	})

	t.Run("calculateWeekStartDate", func(t *testing.T) {
		// 2024년 1월 29일은 월요일
		result := service.calculateWeekStartDate("20240129")
		if result != "2024-01-28" { // 일요일
			t.Errorf("Expected week start date 2024-01-28 for 20240129, got %s", result)
		}
	})
}

func TestForeignWeekChartService_AnalysisFunctions(t *testing.T) {
	service := &ForeignWeekChartService{}

	testData := []models.ForeignWeekChartData{
		{
			High: 200, Low: 180, Close: 195, Volume: 1000,
			WeeklyRange: 20, WeeklyRangeRate: 11.11, ChangeRate: 5.0,
		},
		{
			High: 190, Low: 175, Close: 185, Volume: 1200,
			WeeklyRange: 15, WeeklyRangeRate: 8.57, ChangeRate: -2.5,
		},
		{
			High: 185, Low: 170, Close: 180, Volume: 800,
			WeeklyRange: 15, WeeklyRangeRate: 8.82, ChangeRate: -3.0,
		},
		{
			High: 180, Low: 165, Close: 175, Volume: 1500,
			WeeklyRange: 15, WeeklyRangeRate: 9.09, ChangeRate: -2.0,
		},
	}

	t.Run("GetVolatilityAnalysis", func(t *testing.T) {
		analysis := service.GetVolatilityAnalysis(testData)

		if analysis == nil {
			t.Fatal("Expected volatility analysis, but got nil")
		}

		// 평균 주간 변동폭 확인
		expectedAvgRange := (20.0 + 15.0 + 15.0 + 15.0) / 4.0
		if analysis["avg_weekly_range"] != expectedAvgRange {
			t.Errorf("Expected avg_weekly_range %.2f, got %.2f", expectedAvgRange, analysis["avg_weekly_range"])
		}

		// 평균 주간 변동률 확인
		expectedAvgRangeRate := (11.11 + 8.57 + 8.82 + 9.09) / 4.0
		if analysis["avg_weekly_range_rate"] < expectedAvgRangeRate-0.01 || 
		   analysis["avg_weekly_range_rate"] > expectedAvgRangeRate+0.01 {
			t.Errorf("Expected avg_weekly_range_rate %.2f, got %.2f", expectedAvgRangeRate, analysis["avg_weekly_range_rate"])
		}
	})

	t.Run("Get52WeekHighLow", func(t *testing.T) {
		high, low := service.Get52WeekHighLow(testData)

		if high != 200 {
			t.Errorf("Expected 52-week high 200, got %.2f", high)
		}

		if low != 165 {
			t.Errorf("Expected 52-week low 165, got %.2f", low)
		}
	})

	t.Run("GetTrendAnalysis", func(t *testing.T) {
		// 상승 추세 데이터 (최신 195가 가장 높으므로 상승 추세)
		trend := service.GetTrendAnalysis(testData)
		if trend != "Uptrend" {
			t.Errorf("Expected Uptrend, got %s", trend)
		}

		// 하락 추세 데이터 (최신이 앞에 있고 점점 감소)
		downData := []models.ForeignWeekChartData{
			{Close: 175}, // 최신 (가장 낮음)
			{Close: 180},
			{Close: 185}, 
			{Close: 200}, // 가장 오래된 데이터 (가장 높음)
		}
		trend = service.GetTrendAnalysis(downData)
		if trend != "Downtrend" {
			t.Errorf("Expected Downtrend, got %s", trend)
		}
	})
}

func TestWeekChartPeriod_Methods(t *testing.T) {
	t.Run("FormatDate", func(t *testing.T) {
		period := models.WeekChartPeriod{}

		tests := []struct {
			input    string
			expected string
		}{
			{"2024-01-29", "20240129"},
			{"20240129", "20240129"},
			{"", ""},
			{"2024-1-29", ""}, // 잘못된 형식
		}

		for _, test := range tests {
			result := period.FormatDate(test.input)
			if result != test.expected {
				t.Errorf("FormatDate(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	})

	t.Run("GetFormattedDates", func(t *testing.T) {
		period := models.WeekChartPeriod{
			StartDate: "2023-01-01",
			EndDate:   "2024-02-01",
		}

		startFormatted := period.GetFormattedStartDate()
		endFormatted := period.GetFormattedEndDate()

		if startFormatted != "20230101" {
			t.Errorf("Expected formatted start date 20230101, got %s", startFormatted)
		}

		if endFormatted != "20240201" {
			t.Errorf("Expected formatted end date 20240201, got %s", endFormatted)
		}
	})
}

func TestWeekChartOptions_Methods(t *testing.T) {
	t.Run("GetMarketCode", func(t *testing.T) {
		tests := []struct {
			market   string
			expected string
		}{
			{"NY", models.ForeignMarketNY},
			{"NYSE", models.ForeignMarketNY},
			{"NASDAQ", models.ForeignMarketNASDAQ},
			{"AMEX", models.ForeignMarketAMEX},
			{"unknown", models.ForeignMarketNASDAQ}, // 기본값
		}

		for _, test := range tests {
			options := models.WeekChartOptions{Market: test.market}
			result := options.GetMarketCode()
			if result != test.expected {
				t.Errorf("GetMarketCode(%s) = %s, expected %s", test.market, result, test.expected)
			}
		}
	})

	t.Run("GetAdjustedCode", func(t *testing.T) {
		tests := []struct {
			useAdjusted bool
			expected    string
		}{
			{true, models.AdjustedPriceEnabled},
			{false, models.AdjustedPriceDisabled},
		}

		for _, test := range tests {
			options := models.WeekChartOptions{UseAdjusted: test.useAdjusted}
			result := options.GetAdjustedCode()
			if result != test.expected {
				t.Errorf("GetAdjustedCode(%v) = %s, expected %s", test.useAdjusted, result, test.expected)
			}
		}
	})
}