package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"
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

type ChartVisualizer struct {
	baseDir string
}

func NewChartVisualizer(baseDir string) *ChartVisualizer {
	return &ChartVisualizer{baseDir: baseDir}
}

func (cv *ChartVisualizer) GenerateHTML() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	htmlFile := filepath.Join(cv.baseDir, fmt.Sprintf("chart_report_%s.html", timestamp))
	
	// JSON 파일들 읽기
	files, err := filepath.Glob(filepath.Join(cv.baseDir, "*.json"))
	if err != nil {
		return err
	}
	
	chartData := make(map[string]interface{})
	
	for _, file := range files {
		if strings.Contains(file, "api_results_") {
			continue // 메인 결과 파일은 제외
		}
		
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		
		filename := filepath.Base(file)
		parts := strings.Split(strings.TrimSuffix(filename, ".json"), "_")
		if len(parts) >= 2 {
			chartType := parts[0]
			stockCode := parts[1]
			
			var chartContent interface{}
			json.Unmarshal(data, &chartContent)
			
			key := fmt.Sprintf("%s_%s", chartType, stockCode)
			chartData[key] = chartContent
		}
	}
	
	// HTML 템플릿
	htmlTemplate := `<!DOCTYPE html>
<html lang="ko">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>DB증권 API 차트 분석 리포트</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        body {
            font-family: 'Arial', sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background-color: white;
            padding: 30px;
            border-radius: 10px;
            box-shadow: 0 0 20px rgba(0,0,0,0.1);
        }
        h1 {
            color: #2c3e50;
            text-align: center;
            margin-bottom: 30px;
            border-bottom: 3px solid #3498db;
            padding-bottom: 15px;
        }
        h2 {
            color: #34495e;
            margin-top: 40px;
            margin-bottom: 20px;
        }
        .chart-container {
            position: relative;
            height: 400px;
            margin: 30px 0;
            padding: 20px;
            border: 1px solid #ddd;
            border-radius: 8px;
            background-color: #fafafa;
        }
        .stock-info {
            background-color: #ecf0f1;
            padding: 15px;
            margin: 20px 0;
            border-radius: 5px;
            border-left: 5px solid #3498db;
        }
        .data-table {
            width: 100%;
            border-collapse: collapse;
            margin: 20px 0;
        }
        .data-table th, .data-table td {
            border: 1px solid #ddd;
            padding: 12px;
            text-align: left;
        }
        .data-table th {
            background-color: #3498db;
            color: white;
        }
        .data-table tr:nth-child(even) {
            background-color: #f2f2f2;
        }
        .summary-box {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 20px;
            border-radius: 10px;
            margin: 20px 0;
        }
        .api-status {
            display: inline-block;
            padding: 5px 10px;
            border-radius: 15px;
            font-size: 12px;
            font-weight: bold;
        }
        .success { background-color: #2ecc71; color: white; }
        .failed { background-color: #e74c3c; color: white; }
        .no-data { color: #7f8c8d; font-style: italic; }
    </style>
</head>
<body>
    <div class="container">
        <h1>📈 DB증권 API 차트 분석 리포트</h1>
        <div class="summary-box">
            <h3>🔍 분석 개요</h3>
            <p><strong>생성 시간:</strong> {{.Timestamp}}</p>
            <p><strong>분석 대상:</strong> 해외 주식 차트 데이터 (월차트, 주차트, 일차트)</p>
            <p><strong>주요 종목:</strong> AAPL, MSFT, GOOGL, AMZN, TSLA, NVDA, META</p>
        </div>

        {{range $key, $data := .ChartData}}
        {{if contains $key "MonthChart"}}
        <h2>📊 {{replace $key "MonthChart_" ""}} 월차트 데이터</h2>
        <div class="chart-container">
            <canvas id="chart_{{$key}}"></canvas>
        </div>
        <div class="stock-info">
            <h4>데이터 상세 정보</h4>
            <p><strong>데이터 포인트:</strong> {{len $data}}개</p>
            {{if $data}}
            {{with index $data 0}}
            <p><strong>최신 데이터:</strong> {{.MonthEndDate}} - 종가: ${{printf "%.2f" .Close}}</p>
            <p><strong>시장:</strong> {{.Market}}</p>
            <p><strong>수정주가 적용:</strong> {{if .IsAdjusted}}예{{else}}아니오{{end}}</p>
            {{end}}
            {{end}}
        </div>
        {{end}}
        {{end}}

        {{range $key, $data := .ChartData}}
        {{if contains $key "WeekChart"}}
        <h2>📊 {{replace $key "WeekChart_" ""}} 주차트 데이터</h2>
        <div class="chart-container">
            <canvas id="chart_{{$key}}"></canvas>
        </div>
        {{end}}
        {{end}}

        {{range $key, $data := .ChartData}}
        {{if contains $key "DayChart"}}
        <h2>📊 {{replace $key "DayChart_" ""}} 일차트 데이터</h2>
        <div class="chart-container">
            <canvas id="chart_{{$key}}"></canvas>
        </div>
        {{end}}
        {{end}}

        <h2>📋 원시 데이터 테이블</h2>
        {{range $key, $data := .ChartData}}
        {{if contains $key "MonthChart"}}
        <h3>{{replace $key "MonthChart_" ""}} 월간 데이터</h3>
        <table class="data-table">
            <thead>
                <tr>
                    <th>날짜</th>
                    <th>시가</th>
                    <th>고가</th>
                    <th>저가</th>
                    <th>종가</th>
                    <th>거래량</th>
                    <th>월간 변동률</th>
                </tr>
            </thead>
            <tbody>
                {{range $data}}
                <tr>
                    <td>{{.MonthEndDate}}</td>
                    <td>${{printf "%.2f" .Open}}</td>
                    <td>${{printf "%.2f" .High}}</td>
                    <td>${{printf "%.2f" .Low}}</td>
                    <td>${{printf "%.2f" .Close}}</td>
                    <td>{{.Volume}}</td>
                    <td>{{printf "%.2f" .ChangeRate}}%</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{end}}
        {{end}}

    </div>

<script>
// 차트 생성 스크립트
{{range $key, $data := .ChartData}}
{{if $data}}
// {{$key}} 차트
(function() {
    const ctx = document.getElementById('chart_{{$key}}');
    if (!ctx) return;
    
    const data = {{marshal $data}};
    if (!data || data.length === 0) return;
    
    let labels = [];
    let prices = [];
    let volumes = [];
    
    {{if contains $key "MonthChart"}}
    data.forEach(item => {
        labels.push(item.MonthEndDate || item.month_end_date);
        prices.push(item.Close || item.close);
        volumes.push(item.Volume || item.volume);
    });
    {{else if contains $key "WeekChart"}}
    data.forEach(item => {
        labels.push(item.WeekEndDate || item.week_end_date);
        prices.push(item.Close || item.close);
        volumes.push(item.Volume || item.volume);
    });
    {{else if contains $key "DayChart"}}
    data.forEach(item => {
        labels.push(item.Date || item.date);
        prices.push(item.Close || item.close);
        volumes.push(item.Volume || item.volume);
    });
    {{end}}
    
    new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels.reverse(),
            datasets: [{
                label: '종가 ($)',
                data: prices.reverse(),
                borderColor: 'rgb(75, 192, 192)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                tension: 0.1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                title: {
                    display: true,
                    text: '{{replace $key "_" " "}} 가격 추이'
                }
            },
            scales: {
                y: {
                    beginAtZero: false,
                    title: {
                        display: true,
                        text: '가격 ($)'
                    }
                },
                x: {
                    title: {
                        display: true,
                        text: '날짜'
                    }
                }
            }
        }
    });
})();
{{end}}
{{end}}
</script>

</body>
</html>`

	// 템플릿 함수 정의
	funcMap := template.FuncMap{
		"contains": strings.Contains,
		"replace":  strings.ReplaceAll,
		"marshal": func(v interface{}) template.JS {
			data, _ := json.Marshal(v)
			return template.JS(data)
		},
	}

	tmpl, err := template.New("chart").Funcs(funcMap).Parse(htmlTemplate)
	if err != nil {
		return err
	}

	// HTML 파일 생성
	file, err := os.Create(htmlFile)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		Timestamp string
		ChartData map[string]interface{}
	}{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		ChartData: chartData,
	}

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	fmt.Printf("📊 차트 HTML 리포트 생성: %s\n", htmlFile)
	return nil
}

