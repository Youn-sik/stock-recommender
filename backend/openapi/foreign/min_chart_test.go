package foreign

import (
	"testing"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
	"stock-recommender/backend/openapi/utils"
)

func TestForeignMinChartService_GetMinChart(t *testing.T) {
	// 모의 데이터 생성
	mockData := []models.ForeignMinChartOutput{
		{
			Hour:    "163000",
			Date:    "20240205",
			Prpr:    "187.5700",
			Oprc:    "187.8300",
			Hprc:    "187.8500",
			Lprc:    "187.5150",
			CntgVol: "14162",
		},
		{
			Hour:    "162000",
			Date:    "20240205",
			Prpr:    "187.8300",
			Oprc:    "187.9700",
			Hprc:    "188.0300",
			Lprc:    "187.6800",
			CntgVol: "26443",
		},
	}

	// 모의 서버 생성
	handler := utils.CreateForeignMinChartMockHandler(t, models.PathForeignStockMinChart, "AAPL", mockData)
	mockServer := utils.NewMockServer(t, handler)
	defer mockServer.Close()

	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignMinChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	// 테스트 실행
	t.Run("GetMinChart", func(t *testing.T) {
		period := models.ChartPeriod{
			StartDate: "20240205",
			EndDate:   "20240205",
			IsRange:   true,
		}

		options := models.ChartOptions{
			Interval:    "1min",
			UseAdjusted: true,
			Market:      "NASDAQ",
			DataCount:   100,
		}

		data, err := service.GetMinChart("AAPL", period, options)
		if err != nil {
			t.Fatalf("Failed to get min chart: %v", err)
		}

		// 데이터 검증
		if len(data) == 0 {
			t.Error("Expected chart data, but got empty result")
		}

		// 첫 번째 데이터 검증
		if len(data) > 0 {
			firstData := data[0]
			utils.AssertStringEqual(t, "AAPL", firstData.StockCode, "Stock code")
			utils.AssertFloatEqual(t, 187.57, firstData.Close, "Close price")
			utils.AssertIntEqual(t, 14162, firstData.Volume, "Volume")
		}
	})

	t.Run("GetNASDAQMinChart", func(t *testing.T) {
		data, err := service.GetNASDAQMinChart("AAPL", "1min", 1)
		if err != nil {
			t.Fatalf("Failed to get NASDAQ min chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected chart data for NASDAQ")
		}
	})

	t.Run("GetNYMinChart", func(t *testing.T) {
		data, err := service.GetNYMinChart("IBM", "5min", 1)
		if err != nil {
			t.Fatalf("Failed to get NY min chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected chart data for NY")
		}
	})

	t.Run("GetLatestMinChart", func(t *testing.T) {
		data, err := service.GetLatestMinChart("MSFT", "NASDAQ", "1min", 50)
		if err != nil {
			t.Fatalf("Failed to get latest min chart: %v", err)
		}

		if len(data) == 0 {
			t.Error("Expected latest chart data")
		}
	})
}

func TestForeignMinChartService_GetPopularStocksMinChart(t *testing.T) {
	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignMinChartService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	results, err := service.GetPopularStocksMinChart("1min", 1)
	if err != nil {
		t.Fatalf("Failed to get popular stocks min chart: %v", err)
	}

	// 최소 하나 이상의 결과가 있어야 함
	if len(results) == 0 {
		t.Error("Expected at least one popular stock chart result")
	}

	// 각 종목 데이터 확인
	for stockCode, data := range results {
		if len(data) == 0 {
			t.Errorf("Expected chart data for stock %s", stockCode)
		}

		// 첫 번째 데이터의 종목 코드 확인
		if len(data) > 0 && data[0].StockCode != stockCode {
			t.Errorf("Expected stock code %s, got %s", stockCode, data[0].StockCode)
		}
	}
}

