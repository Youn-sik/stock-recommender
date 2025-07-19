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

func TestForeignStockTickerService_GetForeignStockTickers(t *testing.T) {
	// 테스트 서버 생성
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 경로 확인
		if r.URL.Path != models.PathForeignStockTicker {
			t.Errorf("Expected path %s, got %s", models.PathForeignStockTicker, r.URL.Path)
		}

		// 메소드 확인
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 헤더 확인
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// 요청 본문 파싱
		var req models.ForeignStockTickerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// 요청 데이터 검증
		if req.In.InputDataCode != models.ExchangeNASDAQ {
			t.Errorf("Expected exchange code %s, got %s", models.ExchangeNASDAQ, req.In.InputDataCode)
		}

		// 응답 생성
		response := models.ForeignStockTickerResponse{
			Out: []models.ForeignStockTickerOutput{
				{
					Iscd:         "AAPL",
					KorIsnm:      "애플",
					BstpLargName: "IT",
					ExchClsCode2: "FN",
					SelnVolUnit:  "1",
					ShnuVolUnit:  "1",
				},
				{
					Iscd:         "MSFT",
					KorIsnm:      "마이크로소프트",
					BstpLargName: "IT",
					ExchClsCode2: "FN",
					SelnVolUnit:  "1",
					ShnuVolUnit:  "1",
				},
				{
					Iscd:         "GOOGL",
					KorIsnm:      "알파벳 A",
					BstpLargName: "IT",
					ExchClsCode2: "FN",
					SelnVolUnit:  "1",
					ShnuVolUnit:  "1",
				},
			},
			RspCd:  "00000",
			RspMsg: "정상 처리 되었습니다.",
		}

		// 응답 헤더 설정
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("cont_yn", "N")
		w.Header().Set("cont_key", "")

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
	service := NewForeignStockTickerService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	if !apiClient.HasValidCredentials() {
		t.Skip("API credentials not configured")
	}

	// 테스트 실행
	t.Run("GetNASDAQStocks", func(t *testing.T) {
		stocks, err := service.GetNASDAQStocks()
		if err != nil {
			t.Fatalf("Failed to get NASDAQ stocks: %v", err)
		}

		if len(stocks) == 0 {
			t.Error("Expected at least one stock, got none")
		}

		// 첫 번째 종목 검증
		if len(stocks) > 0 {
			stock := stocks[0]
			if stock.StockCode == "" {
				t.Error("Expected stock code, got empty")
			}
			if stock.KoreanName == "" {
				t.Error("Expected Korean name, got empty")
			}
			if stock.Exchange != "나스닥" {
				t.Errorf("Expected exchange 나스닥, got %s", stock.Exchange)
			}
		}
	})

	t.Run("GetNYStocks", func(t *testing.T) {
		stocks, err := service.GetNYStocks()
		if err != nil {
			t.Fatalf("Failed to get NY stocks: %v", err)
		}

		// NY 주식이 있는 경우 검증
		if len(stocks) > 0 {
			stock := stocks[0]
			if stock.StockCode == "" {
				t.Error("Expected stock code, got empty")
			}
			if stock.Exchange != "뉴욕" {
				t.Errorf("Expected exchange 뉴욕, got %s", stock.Exchange)
			}
		}
	})

	t.Run("GetTechStocks", func(t *testing.T) {
		stocks, err := service.GetTechStocks()
		if err != nil {
			t.Fatalf("Failed to get tech stocks: %v", err)
		}

		// IT 업종 종목이 있는 경우 검증
		for _, stock := range stocks {
			if stock.SectorName != "IT" {
				t.Errorf("Expected sector IT, got %s", stock.SectorName)
			}
		}
	})
}

