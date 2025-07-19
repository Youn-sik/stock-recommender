package foreign

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"stock-recommender/backend/config"
	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
)

func TestForeignCurrentPriceService_GetForeignCurrentPrice(t *testing.T) {
	// 테스트 서버 생성
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 경로 확인
		if r.URL.Path != models.PathForeignStockCurrentPrice {
			t.Errorf("Expected path %s, got %s", models.PathForeignStockCurrentPrice, r.URL.Path)
		}

		// 메소드 확인
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 요청 본문 파싱
		var req models.ForeignCurrentPriceRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// 요청 데이터 검증
		if req.In.InputIscd1 != "TSLA" {
			t.Errorf("Expected stock code TSLA, got %s", req.In.InputIscd1)
		}
		if req.In.InputCondMrktDivCode != models.ForeignMarketNASDAQ {
			t.Errorf("Expected market div %s, got %s", models.ForeignMarketNASDAQ, req.In.InputCondMrktDivCode)
		}

		// 응답 생성
		response := models.ForeignCurrentPriceResponse{
			Out: models.ForeignCurrentPriceOutput{
				Sdpr:             "207.8200",
				Prpr:             "207.8200",
				Mxpr:             "0.0000",
				Llam:             "0.0000",
				Oprc:             "207.8200",
				Hprc:             "207.8200",
				Lprc:             "207.8200",
				PrdyVrss:         "0.0000",
				PrdyCtrt:         "0.00",
				Per:              "32.430",
				AcmlTrPbmn:       "0",
				AcmlVol:          "0",
				PrdyVol:          "78788867",
				Bidp1:            "0.0000",
				Askp1:            "0.0000",
				SdprVrssMrktRate: "0.00",
				PrprVrssOprcRate: "",
				SdprVrssHgprRate: "0.00",
				PrprVrssHgprRate: "",
				SdprVrssLwprRate: "0.00",
				PrprVrssLwprRate: "",
			},
			RspCd:  "00000",
			RspMsg: "정상 처리 되었습니다.",
		}

		// 응답 헤더 설정
		w.Header().Set("Content-Type", "application/json")

		// 응답 작성
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 테스트용 클라이언트 생성
	cfg := &config.Config{
		API: config.APIConfig{
			DBSecAppKey:    "test-key",
			DBSecAppSecret: "test-secret",
		},
	}
	apiClient := client.NewDBSecClient(cfg)
	service := NewForeignCurrentPriceService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	if !apiClient.HasValidCredentials() {
		t.Skip("API credentials not configured")
	}

	// 테스트 실행
	t.Run("GetNASDAQStockPrice", func(t *testing.T) {
		data, err := service.GetNASDAQStockPrice("TSLA")
		if err != nil {
			t.Fatalf("Failed to get NASDAQ stock price: %v", err)
		}

		// 데이터 검증
		if data.StockCode != "TSLA" {
			t.Errorf("Expected stock code TSLA, got %s", data.StockCode)
		}
		if data.Market != "나스닥" {
			t.Errorf("Expected market 나스닥, got %s", data.Market)
		}
		if data.Currency != "USD" {
			t.Errorf("Expected currency USD, got %s", data.Currency)
		}
		if data.CurrentPrice == 0 {
			t.Error("Expected non-zero current price")
		}
	})

	t.Run("GetNYStockPrice", func(t *testing.T) {
		data, err := service.GetNYStockPrice("JPM")
		if err != nil {
			t.Fatalf("Failed to get NY stock price: %v", err)
		}

		if data.Market != "뉴욕" {
			t.Errorf("Expected market 뉴욕, got %s", data.Market)
		}
	})

	t.Run("GetUSStockPrice", func(t *testing.T) {
		data, err := service.GetUSStockPrice("AAPL")
		if err != nil {
			t.Fatalf("Failed to get US stock price: %v", err)
		}

		if data.StockCode != "AAPL" {
			t.Errorf("Expected stock code AAPL, got %s", data.StockCode)
		}
		if data.Currency != "USD" {
			t.Errorf("Expected currency USD, got %s", data.Currency)
		}
	})

	t.Run("GetPopularStockPrices", func(t *testing.T) {
		prices, err := service.GetPopularStockPrices()
		if err != nil {
			t.Fatalf("Failed to get popular stock prices: %v", err)
		}

		// 최소 하나 이상의 결과가 있어야 함
		if len(prices) == 0 {
			t.Error("Expected at least one price result")
		}

		// 각 종목 확인
		for code, data := range prices {
			if data.StockCode != code {
				t.Errorf("Expected stock code %s, got %s", code, data.StockCode)
			}
			if data.Currency != "USD" {
				t.Errorf("Expected currency USD, got %s", data.Currency)
			}
		}
	})

	t.Run("GetTechGiantsPrices", func(t *testing.T) {
		prices, err := service.GetTechGiantsPrices()
		if err != nil {
			t.Fatalf("Failed to get tech giants prices: %v", err)
		}

		// 빅테크 종목들이 포함되어 있는지 확인
		expectedStocks := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "META", "NVDA"}
		for _, stock := range expectedStocks {
			if data, exists := prices[stock]; exists {
				if data.StockCode != stock {
					t.Errorf("Expected stock code %s, got %s", stock, data.StockCode)
				}
			}
		}
	})
}

