package handlers

import (
	"net/http"
	"stock-recommender/backend/config"
	"stock-recommender/backend/models"
	"stock-recommender/backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db            *gorm.DB
	dataCollector *services.DataCollectorService
	config        *config.Config
}

func NewAdminHandler(db *gorm.DB, cfg *config.Config) *AdminHandler {
	return &AdminHandler{
		db:            db,
		dataCollector: services.NewDataCollectorService(db, cfg),
		config:        cfg,
	}
}

// 종목 등록
func (h *AdminHandler) CreateStock(c *gin.Context) {
	var req struct {
		Symbol   string `json:"symbol" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Market   string `json:"market" binding:"required"`
		Exchange string `json:"exchange"`
		Sector   string `json:"sector"`
		Industry string `json:"industry"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 중복 체크
	var existing models.Stock
	result := h.db.Where("symbol = ? AND market = ?", req.Symbol, req.Market).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Stock already exists"})
		return
	}

	stock := models.Stock{
		Symbol:   req.Symbol,
		Name:     req.Name,
		Market:   req.Market,
		Exchange: req.Exchange,
		Sector:   req.Sector,
		Industry: req.Industry,
		IsActive: true,
	}

	if err := h.db.Create(&stock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Stock created successfully",
		"stock":   stock,
	})
}

// 특정 종목 데이터 수집 트리거
func (h *AdminHandler) TriggerDataCollection(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol is required"})
		return
	}

	// 종목 정보 조회
	var stock models.Stock
	if err := h.db.Where("symbol = ?", symbol).First(&stock).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	// 데이터 수집 실행
	err := h.dataCollector.CollectStockData(stock.Symbol, stock.Market)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to collect data",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data collection triggered successfully",
		"symbol":  symbol,
	})
}

// 전체 종목 데이터 수집 트리거
func (h *AdminHandler) TriggerAllDataCollection(c *gin.Context) {
	go func() {
		err := h.dataCollector.CollectAllStocks()
		if err != nil {
			// 로그에 기록 (비동기 처리이므로 응답으로는 보내지 않음)
			// log.Printf("Batch data collection failed: %v", err)
		}
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Batch data collection started",
	})
}

// API 상태 확인
func (h *AdminHandler) GetAPIStatus(c *gin.Context) {
	status := h.dataCollector.GetAPIStatus()
	c.JSON(http.StatusOK, gin.H{
		"api_status": status,
	})
}

// 주요 종목 초기화
func (h *AdminHandler) InitializeMajorStocks(c *gin.Context) {
	err := h.dataCollector.InitializeMajorStocks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to initialize major stocks",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Major stocks initialized successfully",
	})
}

// 종목 목록 조회 (관리용)
func (h *AdminHandler) GetAllStocks(c *gin.Context) {
	var stocks []models.Stock
	if err := h.db.Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stocks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stocks": stocks,
		"count":  len(stocks),
	})
}

// 종목 활성화/비활성화
func (h *AdminHandler) UpdateStockStatus(c *gin.Context) {
	symbol := c.Param("symbol")
	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.db.Model(&models.Stock{}).Where("symbol = ?", symbol).Update("is_active", req.IsActive)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock status"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Stock status updated",
		"symbol":    symbol,
		"is_active": req.IsActive,
	})
}

// 종목 삭제
func (h *AdminHandler) DeleteStock(c *gin.Context) {
	symbol := c.Param("symbol")

	result := h.db.Where("symbol = ?", symbol).Delete(&models.Stock{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete stock"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock deleted successfully",
		"symbol":  symbol,
	})
}

// 데이터베이스 통계
func (h *AdminHandler) GetDatabaseStats(c *gin.Context) {
	var stats struct {
		TotalStocks      int64 `json:"total_stocks"`
		ActiveStocks     int64 `json:"active_stocks"`
		TotalPricePoints int64 `json:"total_price_points"`
		TotalSignals     int64 `json:"total_signals"`
		LastUpdate       string `json:"last_update"`
	}

	h.db.Model(&models.Stock{}).Count(&stats.TotalStocks)
	h.db.Model(&models.Stock{}).Where("is_active = ?", true).Count(&stats.ActiveStocks)
	h.db.Model(&models.StockPrice{}).Count(&stats.TotalPricePoints)
	h.db.Model(&models.TradingSignal{}).Count(&stats.TotalSignals)

	var lastPrice models.StockPrice
	h.db.Order("timestamp DESC").First(&lastPrice)
	if lastPrice.ID != 0 {
		stats.LastUpdate = lastPrice.Timestamp.Format("2006-01-02 15:04:05")
	} else {
		stats.LastUpdate = "No data"
	}

	c.JSON(http.StatusOK, stats)
}