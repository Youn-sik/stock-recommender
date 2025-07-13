package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"stock-recommender/backend/config"
	"stock-recommender/backend/database"
	"stock-recommender/backend/models"
	"stock-recommender/backend/router"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
	cfg    *config.Config
}

func (suite *IntegrationTestSuite) SetupSuite() {
	// Test configuration
	suite.cfg = &config.Config{
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "stockuser",
			Password: "stockpass",
			Name:     "stockdb_test",
		},
		Redis: config.RedisConfig{
			Host: "localhost",
			Port: "6379",
		},
		API: config.APIConfig{
			DBSecAPIKey:  "",
			AIServiceURL: "http://localhost:8001",
		},
	}

	// Initialize test database
	var err error
	suite.db, err = database.Initialize(suite.cfg)
	suite.Require().NoError(err)

	// Setup router
	suite.router = router.Setup(suite.db, suite.cfg)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	// Cleanup test database
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *IntegrationTestSuite) SetupTest() {
	// Clean up test data before each test
	suite.db.Exec("TRUNCATE TABLE stocks, stock_prices, technical_indicators, trading_signals, news_articles RESTART IDENTITY CASCADE")
}

func (suite *IntegrationTestSuite) TestHealthCheck() {
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ok", response["status"])
}

func (suite *IntegrationTestSuite) TestCreateAndGetStock() {
	// Create a test stock
	stock := models.Stock{
		Symbol:   "TEST001",
		Name:     "Test Company",
		Market:   "KR",
		Exchange: "KOSPI",
		Sector:   "Technology",
		Industry: "Software",
		IsActive: true,
	}

	stockJSON, _ := json.Marshal(stock)
	req, _ := http.NewRequest("POST", "/api/v1/admin/stocks", bytes.NewBuffer(stockJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// Get the created stock
	req, _ = http.NewRequest("GET", "/api/v1/stocks/TEST001", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	stockData := response["stock"].(map[string]interface{})
	assert.Equal(suite.T(), "TEST001", stockData["symbol"])
	assert.Equal(suite.T(), "Test Company", stockData["name"])
}

func (suite *IntegrationTestSuite) TestStockPriceOperations() {
	// First create a stock
	stock := models.Stock{
		Symbol:   "TEST002",
		Name:     "Test Company 2",
		Market:   "US",
		Exchange: "NASDAQ",
		IsActive: true,
	}
	suite.db.Create(&stock)

	// Create price data
	priceData := models.StockPrice{
		Symbol:     "TEST002",
		Market:     "US",
		OpenPrice:  100.00,
		HighPrice:  105.00,
		LowPrice:   98.00,
		ClosePrice: 103.50,
		Volume:     1000000,
		Timestamp:  time.Now(),
	}
	suite.db.Create(&priceData)

	// Get price data via API
	req, _ := http.NewRequest("GET", "/api/v1/stocks/TEST002/price", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	price := response["price"].(map[string]interface{})
	assert.Equal(suite.T(), "TEST002", price["symbol"])
	assert.Equal(suite.T(), 103.5, price["close_price"])
}

func (suite *IntegrationTestSuite) TestTechnicalIndicators() {
	// Create stock and price history
	stock := models.Stock{
		Symbol:   "TEST003",
		Name:     "Test Company 3",
		Market:   "KR",
		IsActive: true,
	}
	suite.db.Create(&stock)

	// Create sample technical indicator
	indicator := models.TechnicalIndicator{
		Symbol:         "TEST003",
		IndicatorName:  "RSI",
		IndicatorValue: `{"value": 65.5, "period": 14}`,
		CalculatedAt:   time.Now(),
	}
	suite.db.Create(&indicator)

	// Get indicators via API
	req, _ := http.NewRequest("GET", "/api/v1/stocks/TEST003/indicators", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	indicators := response["indicators"].([]interface{})
	assert.Len(suite.T(), indicators, 1)
	
	firstIndicator := indicators[0].(map[string]interface{})
	assert.Equal(suite.T(), "RSI", firstIndicator["indicator_name"])
}

func (suite *IntegrationTestSuite) TestTradingSignals() {
	// Create stock
	stock := models.Stock{
		Symbol:   "TEST004",
		Name:     "Test Company 4",
		Market:   "KR",
		IsActive: true,
	}
	suite.db.Create(&stock)

	// Create trading signal
	signal := models.TradingSignal{
		Symbol:     "TEST004",
		SignalType: "BUY",
		Strength:   0.8,
		Confidence: 0.75,
		Reasons:    `["RSI oversold", "MACD positive"]`,
		Source:     "AI",
	}
	suite.db.Create(&signal)

	// Get signals via API
	req, _ := http.NewRequest("GET", "/api/v1/signals/TEST004", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	signals := response["signals"].([]interface{})
	assert.Len(suite.T(), signals, 1)
	
	firstSignal := signals[0].(map[string]interface{})
	assert.Equal(suite.T(), "BUY", firstSignal["signal_type"])
	assert.Equal(suite.T(), 0.8, firstSignal["strength"])
}

func (suite *IntegrationTestSuite) TestGetAllStocks() {
	// Create multiple stocks
	stocks := []models.Stock{
		{Symbol: "KR001", Name: "Korean Stock 1", Market: "KR", IsActive: true},
		{Symbol: "US001", Name: "US Stock 1", Market: "US", IsActive: true},
		{Symbol: "KR002", Name: "Korean Stock 2", Market: "KR", IsActive: false},
	}
	
	for _, stock := range stocks {
		suite.db.Create(&stock)
	}

	// Get all stocks
	req, _ := http.NewRequest("GET", "/api/v1/stocks", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	stocks_response := response["stocks"].([]interface{})
	assert.Len(suite.T(), stocks_response, 2) // Only active stocks

	// Test market filter
	req, _ = http.NewRequest("GET", "/api/v1/stocks?market=KR", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	stocks_response = response["stocks"].([]interface{})
	assert.Len(suite.T(), stocks_response, 1) // Only active KR stocks
}

func (suite *IntegrationTestSuite) TestGetAllSignals() {
	// Create stocks and signals
	stocks := []models.Stock{
		{Symbol: "SIG001", Name: "Signal Test 1", Market: "KR", IsActive: true},
		{Symbol: "SIG002", Name: "Signal Test 2", Market: "US", IsActive: true},
	}
	
	for _, stock := range stocks {
		suite.db.Create(&stock)
	}

	signals := []models.TradingSignal{
		{Symbol: "SIG001", SignalType: "BUY", Confidence: 0.8, Source: "AI"},
		{Symbol: "SIG002", SignalType: "SELL", Confidence: 0.7, Source: "AI"},
		{Symbol: "SIG001", SignalType: "HOLD", Confidence: 0.5, Source: "RULE"},
	}
	
	for _, signal := range signals {
		suite.db.Create(&signal)
	}

	// Get all signals
	req, _ := http.NewRequest("GET", "/api/v1/signals", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	signals_response := response["signals"].([]interface{})
	assert.Len(suite.T(), signals_response, 3)

	// Test signal type filter
	req, _ = http.NewRequest("GET", "/api/v1/signals?signal_type=BUY", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	
	signals_response = response["signals"].([]interface{})
	assert.Len(suite.T(), signals_response, 1)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}