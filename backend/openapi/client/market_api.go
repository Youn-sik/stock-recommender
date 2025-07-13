package client

import (
	"encoding/json"
	"fmt"
	"stock-recommender/backend/openapi/models"
	"time"
)

// 국내주식 현재가 조회
func (c *DBSecClient) GetDomesticStockPrice(symbol string) (*models.ParsedStockPrice, error) {
	params := map[string]string{
		"fid_cond_mrkt_div_code": "J", // KOSPI
		"fid_input_iscd":         symbol,
	}

	respBody, err := c.makeRequest("GET", models.PathDomesticStockPrice, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get domestic stock price: %w", err)
	}

	var response models.DomesticStockPrice
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse stock price response: %w", err)
	}

	// API 에러 체크
	if response.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s - %s", response.ResultCode, response.ResultMessage)
	}

	// 응답 데이터를 내부 구조체로 변환
	return &models.ParsedStockPrice{
		Symbol:         symbol,
		Name:           response.Output.StockNameKor,
		CurrentPrice:   c.parseFloat(response.Output.StockPrice),
		OpenPrice:      c.parseFloat(response.Output.OpenPrice),
		HighPrice:      c.parseFloat(response.Output.HighPrice),
		LowPrice:       c.parseFloat(response.Output.LowPrice),
		PrevClosePrice: c.parseFloat(response.Output.PrevClosePrice),
		Change:         c.parseFloat(response.Output.PrevDayDiff),
		ChangeRate:     c.parseFloat(response.Output.PrevDayDiffRate),
		Volume:         c.parseInt(response.Output.AccTradeVolume),
		TradeAmount:    c.parseInt(response.Output.AccTradePrice),
		Timestamp:      time.Now(),
		Market:         "KR",
	}, nil
}

// 국내주식 호가 정보 조회
func (c *DBSecClient) GetDomesticStockAskingPrice(symbol string) (*models.ParsedAskingPrice, error) {
	params := map[string]string{
		"fid_cond_mrkt_div_code": "J",
		"fid_input_iscd":         symbol,
	}

	respBody, err := c.makeRequest("GET", models.PathDomesticStockAsking, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get asking price: %w", err)
	}

	var response models.DomesticStockAskingPrice
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse asking price response: %w", err)
	}

	if response.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s - %s", response.ResultCode, response.ResultMessage)
	}

	if len(response.Output1) == 0 || len(response.Output2) == 0 {
		return nil, fmt.Errorf("invalid asking price data")
	}

	output1 := response.Output1[0]
	output2 := response.Output2[0]

	return &models.ParsedAskingPrice{
		Symbol: symbol,
		AskPrices: [5]float64{
			c.parseFloat(output1.AskPrice1),
			c.parseFloat(output1.AskPrice2),
			c.parseFloat(output1.AskPrice3),
			c.parseFloat(output1.AskPrice4),
			c.parseFloat(output1.AskPrice5),
		},
		BidPrices: [5]float64{
			c.parseFloat(output1.BidPrice1),
			c.parseFloat(output1.BidPrice2),
			c.parseFloat(output1.BidPrice3),
			c.parseFloat(output1.BidPrice4),
			c.parseFloat(output1.BidPrice5),
		},
		AskVolumes: [5]int64{
			c.parseInt(output1.AskVolume1),
			c.parseInt(output1.AskVolume2),
			c.parseInt(output1.AskVolume3),
			c.parseInt(output1.AskVolume4),
			c.parseInt(output1.AskVolume5),
		},
		BidVolumes: [5]int64{
			c.parseInt(output1.BidVolume1),
			c.parseInt(output1.BidVolume2),
			c.parseInt(output1.BidVolume3),
			c.parseInt(output1.BidVolume4),
			c.parseInt(output1.BidVolume5),
		},
		TotalAskVol: c.parseInt(output2.TotalAskVolume),
		TotalBidVol: c.parseInt(output2.TotalBidVolume),
		Timestamp:   time.Now(),
	}, nil
}

