package models

import (
	"time"
)

// 공통 응답 구조
type BaseResponse struct {
	ResultCode    string `json:"rt_cd"`      // 결과코드
	ResultMessage string `json:"msg_cd"`     // 메시지코드  
	ResultData    string `json:"msg1"`       // 결과메시지
}

// 국내주식 현재가 응답
type DomesticStockPrice struct {
	BaseResponse
	Output StockPriceOutput `json:"output"`
}

type StockPriceOutput struct {
	StockCode        string `json:"mksc_shrn_iscd"`    // 유가증권 단축 종목코드
	StockNameKor     string `json:"hts_kor_isnm"`      // HTS 한글 종목명
	StockPrice       string `json:"stck_prpr"`         // 주식 현재가
	PrevDayDiff      string `json:"prdy_vrss"`         // 전일 대비
	PrevDayDiffSign  string `json:"prdy_vrss_sign"`    // 전일 대비 부호
	PrevDayDiffRate  string `json:"prdy_ctrt"`         // 전일 대비율
	AccTradeVolume   string `json:"acml_vol"`          // 누적 거래량
	AccTradePrice    string `json:"acml_tr_pbmn"`      // 누적 거래대금
	HtsAvgPrice      string `json:"hts_avls"`          // HTS 시가총액
	PER              string `json:"per"`               // PER
	PBR              string `json:"pbr"`               // PBR
	OpenPrice        string `json:"stck_oprc"`         // 주식 시가
	HighPrice        string `json:"stck_hgpr"`         // 주식 최고가
	LowPrice         string `json:"stck_lwpr"`         // 주식 최저가
	PrevClosePrice   string `json:"stck_sdpr"`         // 주식 기준가(전일종가)
	Volume           string `json:"hts_deal_qty_unit_val"` // HTS 거래량 단위값
}

// 국내주식 호가 정보
type DomesticStockAskingPrice struct {
	BaseResponse
	Output1 []AskingPriceOutput1 `json:"output1"`
	Output2 []AskingPriceOutput2 `json:"output2"`
}

type AskingPriceOutput1 struct {
	AskPrice1   string `json:"askp1"`      // 매도호가1
	AskPrice2   string `json:"askp2"`      // 매도호가2
	AskPrice3   string `json:"askp3"`      // 매도호가3
	AskPrice4   string `json:"askp4"`      // 매도호가4
	AskPrice5   string `json:"askp5"`      // 매도호가5
	BidPrice1   string `json:"bidp1"`      // 매수호가1
	BidPrice2   string `json:"bidp2"`      // 매수호가2
	BidPrice3   string `json:"bidp3"`      // 매수호가3
	BidPrice4   string `json:"bidp4"`      // 매수호가4
	BidPrice5   string `json:"bidp5"`      // 매수호가5
	AskVolume1  string `json:"askp_rsqn1"` // 매도호가 잔량1
	AskVolume2  string `json:"askp_rsqn2"` // 매도호가 잔량2
	AskVolume3  string `json:"askp_rsqn3"` // 매도호가 잔량3
	AskVolume4  string `json:"askp_rsqn4"` // 매도호가 잔량4
	AskVolume5  string `json:"askp_rsqn5"` // 매도호가 잔량5
	BidVolume1  string `json:"bidp_rsqn1"` // 매수호가 잔량1
	BidVolume2  string `json:"bidp_rsqn2"` // 매수호가 잔량2
	BidVolume3  string `json:"bidp_rsqn3"` // 매수호가 잔량3
	BidVolume4  string `json:"bidp_rsqn4"` // 매수호가 잔량4
	BidVolume5  string `json:"bidp_rsqn5"` // 매수호가 잔량5
}

type AskingPriceOutput2 struct {
	TotalAskVolume string `json:"total_askp_rsqn"` // 총 매도호가 잔량
	TotalBidVolume string `json:"total_bidp_rsqn"` // 총 매수호가 잔량
	OverAskVolume  string `json:"ovtm_total_askp_rsqn"` // 시간외 총 매도호가 잔량
	OverBidVolume  string `json:"ovtm_total_bidp_rsqn"` // 시간외 총 매수호가 잔량
}