func TestForeignStockTickerService_GetAllUSStocks(t *testing.T) {
	// 모의 서버 생성
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++

		var req models.ForeignStockTickerRequest
		json.NewDecoder(r.Body).Decode(&req)

		var response models.ForeignStockTickerResponse

		// 요청된 거래소에 따라 다른 응답 생성
		switch req.In.InputDataCode {
		case models.ExchangeNY:
			response = models.ForeignStockTickerResponse{
				Out: []models.ForeignStockTickerOutput{
					{Iscd: "JPM", KorIsnm: "JP모건체이스", BstpLargName: "금융", ExchClsCode2: "FN", SelnVolUnit: "1", ShnuVolUnit: "1"},
				},
				RspCd: "00000", RspMsg: "정상 처리 되었습니다.",
			}
		case models.ExchangeNASDAQ:
			response = models.ForeignStockTickerResponse{
				Out: []models.ForeignStockTickerOutput{
					{Iscd: "AAPL", KorIsnm: "애플", BstpLargName: "IT", ExchClsCode2: "FN", SelnVolUnit: "1", ShnuVolUnit: "1"},
				},
				RspCd: "00000", RspMsg: "정상 처리 되었습니다.",
			}
		case models.ExchangeAMEX:
			response = models.ForeignStockTickerResponse{
				Out: []models.ForeignStockTickerOutput{
					{Iscd: "SPY", KorIsnm: "SPDR S&P500 ETF", BstpLargName: "ETF", ExchClsCode2: "FN", SelnVolUnit: "1", ShnuVolUnit: "1"},
				},
				RspCd: "00000", RspMsg: "정상 처리 되었습니다.",
			}
		}

		w.Header().Set("cont_yn", "N")
		w.Header().Set("cont_key", "")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 테스트는 실제 구현 시 진행
	t.Skip("All US stocks test requires mock client setup")
}

func TestForeignStockTickerService_DataConversion(t *testing.T) {
	service := &ForeignStockTickerService{}

	// 테스트 데이터
	output := &models.ForeignStockTickerOutput{
		Iscd:         "AAPL",
		KorIsnm:      "애플",
		BstpLargName: "IT",
		ExchClsCode2: "FN",
		SelnVolUnit:  "1",
		ShnuVolUnit:  "1",
	}

	// 변환 테스트
	data := service.convertToForeignStockData(models.ExchangeNASDAQ, output)

	// 검증
	if data.StockCode != "AAPL" {
		t.Errorf("Expected stock code AAPL, got %s", data.StockCode)
	}
	if data.KoreanName != "애플" {
		t.Errorf("Expected Korean name 애플, got %s", data.KoreanName)
	}
	if data.SectorName != "IT" {
		t.Errorf("Expected sector IT, got %s", data.SectorName)
	}
	if data.Exchange != "나스닥" {
		t.Errorf("Expected exchange 나스닥, got %s", data.Exchange)
	}
	if data.SellUnit != 1 {
		t.Errorf("Expected sell unit 1, got %d", data.SellUnit)
	}
	if data.BuyUnit != 1 {
		t.Errorf("Expected buy unit 1, got %d", data.BuyUnit)
	}
}

func TestForeignStockTickerService_UtilityFunctions(t *testing.T) {
	service := &ForeignStockTickerService{}

	t.Run("getExchangeName", func(t *testing.T) {
		tests := []struct {
			code     string
			expected string
		}{
			{models.ExchangeNY, "뉴욕"},
			{models.ExchangeNASDAQ, "나스닥"},
			{models.ExchangeAMEX, "아멕스"},
			{"UNKNOWN", "UNKNOWN"},
		}

		for _, test := range tests {
			result := service.getExchangeName(test.code)
			if result != test.expected {
				t.Errorf("getExchangeName(%s) = %s, expected %s", test.code, result, test.expected)
			}
		}
	})

	t.Run("parseInt", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{"1", 1},
			{"100", 100},
			{"  50  ", 50},
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

	t.Run("getContYn", func(t *testing.T) {
		tests := []struct {
			contKey  string
			expected string
		}{
			{"", "N"},
			{"some_key", "Y"},
		}

		for _, test := range tests {
			result := service.getContYn(test.contKey)
			if result != test.expected {
				t.Errorf("getContYn(%s) = %s, expected %s", test.contKey, result, test.expected)
			}
		}
	})
}