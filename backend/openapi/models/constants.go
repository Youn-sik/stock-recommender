package models

// API 경로 상수
const (
	// 국내주식
	PathDomesticStockPrice  = "/api/v1/quote/kr-stock/stocks/{symbol}/price"
	PathDomesticStockAsking = "/api/v1/quote/kr-stock/stocks/{symbol}/asking-price"
	PathDomesticStockDaily  = "/api/v1/quote/kr-stock/stocks/{symbol}/days"
	PathDomesticStockList   = "/api/v1/quote/kr-stock/list"
	PathDomesticStockTicker = "/api/v1/quote/kr-stock/inquiry/stock-ticker"

	// 해외주식
	PathForeignStockPrice = "/api/v1/quote/foreign-stock/price"
	PathForeignStockDaily = "/api/v1/quote/foreign-stock/daily-price"

	// 지수
	PathIndexPrice = "/api/v1/quote/index/price"
)

// 시장 구분 코드
const (
	MarketDivStock = "J"  // 주식
	MarketDivETF   = "E"  // ETF
	MarketDivETN   = "EN" // ETN
)

// 시장분류구분코드
const (
	MarketClassKosdaq = "1" // 코스닥
	MarketClassKospi  = "4" // 코스피
)

// 트랜잭션 ID
const (
	TrIdStockTicker = "JCODES" // 주식종목 조회
)