// 국내주식 일봉차트 조회
type DomesticStockDaily struct {
	BaseResponse
	Output1 StockDailyOutput1   `json:"output1"`
	Output2 []StockDailyOutput2 `json:"output2"`
}

type StockDailyOutput1 struct {
	PrdtTypeCd   string `json:"prdt_type_cd"`   // 상품유형코드
	NextTrCont   string `json:"next_tr_cont"`   // 연속조회검색조건
	NextTrKey    string `json:"next_tr_key"`    // 연속조회키
}

type StockDailyOutput2 struct {
	StockDate      string `json:"stck_bsop_date"` // 주식 영업 일자
	StockClosePrice string `json:"stck_clpr"`      // 주식 종가
	StockOpenPrice  string `json:"stck_oprc"`      // 주식 시가
	StockHighPrice  string `json:"stck_hgpr"`      // 주식 최고가
	StockLowPrice   string `json:"stck_lwpr"`      // 주식 최저가
	AccTradeVolume  string `json:"acml_vol"`       // 누적 거래량
	AccTradePrice   string `json:"acml_tr_pbmn"`   // 누적 거래대금
	FluctuationRate string `json:"flng_cls_code"`  // 등락구분코드
	PrevDayDiff     string `json:"prdy_vrss_vol_rate"` // 전일 대비 거래량 비율
	PersonalNetBuy  string `json:"pers_ntby_qty"`      // 개인 순매수 수량
	ForeignNetBuy   string `json:"frgn_ntby_qty"`      // 외국인 순매수 수량
	OrganNetBuy     string `json:"orgn_ntby_qty"`      // 기관계 순매수 수량
}

// 해외주식 현재가 조회
type ForeignStockPrice struct {
	BaseResponse
	Output ForeignStockPriceOutput `json:"output"`
}

type ForeignStockPriceOutput struct {
	SymbolCode      string `json:"symb"`           // 심볼
	SecurityName    string `json:"rsym"`           // 실시간 심볼
	ZdivCode        string `json:"zdiv"`           // 소수점 자리수
	BasePrice       string `json:"base"`           // 전일종가
	PrevPrice       string `json:"pvol"`           // 전일거래량
	LastPrice       string `json:"last"`           // 현재가
	DiffSign        string `json:"sign"`           // 대비구분
	Change          string `json:"diff"`           // 대비
	Rate            string `json:"rate"`           // 등락율
	Volume          string `json:"tvol"`           // 거래량
	TradePrice      string `json:"tamt"`           // 거래대금
	OpenPrice       string `json:"open"`           // 시가
	HighPrice       string `json:"high"`           // 고가
	LowPrice        string `json:"low"`            // 저가
}

// 지수 현재가 조회
type IndexPrice struct {
	BaseResponse
	Output IndexPriceOutput `json:"output"`
}

type IndexPriceOutput struct {
	IndexName       string `json:"hts_kor_isnm"`  // HTS 한글 종목명
	IndexValue      string `json:"bstp_nmix_prpr"` // 업종지수 현재가
	IndexChange     string `json:"bstp_nmix_prdy_vrss"` // 업종지수 전일대비
	IndexChangeSign string `json:"prdy_vrss_sign"` // 전일 대비 부호
	IndexChangeRate string `json:"bstp_nmix_prdy_ctrt"` // 업종지수 전일대비율
	AccTradeVolume  string `json:"acml_vol"`       // 누적 거래량
	AccTradePrice   string `json:"acml_tr_pbmn"`   // 누적 거래대금
	OpenValue       string `json:"bstp_nmix_oprc"` // 업종지수 시가
	HighValue       string `json:"bstp_nmix_hgpr"` // 업종지수 최고가
	LowValue        string `json:"bstp_nmix_lwpr"` // 업종지수 최저가
}

// KOSPI/KOSDAQ 종목 리스트
type StockList struct {
	BaseResponse
	Output []StockListOutput `json:"output"`
}

