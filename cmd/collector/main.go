package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"stock-recommender/backend/config"
	"stock-recommender/backend/database"
	"stock-recommender/backend/models"
	"stock-recommender/backend/services"
	"syscall"
	"time"

	"gorm.io/gorm"
)

func main() {
	log.Println("Starting Stock Data Collector Service")

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	apiClient := services.NewDBSecAPIClient(cfg)
	cacheService := services.NewCacheService(cfg)
	
	// Initialize queue service (optional for collector)
	queueService, err := services.NewQueueService(cfg)
	if err != nil {
		log.Printf("Warning: Failed to initialize queue service: %v", err)
		queueService = nil
	}

	// Create collector
	collector := NewDataCollector(db, apiClient, cacheService, queueService)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start collection scheduler
	go collector.StartScheduler()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down collector service...")

	// Cleanup
	collector.Stop()
	if queueService != nil {
		queueService.Close()
	}

	log.Println("Collector service stopped")
}

type DataCollector struct {
	db           *gorm.DB
	apiClient    *services.DBSecAPIClient
	cacheService *services.CacheService
	queueService *services.QueueService
	stopChan     chan bool
}

func NewDataCollector(
	db *gorm.DB,
	apiClient *services.DBSecAPIClient,
	cacheService *services.CacheService,
	queueService *services.QueueService,
) *DataCollector {
	return &DataCollector{
		db:           db,
		apiClient:    apiClient,
		cacheService: cacheService,
		queueService: queueService,
		stopChan:     make(chan bool),
	}
}

func (dc *DataCollector) StartScheduler() {
	log.Println("Data collector scheduler started")

	// Initial collection
	dc.collectAllStocks()

	// Schedule regular collections
	ticker := time.NewTicker(5 * time.Minute) // 5분마다 수집
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dc.collectAllStocks()
		case <-dc.stopChan:
			log.Println("Collector scheduler stopped")
			return
		}
	}
}

func (dc *DataCollector) collectAllStocks() {
	log.Println("Starting stock data collection cycle")

	// Get active stocks
	var stocks []models.Stock
	err := dc.db.Where("is_active = ?", true).Find(&stocks).Error
	if err != nil {
		log.Printf("Failed to fetch active stocks: %v", err)
		return
	}

	successCount := 0
	errorCount := 0

	for _, stock := range stocks {
		err := dc.collectStockData(stock.Symbol, stock.Market)
		if err != nil {
			log.Printf("Failed to collect data for %s: %v", stock.Symbol, err)
			errorCount++
		} else {
			successCount++
		}

		// Rate limiting - 1초 대기
		time.Sleep(1 * time.Second)
	}

	log.Printf("Collection cycle completed: %d success, %d errors", successCount, errorCount)
}

func (dc *DataCollector) collectStockData(symbol, market string) error {
	log.Printf("Collecting data for %s (%s)", symbol, market)

	// Check if using real API or mock data
	var stockPrice *models.StockPrice
	var err error

	if dc.apiClient.HasValidAPIKey() {
		// Use real API
		stockPrice, err = dc.apiClient.FetchStockPrice(symbol, market)
	} else {
		// Use mock data for development
		log.Printf("Using mock data for %s (no API key configured)", symbol)
		stockPrice = dc.apiClient.GenerateMockData(symbol, market)
	}

	if err != nil {
		return fmt.Errorf("failed to fetch stock price: %w", err)
	}

	if stockPrice == nil {
		return fmt.Errorf("received nil stock price data")
	}

	// Save to database
	err = dc.db.Create(stockPrice).Error
	if err != nil {
		return fmt.Errorf("failed to save stock price: %w", err)
	}

	// Update cache
	dc.cacheService.SetStockPrice(symbol, stockPrice)

	// Publish to message queue
	if dc.queueService != nil {
		dc.queueService.PublishPriceUpdate(symbol, market, stockPrice)
	}

	log.Printf("Successfully collected and saved data for %s", symbol)
	return nil
}

func (dc *DataCollector) Stop() {
	close(dc.stopChan)
}