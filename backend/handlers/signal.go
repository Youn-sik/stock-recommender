package handlers

import (
	"net/http"
	"stock-recommender/backend/config"
	"stock-recommender/backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SignalHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewSignalHandler(db *gorm.DB, cfg *config.Config) *SignalHandler {
	return &SignalHandler{db: db, cfg: cfg}
}

func (h *SignalHandler) GetSignals(c *gin.Context) {
	var signals []models.TradingSignal
	
	// Query parameters
	signalType := c.Query("signal_type") // BUY, SELL, HOLD
	market := c.Query("market")          // KR, US
	limitStr := c.DefaultQuery("limit", "50")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	
	query := h.db.Model(&models.TradingSignal{})
	
	if signalType != "" {
		query = query.Where("signal_type = ?", signalType)
	}
	
	// Join with stock to filter by market
	if market != "" {
		query = query.Joins("JOIN stocks ON stocks.symbol = trading_signals.symbol").
			Where("stocks.market = ?", market)
	}
	
	if err := query.Order("created_at desc").
		Limit(limit).
		Find(&signals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch signals"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"signals": signals,
		"total":   len(signals),
	})
}

func (h *SignalHandler) GetSignalsBySymbol(c *gin.Context) {
	symbol := c.Param("symbol")
	limitStr := c.DefaultQuery("limit", "20")
	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	
	var signals []models.TradingSignal
	if err := h.db.Where("symbol = ?", symbol).
		Order("created_at desc").
		Limit(limit).
		Find(&signals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch signals"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"symbol":  symbol,
		"signals": signals,
		"total":   len(signals),
	})
}