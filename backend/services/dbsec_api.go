package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stock-recommender/backend/config"
	"stock-recommender/backend/models"
	"time"
)

type DBSecAPIClient struct {
	apiKey     string
	baseURL    string
	client     *http.Client
	rateLimiter chan struct{}
}

// DB증권 API 응답 구조체
type DBSecResponse struct {
	Symbol    string  `json:"symbol"`
	Name      string  `json:"name"`
	Market    string  `json:"market"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
	Timestamp string  `json:"timestamp"`
}

func NewDBSecAPIClient(cfg *config.Config) *DBSecAPIClient {
	// Rate limiter: 초당 10 요청으로 제한
	rateLimiter := make(chan struct{}, 10)
	go func() {
		for {
			time.Sleep(100 * time.Millisecond) // 100ms마다 1개 허용
			select {
			case rateLimiter <- struct{}{}:
			default:
			}
		}
	}()

	return &DBSecAPIClient{
		apiKey:      cfg.API.DBSecAPIKey,
		baseURL:     "https://openapi.dbsec.co.kr/v1",
		client:      &http.Client{Timeout: 30 * time.Second},
		rateLimiter: rateLimiter,
	}
}

func (c *DBSecAPIClient) FetchStockPrice(symbol string, market string) (*models.StockPrice, error) {
	// Rate limiting
	<-c.rateLimiter

	url := fmt.Sprintf("%s/quote/%s", c.baseURL, symbol)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// API 키 설정
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse DBSecResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 시간 파싱
	timestamp, err := time.Parse(time.RFC3339, apiResponse.Timestamp)
	if err != nil {
		timestamp = time.Now()
	}

	stockPrice := &models.StockPrice{
		Symbol:     apiResponse.Symbol,
		Market:     market,
		OpenPrice:  apiResponse.Open,
		HighPrice:  apiResponse.High,
		LowPrice:   apiResponse.Low,
		ClosePrice: apiResponse.Close,
		Volume:     apiResponse.Volume,
		Timestamp:  timestamp,
		CreatedAt:  time.Now(),
	}

	return stockPrice, nil
}

func (c *DBSecAPIClient) FetchHistoricalData(symbol string, market string, days int) ([]*models.StockPrice, error) {
	// Rate limiting
	<-c.rateLimiter

	url := fmt.Sprintf("%s/history/%s?days=%d", c.baseURL, symbol, days)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponses []DBSecResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponses); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var stockPrices []*models.StockPrice
	for _, apiResp := range apiResponses {
		timestamp, err := time.Parse(time.RFC3339, apiResp.Timestamp)
		if err != nil {
			timestamp = time.Now()
		}

		stockPrice := &models.StockPrice{
			Symbol:     apiResp.Symbol,
			Market:     market,
			OpenPrice:  apiResp.Open,
			HighPrice:  apiResp.High,
			LowPrice:   apiResp.Low,
			ClosePrice: apiResp.Close,
			Volume:     apiResp.Volume,
			Timestamp:  timestamp,
			CreatedAt:  time.Now(),
		}
		stockPrices = append(stockPrices, stockPrice)
	}

	return stockPrices, nil
}

// Mock 데이터 생성 (개발용)
func (c *DBSecAPIClient) GenerateMockData(symbol string, market string) *models.StockPrice {
	basePrice := 100000.0
	if market == "US" {
		basePrice = 150.0
	}

	// 간단한 랜덤 변동 시뮬레이션
	variation := 0.05 // 5% 변동
	randFactor1 := float64(time.Now().Unix()%100)/100.0
	randFactor2 := float64(time.Now().Nanosecond()%100)/100.0
	open := basePrice * (1 + (variation * (2*randFactor1 - 1)))
	high := open * (1 + variation/2)
	low := open * (1 - variation/2)
	close := open * (1 + (variation * (2*randFactor2 - 1)))

	return &models.StockPrice{
		Symbol:     symbol,
		Market:     market,
		OpenPrice:  open,
		HighPrice:  high,
		LowPrice:   low,
		ClosePrice: close,
		Volume:     int64(1000000 + time.Now().Unix()%500000),
		Timestamp:  time.Now(),
		CreatedAt:  time.Now(),
	}
}

// HasValidAPIKey checks if API key is configured
func (c *DBSecAPIClient) HasValidAPIKey() bool {
	return c.apiKey != ""
}

func (c *DBSecAPIClient) HealthCheck() error {
	if !c.HasValidAPIKey() {
		return fmt.Errorf("no API key configured")
	}

	url := fmt.Sprintf("%s/health", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("health check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode)
	}

	return nil
}