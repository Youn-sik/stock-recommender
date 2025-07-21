package models

import (
	"fmt"
	
	"stock-recommender/backend/openapi/utils"
)

// ForeignMinChartRequest 해외주식 분차트조회 요청
type ForeignMinChartRequest struct {
	In ForeignMinChartInput `json:"In"`
}

// ForeignMinChartInput 해외주식 분차트조회 입력
type ForeignMinChartInput struct {
	InputCondMrktDivCode string `json:"InputCondMrktDivCode"` // 시장분류코드 (FY:뉴욕, FN:나스닥, FA:아멕스)
	InputIscd1           string `json:"InputIscd1"`           // 해외주식종목코드
	InputDate1           string `json:"InputDate1"`           // 시작날짜 (YYYYMMDD)
	InputDate2           string `json:"InputDate2"`           // 종료날짜 (YYYYMMDD)
	InputHourClsCode     string `json:"InputHourClsCode"`     // 시간구분코드 (고정값: "0")
	InputDivXtick        string `json:"InputDivXtick"`        // 분일별구분코드 (30:30초, 60:1분, 600:10분, 3600:60분)
	InputPwDataIncuYn    string `json:"InputPwDataIncuYn"`    // 기간지정여부코드 (Y:기간지정, N:기간미지정)
	InputOrgAdjPrc       string `json:"InputOrgAdjPrc"`       // 수정주가사용여부 (0:미사용, 1:사용)
	DataCnt              string `json:"dataCnt"`              // 호출건수 (1~2000, 공백시 기본400개)
}

// ForeignMinChartResponse 해외주식 분차트조회 응답
type ForeignMinChartResponse struct {
	utils.BaseAPIResponse
	Out []ForeignMinChartOutput `json:"Out"`
}

// ForeignMinChartOutput 해외주식 분차트조회 출력
type ForeignMinChartOutput struct {
	Hour    string `json:"Hour"`    // 시간 (HHMMSS)
	Date    string `json:"Date"`    // 일자 (YYYYMMDD)
	Prpr    string `json:"Prpr"`    // 현재가
	Oprc    string `json:"Oprc"`    // 시가
	Hprc    string `json:"Hprc"`    // 고가
	Lprc    string `json:"Lprc"`    // 저가
	CntgVol string `json:"CntgVol"` // 체결거래량
}

// ForeignMinChartData 해외주식 분차트 비즈니스 모델
type ForeignMinChartData struct {
	StockCode     string  `json:"stock_code"`     // 종목코드
	DateTime      string  `json:"date_time"`      // 일시 (YYYY-MM-DD HH:MM:SS)
	Date          string  `json:"date"`           // 일자 (YYYY-MM-DD)
	Time          string  `json:"time"`           // 시간 (HH:MM:SS)
	Open          float64 `json:"open"`           // 시가
	High          float64 `json:"high"`           // 고가
	Low           float64 `json:"low"`            // 저가
	Close         float64 `json:"close"`          // 종가(현재가)
	Volume        int64   `json:"volume"`         // 거래량
	Market        string  `json:"market"`         // 시장명
	MarketCode    string  `json:"market_code"`    // 시장코드
	Interval      string  `json:"interval"`       // 시간간격
	IntervalCode  string  `json:"interval_code"`  // 시간간격코드
	IsAdjusted    bool    `json:"is_adjusted"`    // 수정주가 적용여부
}

// ChartPeriod 차트 조회 기간 설정
type ChartPeriod struct {
	StartDate string `json:"start_date"` // 시작일 (YYYY-MM-DD)
	EndDate   string `json:"end_date"`   // 종료일 (YYYY-MM-DD)
	IsRange   bool   `json:"is_range"`   // 기간 지정 여부
}

// ChartOptions 차트 조회 옵션
type ChartOptions struct {
	Interval     string `json:"interval"`      // 시간간격 (30sec, 1min, 5min, 10min, 60min)
	UseAdjusted  bool   `json:"use_adjusted"`  // 수정주가 사용여부
	DataCount    int    `json:"data_count"`    // 조회 건수 (1~2000)
	Market       string `json:"market"`        // 시장 (NY, NASDAQ, AMEX)
}

