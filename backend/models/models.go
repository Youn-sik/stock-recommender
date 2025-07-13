package models

import (
	"time"

	"gorm.io/gorm"
)

// Stock represents a stock symbol information
type Stock struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Symbol      string         `gorm:"uniqueIndex;size:20;not null" json:"symbol"`
	Name        string         `gorm:"size:100" json:"name"`
	Market      string         `gorm:"size:5;not null" json:"market"` // KR or US
	Exchange    string         `gorm:"size:20" json:"exchange"`       // KOSPI, NASDAQ, etc.
	Sector      string         `gorm:"size:50" json:"sector"`
	Industry    string         `gorm:"size:50" json:"industry"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// StockPrice represents historical and real-time stock price data
type StockPrice struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Symbol    string    `gorm:"index:idx_symbol_timestamp;size:20;not null" json:"symbol"`
	Market    string    `gorm:"size:5;not null" json:"market"`
	OpenPrice float64   `gorm:"type:decimal(12,4)" json:"open_price"`
	HighPrice float64   `gorm:"type:decimal(12,4)" json:"high_price"`
	LowPrice  float64   `gorm:"type:decimal(12,4)" json:"low_price"`
	ClosePrice float64  `gorm:"type:decimal(12,4)" json:"close_price"`
	Volume    int64     `json:"volume"`
	Timestamp time.Time `gorm:"index:idx_symbol_timestamp;not null" json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
}

// TechnicalIndicator represents calculated technical indicators
type TechnicalIndicator struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Symbol        string    `gorm:"index:idx_symbol_calculated;size:20;not null" json:"symbol"`
	IndicatorName string    `gorm:"size:50;not null" json:"indicator_name"`
	IndicatorValue string   `gorm:"type:jsonb" json:"indicator_value"` // JSON for flexible data
	CalculatedAt  time.Time `gorm:"index:idx_symbol_calculated;not null" json:"calculated_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// TradingSignal represents buy/sell/hold signals
type TradingSignal struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	Symbol     string    `gorm:"index:idx_symbol_created;size:20;not null" json:"symbol"`
	SignalType string    `gorm:"size:10;not null" json:"signal_type"` // BUY, SELL, HOLD
	Strength   float64   `gorm:"type:decimal(3,2)" json:"strength"`   // 0.0 ~ 1.0
	Confidence float64   `gorm:"type:decimal(3,2)" json:"confidence"` // 0.0 ~ 1.0
	Reasons    string    `gorm:"type:jsonb" json:"reasons"`           // JSON array of reasons
	Source     string    `gorm:"size:20" json:"source"`               // AI, RULE, MANUAL
	CreatedAt  time.Time `gorm:"index:idx_symbol_created" json:"created_at"`
}

// NewsArticle represents news articles for sentiment analysis
type NewsArticle struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Title          string    `gorm:"type:text;not null" json:"title"`
	Content        string    `gorm:"type:text" json:"content"`
	URL            string    `gorm:"type:text" json:"url"`
	Source         string    `gorm:"size:100" json:"source"`
	SentimentScore float64   `gorm:"type:decimal(3,2)" json:"sentiment_score"` // -1.0 ~ 1.0
	Keywords       string    `gorm:"type:jsonb" json:"keywords"`               // JSON array
	RelatedSymbols string    `gorm:"type:jsonb" json:"related_symbols"`        // JSON array
	PublishedAt    time.Time `json:"published_at"`
	CreatedAt      time.Time `json:"created_at"`
}

// AIDecisionRequest represents data sent to AI service
type AIDecisionRequest struct {
	Symbol      string                 `json:"symbol"`
	Market      string                 `json:"market"`
	Price       StockPrice            `json:"price"`
	Indicators  map[string]float64    `json:"indicators"`
	NewsScore   float64               `json:"news_score,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AIDecisionResponse represents response from AI service
type AIDecisionResponse struct {
	Symbol     string    `json:"symbol"`
	Decision   string    `json:"decision"`   // BUY/SELL/HOLD
	Confidence float64   `json:"confidence"` // 0.0 ~ 1.0
	Reasoning  []string  `json:"reasoning"`
	Timestamp  time.Time `json:"timestamp"`
}