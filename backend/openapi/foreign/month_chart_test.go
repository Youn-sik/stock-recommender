package foreign

import (
	"testing"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
	"stock-recommender/backend/openapi/utils"
)

func TestForeignMonthChartService_GetMonthChart(t *testing.T) {
	// 모의 데이터 생성
	mockData := []models.ForeignMonthChartOutput{
		{
			Hour:    "",
			Date:    "20240131",
			Prpr:    "187.9100",
			Oprc:    "185.6300",
			Hprc:    "196.3593",
			Lprc:    "182.0000",
			AcmlVol: "2394275082",
		},
		{
			Hour:    "",
			Date:    "20231230",
			Prpr:    "248.4800",
			Oprc:    "252.7400",
			Hprc:    "271.0000",
			Lprc:    "235.0000",
			AcmlVol: "3443091887",
		},
	}

	// 모의 서버 생성
	handler := utils.CreateForeignMonthChartMockHandler(t, models.PathForeignStockMonthChart, "TSLA", mockData)
	mockServer := utils.NewMockServer(t, handler)
	defer mockServer.Close()

	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignMonthChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	// 테스트 실행
	t.Run("GetMonthChart", func(t *testing.T) {
		period := models.MonthChartPeriod{
			StartDate: "20230101",
			EndDate:   "20240201",
		}

		options := models.MonthChartOptions{
			UseAdjusted: true,
			Market:      "NASDAQ",
		}

		data, err := service.GetMonthChart("TSLA", period, options)
		if err != nil {
			t.Fatalf("Failed to get month chart: %v", err)
		}

		// 데이터 검증
		if len(data) == 0 {
			t.Error("Expected month chart data, but got empty result")
		}

		// 첫 번째 데이터 검증
		if len(data) > 0 {
			firstData := data[0]
			utils.AssertStringEqual(t, "TSLA", firstData.StockCode, "Stock code")
			utils.AssertFloatEqual(t, 187.91, firstData.Close, "Close price")
			utils.AssertStringEqual(t, "2024-01-31", firstData.MonthEndDate, "Month end date format")
			
			// 월간 변동폭 계산 확인
			expectedRange := 196.3593 - 182.0000
			if firstData.MonthlyRange < expectedRange-0.01 || firstData.MonthlyRange > expectedRange+0.01 {
				t.Errorf("Expected monthly range %.4f, got %.4f", expectedRange, firstData.MonthlyRange)
			}
			
			// 연도와 월 확인
			if firstData.Year != 2024 || firstData.Month != 1 {
				t.Errorf("Expected year 2024, month 1, got year %d, month %d", firstData.Year, firstData.Month)
			}
		}
	})

	t.Run("GetNASDAQMonthChart", func(t *testing.T) {
		data, err := service.GetNASDAQMonthChart("AAPL", 12) // 12개월
		if err != nil {
			t.Fatalf("Failed to get NASDAQ month chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected month chart data for NASDAQ")
		}
	})

	t.Run("GetNYMonthChart", func(t *testing.T) {
		data, err := service.GetNYMonthChart("IBM", 12)
		if err != nil {
			t.Fatalf("Failed to get NY month chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected month chart data for NY")
		}
	})

	t.Run("GetRecentMonthChart", func(t *testing.T) {
		data, err := service.GetRecentMonthChart("MSFT", "NASDAQ", 6) // 최근 6개월
		if err != nil {
			t.Fatalf("Failed to get recent month chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected recent month chart data")
		}
	})
}

func TestForeignMonthChartService_PeriodMethods(t *testing.T) {
	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignMonthChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	t.Run("Get12MonthChart", func(t *testing.T) {
		data, err := service.Get12MonthChart("AAPL", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get 12 month chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected 12 month chart data")
		}
	})

	t.Run("Get24MonthChart", func(t *testing.T) {
		data, err := service.Get24MonthChart("TSLA", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get 24 month chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected 24 month chart data")
		}
	})

	t.Run("Get36MonthChart", func(t *testing.T) {
		data, err := service.Get36MonthChart("NVDA", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get 36 month chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected 36 month chart data")
		}
	})

	t.Run("Get60MonthChart", func(t *testing.T) {
		data, err := service.Get60MonthChart("GOOGL", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get 60 month chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected 60 month chart data")
		}
	})
}

func TestForeignMonthChartService_GetTechGiantsMonthChart(t *testing.T) {
	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignMonthChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	results, err := service.GetTechGiantsMonthChart(12) // 12개월
	if err != nil {
		t.Fatalf("Failed to get tech giants month chart: %v", err)
	}

	// 최소 하나 이상의 결과가 있어야 함
	if len(results) == 0 {
		t.Error("Expected at least one tech stock chart result")
	}

	// 각 종목 데이터 확인
	for stockCode, data := range results {
		if len(data) == 0 {
			t.Errorf("Expected month chart data for stock %s", stockCode)
		}

		// 첫 번째 데이터의 종목 코드 확인
		if len(data) > 0 && data[0].StockCode != stockCode {
			t.Errorf("Expected stock code %s, got %s", stockCode, data[0].StockCode)
		}
	}
}

