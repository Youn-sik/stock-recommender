package models

import "stock-recommender/backend/openapi/utils"

// CurrentPriceRequest 현재가조회 요청
type CurrentPriceRequest struct {
	In CurrentPriceInput `json:"In"`
}

// CurrentPriceInput 현재가조회 입력
type CurrentPriceInput struct {
	InputCondMrktDivCode string `json:"InputCondMrktDivCode"` // 시장분류코드 (J: 주식, E: ETF, EN: ETN, W: ELW, U: 업종&지수)
	InputIscd1           string `json:"InputIscd1"`           // 종목코드 (6자리) 또는 지수코드
}

// CurrentPriceResponse 현재가조회 응답
type CurrentPriceResponse struct {
	utils.BaseAPIResponse
	Out CurrentPriceOutput `json:"Out"`
}

// CurrentPriceOutput 현재가조회 출력
type CurrentPriceOutput struct {
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
	Pbr                  string `json:"Pbr"`                  // PBR
	AcmlTrPbmn           string `json:"AcmlTrPbmn"`           // 거래대금
	AcmlVol              string `json:"AcmlVol"`              // 거래량
	PrdyVol              string `json:"PrdyVol"`              // 전일거래량
	Bidp1                string `json:"Bidp1"`                // 매수호가
	Askp1                string `json:"Askp1"`                // 매도호가
	SdprVrssMrktRate     string `json:"SdprVrssMrktRate"`     // 기준가대비시가비율
	PrprVrssOprcRate     string `json:"PrprVrssOprcRate"`     // 현재가대비시가비율
	SdprVrssHgprRate     string `json:"SdprVrssHgprRate"`     // 기준가대비고가비율
	PrprVrssHgprRate     string `json:"PrprVrssHgprRate"`     // 현재가대비고가비율
	SdprVrssLwprRate     string `json:"SdprVrssLwprRate"`     // 기준가대비저가비율
	PrprVrssLwprRate     string `json:"PrprVrssLwprRate"`     // 현재가대비저가비율
	HtsOtstStplQty       string `json:"HtsOtstStplQty"`       // 미결제약정수량
	OtstStplQtyIcdc      string `json:"OtstStplQtyIcdc"`      // 미결제증감
}

// CurrentPriceData 현재가 데이터 (변환된 형식)
type CurrentPriceData struct {
	StockCode        string  `json:"stock_code"`         // 종목코드
	BasePrice        float64 `json:"base_price"`         // 기준가
	CurrentPrice     float64 `json:"current_price"`      // 현재가
	UpperLimit       float64 `json:"upper_limit"`        // 상한가
	LowerLimit       float64 `json:"lower_limit"`        // 하한가
	OpenPrice        float64 `json:"open_price"`         // 시가
	HighPrice        float64 `json:"high_price"`         // 고가
	LowPrice         float64 `json:"low_price"`          // 저가
	PriceChange      float64 `json:"price_change"`       // 전일대비
	PriceChangeRate  float64 `json:"price_change_rate"`  // 전일대비율 (%)
	PER              float64 `json:"per"`                // PER
	PBR              float64 `json:"pbr"`                // PBR
	TradingValue     float64 `json:"trading_value"`      // 거래대금
	TradingVolume    int64   `json:"trading_volume"`     // 거래량
	YesterdayVolume  int64   `json:"yesterday_volume"`   // 전일거래량
	BidPrice         float64 `json:"bid_price"`          // 매수호가
	AskPrice         float64 `json:"ask_price"`          // 매도호가
	MarketOpenRate   float64 `json:"market_open_rate"`   // 기준가대비시가비율
	CurrentOpenRate  float64 `json:"current_open_rate"`  // 현재가대비시가비율
	MarketHighRate   float64 `json:"market_high_rate"`   // 기준가대비고가비율
	CurrentHighRate  float64 `json:"current_high_rate"`  // 현재가대비고가비율
	MarketLowRate    float64 `json:"market_low_rate"`    // 기준가대비저가비율
	CurrentLowRate   float64 `json:"current_low_rate"`   // 현재가대비저가비율
}