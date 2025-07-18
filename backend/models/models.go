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
	ID             uint      `gorm:"primarykey" json:"id"`
	Symbol         string    `gorm:"index:idx_symbol_timestamp;size:20;not null" json:"symbol"`
	Market         string    `gorm:"size:5;not null" json:"market"`
	OpenPrice      float64   `gorm:"type:decimal(12,4)" json:"open_price"`
	HighPrice      float64   `gorm:"type:decimal(12,4)" json:"high_price"`
	LowPrice       float64   `gorm:"type:decimal(12,4)" json:"low_price"`
	ClosePrice     float64   `gorm:"type:decimal(12,4)" json:"close_price"`
	Volume         int64     `json:"volume"`
	TradeAmount    int64     `json:"trade_amount"`
	PrevClosePrice float64   `gorm:"type:decimal(12,4)" json:"prev_close_price"`
	Change         float64   `gorm:"type:decimal(12,4)" json:"change"`
	ChangeRate     float64   `gorm:"type:decimal(5,2)" json:"change_rate"`
	Timestamp      time.Time `gorm:"index:idx_symbol_timestamp;not null" json:"timestamp"`
	CreatedAt      time.Time `json:"created_at"`
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

// AskingPrice represents bid/ask price information
type AskingPrice struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Symbol      string    `gorm:"index;size:20;not null" json:"symbol"`
	AskPrice1   float64   `gorm:"type:decimal(12,4)" json:"ask_price_1"`
	AskPrice2   float64   `gorm:"type:decimal(12,4)" json:"ask_price_2"`
	AskPrice3   float64   `gorm:"type:decimal(12,4)" json:"ask_price_3"`
	AskPrice4   float64   `gorm:"type:decimal(12,4)" json:"ask_price_4"`
	AskPrice5   float64   `gorm:"type:decimal(12,4)" json:"ask_price_5"`
	BidPrice1   float64   `gorm:"type:decimal(12,4)" json:"bid_price_1"`
	BidPrice2   float64   `gorm:"type:decimal(12,4)" json:"bid_price_2"`
	BidPrice3   float64   `gorm:"type:decimal(12,4)" json:"bid_price_3"`
	BidPrice4   float64   `gorm:"type:decimal(12,4)" json:"bid_price_4"`
	BidPrice5   float64   `gorm:"type:decimal(12,4)" json:"bid_price_5"`
	AskVolume1  int64     `json:"ask_volume_1"`
	AskVolume2  int64     `json:"ask_volume_2"`
	AskVolume3  int64     `json:"ask_volume_3"`
	AskVolume4  int64     `json:"ask_volume_4"`
	AskVolume5  int64     `json:"ask_volume_5"`
	BidVolume1  int64     `json:"bid_volume_1"`
	BidVolume2  int64     `json:"bid_volume_2"`
	BidVolume3  int64     `json:"bid_volume_3"`
	BidVolume4  int64     `json:"bid_volume_4"`
	BidVolume5  int64     `json:"bid_volume_5"`
	TotalAskVol int64     `json:"total_ask_volume"`
	TotalBidVol int64     `json:"total_bid_volume"`
	Timestamp   time.Time `gorm:"not null" json:"timestamp"`
	CreatedAt   time.Time `json:"created_at"`
}