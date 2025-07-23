package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"stock-recommender/backend/config"
	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/foreign"
	"stock-recommender/backend/openapi/models"
)

type APICallResult struct {
	Timestamp    string      `json:"timestamp"`
	API          string      `json:"api"`
	StockCode    string      `json:"stock_code"`
	Success      bool        `json:"success"`
	DataCount    int         `json:"data_count"`
	ResponseTime string      `json:"response_time"`
	Error        string      `json:"error,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type APIAnalyzer struct {
	client  *client.DBSecClient
	results []APICallResult
	baseDir string
}

func NewAPIAnalyzer(client *client.DBSecClient) *APIAnalyzer {
	// 결과 저장 디렉토리 생성 (프로젝트 루트의 results 디렉토리)
	baseDir := "../results"
	os.MkdirAll(baseDir, 0755)
	
	return &APIAnalyzer{
		client:  client,
		results: make([]APICallResult, 0),
		baseDir: baseDir,
	}
}

func (a *APIAnalyzer) recordCall(api, stockCode string, success bool, dataCount int, responseTime time.Duration, err error, data interface{}) {
	result := APICallResult{
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
		API:          api,
		StockCode:    stockCode,
		Success:      success,
		DataCount:    dataCount,
		ResponseTime: responseTime.String(),
		Data:         data,
	}
	
	if err != nil {
		result.Error = err.Error()
	}
	
	a.results = append(a.results, result)
}

func (a *APIAnalyzer) TestCurrentPrices() {
	fmt.Println("📊 Testing Current Prices...")
	service := foreign.NewForeignCurrentPriceService(a.client)
	
	stocks := []string{"AAPL", "MSFT", "GOOGL", "AMZN", "TSLA", "NVDA", "META"}
	
	for _, stock := range stocks {
		start := time.Now()
		data, err := service.GetNASDAQStockPrice(stock)
		duration := time.Since(start)
		
		if err != nil {
			a.recordCall("CurrentPrice", stock, false, 0, duration, err, nil)
			fmt.Printf("   ❌ %s: %v\n", stock, err)
		} else {
			a.recordCall("CurrentPrice", stock, true, 1, duration, nil, data)
			fmt.Printf("   ✅ %s: $%.2f\n", stock, data.CurrentPrice)
		}
		
		// Rate limiting 고려
		time.Sleep(100 * time.Millisecond)
	}
}

func (a *APIAnalyzer) TestMonthCharts() {
	fmt.Println("📊 Testing Month Charts...")
	service := foreign.NewForeignMonthChartService(a.client)
	
	stocks := []string{"AAPL", "TSLA", "NVDA"}
	
	period := models.MonthChartPeriod{
		StartDate: "2024-01-01",
		EndDate:   "2024-07-23",
	}
	
	options := models.MonthChartOptions{
		UseAdjusted: true,
		Market:      "NASDAQ",
	}
	
	for _, stock := range stocks {
		start := time.Now()
		data, err := service.GetMonthChart(stock, period, options)
		duration := time.Since(start)
		
		if err != nil {
			a.recordCall("MonthChart", stock, false, 0, duration, err, nil)
			fmt.Printf("   ❌ %s: %v\n", stock, err)
		} else {
			a.recordCall("MonthChart", stock, true, len(data), duration, nil, data)
			fmt.Printf("   ✅ %s: %d months\n", stock, len(data))
		}
		
		time.Sleep(100 * time.Millisecond)
	}
}

func (a *APIAnalyzer) TestWeekCharts() {
	fmt.Println("📊 Testing Week Charts...")
	service := foreign.NewForeignWeekChartService(a.client)
	
	stocks := []string{"AAPL", "MSFT", "GOOGL"}
	
	for _, stock := range stocks {
		start := time.Now()
		data, err := service.GetRecentWeekChart(stock, "NASDAQ", 8)
		duration := time.Since(start)
		
		if err != nil {
			a.recordCall("WeekChart", stock, false, 0, duration, err, nil)
			fmt.Printf("   ❌ %s: %v\n", stock, err)
		} else {
			a.recordCall("WeekChart", stock, true, len(data), duration, nil, data)
			fmt.Printf("   ✅ %s: %d weeks\n", stock, len(data))
		}
		
		time.Sleep(100 * time.Millisecond)
	}
}

func (a *APIAnalyzer) TestDayCharts() {
	fmt.Println("📊 Testing Day Charts...")
	service := foreign.NewForeignDayChartService(a.client)
	
	stocks := []string{"AAPL", "TSLA"}
	
	for _, stock := range stocks {
		start := time.Now()
		data, err := service.GetRecentDayChart(stock, "NASDAQ", 10)
		duration := time.Since(start)
		
		if err != nil {
			a.recordCall("DayChart", stock, false, 0, duration, err, nil)
			fmt.Printf("   ❌ %s: %v\n", stock, err)
		} else {
			a.recordCall("DayChart", stock, true, len(data), duration, nil, data)
			fmt.Printf("   ✅ %s: %d days\n", stock, len(data))
		}
		
		time.Sleep(100 * time.Millisecond)
	}
}

func (a *APIAnalyzer) SaveResults() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	
	// JSON 결과 저장
	jsonFile := filepath.Join(a.baseDir, fmt.Sprintf("api_results_%s.json", timestamp))
	jsonData, err := json.MarshalIndent(a.results, "", "  ")
	if err != nil {
		return err
	}
	
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return err
	}
	
	// CSV 결과 저장
	csvFile := filepath.Join(a.baseDir, fmt.Sprintf("api_results_%s.csv", timestamp))
	file, err := os.Create(csvFile)
	if err != nil {
		return err
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// CSV 헤더
	writer.Write([]string{"Timestamp", "API", "StockCode", "Success", "DataCount", "ResponseTime", "Error"})
	
	// CSV 데이터
	for _, result := range a.results {
		writer.Write([]string{
			result.Timestamp,
			result.API,
			result.StockCode,
			fmt.Sprintf("%t", result.Success),
			fmt.Sprintf("%d", result.DataCount),
			result.ResponseTime,
			result.Error,
		})
	}
	
	fmt.Printf("📁 Results saved to:\n")
	fmt.Printf("   JSON: %s\n", jsonFile)
	fmt.Printf("   CSV:  %s\n", csvFile)
	
	return nil
}

func (a *APIAnalyzer) SaveDetailedChartData() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	
	// 각 API별로 성공한 데이터만 별도 저장
	for _, result := range a.results {
		if !result.Success || result.Data == nil {
			continue
		}
		
		filename := fmt.Sprintf("%s_%s_%s.json", result.API, result.StockCode, timestamp)
		filepath := filepath.Join(a.baseDir, filename)
		
		data, err := json.MarshalIndent(result.Data, "", "  ")
		if err != nil {
			continue
		}
		
		os.WriteFile(filepath, data, 0644)
	}
	
	return nil
}

func (a *APIAnalyzer) GenerateReport() {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📈 DB증권 API 호출 결과 분석 리포트")
	fmt.Println(strings.Repeat("=", 60))
	
	totalCalls := len(a.results)
	successCalls := 0
	failedCalls := 0
	totalDataPoints := 0
	
	apiStats := make(map[string]int)
	errorStats := make(map[string]int)
	
	for _, result := range a.results {
		if result.Success {
			successCalls++
			totalDataPoints += result.DataCount
		} else {
			failedCalls++
			errorStats[result.Error]++
		}
		apiStats[result.API]++
	}
	
	fmt.Printf("📊 전체 호출 통계:\n")
	fmt.Printf("   총 호출 수: %d\n", totalCalls)
	fmt.Printf("   성공: %d (%.1f%%)\n", successCalls, float64(successCalls)/float64(totalCalls)*100)
	fmt.Printf("   실패: %d (%.1f%%)\n", failedCalls, float64(failedCalls)/float64(totalCalls)*100)
	fmt.Printf("   총 데이터 포인트: %d\n", totalDataPoints)
	
	fmt.Printf("\n📋 API별 호출 현황:\n")
	for api, count := range apiStats {
		fmt.Printf("   %s: %d회\n", api, count)
	}
	
	if len(errorStats) > 0 {
		fmt.Printf("\n❌ 에러 통계:\n")
		for errMsg, count := range errorStats {
			// 에러 메시지를 짧게 요약
			if strings.Contains(errMsg, "호출 거래건수를 초과") {
				fmt.Printf("   API 호출 한도 초과: %d회\n", count)
			} else if strings.Contains(errMsg, "authentication failed") {
				fmt.Printf("   인증 실패: %d회\n", count)
			} else {
				fmt.Printf("   기타 에러: %d회\n", count)
			}
		}
	}
	
	fmt.Printf("\n💡 분석 결과:\n")
	if failedCalls > 0 && strings.Contains(fmt.Sprintf("%v", errorStats), "호출 거래건수를 초과") {
		fmt.Printf("   ⚠️  API 일일 호출 한도에 도달했습니다.\n")
		fmt.Printf("   ⚠️  DB증권 API는 일일 호출 제한이 있는 것으로 보입니다.\n")
	}
	
	if successCalls > 0 {
		fmt.Printf("   ✅ 성공한 API 호출에서 총 %d개의 데이터 포인트를 수집했습니다.\n", totalDataPoints)
	}
}

func main() {
	fmt.Println("🔍 DB증권 API 분석 도구 시작")
	fmt.Println("Current working directory:", func() string { wd, _ := os.Getwd(); return wd }())
	
	// 환경변수 설정
	os.Setenv("DBSEC_APP_KEY", "PSxUUPVxVizXuOpUaL6P9Dk0mHGK2a8TNqS6")
	os.Setenv("DBSEC_APP_SECRET", "2UCcBWcVoHWvx1eAuBLusNdEPCoAOedw")
	
	// 설정 로드
	cfg := config.Load()
	apiClient := client.NewDBSecClient(cfg)
	
	// 분석기 생성
	analyzer := NewAPIAnalyzer(apiClient)
	
	// 인증 테스트
	fmt.Println("🔐 인증 테스트...")
	if err := apiClient.HealthCheck(); err != nil {
		fmt.Printf("❌ 인증 실패: %v\n", err)
		analyzer.recordCall("Authentication", "", false, 0, 0, err, nil)
	} else {
		fmt.Println("✅ 인증 성공!")
		analyzer.recordCall("Authentication", "", true, 1, 0, nil, nil)
	}
	
	// 각종 API 테스트 (호출 한도까지 테스트)
	analyzer.TestCurrentPrices()
	analyzer.TestMonthCharts()
	analyzer.TestWeekCharts()
	analyzer.TestDayCharts()
	
	// 결과 저장
	if err := analyzer.SaveResults(); err != nil {
		fmt.Printf("❌ 결과 저장 실패: %v\n", err)
	}
	
	if err := analyzer.SaveDetailedChartData(); err != nil {
		fmt.Printf("❌ 상세 데이터 저장 실패: %v\n", err)
	}
	
	// 분석 리포트 생성
	analyzer.GenerateReport()
	
	fmt.Println("\n🎉 분석 완료!")
}