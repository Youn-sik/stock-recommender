package models

import "stock-recommender/backend/openapi/utils"

// StockTickerRequest 주식종목 조회 요청
type StockTickerRequest struct {
	In StockTickerInput `json:"In"`
}

// StockTickerInput 주식종목 조회 입력
type StockTickerInput struct {
	InputCondMrktDivCode string `json:"InputCondMrktDivCode"` // 시장분류코드 (J: 주식, E: ETF, EN: ETN)
}

// StockTickerResponse 주식종목 조회 응답
type StockTickerResponse struct {
	utils.BaseAPIResponse
	Out []StockTickerOutput `json:"Out"`
}

// StockTickerOutput 주식종목 조회 출력
type StockTickerOutput struct {
	Iscd        string `json:"Iscd"`        // 종목코드 (9자리)
	StndIscd    string `json:"StndIscd"`    // 표준종목코드 (12자리)
	KorIsnm     string `json:"KorIsnm"`     // 한글종목명 (40자)
	MrktClsCode string `json:"MrktClsCode"` // 시장분류구분코드 (1: 코스닥, 4: 코스피)
}

// StockTickerHeader 주식종목 조회 헤더
type StockTickerHeader struct {
	ContentType string `json:"content-type"`
	Authorization string `json:"authorization"` 
	ContYn      string `json:"cont_yn"`       // 연속거래 여부 (Y/N)
	ContKey     string `json:"cont_key"`      // 연속키 값 (최대 70자)
	MacAddress  string `json:"mac_address"`   // MAC 주소 (법인용, 12자)
}