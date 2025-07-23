package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("📊 리포트 생성 도구 시작")
	
	// 차트 시각화 도구 생성
	visualizer := NewChartVisualizer("../results")
	
	// HTML 리포트 생성
	fmt.Println("🌐 HTML 차트 리포트 생성 중...")
	if err := visualizer.GenerateHTML(); err != nil {
		log.Printf("HTML 생성 실패: %v", err)
	}
	
	// 마크다운 리포트 생성
	fmt.Println("📝 마크다운 분석 리포트 생성 중...")
	if err := visualizer.GenerateMarkdownReport(); err != nil {
		log.Printf("마크다운 생성 실패: %v", err)
	}
	
	fmt.Println("✅ 모든 리포트 생성 완료!")
}