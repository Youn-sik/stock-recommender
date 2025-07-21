package foreign

import (
	"testing"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
	"stock-recommender/backend/openapi/utils"
)

func TestForeignDayChartService_GetDayChart(t *testing.T) {
	// 모의 데이터 생성
	mockData := []models.ForeignDayChartOutput{
		{
			Hour:    "",
			Date:    "20250711",
			Prpr:    "313.5100",
			Oprc:    "307.8900",
			Hprc:    "314.0900",
			Lprc:    "305.6500",
			AcmlVol: "79236442",
		},
		{
			Hour:    "",
			Date:    "20250710",
			Prpr:    "309.8700",
			Oprc:    "300.0500",
			Hprc:    "310.4800",
			Lprc:    "300.0000",
			AcmlVol: "104365271",
		},
	}

	// 모의 서버 생성
	handler := utils.CreateForeignDayChartMockHandler(t, models.PathForeignStockDayChart, "TSLA", mockData)
	mockServer := utils.NewMockServer(t, handler)
	defer mockServer.Close()

	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignDayChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	// 테스트 실행
	t.Run("GetDayChart", func(t *testing.T) {
		period := models.DayChartPeriod{
			StartDate: "20250701",
			EndDate:   "20250714",
		}

		options := models.DayChartOptions{
			UseAdjusted: true,
			Market:      "NASDAQ",
		}

		data, err := service.GetDayChart("TSLA", period, options)
		if err != nil {
			t.Fatalf("Failed to get day chart: %v", err)
		}

		// 데이터 검증
		if len(data) == 0 {
			t.Error("Expected day chart data, but got empty result")
		}

		// 첫 번째 데이터 검증
		if len(data) > 0 {
			firstData := data[0]
			utils.AssertStringEqual(t, "TSLA", firstData.StockCode, "Stock code")
			utils.AssertFloatEqual(t, 313.51, firstData.Close, "Close price")
			utils.AssertIntEqual(t, 79236442, firstData.Volume, "Volume")
			utils.AssertStringEqual(t, "2025-07-11", firstData.Date, "Date format")
		}
	})

	t.Run("GetNASDAQDayChart", func(t *testing.T) {
		data, err := service.GetNASDAQDayChart("AAPL", 30)
		if err != nil {
			t.Fatalf("Failed to get NASDAQ day chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected day chart data for NASDAQ")
		}
	})

	t.Run("GetNYDayChart", func(t *testing.T) {
		data, err := service.GetNYDayChart("IBM", 30)
		if err != nil {
			t.Fatalf("Failed to get NY day chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected day chart data for NY")
		}
	})

	t.Run("GetRecentDayChart", func(t *testing.T) {
		data, err := service.GetRecentDayChart("MSFT", "NASDAQ", 7)
		if err != nil {
			t.Fatalf("Failed to get recent day chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected recent day chart data")
		}
	})
}

func TestForeignDayChartService_PeriodMethods(t *testing.T) {
	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignDayChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	t.Run("GetYearChart", func(t *testing.T) {
		data, err := service.GetYearChart("AAPL", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get year chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected year chart data")
		}
	})

	t.Run("GetMonthChart", func(t *testing.T) {
		data, err := service.GetMonthChart("TSLA", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get month chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected month chart data")
		}
	})

	t.Run("GetWeekChart", func(t *testing.T) {
		data, err := service.GetWeekChart("NVDA", "NASDAQ")
		if err != nil {
			t.Fatalf("Failed to get week chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected week chart data")
		}
	})
}

func TestForeignDayChartService_GetTechGiantsDayChart(t *testing.T) {
	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignDayChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	results, err := service.GetTechGiantsDayChart(7)
	if err != nil {
		t.Fatalf("Failed to get tech giants day chart: %v", err)
	}

	// 최소 하나 이상의 결과가 있어야 함
	if len(results) == 0 {
		t.Error("Expected at least one tech stock chart result")
	}

	// 각 종목 데이터 확인
	for stockCode, data := range results {
		if len(data) == 0 {
			t.Errorf("Expected day chart data for stock %s", stockCode)
		}

		// 첫 번째 데이터의 종목 코드 확인
		if len(data) > 0 && data[0].StockCode != stockCode {
			t.Errorf("Expected stock code %s, got %s", stockCode, data[0].StockCode)
		}
	}
}

