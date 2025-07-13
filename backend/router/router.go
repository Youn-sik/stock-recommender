package router

import (
	"stock-recommender/backend/config"
	"stock-recommender/backend/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	
	r := gin.Default()

	// Middleware
	r.Use(CORSMiddleware())
	r.Use(LoggingMiddleware())

	// Initialize handlers
	stockHandler := handlers.NewStockHandler(db, cfg)
	signalHandler := handlers.NewSignalHandler(db, cfg)
	healthHandler := handlers.NewHealthHandler(db)
	adminHandler := handlers.NewAdminHandler(db, cfg)

	// Health check
	r.GET("/health", healthHandler.HealthCheck)

	// API routes
	api := r.Group("/api/v1")
	{
		// Stock endpoints
		stocks := api.Group("/stocks")
		{
			stocks.GET("/", stockHandler.GetStocks)
			stocks.GET("/:symbol", stockHandler.GetStock)
			stocks.GET("/:symbol/price", stockHandler.GetStockPrice)
			stocks.GET("/:symbol/indicators", stockHandler.GetIndicators)
		}

		// Signal endpoints
		signals := api.Group("/signals")
		{
			signals.GET("/", signalHandler.GetSignals)
			signals.GET("/:symbol", signalHandler.GetSignalsBySymbol)
		}

		// Admin endpoints (for testing and management)
		admin := api.Group("/admin")
		{
			// Stock management
			admin.POST("/stocks", adminHandler.CreateStock)
			admin.GET("/stocks", adminHandler.GetAllStocks)
			admin.PUT("/stocks/:symbol/status", adminHandler.UpdateStockStatus)
			admin.DELETE("/stocks/:symbol", adminHandler.DeleteStock)
			
			// Data collection
			admin.POST("/collect/:symbol", adminHandler.TriggerDataCollection)
			admin.POST("/collect/all", adminHandler.TriggerAllDataCollection)
			admin.POST("/initialize/major-stocks", adminHandler.InitializeMajorStocks)
			
			// System status
			admin.GET("/api-status", adminHandler.GetAPIStatus)
			admin.GET("/database/stats", adminHandler.GetDatabaseStats)
		}
	}

	return r
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	})
}