// GetIntervalCode 시간간격 문자열을 코드로 변환
func (opts *ChartOptions) GetIntervalCode() string {
	switch opts.Interval {
	case "30sec":
		return ChartInterval30Sec
	case "1min":
		return ChartInterval1Min
	case "2min":
		return ChartInterval2Min
	case "5min":
		return ChartInterval5Min
	case "10min":
		return ChartInterval10Min
	case "60min":
		return ChartInterval60Min
	default:
		return ChartInterval1Min // 기본값: 1분
	}
}

// GetMarketCode 시장명을 코드로 변환
func (opts *ChartOptions) GetMarketCode() string {
	switch opts.Market {
	case "NY", "NYSE":
		return ForeignMarketNY
	case "NASDAQ":
		return ForeignMarketNASDAQ
	case "AMEX":
		return ForeignMarketAMEX
	default:
		return ForeignMarketNASDAQ // 기본값: 나스닥
	}
}

// GetAdjustedCode 수정주가 사용여부를 코드로 변환
func (opts *ChartOptions) GetAdjustedCode() string {
	if opts.UseAdjusted {
		return AdjustedPriceEnabled
	}
	return AdjustedPriceDisabled
}

// GetDataCountString 조회건수를 문자열로 변환
func (opts *ChartOptions) GetDataCountString() string {
	if opts.DataCount <= 0 || opts.DataCount > 2000 {
		return "" // 기본값 사용 (400개)
	}
	return fmt.Sprintf("%d", opts.DataCount)
}

// =============================================================================
// 해외주식 일차트 조회 관련 모델
// =============================================================================

// ForeignDayChartRequest 해외주식 일차트조회 요청
type ForeignDayChartRequest struct {
	In ForeignDayChartInput `json:"In"`
}

// ForeignDayChartInput 해외주식 일차트조회 입력
type ForeignDayChartInput struct {
	InputCondMrktDivCode string `json:"InputCondMrktDivCode"` // 시장분류코드 (FY:뉴욕, FN:나스닥, FA:아멕스)
	InputOrgAdjPrc       string `json:"InputOrgAdjPrc"`       // 수정주가사용여부 (0:미사용, 1:사용)
	InputIscd1           string `json:"InputIscd1"`           // 해외주식종목코드
	InputDate1           string `json:"InputDate1"`           // 시작날짜 (YYYYMMDD)
	InputDate2           string `json:"InputDate2"`           // 종료날짜 (YYYYMMDD)
}

// ForeignDayChartResponse 해외주식 일차트조회 응답
type ForeignDayChartResponse struct {
	utils.BaseAPIResponse
	Out []ForeignDayChartOutput `json:"Out"`
}

// ForeignDayChartOutput 해외주식 일차트조회 출력
type ForeignDayChartOutput struct {
	Hour    string `json:"Hour"`    // 시간 (일차트에서는 빈 값)
	Date    string `json:"Date"`    // 일자 (YYYYMMDD)
	Prpr    string `json:"Prpr"`    // 현재가(종가)
	Oprc    string `json:"Oprc"`    // 시가
	Hprc    string `json:"Hprc"`    // 고가
	Lprc    string `json:"Lprc"`    // 저가
	AcmlVol string `json:"AcmlVol"` // 누적거래량 (일차트에서는 AcmlVol 사용)
}

// ForeignDayChartData 해외주식 일차트 비즈니스 모델
type ForeignDayChartData struct {
	StockCode    string  `json:"stock_code"`    // 종목코드
	Date         string  `json:"date"`          // 일자 (YYYY-MM-DD)
	Open         float64 `json:"open"`          // 시가
	High         float64 `json:"high"`          // 고가
	Low          float64 `json:"low"`           // 저가
	Close        float64 `json:"close"`         // 종가
	Volume       int64   `json:"volume"`        // 거래량
	Market       string  `json:"market"`        // 시장명
	MarketCode   string  `json:"market_code"`   // 시장코드
	IsAdjusted   bool    `json:"is_adjusted"`   // 수정주가 적용여부
	WeekDay      string  `json:"week_day"`      // 요일
	PriceChange  float64 `json:"price_change"`  // 전일대비 가격 변화 (계산된 값)
	ChangeRate   float64 `json:"change_rate"`   // 전일대비 변화율 (계산된 값)
}

// DayChartPeriod 일차트 조회 기간 설정
type DayChartPeriod struct {
	StartDate string `json:"start_date"` // 시작일 (YYYY-MM-DD 또는 YYYYMMDD)
	EndDate   string `json:"end_date"`   // 종료일 (YYYY-MM-DD 또는 YYYYMMDD)
}