func TestForeignMonthChartService_DataConversion(t *testing.T) {
	service := &ForeignMonthChartService{}

	// 테스트 데이터
	outputs := []models.ForeignMonthChartOutput{
		{
			Hour:    "",
			Date:    "20240131",
			Prpr:    "187.9100",
			Oprc:    "185.6300",
			Hprc:    "196.3593",
			Lprc:    "182.0000",
			AcmlVol: "",
		},
		{
			Hour:    "",
			Date:    "20231231",
			Prpr:    "248.4800",
			Oprc:    "252.7400",
			Hprc:    "271.0000",
			Lprc:    "235.0000",
			AcmlVol: "1000000",
		},
	}

	options := models.MonthChartOptions{
		UseAdjusted: true,
		Market:      "NASDAQ",
	}

	// 변환 테스트
	data := service.convertToChartData("TSLA", outputs, options)

	// 검증
	if len(data) != 2 {
		t.Errorf("Expected 2 month chart data, got %d", len(data))
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
		if chartData.MonthEndDate != "2024-01-31" {
			t.Errorf("Expected month end date 2024-01-31, got %s", chartData.MonthEndDate)
		}
		if chartData.MonthStartDate != "2024-01-01" {
			t.Errorf("Expected month start date 2024-01-01, got %s", chartData.MonthStartDate)
		}
		if chartData.Year != 2024 || chartData.Month != 1 {
			t.Errorf("Expected year 2024, month 1, got year %d, month %d", chartData.Year, chartData.Month)
		}
		
		// 월간 변동폭 검증
		expectedRange := 196.3593 - 182.0000
		if chartData.MonthlyRange < expectedRange-0.01 || chartData.MonthlyRange > expectedRange+0.01 {
			t.Errorf("Expected monthly range %.4f, got %.4f", expectedRange, chartData.MonthlyRange)
		}
	}
	
	// 두 번째 데이터의 거래량 확인
	if len(data) > 1 && data[1].Volume != 1000000 {
		t.Errorf("Expected volume 1000000, got %d", data[1].Volume)
	}
}

func TestForeignMonthChartService_UtilityFunctions(t *testing.T) {
	service := &ForeignMonthChartService{}

	t.Run("formatDate", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"20240131", "2024-01-31"},
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

	t.Run("getYearMonth", func(t *testing.T) {
		tests := []struct {
			input         string
			expectedYear  int
			expectedMonth int
		}{
			{"20240131", 2024, 1},  // 2024년 1월
			{"20231225", 2023, 12}, // 2023년 12월
			{"", 0, 0},
			{"2024012", 0, 0},
		}

		for _, test := range tests {
			year, month := service.getYearMonth(test.input)
			if year != test.expectedYear || month != test.expectedMonth {
				t.Errorf("getYearMonth(%s) = (%d, %d), expected (%d, %d)", 
					test.input, year, month, test.expectedYear, test.expectedMonth)
			}
		}
	})

	t.Run("calculateMonthStartDate", func(t *testing.T) {
		// 2024년 1월 31일의 월 시작일은 2024년 1월 1일
		result := service.calculateMonthStartDate("20240131")
		if result != "2024-01-01" {
			t.Errorf("Expected month start date 2024-01-01 for 20240131, got %s", result)
		}
		
		// 2023년 12월 25일의 월 시작일은 2023년 12월 1일
		result = service.calculateMonthStartDate("20231225")
		if result != "2023-12-01" {
			t.Errorf("Expected month start date 2023-12-01 for 20231225, got %s", result)
		}
	})
}