func TestForeignMinChartService_DataConversion(t *testing.T) {
	service := &ForeignMinChartService{}

	// 테스트 데이터
	outputs := []models.ForeignMinChartOutput{
		{
			Hour:    "163000",
			Date:    "20240205",
			Prpr:    "187.5700",
			Oprc:    "187.8300",
			Hprc:    "187.8500",
			Lprc:    "187.5150",
			CntgVol: "14162",
		},
	}

	options := models.ChartOptions{
		Interval:    "1min",
		UseAdjusted: true,
		Market:      "NASDAQ",
	}

	// 변환 테스트
	data := service.convertToChartData("AAPL", outputs, options)

	// 검증
	if len(data) != 1 {
		t.Errorf("Expected 1 chart data, got %d", len(data))
	}

	if len(data) > 0 {
		chartData := data[0]
		if chartData.StockCode != "AAPL" {
			t.Errorf("Expected stock code AAPL, got %s", chartData.StockCode)
		}
		if chartData.Open != 187.83 {
			t.Errorf("Expected open price 187.83, got %.2f", chartData.Open)
		}
		if chartData.High != 187.85 {
			t.Errorf("Expected high price 187.85, got %.2f", chartData.High)
		}
		if chartData.Low != 187.515 {
			t.Errorf("Expected low price 187.515, got %.3f", chartData.Low)
		}
		if chartData.Close != 187.57 {
			t.Errorf("Expected close price 187.57, got %.2f", chartData.Close)
		}
		if chartData.Volume != 14162 {
			t.Errorf("Expected volume 14162, got %d", chartData.Volume)
		}
		if chartData.Market != "NASDAQ" {
			t.Errorf("Expected market NASDAQ, got %s", chartData.Market)
		}
		if chartData.DateTime != "2024-02-05 16:30:00" {
			t.Errorf("Expected datetime 2024-02-05 16:30:00, got %s", chartData.DateTime)
		}
	}
}

func TestForeignMinChartService_UtilityFunctions(t *testing.T) {
	service := &ForeignMinChartService{}

	t.Run("formatDateTime", func(t *testing.T) {
		tests := []struct {
			date     string
			hour     string
			expected string
		}{
			{"20240205", "163000", "2024-02-05 16:30:00"},
			{"20231225", "093000", "2023-12-25 09:30:00"},
			{"", "163000", ""},
			{"20240205", "", ""},
		}

		for _, test := range tests {
			result := service.formatDateTime(test.date, test.hour)
			if result != test.expected {
				t.Errorf("formatDateTime(%s, %s) = %s, expected %s", test.date, test.hour, result, test.expected)
			}
		}
	})

	t.Run("formatDate", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"20240205", "2024-02-05"},
			{"20231225", "2023-12-25"},
			{"", ""},
			{"2024020", ""},
		}

		for _, test := range tests {
			result := service.formatDate(test.input)
			if result != test.expected {
				t.Errorf("formatDate(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	})

	t.Run("formatTime", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"163000", "16:30:00"},
			{"093000", "09:30:00"},
			{"", ""},
			{"16300", ""},
		}

		for _, test := range tests {
			result := service.formatTime(test.input)
			if result != test.expected {
				t.Errorf("formatTime(%s) = %s, expected %s", test.input, result, test.expected)
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

	t.Run("GetIntervalDescription", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{models.ChartInterval30Sec, "30초"},
			{models.ChartInterval1Min, "1분"},
			{models.ChartInterval5Min, "5분"},
			{models.ChartInterval10Min, "10분"},
			{models.ChartInterval60Min, "60분"},
			{"UNKNOWN", "알 수 없음"},
		}

		for _, test := range tests {
			result := service.GetIntervalDescription(test.input)
			if result != test.expected {
				t.Errorf("GetIntervalDescription(%s) = %s, expected %s", test.input, result, test.expected)
			}
		}
	})
}

func TestChartOptions_Methods(t *testing.T) {
	t.Run("GetIntervalCode", func(t *testing.T) {
		tests := []struct {
			interval string
			expected string
		}{
			{"30sec", models.ChartInterval30Sec},
			{"1min", models.ChartInterval1Min},
			{"5min", models.ChartInterval5Min},
			{"10min", models.ChartInterval10Min},
			{"60min", models.ChartInterval60Min},
			{"unknown", models.ChartInterval1Min}, // 기본값
		}

		for _, test := range tests {
			options := models.ChartOptions{Interval: test.interval}
			result := options.GetIntervalCode()
			if result != test.expected {
				t.Errorf("GetIntervalCode(%s) = %s, expected %s", test.interval, result, test.expected)
			}
		}
	})

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
			options := models.ChartOptions{Market: test.market}
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
			options := models.ChartOptions{UseAdjusted: test.useAdjusted}
			result := options.GetAdjustedCode()
			if result != test.expected {
				t.Errorf("GetAdjustedCode(%v) = %s, expected %s", test.useAdjusted, result, test.expected)
			}
		}
	})

	t.Run("GetDataCountString", func(t *testing.T) {
		tests := []struct {
			dataCount int
			expected  string
		}{
			{100, "100"},
			{500, "500"},
			{0, ""},     // 기본값
			{-1, ""},    // 기본값
			{3000, ""},  // 기본값 (범위 초과)
		}

		for _, test := range tests {
			options := models.ChartOptions{DataCount: test.dataCount}
			result := options.GetDataCountString()
			if result != test.expected {
				t.Errorf("GetDataCountString(%d) = %s, expected %s", test.dataCount, result, test.expected)
			}
		}
	})
}