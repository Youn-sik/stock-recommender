package main

import (
	"fmt"
	"os"

	"stock-recommender/backend/config"
	"stock-recommender/backend/openapi/client"
	"stock-recommender/backend/openapi/foreign"
	"stock-recommender/backend/openapi/models"
)

func main() {
	fmt.Println("=== DB Securities API Integration Test ===")
	
	// .env 파일 로드를 위해 상위 디렉토리로 이동
	if err := os.Chdir("../"); err != nil {
		fmt.Printf("Failed to change directory: %v\n", err)
		return
	}
	
	// 환경변수 확인
	appKey := os.Getenv("DBSEC_APP_KEY")
	appSecret := os.Getenv("DBSEC_APP_SECRET")
	
	if appKey == "" || appSecret == "" {
		fmt.Println("❌ Environment variables not loaded. Trying to load from file...")
		// 수동으로 .env 파일 읽기
		return
	}
	
	// 환경변수 출력 (키는 마스킹)
	if len(appKey) > 10 {
		fmt.Printf("DBSEC_APP_KEY: %s***\n", appKey[:10])
	} else {
		fmt.Printf("DBSEC_APP_KEY: %s\n", appKey)
	}
	if len(appSecret) > 10 {
		fmt.Printf("DBSEC_APP_SECRET: %s***\n", appSecret[:10])
	} else {
		fmt.Printf("DBSEC_APP_SECRET: %s\n", appSecret)
	}
	
	// 설정 로드
	cfg := config.Load()
	fmt.Printf("Config loaded - AppKey: %s***, AppSecret: %s***\n", 
		cfg.API.DBSecAppKey[:10], cfg.API.DBSecAppSecret[:10])
	
	// DB Securities 클라이언트 생성
	apiClient := client.NewDBSecClient(cfg)
	
	// 1. 인증 테스트
	fmt.Println("\n1. Testing Authentication...")
	if err := apiClient.HealthCheck(); err != nil {
		fmt.Printf("❌ Authentication failed: %v\n", err)
		return
	} else {
		fmt.Println("✅ Authentication successful!")
	}
	
	// 2. 해외 주식 현재가 조회 테스트
	fmt.Println("\n2. Testing Foreign Stock Current Price (AAPL)...")
	currentPriceService := foreign.NewForeignCurrentPriceService(apiClient)
	currentPrice, err := currentPriceService.GetNASDAQStockPrice("AAPL")
	if err != nil {
		fmt.Printf("❌ Current price query failed: %v\n", err)
	} else {
		fmt.Printf("✅ AAPL Current Price: $%.2f\n", currentPrice.CurrentPrice)
	}
	
	// 3. 해외 주식 월차트 조회 테스트
	fmt.Println("\n3. Testing Foreign Stock Month Chart (TSLA)...")
	monthChartService := foreign.NewForeignMonthChartService(apiClient)
	
	period := models.MonthChartPeriod{
		StartDate: "2024-01-01",
		EndDate:   "2024-07-23",
	}
	
	options := models.MonthChartOptions{
		UseAdjusted: true,
		Market:      "NASDAQ",
	}
	
	monthData, err := monthChartService.GetMonthChart("TSLA", period, options)
	if err != nil {
		fmt.Printf("❌ Month chart query failed: %v\n", err)
	} else {
		fmt.Printf("✅ TSLA Month Chart: Retrieved %d months of data\n", len(monthData))
		if len(monthData) > 0 {
			latest := monthData[0]
			fmt.Printf("   Latest: %s - Close: $%.2f, Volume: %d\n", 
				latest.MonthEndDate, latest.Close, latest.Volume)
		}
	}
	
	// 4. 해외 주식 주차트 조회 테스트
	fmt.Println("\n4. Testing Foreign Stock Week Chart (NVDA)...")
	weekChartService := foreign.NewForeignWeekChartService(apiClient)
	weekData, err := weekChartService.GetRecentWeekChart("NVDA", "NASDAQ", 4)
	if err != nil {
		fmt.Printf("❌ Week chart query failed: %v\n", err)
	} else {
		fmt.Printf("✅ NVDA Week Chart: Retrieved %d weeks of data\n", len(weekData))
		if len(weekData) > 0 {
			latest := weekData[0]
			fmt.Printf("   Latest: %s - Close: $%.2f\n", 
				latest.WeekEndDate, latest.Close)
		}
	}
	
	// 5. 해외 주식 일차트 조회 테스트
	fmt.Println("\n5. Testing Foreign Stock Day Chart (MSFT)...")
	dayChartService := foreign.NewForeignDayChartService(apiClient)
	dayData, err := dayChartService.GetRecentDayChart("MSFT", "NASDAQ", 5)
	if err != nil {
		fmt.Printf("❌ Day chart query failed: %v\n", err)
	} else {
		fmt.Printf("✅ MSFT Day Chart: Retrieved %d days of data\n", len(dayData))
		if len(dayData) > 0 {
			latest := dayData[0]
			fmt.Printf("   Latest: %s - Close: $%.2f\n", 
				latest.Date, latest.Close)
		}
	}
	
	// 6. 해외 주식 분차트 조회 테스트
	fmt.Println("\n6. Testing Foreign Stock Min Chart (GOOGL)...")
	minChartService := foreign.NewForeignMinChartService(apiClient)
	
	minPeriod := models.ChartPeriod{
		StartDate: "20240722",
		EndDate:   "20240723",
		IsRange:   true,
	}
	
	minOptions := models.ChartOptions{
		Interval:    "1min",
		UseAdjusted: true,
		DataCount:   50,
		Market:      "NASDAQ",
	}
	
	minData, err := minChartService.GetMinChart("GOOGL", minPeriod, minOptions)
	if err != nil {
		fmt.Printf("❌ Min chart query failed: %v\n", err)
	} else {
		fmt.Printf("✅ GOOGL Min Chart: Retrieved %d data points\n", len(minData))
		if len(minData) > 0 {
			latest := minData[0]
			fmt.Printf("   Latest: %s - Close: $%.2f\n", 
				latest.DateTime, latest.Close)
		}
	}
	
	fmt.Println("\n=== Integration Test Complete ===")
}