// 국내주식 일봉차트 조회
func (c *DBSecClient) GetDomesticStockDaily(symbol string, startDate, endDate string) ([]models.ParsedDailyData, error) {
	params := map[string]string{
		"fid_cond_mrkt_div_code": "J",
		"fid_input_iscd":         symbol,
		"fid_input_date_1":       startDate, // YYYYMMDD
		"fid_input_date_2":       endDate,   // YYYYMMDD
		"fid_period_div_code":    "D",       // 일봉
		"fid_org_adj_prc":        "0",       // 수정주가
	}

	respBody, err := c.makeRequest("GET", models.PathDomesticStockDaily, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get daily data: %w", err)
	}

	var response models.DomesticStockDaily
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse daily data response: %w", err)
	}

	if response.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s - %s", response.ResultCode, response.ResultMessage)
	}

	var result []models.ParsedDailyData
	for _, item := range response.Output2 {
		result = append(result, models.ParsedDailyData{
			Symbol:      symbol,
			Date:        c.parseDate(item.StockDate),
			OpenPrice:   c.parseFloat(item.StockOpenPrice),
			HighPrice:   c.parseFloat(item.StockHighPrice),
			LowPrice:    c.parseFloat(item.StockLowPrice),
			ClosePrice:  c.parseFloat(item.StockClosePrice),
			Volume:      c.parseInt(item.AccTradeVolume),
			TradeAmount: c.parseInt(item.AccTradePrice),
		})
	}

	return result, nil
}

// 해외주식 현재가 조회
func (c *DBSecClient) GetForeignStockPrice(symbol, exchange string) (*models.ParsedStockPrice, error) {
	params := map[string]string{
		"AUTH": exchange, // NAS, NYS, AMS
		"SYMB": symbol,
	}

	respBody, err := c.makeRequest("GET", models.PathForeignStockPrice, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get foreign stock price: %w", err)
	}

	var response models.ForeignStockPrice
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse foreign stock price response: %w", err)
	}

	if response.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s - %s", response.ResultCode, response.ResultMessage)
	}

	return &models.ParsedStockPrice{
		Symbol:         symbol,
		Name:           response.Output.SecurityName,
		CurrentPrice:   c.parseFloat(response.Output.LastPrice),
		OpenPrice:      c.parseFloat(response.Output.OpenPrice),
		HighPrice:      c.parseFloat(response.Output.HighPrice),
		LowPrice:       c.parseFloat(response.Output.LowPrice),
		PrevClosePrice: c.parseFloat(response.Output.BasePrice),
		Change:         c.parseFloat(response.Output.Change),
		ChangeRate:     c.parseFloat(response.Output.Rate),
		Volume:         c.parseInt(response.Output.Volume),
		TradeAmount:    c.parseInt(response.Output.TradePrice),
		Timestamp:      time.Now(),
		Market:         "US",
	}, nil
}

// 지수 현재가 조회
func (c *DBSecClient) GetIndexPrice(indexCode string) (*models.ParsedStockPrice, error) {
	params := map[string]string{
		"fid_cond_mrkt_div_code": "U", // 업종지수
		"fid_input_iscd":         indexCode,
	}

	respBody, err := c.makeRequest("GET", models.PathIndexPrice, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get index price: %w", err)
	}

	var response models.IndexPrice
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse index price response: %w", err)
	}

	if response.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s - %s", response.ResultCode, response.ResultMessage)
	}

	return &models.ParsedStockPrice{
		Symbol:         indexCode,
		Name:           response.Output.IndexName,
		CurrentPrice:   c.parseFloat(response.Output.IndexValue),
		OpenPrice:      c.parseFloat(response.Output.OpenValue),
		HighPrice:      c.parseFloat(response.Output.HighValue),
		LowPrice:       c.parseFloat(response.Output.LowValue),
		PrevClosePrice: 0, // 지수는 전일종가 정보가 없음
		Change:         c.parseFloat(response.Output.IndexChange),
		ChangeRate:     c.parseFloat(response.Output.IndexChangeRate),
		Volume:         c.parseInt(response.Output.AccTradeVolume),
		TradeAmount:    c.parseInt(response.Output.AccTradePrice),
		Timestamp:      time.Now(),
		Market:         "INDEX",
	}, nil
}