func TestForeignMonthChartService_AnalysisFunctions(t *testing.T) {
	service := &ForeignMonthChartService{}

	testData := []models.ForeignMonthChartData{
		{
			High: 200, Low: 180, Close: 195, Volume: 1000,
			MonthlyRange: 20, MonthlyRangeRate: 11.11, ChangeRate: 5.0,
			Year: 2024, Month: 1,
		},
		{
			High: 190, Low: 175, Close: 185, Volume: 1200,
			MonthlyRange: 15, MonthlyRangeRate: 8.57, ChangeRate: -2.5,
			Year: 2023, Month: 12,
		},
		{
			High: 185, Low: 170, Close: 180, Volume: 800,
			MonthlyRange: 15, MonthlyRangeRate: 8.82, ChangeRate: -3.0,
			Year: 2023, Month: 11,
		},
		{
			High: 180, Low: 165, Close: 175, Volume: 1500,
			MonthlyRange: 15, MonthlyRangeRate: 9.09, ChangeRate: -2.0,
			Year: 2023, Month: 10,
		},
		{
			High: 175, Low: 160, Close: 170, Volume: 1300,
			MonthlyRange: 15, MonthlyRangeRate: 9.37, ChangeRate: -1.5,
			Year: 2023, Month: 9,
		},
		{
			High: 170, Low: 155, Close: 165, Volume: 1100,
			MonthlyRange: 15, MonthlyRangeRate: 9.67, ChangeRate: -1.0,
			Year: 2023, Month: 8,
		},
	}

	t.Run("GetVolatilityAnalysis", func(t *testing.T) {
		analysis := service.GetVolatilityAnalysis(testData)

		if analysis == nil {
			t.Fatal("Expected volatility analysis, but got nil")
		}

		// 평균 월간 변동폭 확인
		expectedAvgRange := (20.0 + 15.0 + 15.0 + 15.0 + 15.0 + 15.0) / 6.0
		if analysis["avg_monthly_range"] < expectedAvgRange-0.01 || analysis["avg_monthly_range"] > expectedAvgRange+0.01 {
			t.Errorf("Expected avg_monthly_range %.2f, got %.2f", expectedAvgRange, analysis["avg_monthly_range"])
		}

		// 평균 월간 변동률 확인
		expectedAvgRangeRate := (11.11 + 8.57 + 8.82 + 9.09 + 9.37 + 9.67) / 6.0
		if analysis["avg_monthly_range_rate"] < expectedAvgRangeRate-0.01 || 
		   analysis["avg_monthly_range_rate"] > expectedAvgRangeRate+0.01 {
			t.Errorf("Expected avg_monthly_range_rate %.2f, got %.2f", expectedAvgRangeRate, analysis["avg_monthly_range_rate"])
		}
	})

	t.Run("Get12MonthHighLow", func(t *testing.T) {
		high, low := service.Get12MonthHighLow(testData)

		if high != 200 {
			t.Errorf("Expected 12-month high 200, got %.2f", high)
		}

		if low != 155 {
			t.Errorf("Expected 12-month low 155, got %.2f", low)
		}
	})

	t.Run("GetLongTermTrend", func(t *testing.T) {
		// 상승 추세 데이터 (최신 195가 가장 높으므로 상승 추세)
		trend := service.GetLongTermTrend(testData)
		if trend != "Long-term Uptrend" {
			t.Errorf("Expected Long-term Uptrend, got %s", trend)
		}

		// 하락 추세 데이터 (최신이 앞에 있고 점점 감소)
		downData := []models.ForeignMonthChartData{
			{Close: 175}, // 최신 (가장 낮음)
			{Close: 180},
			{Close: 185},
			{Close: 190},
			{Close: 195},
			{Close: 200}, // 가장 오래된 데이터 (가장 높음)
		}
		trend = service.GetLongTermTrend(downData)
		if trend != "Long-term Downtrend" {
			t.Errorf("Expected Long-term Downtrend, got %s", trend)
		}
	})

	t.Run("GetSeasonalAnalysis", func(t *testing.T) {
		seasonalData := service.GetSeasonalAnalysis(testData)

		// 1월, 10월, 11월, 12월 데이터가 있어야 함
		if len(seasonalData) == 0 {
			t.Error("Expected seasonal analysis data")
		}

		// 1월 평균 수익률 확인 (5.0%)
		if seasonalData[1] != 5.0 {
			t.Errorf("Expected seasonal data for January 5.0, got %.2f", seasonalData[1])
		}
	})
}

func TestMonthChartPeriod_Methods(t *testing.T) {
	t.Run("FormatDate", func(t *testing.T) {
		period := models.MonthChartPeriod{}

		tests := []struct {
			input    string
			expected string
		}{
			{"2024-01-31", "20240131"},
			{"20240131", "20240131"},
			{"", ""},
			{"2024-1-31", ""}, // 잘못된 형식
		}

		for _, test := range tests {
			result := period.FormatDate(test.input)
			if result != test.expected {
				t.Errorf("FormatDate(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	})

	t.Run("GetFormattedDates", func(t *testing.T) {
		period := models.MonthChartPeriod{
			StartDate: "2023-01-01",
			EndDate:   "2024-01-31",
		}

		startFormatted := period.GetFormattedStartDate()
		endFormatted := period.GetFormattedEndDate()

		if startFormatted != "20230101" {
			t.Errorf("Expected formatted start date 20230101, got %s", startFormatted)
		}

		if endFormatted != "20240131" {
			t.Errorf("Expected formatted end date 20240131, got %s", endFormatted)
		}
	})
}

func TestMonthChartOptions_Methods(t *testing.T) {
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
			options := models.MonthChartOptions{Market: test.market}
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
			options := models.MonthChartOptions{UseAdjusted: test.useAdjusted}
			result := options.GetAdjustedCode()
			if result != test.expected {
				t.Errorf("GetAdjustedCode(%v) = %s, expected %s", test.useAdjusted, result, test.expected)
			}
		}
	})
}