// DayChartOptions 일차트 조회 옵션
type DayChartOptions struct {
	UseAdjusted bool   `json:"use_adjusted"` // 수정주가 사용여부
	Market      string `json:"market"`       // 시장 (NY, NASDAQ, AMEX)
}

// GetMarketCode 시장명을 코드로 변환
func (opts *DayChartOptions) GetMarketCode() string {
	switch opts.Market {
	case "NY", "NYSE":
		return ForeignMarketNY
	case "NASDAQ":
		return ForeignMarketNASDAQ
	case "AMEX":
		return ForeignMarketAMEX
	default:
		return ForeignMarketNASDAQ // 기본값: 나스닥
	}
}

// GetAdjustedCode 수정주가 사용여부를 코드로 변환
func (opts *DayChartOptions) GetAdjustedCode() string {
	if opts.UseAdjusted {
		return AdjustedPriceEnabled
	}
	return AdjustedPriceDisabled
}

// FormatDate 날짜를 YYYYMMDD 형식으로 변환
func (p *DayChartPeriod) FormatDate(date string) string {
	// YYYY-MM-DD 형식을 YYYYMMDD로 변환
	if len(date) == 10 && date[4] == '-' && date[7] == '-' {
		return date[:4] + date[5:7] + date[8:10]
	}
	// 이미 YYYYMMDD 형식이면 그대로 반환
	if len(date) == 8 {
		return date
	}
	return ""
}

// GetFormattedStartDate 포맷된 시작일 반환
func (p *DayChartPeriod) GetFormattedStartDate() string {
	return p.FormatDate(p.StartDate)
}

// GetFormattedEndDate 포맷된 종료일 반환
func (p *DayChartPeriod) GetFormattedEndDate() string {
	return p.FormatDate(p.EndDate)
}

// =============================================================================
// 해외주식 주차트 조회 관련 모델
// =============================================================================

// ForeignWeekChartRequest 해외주식 주차트조회 요청
type ForeignWeekChartRequest struct {
	In ForeignWeekChartInput `json:"In"`
}

// ForeignWeekChartInput 해외주식 주차트조회 입력
type ForeignWeekChartInput struct {
	InputCondMrktDivCode string `json:"InputCondMrktDivCode"` // 시장분류코드 (FY:뉴욕, FN:나스닥, FA:아멕스)
	InputOrgAdjPrc       string `json:"InputOrgAdjPrc"`       // 수정주가사용여부 (0:미사용, 1:사용)
	InputIscd1           string `json:"InputIscd1"`           // 해외주식종목코드
	InputDate1           string `json:"InputDate1"`           // 시작날짜 (YYYYMMDD)
	InputDate2           string `json:"InputDate2"`           // 종료날짜 (YYYYMMDD)
	InputPeriodDivCode   string `json:"InputPeriodDivCode"`   // 기간구분코드 (W:주간)
}

// ForeignWeekChartResponse 해외주식 주차트조회 응답
type ForeignWeekChartResponse struct {
	utils.BaseAPIResponse
	Out []ForeignWeekChartOutput `json:"Out"`
}

// ForeignWeekChartOutput 해외주식 주차트조회 출력
type ForeignWeekChartOutput struct {
	Hour    string `json:"Hour"`    // 시간 (주차트에서는 빈 값)
	Date    string `json:"Date"`    // 일자 (YYYYMMDD) - 주의 마지막 날짜
	Prpr    string `json:"Prpr"`    // 현재가(주간종가)
	Oprc    string `json:"Oprc"`    // 시가
	Hprc    string `json:"Hprc"`    // 고가
	Lprc    string `json:"Lprc"`    // 저가
	CntgVol string `json:"CntgVol"` // 체결거래량 (주차트에서는 빈 값인 경우가 많음)
}

