package models

// ForeignStockTickerRequest 해외주식종목 조회 요청
type ForeignStockTickerRequest struct {
	In ForeignStockTickerInput `json:"In"`
}

// ForeignStockTickerInput 해외주식종목 조회 입력
type ForeignStockTickerInput struct {
	InputDataCode string `json:"InputDataCode"` // 해외증시구분코드 (NY: 뉴욕, NA: 나스닥, AM: 아멕스)
}

// ForeignStockTickerResponse 해외주식종목 조회 응답
type ForeignStockTickerResponse struct {
	Out    []ForeignStockTickerOutput `json:"Out"`
	RspCd  string                     `json:"rsp_cd"`  // 응답코드
	RspMsg string                     `json:"rsp_msg"` // 응답메시지
}

// ForeignStockTickerOutput 해외주식종목 조회 출력
type ForeignStockTickerOutput struct {
	Iscd          string `json:"Iscd"`          // 종목코드 (9자리)
	KorIsnm       string `json:"KorIsnm"`       // 한글종목명 (40자)
	BstpLargName  string `json:"BstpLargName"`  // 업종대분류명 (40자)
	ExchClsCode2  string `json:"ExchClsCode2"`  // 거래소코드2 (4자)
	SelnVolUnit   string `json:"SelnVolUnit"`   // 매도량단위 (9자)
	ShnuVolUnit   string `json:"ShnuVolUnit"`   // 매수량단위 (9자)
}

// ForeignStockTickerHeader 해외주식종목 조회 헤더
type ForeignStockTickerHeader struct {
	ContentType   string `json:"content-type"`
	Authorization string `json:"authorization"`
	ContYn        string `json:"cont_yn"`      // 연속거래 여부 (Y/N)
	ContKey       string `json:"cont_key"`     // 연속키 값 (최대 70자)
	MacAddress    string `json:"mac_address"`  // MAC 주소 (법인용, 12자)
}

// ForeignStockData 해외주식 데이터 (변환된 형식)
type ForeignStockData struct {
	StockCode     string `json:"stock_code"`     // 종목코드
	KoreanName    string `json:"korean_name"`    // 한글종목명
	SectorName    string `json:"sector_name"`    // 업종대분류명
	ExchangeCode  string `json:"exchange_code"`  // 거래소코드
	Exchange      string `json:"exchange"`       // 거래소명 (NY/NASDAQ/AMEX)
	SellUnit      int64  `json:"sell_unit"`      // 매도량단위
	BuyUnit       int64  `json:"buy_unit"`       // 매수량단위
}