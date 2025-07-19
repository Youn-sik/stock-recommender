package models

// API 경로 상수
const (
	// 국내주식
	PathDomesticStockPrice        = "/api/v1/quote/kr-stock/stocks/{symbol}/price"
	PathDomesticStockAsking       = "/api/v1/quote/kr-stock/stocks/{symbol}/asking-price"
	PathDomesticStockDaily        = "/api/v1/quote/kr-stock/stocks/{symbol}/days"
	PathDomesticStockList         = "/api/v1/quote/kr-stock/list"
	PathDomesticStockTicker       = "/api/v1/quote/kr-stock/inquiry/stock-ticker"
	PathDomesticStockCurrentPrice = "/api/v1/quote/kr-stock/inquiry/price"

	// 해외주식
	PathForeignStockPrice        = "/api/v1/quote/foreign-stock/price"
	PathForeignStockDaily        = "/api/v1/quote/foreign-stock/daily-price"
	PathForeignStockTicker       = "/api/v1/quote/overseas-stock/inquiry/stock-ticker"
	PathForeignStockCurrentPrice = "/api/v1/quote/overseas-stock/inquiry/price"

	// 지수
	PathIndexPrice = "/api/v1/quote/index/price"
)

// 시장 구분 코드
const (
	MarketDivStock = "J"  // 주식
	MarketDivETF   = "E"  // ETF
	MarketDivETN   = "EN" // ETN
	MarketDivELW   = "W"  // ELW
	MarketDivIndex = "U"  // 업종&지수
)

// 시장분류구분코드
const (
	MarketClassKosdaq = "1" // 코스닥
	MarketClassKospi  = "4" // 코스피
)

// 트랜잭션 ID
const (
	TrIdStockTicker               = "JCODES"    // 주식종목 조회
	TrIdStockCurrentPrice         = "PRICE"     // 현재가조회
	TrIdForeignStockTicker        = "FSTKCODES" // 해외주식종목 조회
	TrIdForeignStockCurrentPrice  = "FSTKPRICE" // 해외주식현재가조회
)

// 주요 지수 코드
const (
	IndexKOSPI         = "1001" // KOSPI
	IndexKOSDAQ        = "2001" // KOSDAQ
	IndexKOSPI200      = "3001" // KOSPI200
	IndexKOSPILarge    = "1002" // 코스피(대형주)
	IndexKOSPISmall    = "1004" // 코스피(소형주)
	IndexKOSPI50       = "1053" // KOSPI50종합지수
	IndexKOSPI100      = "1054" // KOSPI100종합지수
	IndexKOSPIDiv50    = "1163" // 코스피고배당50
	IndexKOSDAQLarge   = "2002" // 코스닥(대형주)
	IndexKOSDAQSmall   = "2004" // 코스닥(소형주)
	IndexKOSDAQ150     = "2203" // 코스닥 150
	IndexKP200Leverage = "3903" // KP200레버리지지수
	IndexVolatility    = "3907" // 변동성지수
	IndexKRX100        = "0100" // KRX100
	IndexKTOP30        = "0600" // KTOP 30
	IndexKOVIXI00      = "K001" // KOVIXI00
)

// 해외증시구분코드 (종목조회용)
const (
	ExchangeNY     = "NY" // 뉴욕
	ExchangeNASDAQ = "NA" // 나스닥
	ExchangeAMEX   = "AM" // 아멕스
)

// 해외주식 시장분류코드 (현재가조회용)
const (
	ForeignMarketNY     = "FY" // 뉴욕
	ForeignMarketNASDAQ = "FN" // 나스닥
	ForeignMarketAMEX   = "FA" // 아멕스
)