// ForeignWeekChartData 해외주식 주차트 비즈니스 모델
type ForeignWeekChartData struct {
	StockCode       string  `json:"stock_code"`       // 종목코드
	WeekEndDate     string  `json:"week_end_date"`    // 주 종료일 (YYYY-MM-DD)
	WeekStartDate   string  `json:"week_start_date"`  // 주 시작일 (YYYY-MM-DD) - 계산된 값
	Open            float64 `json:"open"`             // 시가
	High            float64 `json:"high"`             // 고가
	Low             float64 `json:"low"`              // 저가
	Close           float64 `json:"close"`            // 종가
	Volume          int64   `json:"volume"`           // 거래량 (0일 수 있음)
	Market          string  `json:"market"`           // 시장명
	MarketCode      string  `json:"market_code"`      // 시장코드
	IsAdjusted      bool    `json:"is_adjusted"`      // 수정주가 적용여부
	WeekNumber      int     `json:"week_number"`      // 연도 내 주차 번호 (계산된 값)
	Year            int     `json:"year"`             // 연도
	PriceChange     float64 `json:"price_change"`     // 전주대비 가격 변화 (계산된 값)
	ChangeRate      float64 `json:"change_rate"`      // 전주대비 변화율 (계산된 값)
	WeeklyRange     float64 `json:"weekly_range"`     // 주간 변동폭 (고가-저가)
	WeeklyRangeRate float64 `json:"weekly_range_rate"` // 주간 변동률 ((고가-저가)/저가*100)
}

// WeekChartPeriod 주차트 조회 기간 설정
type WeekChartPeriod struct {
	StartDate string `json:"start_date"` // 시작일 (YYYY-MM-DD 또는 YYYYMMDD)
	EndDate   string `json:"end_date"`   // 종료일 (YYYY-MM-DD 또는 YYYYMMDD)
}

// WeekChartOptions 주차트 조회 옵션
type WeekChartOptions struct {
	UseAdjusted bool   `json:"use_adjusted"` // 수정주가 사용여부
	Market      string `json:"market"`       // 시장 (NY, NASDAQ, AMEX)
}

// GetMarketCode 시장명을 코드로 변환
func (opts *WeekChartOptions) GetMarketCode() string {
	switch opts.Market {
	case "NY", "NYSE":
		return ForeignMarketNY
	case "NASDAQ":
		return ForeignMarketNASDAQ
	case "AMEX":
		return ForeignMarketAMEX
	default:
		return ForeignMarketNASDAQ // 기본값: 나스닥
	}
}

// GetAdjustedCode 수정주가 사용여부를 코드로 변환
func (opts *WeekChartOptions) GetAdjustedCode() string {
	if opts.UseAdjusted {
		return AdjustedPriceEnabled
	}
	return AdjustedPriceDisabled
}

// FormatDate 날짜를 YYYYMMDD 형식으로 변환
func (p *WeekChartPeriod) FormatDate(date string) string {
	// YYYY-MM-DD 형식을 YYYYMMDD로 변환
	if len(date) == 10 && date[4] == '-' && date[7] == '-' {
		return date[:4] + date[5:7] + date[8:10]
	}
	// 이미 YYYYMMDD 형식이면 그대로 반환
	if len(date) == 8 {
		return date
	}
	return ""
}

// GetFormattedStartDate 포맷된 시작일 반환
func (p *WeekChartPeriod) GetFormattedStartDate() string {
	return p.FormatDate(p.StartDate)
}

// GetFormattedEndDate 포맷된 종료일 반환
func (p *WeekChartPeriod) GetFormattedEndDate() string {
	return p.FormatDate(p.EndDate)
}

// =============================================================================
// 해외주식 월차트 조회 관련 모델
// =============================================================================

// ForeignMonthChartRequest 해외주식 월차트조회 요청
type ForeignMonthChartRequest struct {
	In ForeignMonthChartInput `json:"In"`
}

// ForeignMonthChartInput 해외주식 월차트조회 입력
type ForeignMonthChartInput struct {
	InputOrgAdjPrc       string `json:"InputOrgAdjPrc"`       // 수정주가사용여부 (0:미사용, 1:사용)
	InputCondMrktDivCode string `json:"InputCondMrktDivCode"` // 시장분류코드 (FY:뉴욕, FN:나스닥, FA:아멕스)
	InputIscd1           string `json:"InputIscd1"`           // 해외주식종목코드
	InputDate1           string `json:"InputDate1"`           // 시작날짜 (YYYYMMDD)
	InputDate2           string `json:"InputDate2"`           // 종료날짜 (YYYYMMDD)
}