func TestForeignDayChartService_DataConversion(t *testing.T) {
	service := &ForeignDayChartService{}

	// 테스트 데이터
	outputs := []models.ForeignDayChartOutput{
		{
			Hour:    "",
			Date:    "20250711",
			Prpr:    "313.5100",
			Oprc:    "307.8900",
			Hprc:    "314.0900",
			Lprc:    "305.6500",
			AcmlVol: "79236442",
		},
		{
			Hour:    "",
			Date:    "20250710",
			Prpr:    "309.8700",
			Oprc:    "300.0500",
			Hprc:    "310.4800",
			Lprc:    "300.0000",
			AcmlVol: "104365271",
		},
	}

	options := models.DayChartOptions{
		UseAdjusted: true,
		Market:      "NASDAQ",
	}

	// 변환 테스트
	data := service.convertToChartData("TSLA", outputs, options)

	// 검증
	if len(data) != 2 {
		t.Errorf("Expected 2 day chart data, got %d", len(data))
	}

	if len(data) > 0 {
		chartData := data[0]
		if chartData.StockCode != "TSLA" {
			t.Errorf("Expected stock code TSLA, got %s", chartData.StockCode)
		}
		if chartData.Open != 307.89 {
			t.Errorf("Expected open price 307.89, got %.2f", chartData.Open)
		}
		if chartData.High != 314.09 {
			t.Errorf("Expected high price 314.09, got %.2f", chartData.High)
		}
		if chartData.Low != 305.65 {
			t.Errorf("Expected low price 305.65, got %.2f", chartData.Low)
		}
		if chartData.Close != 313.51 {
			t.Errorf("Expected close price 313.51, got %.2f", chartData.Close)
		}
		if chartData.Volume != 79236442 {
			t.Errorf("Expected volume 79236442, got %d", chartData.Volume)
		}
		if chartData.Market != "NASDAQ" {
			t.Errorf("Expected market NASDAQ, got %s", chartData.Market)
		}
		if chartData.Date != "2025-07-11" {
			t.Errorf("Expected date 2025-07-11, got %s", chartData.Date)
		}
		if chartData.WeekDay == "" {
			t.Error("Expected week day to be set")
		}
	}
}

func TestForeignDayChartService_UtilityFunctions(t *testing.T) {
	service := &ForeignDayChartService{}

	t.Run("formatDate", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"20250711", "2025-07-11"},
			{"20231225", "2023-12-25"},
			{"", ""},
			{"2025071", ""},
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

	t.Run("getWeekDay", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"20250711", "금"}, // 2025년 7월 11일은 금요일
			{"20250710", "목"}, // 2025년 7월 10일은 목요일
			{"", ""},
			{"2025071", ""},
		}

		for _, test := range tests {
			result := service.getWeekDay(test.input)
			if result != test.expected {
				t.Errorf("getWeekDay(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	})
}

func TestForeignDayChartService_GetPriceStatistics(t *testing.T) {
	service := &ForeignDayChartService{}

	testData := []models.ForeignDayChartData{
		{High: 100, Low: 90, Close: 95, Volume: 1000},
		{High: 110, Low: 95, Close: 105, Volume: 1500},
		{High: 105, Low: 100, Close: 102, Volume: 1200},
	}

	stats := service.GetPriceStatistics(testData)

	if stats == nil {
		t.Fatal("Expected statistics, but got nil")
	}

	// 최고가 확인
	if stats["max_high"] != 110 {
		t.Errorf("Expected max_high 110, got %.2f", stats["max_high"])
	}

	// 최저가 확인
	if stats["min_low"] != 90 {
		t.Errorf("Expected min_low 90, got %.2f", stats["min_low"])
	}

	// 평균 종가 확인 (95+105+102)/3 = 100.67
	expectedAvg := (95.0 + 105.0 + 102.0) / 3.0
	if stats["avg_close"] != expectedAvg {
		t.Errorf("Expected avg_close %.2f, got %.2f", expectedAvg, stats["avg_close"])
	}

	// 평균 거래량 확인 (1000+1500+1200)/3 = 1233.33
	expectedAvgVol := (1000.0 + 1500.0 + 1200.0) / 3.0
	if stats["avg_volume"] != expectedAvgVol {
		t.Errorf("Expected avg_volume %.2f, got %.2f", expectedAvgVol, stats["avg_volume"])
	}
}

func TestDayChartPeriod_Methods(t *testing.T) {
	t.Run("FormatDate", func(t *testing.T) {
		period := models.DayChartPeriod{}

		tests := []struct {
			input    string
			expected string
		}{
			{"2025-07-11", "20250711"},
			{"20250711", "20250711"},
			{"", ""},
			{"2025-7-11", ""}, // 잘못된 형식
		}

		for _, test := range tests {
			result := period.FormatDate(test.input)
			if result != test.expected {
				t.Errorf("FormatDate(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	})

	t.Run("GetFormattedDates", func(t *testing.T) {
		period := models.DayChartPeriod{
			StartDate: "2025-07-01",
			EndDate:   "2025-07-11",
		}

		startFormatted := period.GetFormattedStartDate()
		endFormatted := period.GetFormattedEndDate()

		if startFormatted != "20250701" {
			t.Errorf("Expected formatted start date 20250701, got %s", startFormatted)
		}

		if endFormatted != "20250711" {
			t.Errorf("Expected formatted end date 20250711, got %s", endFormatted)
		}
	})
}

func TestDayChartOptions_Methods(t *testing.T) {
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
			options := models.DayChartOptions{Market: test.market}
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
			options := models.DayChartOptions{UseAdjusted: test.useAdjusted}
			result := options.GetAdjustedCode()
			if result != test.expected {
				t.Errorf("GetAdjustedCode(%v) = %s, expected %s", test.useAdjusted, result, test.expected)
			}
		}
	})
}