// 종목 리스트 조회 (KOSPI/KOSDAQ)
func (c *DBSecClient) GetStockList(marketType string) ([]models.StockListOutput, error) {
	params := map[string]string{
		"SCNT":                 "0",      // 검색개수
		"PRDT_TYPE_CD":         "300",    // 상품유형코드
		"MKT_ID_CD":            marketType, // J:KOSPI, Q:KOSDAQ
		"STCK_SHRN_ISCD":       "",       // 주식단축종목코드 (전체조회시 공백)
		"STCK_SCHD_YN":         "Y",      // 주식예약여부
		"STCK_LSTG_STQT_BZTP":  "",       // 주식상장주식수업종구분
		"STCK_PRDT_GRP_CD":     "",       // 주식상품그룹코드
	}

	respBody, err := c.makeRequest("GET", models.PathDomesticStockList, params, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get stock list: %w", err)
	}

	var response models.StockList
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse stock list response: %w", err)
	}

	if response.ResultCode != "0" {
		return nil, fmt.Errorf("API error: %s - %s", response.ResultCode, response.ResultMessage)
	}

	return response.Output, nil
}

// 종목별 최신 데이터 수집 (통합)
func (c *DBSecClient) CollectStockData(symbol, market string) (*models.ParsedStockPrice, *models.ParsedAskingPrice, error) {
	var price *models.ParsedStockPrice
	var asking *models.ParsedAskingPrice
	var err error

	// 시장에 따라 다른 API 호출
	switch market {
	case "KR":
		price, err = c.GetDomesticStockPrice(symbol)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get domestic price: %w", err)
		}
		
		asking, err = c.GetDomesticStockAskingPrice(symbol)
		if err != nil {
			return price, nil, fmt.Errorf("failed to get asking price: %w", err)
		}

	case "US":
		// 미국 주식은 거래소별로 구분 (기본: NASDAQ)
		exchange := models.ExchangeNASDAQ
		if len(symbol) <= 3 {
			exchange = models.ExchangeNYSE // 짧은 심볼은 보통 NYSE
		}
		
		price, err = c.GetForeignStockPrice(symbol, exchange)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get foreign price: %w", err)
		}
		// 해외주식은 호가 정보 없음

	case "INDEX":
		price, err = c.GetIndexPrice(symbol)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get index price: %w", err)
		}
		// 지수는 호가 정보 없음

	default:
		return nil, nil, fmt.Errorf("unsupported market: %s", market)
	}

	return price, asking, nil
}

// 시장별 주요 종목 코드 가져오기
func (c *DBSecClient) GetMajorStocks() map[string][]string {
	return map[string][]string{
		"KR": {
			"005930", // 삼성전자
			"000660", // SK하이닉스
			"051910", // LG화학
			"035420", // NAVER
			"006400", // 삼성SDI
			"035720", // 카카오
			"012330", // 현대모비스
			"028260", // 삼성물산
			"066570", // LG전자
			"096770", // SK이노베이션
		},
		"US": {
			"AAPL",  // Apple
			"MSFT",  // Microsoft
			"GOOGL", // Alphabet
			"AMZN",  // Amazon
			"TSLA",  // Tesla
			"META",  // Meta
			"NVDA",  // NVIDIA
			"AMD",   // AMD
			"INTC",  // Intel
			"ORCL",  // Oracle
		},
		"INDEX": {
			models.IndexKOSPI,
			models.IndexKOSDAQ,
			models.IndexKOSPI200,
		},
	}
}

// API 상태 및 제한 확인
func (c *DBSecClient) GetAPIStatus() map[string]interface{} {
	status := map[string]interface{}{
		"authenticated": c.accessToken != "",
		"base_url":      c.baseURL,
		"rate_limit":    "20 requests/second",
		"timestamp":     time.Now(),
	}

	// 간단한 API 호출로 연결 상태 확인
	_, err := c.GetDomesticStockPrice("005930")
	status["api_available"] = err == nil
	if err != nil {
		status["last_error"] = err.Error()
	}

	return status
}