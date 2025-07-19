package domestic

import (
	"testing"

	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/models"
	"stock-recommender/backend/openapi/utils"
)

func TestCurrentPriceService_GetCurrentPrice(t *testing.T) {
	// 모의 데이터 생성
	mockData := models.CurrentPriceOutput{
		Sdpr:             "53900",
		Prpr:             "55550",
		Mxpr:             "70000",
		Llam:             "37800",
		Oprc:             "54300",
		Hprc:             "55900",
		Lprc:             "54200",
		PrdyVrss:         "1650",
		PrdyCtrt:         "3.06",
		Per:              "10.89",
		Pbr:              "0.93",
		AcmlTrPbmn:       "400303637800",
		AcmlVol:          "7240324",
		PrdyVol:          "13439520",
		Bidp1:            "55500",
		Askp1:            "55600",
		SdprVrssMrktRate: "0.74",
		PrprVrssOprcRate: "-2.25",
		SdprVrssHgprRate: "3.71",
		PrprVrssHgprRate: "0.63",
		SdprVrssLwprRate: "0.56",
		PrprVrssLwprRate: "-2.43",
	}

	// 모의 서버 생성
	handler := utils.CreateCurrentPriceMockHandler(t, models.PathDomesticStockCurrentPrice, "005930", mockData)
	mockServer := utils.NewMockServer(t, handler)
	defer mockServer.Close()

	// 테스트용 클라이언트 생성
	cfg := utils.CreateTestConfig()
	apiClient := client.NewDBSecClient(cfg)
	service := NewCurrentPriceService(apiClient)

	// 실제 API 호출이 설정되어 있지 않은 경우 스킵
	utils.SkipIfNoCredentials(t, apiClient)

	// 테스트 실행
	t.Run("GetStockPrice", func(t *testing.T) {
		data, err := service.GetStockPrice("005930")
		if err != nil {
			t.Fatalf("Failed to get stock price: %v", err)
		}

		// 데이터 검증
		utils.AssertStringEqual(t, "005930", data.StockCode, "Stock code")
		utils.AssertFloatEqual(t, 55550, data.CurrentPrice, "Current price")
		utils.AssertIntEqual(t, 7240324, data.TradingVolume, "Trading volume")
	})

	t.Run("GetKOSPIPrice", func(t *testing.T) {
		data, err := service.GetKOSPIPrice()
		if err != nil {
			t.Fatalf("Failed to get KOSPI price: %v", err)
		}

		if data.StockCode != models.IndexKOSPI {
			t.Errorf("Expected index code %s, got %s", models.IndexKOSPI, data.StockCode)
		}
	})

	t.Run("GetMultipleStockPrices", func(t *testing.T) {
		codes := []string{"005930", "000660", "035720"}
		prices, err := service.GetMultipleStockPrices(codes)
		if err != nil {
			t.Fatalf("Failed to get multiple stock prices: %v", err)
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
		}
	})
}

func TestCurrentPriceService_DataConversion(t *testing.T) {
	service := &CurrentPriceService{}

	// 테스트 데이터
	output := &models.CurrentPriceOutput{
		Sdpr:             "53900",
		Prpr:             "55550",
		Mxpr:             "70000",
		Llam:             "37800",
		Oprc:             "54300",
		Hprc:             "55900",
		Lprc:             "54200",
		PrdyVrss:         "1650",
		PrdyCtrt:         "3.06",
		Per:              "10.89",
		Pbr:              "0.93",
		AcmlTrPbmn:       "400303637800",
		AcmlVol:          "7240324",
		PrdyVol:          "13439520",
		Bidp1:            "55500",
		Askp1:            "55600",
		SdprVrssMrktRate: "0.74",
		PrprVrssOprcRate: "-2.25",
	}

	// 변환 테스트
	data := service.convertToCurrentPriceData("005930", output)

	// 검증
	if data.StockCode != "005930" {
		t.Errorf("Expected stock code 005930, got %s", data.StockCode)
	}
	if data.BasePrice != 53900 {
		t.Errorf("Expected base price 53900, got %.0f", data.BasePrice)
	}
	if data.CurrentPrice != 55550 {
		t.Errorf("Expected current price 55550, got %.0f", data.CurrentPrice)
	}
	if data.PriceChangeRate != 3.06 {
		t.Errorf("Expected price change rate 3.06, got %.2f", data.PriceChangeRate)
	}
	if data.TradingVolume != 7240324 {
		t.Errorf("Expected trading volume 7240324, got %d", data.TradingVolume)
	}
	if data.PER != 10.89 {
		t.Errorf("Expected PER 10.89, got %.2f", data.PER)
	}
	if data.PBR != 0.93 {
		t.Errorf("Expected PBR 0.93, got %.2f", data.PBR)
	}
}

func TestCurrentPriceService_ParseFunctions(t *testing.T) {

	t.Run("parseFloat", func(t *testing.T) {
		tests := []struct {
			input    string
			expected float64
		}{
			{"123.45", 123.45},
			{"  678.90  ", 678.90},
			{"", 0},
			{"invalid", 0},
		}

		for _, test := range tests {
			result := utils.ParseFloat(test.input)
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
			{"12345", 12345},
			{"  67890  ", 67890},
			{"", 0},
			{"invalid", 0},
		}

		for _, test := range tests {
			result := utils.ParseInt(test.input)
			if result != test.expected {
				t.Errorf("parseInt(%s) = %d, expected %d", test.input, result, test.expected)
			}
		}
	})
}