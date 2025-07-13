package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  string    `json:"database"`
	Version   string    `json:"version"`
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Version:   "1.0.0",
	}

	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		response.Status = "error"
		response.Database = "connection failed"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		response.Status = "error"
		response.Database = "ping failed"
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Database = "connected"
	c.JSON(http.StatusOK, response)
}