func TestForeignCurrentPriceService_GetMultipleForeignStockPrices(t *testing.T) {
	// 모의 서버 생성
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		var req models.ForeignCurrentPriceRequest
		json.NewDecoder(r.Body).Decode(&req)

		var response models.ForeignCurrentPriceResponse

		// 요청된 종목에 따라 다른 응답 생성
		switch req.In.InputIscd1 {
		case "AAPL":
			response = models.ForeignCurrentPriceResponse{
				Out: models.ForeignCurrentPriceOutput{
					Sdpr: "150.00", Prpr: "155.50", Per: "28.5",
					AcmlVol: "1000000", PrdyVol: "950000",
				},
				RspCd: "00000", RspMsg: "정상 처리 되었습니다.",
			}
		case "TSLA":
			response = models.ForeignCurrentPriceResponse{
				Out: models.ForeignCurrentPriceOutput{
					Sdpr: "207.82", Prpr: "207.82", Per: "32.43",
					AcmlVol: "0", PrdyVol: "78788867",
				},
				RspCd: "00000", RspMsg: "정상 처리 되었습니다.",
			}
		default:
			response = models.ForeignCurrentPriceResponse{
				Out: models.ForeignCurrentPriceOutput{
					Sdpr: "100.00", Prpr: "100.00", Per: "20.0",
				},
				RspCd: "00000", RspMsg: "정상 처리 되었습니다.",
			}
		}

		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 테스트는 실제 구현 시 진행
	t.Skip("Multiple foreign stock prices test requires mock client setup")
}

func TestForeignCurrentPriceService_DataConversion(t *testing.T) {
	service := &ForeignCurrentPriceService{}

	// 테스트 데이터
	output := &models.ForeignCurrentPriceOutput{
		Sdpr:             "207.8200",
		Prpr:             "207.8200",
		Mxpr:             "0.0000",
		Llam:             "0.0000",
		Oprc:             "207.8200",
		Hprc:             "207.8200",
		Lprc:             "207.8200",
		PrdyVrss:         "0.0000",
		PrdyCtrt:         "0.00",
		Per:              "32.430",
		AcmlTrPbmn:       "0",
		AcmlVol:          "0",
		PrdyVol:          "78788867",
		Bidp1:            "0.0000",
		Askp1:            "0.0000",
		SdprVrssMrktRate: "0.00",
	}

	// 변환 테스트
	data := service.convertToForeignCurrentPriceData("TSLA", models.ForeignMarketNASDAQ, output)

	// 검증
	if data.StockCode != "TSLA" {
		t.Errorf("Expected stock code TSLA, got %s", data.StockCode)
	}
	if data.Market != "나스닥" {
		t.Errorf("Expected market 나스닥, got %s", data.Market)
	}
	if data.BasePrice != 207.82 {
		t.Errorf("Expected base price 207.82, got %.2f", data.BasePrice)
	}
	if data.CurrentPrice != 207.82 {
		t.Errorf("Expected current price 207.82, got %.2f", data.CurrentPrice)
	}
	if data.PER != 32.43 {
		t.Errorf("Expected PER 32.43, got %.2f", data.PER)
	}
	if data.YesterdayVolume != 78788867 {
		t.Errorf("Expected yesterday volume 78788867, got %d", data.YesterdayVolume)
	}
	if data.Currency != "USD" {
		t.Errorf("Expected currency USD, got %s", data.Currency)
	}
}

func TestForeignCurrentPriceService_UtilityFunctions(t *testing.T) {
	service := &ForeignCurrentPriceService{}

	t.Run("getMarketName", func(t *testing.T) {
		tests := []struct {
			code     string
			expected string
		}{
			{models.ForeignMarketNY, "뉴욕"},
			{models.ForeignMarketNASDAQ, "나스닥"},
			{models.ForeignMarketAMEX, "아멕스"},
			{"UNKNOWN", "UNKNOWN"},
		}

		for _, test := range tests {
			result := service.getMarketName(test.code)
			if result != test.expected {
				t.Errorf("getMarketName(%s) = %s, expected %s", test.code, result, test.expected)
			}
		}
	})

	t.Run("parseFloat", func(t *testing.T) {
		tests := []struct {
			input    string
			expected float64
		}{
			{"207.8200", 207.82},
			{"  32.430  ", 32.43},
			{"", 0},
			{"invalid", 0},
		}

		for _, test := range tests {
			result := service.parseFloat(test.input)
			if result != test.expected {
				t.Errorf("parseFloat(%s) = %f, expected %f", test.input, result, test.expected)
			}
		}
	})

	t.Run("parseInt", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{"78788867", 78788867},
			{"  1000000  ", 1000000},
			{"", 0},
			{"invalid", 0},
		}

		for _, test := range tests {
			result := service.parseInt(test.input)
			if result != test.expected {
				t.Errorf("parseInt(%s) = %d, expected %d", test.input, result, test.expected)
			}
		}
	})
}