type StockListOutput struct {
	StockCode       string `json:"mksc_shrn_iscd"`     // 유가증권 단축 종목코드
	StandardCode    string `json:"stnd_iscd"`          // 표준 종목코드
	StockNameKor    string `json:"hts_kor_isnm"`       // HTS 한글 종목명
	StockNameEng    string `json:"lst_stck_vl_100"`    // 상장 주식 값 100
	MarketWarning   string `json:"mktw_cls_code"`      // 시장경고코드
	GroupCode       string `json:"grp_code"`           // 그룹코드
	StockRank       string `json:"stck_kind"`          // 주식 종류
	MfDate          string `json:"mf_rate_cls_code"`   // 제조업 구분 코드
	BusinessType    string `json:"bztp_cls_code"`      // 업종 구분 코드
	ListedShares    string `json:"lstg_stqt"`          // 상장 주식수
}

// API 요청용 구조체들
type StockPriceRequest struct {
	StockCode string `json:"fid_cond_mrkt_div_code"` // 시장분류코드
	Symbol    string `json:"fid_input_iscd"`         // 종목코드
}

type StockDailyRequest struct {
	StockCode     string `json:"fid_cond_mrkt_div_code"` // 시장분류코드
	Symbol        string `json:"fid_input_iscd"`         // 종목코드
	StartDate     string `json:"fid_input_date_1"`       // 조회시작일자
	EndDate       string `json:"fid_input_date_2"`       // 조회종료일자
	PeriodDivCode string `json:"fid_period_div_code"`    // 기간분류코드
	OrgAdjPrice   string `json:"fid_org_adj_prc"`        // 수정주가 원주가 가격
}

type ForeignStockRequest struct {
	ExchangeCode string `json:"AUTH"`       // 거래소코드
	Symbol       string `json:"SYMB"`       // 심볼
}

// 내부 변환용 구조체
type ParsedStockPrice struct {
	Symbol         string
	Name           string
	CurrentPrice   float64
	OpenPrice      float64
	HighPrice      float64
	LowPrice       float64
	PrevClosePrice float64
	Change         float64
	ChangeRate     float64
	Volume         int64
	TradeAmount    int64
	Timestamp      time.Time
	Market         string
}

type ParsedAskingPrice struct {
	Symbol      string
	AskPrices   [5]float64  // 매도호가 1~5
	BidPrices   [5]float64  // 매수호가 1~5
	AskVolumes  [5]int64    // 매도 잔량 1~5
	BidVolumes  [5]int64    // 매수 잔량 1~5
	TotalAskVol int64       // 총 매도 잔량
	TotalBidVol int64       // 총 매수 잔량
	Timestamp   time.Time
}

type ParsedDailyData struct {
	Symbol      string
	Date        time.Time
	OpenPrice   float64
	HighPrice   float64
	LowPrice    float64
	ClosePrice  float64
	Volume      int64
	TradeAmount int64
}

// 거래소 코드 상수
const (
	MarketKOSPI   = "J"  // KOSPI
	MarketKOSDAQ  = "Q"  // KOSDAQ
	MarketKONEX   = "K"  // KONEX
	
	ExchangeNASDAQ = "NAS" // NASDAQ
	ExchangeNYSE   = "NYS" // NYSE
	ExchangeAMEX   = "AMS" // AMEX
)

// 주요 지수 코드
const (
	IndexKOSPI     = "0001"  // KOSPI
	IndexKOSDAQ    = "1001"  // KOSDAQ
	IndexKOSPI200  = "1028"  // KOSPI200
	IndexDow       = "DJI"   // 다우존스
	IndexNASDAQ    = "IXIC"  // 나스닥
	IndexSP500     = "SPX"   // S&P500
)

// API 경로 상수
const (
	// 국내주식
	PathDomesticStockPrice    = "/uapi/domestic-stock/v1/quotations/inquire-price"
	PathDomesticStockAsking   = "/uapi/domestic-stock/v1/quotations/inquire-asking-price-exp-ccn"
	PathDomesticStockDaily    = "/uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice"
	PathDomesticStockList     = "/uapi/domestic-stock/v1/quotations/search-stock-info"
	
	// 해외주식
	PathForeignStockPrice     = "/uapi/overseas-price/v1/quotations/price"
	PathForeignStockDaily     = "/uapi/overseas-price/v1/quotations/dailyprice"
	
	// 지수
	PathIndexPrice            = "/uapi/domestic-stock/v1/quotations/inquire-index-price"
)