func (cv *ChartVisualizer) GenerateMarkdownReport() error {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	mdFile := filepath.Join(cv.baseDir, fmt.Sprintf("analysis_report_%s.md", timestamp))
	
	// API 결과 파일 찾기
	resultFiles, err := filepath.Glob(filepath.Join(cv.baseDir, "api_results_*.json"))
	if err != nil || len(resultFiles) == 0 {
		return fmt.Errorf("API 결과 파일을 찾을 수 없습니다")
	}
	
	// 최신 결과 파일 읽기
	latestFile := resultFiles[len(resultFiles)-1]
	data, err := os.ReadFile(latestFile)
	if err != nil {
		return err
	}
	
	var results []APICallResult
	if err := json.Unmarshal(data, &results); err != nil {
		return err
	}
	
	// 마크다운 리포트 생성
	md := strings.Builder{}
	md.WriteString("# 📈 DB증권 API 분석 리포트\n\n")
	md.WriteString(fmt.Sprintf("**생성 시간:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	
	// 통계 계산
	totalCalls := len(results)
	successCalls := 0
	failedCalls := 0
	apiStats := make(map[string]int)
	
	for _, result := range results {
		if result.Success {
			successCalls++
		} else {
			failedCalls++
		}
		apiStats[result.API]++
	}
	
	md.WriteString("## 📊 전체 통계\n\n")
	md.WriteString(fmt.Sprintf("- **총 호출 수:** %d\n", totalCalls))
	md.WriteString(fmt.Sprintf("- **성공 호출:** %d (%.1f%%)\n", successCalls, float64(successCalls)/float64(totalCalls)*100))
	md.WriteString(fmt.Sprintf("- **실패 호출:** %d (%.1f%%)\n\n", failedCalls, float64(failedCalls)/float64(totalCalls)*100))
	
	md.WriteString("## 🔍 API별 호출 현황\n\n")
	md.WriteString("| API | 호출 수 | 상태 |\n")
	md.WriteString("|-----|---------|------|\n")
	
	for api, count := range apiStats {
		status := "✅"
		for _, result := range results {
			if result.API == api && !result.Success {
				status = "❌"
				break
			}
		}
		md.WriteString(fmt.Sprintf("| %s | %d | %s |\n", api, count, status))
	}
	
	md.WriteString("\n## 📋 상세 호출 로그\n\n")
	md.WriteString("| 시간 | API | 종목 | 성공 | 데이터 수 | 응답시간 | 에러 |\n")
	md.WriteString("|------|-----|------|------|-----------|----------|-------|\n")
	
	for _, result := range results {
		status := "✅"
		if !result.Success {
			status = "❌"
		}
		
		errorMsg := result.Error
		if len(errorMsg) > 50 {
			errorMsg = errorMsg[:50] + "..."
		}
		
		md.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %d | %s | %s |\n",
			result.Timestamp, result.API, result.StockCode, status, 
			result.DataCount, result.ResponseTime, errorMsg))
	}
	
	// 한도 분석
	hasLimitError := false
	for _, result := range results {
		if strings.Contains(result.Error, "호출 거래건수를 초과") {
			hasLimitError = true
			break
		}
	}
	
	md.WriteString("\n## 💡 분석 결과\n\n")
	if hasLimitError {
		md.WriteString("⚠️ **API 호출 한도 감지**\n")
		md.WriteString("- DB증권 API는 일일 호출 제한이 있는 것으로 확인됩니다.\n")
		md.WriteString("- 'IGW00201: 호출 거래건수를 초과하였습니다' 에러 발생\n")
		md.WriteString("- API 사용량 관리가 필요합니다.\n\n")
	}
	
	if successCalls > 0 {
		md.WriteString("✅ **성공적인 데이터 수집**\n")
		md.WriteString(fmt.Sprintf("- %d개의 성공적인 API 호출 확인\n", successCalls))
		md.WriteString("- 실시간 주식 데이터 수집 기능 정상 작동\n\n")
	}
	
	// 파일 저장
	if err := os.WriteFile(mdFile, []byte(md.String()), 0644); err != nil {
		return err
	}
	
	fmt.Printf("📝 마크다운 리포트 생성: %s\n", mdFile)
	return nil
}