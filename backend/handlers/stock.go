package handlers

import (
	"net/http"
	"stock-recommender/backend/config"
	"stock-recommender/backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StockHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewStockHandler(db *gorm.DB, cfg *config.Config) *StockHandler {
	return &StockHandler{db: db, cfg: cfg}
}

func (h *StockHandler) GetStocks(c *gin.Context) {
	var stocks []models.Stock
	
	market := c.Query("market") // KR or US
	query := h.db.Where("is_active = ?", true)
	
	if market != "" {
		query = query.Where("market = ?", market)
	}
	
	if err := query.Find(&stocks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stocks"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"stocks": stocks})
}

func (h *StockHandler) GetStock(c *gin.Context) {
	symbol := c.Param("symbol")
	
	var stock models.Stock
	if err := h.db.Where("symbol = ?", symbol).First(&stock).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stock not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"stock": stock})
}

func (h *StockHandler) GetStockPrice(c *gin.Context) {
	symbol := c.Param("symbol")
	
	var price models.StockPrice
	if err := h.db.Where("symbol = ?", symbol).
		Order("timestamp desc").
		First(&price).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Price data not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"price": price})
}

func (h *StockHandler) GetIndicators(c *gin.Context) {
	symbol := c.Param("symbol")
	
	var indicators []models.TechnicalIndicator
	if err := h.db.Where("symbol = ?", symbol).
		Order("calculated_at desc").
		Limit(50).
		Find(&indicators).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch indicators"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"indicators": indicators})
}

func (h *StockHandler) CreateStock(c *gin.Context) {
	var stock models.Stock
	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set default values
	stock.IsActive = true
	stock.CreatedAt = time.Now()
	stock.UpdatedAt = time.Now()
	
	if err := h.db.Create(&stock).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"stock": stock})
}

func (h *StockHandler) TriggerCollection(c *gin.Context) {
	symbol := c.Param("symbol")
	
	// TODO: Trigger data collection for specific symbol
	// This will be implemented when we add RabbitMQ
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Collection triggered",
		"symbol":  symbol,
		"status":  "pending",
	})
}