// ForeignMonthChartResponse 해외주식 월차트조회 응답
type ForeignMonthChartResponse struct {
	utils.BaseAPIResponse
	Out []ForeignMonthChartOutput `json:"Out"`
}

// ForeignMonthChartOutput 해외주식 월차트조회 출력
type ForeignMonthChartOutput struct {
	Hour    string `json:"Hour"`    // 시간 (월차트에서는 빈 값)
	Date    string `json:"Date"`    // 일자 (YYYYMMDD) - 월의 마지막 날짜
	Prpr    string `json:"Prpr"`    // 현재가(월간종가)
	Oprc    string `json:"Oprc"`    // 시가
	Hprc    string `json:"Hprc"`    // 고가
	Lprc    string `json:"Lprc"`    // 저가
	AcmlVol string `json:"AcmlVol"` // 누적체결거래량
}

// ForeignMonthChartData 해외주식 월차트 비즈니스 모델
type ForeignMonthChartData struct {
	StockCode         string  `json:"stock_code"`         // 종목코드
	MonthEndDate      string  `json:"month_end_date"`     // 월 종료일 (YYYY-MM-DD)
	MonthStartDate    string  `json:"month_start_date"`   // 월 시작일 (YYYY-MM-DD) - 계산된 값
	Open              float64 `json:"open"`               // 시가
	High              float64 `json:"high"`               // 고가
	Low               float64 `json:"low"`                // 저가
	Close             float64 `json:"close"`              // 종가
	Volume            int64   `json:"volume"`             // 거래량
	Market            string  `json:"market"`             // 시장명
	MarketCode        string  `json:"market_code"`        // 시장코드
	IsAdjusted        bool    `json:"is_adjusted"`        // 수정주가 적용여부
	Year              int     `json:"year"`               // 연도
	Month             int     `json:"month"`              // 월
	PriceChange       float64 `json:"price_change"`       // 전월대비 가격 변화 (계산된 값)
	ChangeRate        float64 `json:"change_rate"`        // 전월대비 변화율 (계산된 값)
	MonthlyRange      float64 `json:"monthly_range"`      // 월간 변동폭 (고가-저가)
	MonthlyRangeRate  float64 `json:"monthly_range_rate"` // 월간 변동률 ((고가-저가)/저가*100)
}

// MonthChartPeriod 월차트 조회 기간 설정
type MonthChartPeriod struct {
	StartDate string `json:"start_date"` // 시작일 (YYYY-MM-DD 또는 YYYYMMDD)
	EndDate   string `json:"end_date"`   // 종료일 (YYYY-MM-DD 또는 YYYYMMDD)
}

// MonthChartOptions 월차트 조회 옵션
type MonthChartOptions struct {
	UseAdjusted bool   `json:"use_adjusted"` // 수정주가 사용여부
	Market      string `json:"market"`       // 시장 (NY, NASDAQ, AMEX)
}

// GetMarketCode 시장명을 코드로 변환
func (opts *MonthChartOptions) GetMarketCode() string {
	switch opts.Market {
	case "NY", "NYSE":
		return ForeignMarketNY
	case "NASDAQ":
		return ForeignMarketNASDAQ
	case "AMEX":
		return ForeignMarketAMEX
	default:
		return ForeignMarketNASDAQ // 기본값: 나스닥
	}
}

// GetAdjustedCode 수정주가 사용여부를 코드로 변환
func (opts *MonthChartOptions) GetAdjustedCode() string {
	if opts.UseAdjusted {
		return AdjustedPriceEnabled
	}
	return AdjustedPriceDisabled
}

// FormatDate 날짜를 YYYYMMDD 형식으로 변환
func (p *MonthChartPeriod) FormatDate(date string) string {
	// YYYY-MM-DD 형식을 YYYYMMDD로 변환
	if len(date) == 10 && date[4] == '-' && date[7] == '-' {
		return date[:4] + date[5:7] + date[8:10]
	}
	// 이미 YYYYMMDD 형식이면 그대로 반환
	if len(date) == 8 {
		return date
	}
	return ""
}

// GetFormattedStartDate 포맷된 시작일 반환
func (p *MonthChartPeriod) GetFormattedStartDate() string {
	return p.FormatDate(p.StartDate)
}

// GetFormattedEndDate 포맷된 종료일 반환
func (p *MonthChartPeriod) GetFormattedEndDate() string {
	return p.FormatDate(p.EndDate)
}