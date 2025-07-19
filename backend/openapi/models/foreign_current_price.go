package models

// ForeignCurrentPriceRequest 해외주식현재가조회 요청
type ForeignCurrentPriceRequest struct {
	In ForeignCurrentPriceInput `json:"In"`
}

// ForeignCurrentPriceInput 해외주식현재가조회 입력
type ForeignCurrentPriceInput struct {
	InputCondMrktDivCode string `json:"InputCondMrktDivCode"` // 시장분류코드 (FY: 뉴욕, FN: 나스닥, FA: 아멕스)
	InputIscd1           string `json:"InputIscd1"`           // 해외주식종목코드 (예: TSLA, AAPL)
}

// ForeignCurrentPriceResponse 해외주식현재가조회 응답
type ForeignCurrentPriceResponse struct {
	Out    ForeignCurrentPriceOutput `json:"Out"`
	RspCd  string                    `json:"rsp_cd"`  // 응답코드
	RspMsg string                    `json:"rsp_msg"` // 응답메시지
}

// ForeignCurrentPriceOutput 해외주식현재가조회 출력
type ForeignCurrentPriceOutput struct {
	Sdpr                 string `json:"Sdpr"`                 // 기준가
	Prpr                 string `json:"Prpr"`                 // 현재가
	Mxpr                 string `json:"Mxpr"`                 // 상한가
	Llam                 string `json:"Llam"`                 // 하한가
	Oprc                 string `json:"Oprc"`                 // 시가
	Hprc                 string `json:"Hprc"`                 // 고가
	Lprc                 string `json:"Lprc"`                 // 저가
	PrdyVrss             string `json:"PrdyVrss"`             // 전일대비
	PrdyCtrt             string `json:"PrdyCtrt"`             // 전일대비율
	Per                  string `json:"Per"`                  // PER
	AcmlTrPbmn           string `json:"AcmlTrPbmn"`           // 거래대금
	AcmlVol              string `json:"AcmlVol"`              // 거래량
	PrdyVol              string `json:"prdyVol"`              // 전일거래량 (소문자 p)
	Bidp1                string `json:"bidp1"`                // 매수호가 (소문자)
	Askp1                string `json:"askp1"`                // 매도호가 (소문자)
	SdprVrssMrktRate     string `json:"SdprVrssMrktRate"`     // 기준가대비시가비율
	PrprVrssOprcRate     string `json:"PrprVrssOprcRate"`     // 현재가대비시가비율
	SdprVrssHgprRate     string `json:"SdprVrssHgprRate"`     // 기준가대비고가비율
	PrprVrssHgprRate     string `json:"PrprVrssHgprRate"`     // 현재가대비고가비율
	SdprVrssLwprRate     string `json:"SdprVrssLwprRate"`     // 기준가대비저가비율
	PrprVrssLwprRate     string `json:"PrprVrssLwprRate"`     // 현재가대비저가비율
}

// ForeignCurrentPriceData 해외주식 현재가 데이터 (변환된 형식)
type ForeignCurrentPriceData struct {
	StockCode        string  `json:"stock_code"`         // 종목코드
	Market           string  `json:"market"`             // 시장 (뉴욕/나스닥/아멕스)
	BasePrice        float64 `json:"base_price"`         // 기준가 (USD)
	CurrentPrice     float64 `json:"current_price"`      // 현재가 (USD)
	UpperLimit       float64 `json:"upper_limit"`        // 상한가 (USD)
	LowerLimit       float64 `json:"lower_limit"`        // 하한가 (USD)
	OpenPrice        float64 `json:"open_price"`         // 시가 (USD)
	HighPrice        float64 `json:"high_price"`         // 고가 (USD)
	LowPrice         float64 `json:"low_price"`          // 저가 (USD)
	PriceChange      float64 `json:"price_change"`       // 전일대비 (USD)
	PriceChangeRate  float64 `json:"price_change_rate"`  // 전일대비율 (%)
	PER              float64 `json:"per"`                // PER
	TradingValue     float64 `json:"trading_value"`      // 거래대금 (USD)
	TradingVolume    int64   `json:"trading_volume"`     // 거래량
	YesterdayVolume  int64   `json:"yesterday_volume"`   // 전일거래량
	BidPrice         float64 `json:"bid_price"`          // 매수호가 (USD)
	AskPrice         float64 `json:"ask_price"`          // 매도호가 (USD)
	MarketOpenRate   float64 `json:"market_open_rate"`   // 기준가대비시가비율
	CurrentOpenRate  float64 `json:"current_open_rate"`  // 현재가대비시가비율
	MarketHighRate   float64 `json:"market_high_rate"`   // 기준가대비고가비율
	CurrentHighRate  float64 `json:"current_high_rate"`  // 현재가대비고가비율
	MarketLowRate    float64 `json:"market_low_rate"`    // 기준가대비저가비율
	CurrentLowRate   float64 `json:"current_low_rate"`   // 현재가대비저가비율
	Currency         string  `json:"currency"`           // 통화 (USD)
}