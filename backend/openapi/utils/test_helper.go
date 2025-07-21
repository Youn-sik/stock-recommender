package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"stock-recommender/backend/config"
)

// MockServer 테스트용 모의 서버 구조체
type MockServer struct {
	server *httptest.Server
	t      *testing.T
}

// NewMockServer 새로운 모의 서버 생성
func NewMockServer(t *testing.T, handler http.HandlerFunc) *MockServer {
	server := httptest.NewServer(handler)
	return &MockServer{
		server: server,
		t:      t,
	}
}

// Close 모의 서버 종료
func (m *MockServer) Close() {
	m.server.Close()
}

// URL 모의 서버 URL 반환
func (m *MockServer) URL() string {
	return m.server.URL
}

// ClientInterface 클라이언트 인터페이스 (순환 import 방지)
type ClientInterface interface {
	HasValidCredentials() bool
}

// CreateTestConfig 테스트용 설정 생성
func CreateTestConfig() *config.Config {
	return &config.Config{
		API: config.APIConfig{
			DBSecAppKey:    "test-key",
			DBSecAppSecret: "test-secret",
		},
	}
}

// SkipIfNoCredentials API 자격증명이 없으면 테스트 스킵
func SkipIfNoCredentials(t *testing.T, client ClientInterface) {
	if !client.HasValidCredentials() {
		t.Skip("API credentials not configured")
	}
}

// MockAPIResponse 모의 API 응답 생성 헬퍼
type MockAPIResponse struct {
	ResponseCode string
	ResponseMsg  string
	Data         interface{}
}

// CreateMockResponse 모의 응답 생성
func CreateMockResponse(mockResp MockAPIResponse) []byte {
	response := map[string]interface{}{
		"rsp_cd":  mockResp.ResponseCode,
		"rsp_msg": mockResp.ResponseMsg,
	}
	
	if mockResp.Data != nil {
		response["Out"] = mockResp.Data
	}
	
	respBody, _ := json.Marshal(response)
	return respBody
}

// CreateStockTickerMockHandler 종목 조회용 모의 핸들러 생성
func CreateStockTickerMockHandler(t *testing.T, expectedPath string, expectedMarketDiv string, mockData interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 경로 확인
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 메소드 확인
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 요청 본문 파싱
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// 요청 데이터 검증
		if in, ok := req["In"].(map[string]interface{}); ok {
			if marketDiv, ok := in["InputCondMrktDivCode"].(string); ok {
				if marketDiv != expectedMarketDiv {
					t.Errorf("Expected market div %s, got %s", expectedMarketDiv, marketDiv)
				}
			}
		}

		// 응답 생성
		response := map[string]interface{}{
			"rsp_cd":  "00000",
			"rsp_msg": "정상 처리 되었습니다.",
			"Out":     mockData,
		}

		// 응답 헤더 설정
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("cont_yn", "N")
		w.Header().Set("cont_key", "")

		// 응답 작성
		json.NewEncoder(w).Encode(response)
	}
}

// CreateCurrentPriceMockHandler 현재가 조회용 모의 핸들러 생성
func CreateCurrentPriceMockHandler(t *testing.T, expectedPath string, expectedStockCode string, mockData interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 경로 확인
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 메소드 확인
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 요청 본문 파싱
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// 요청 데이터 검증
		if in, ok := req["In"].(map[string]interface{}); ok {
			if stockCode, ok := in["InputIscd1"].(string); ok {
				if stockCode != expectedStockCode {
					t.Errorf("Expected stock code %s, got %s", expectedStockCode, stockCode)
				}
			}
		}

		// 응답 생성
		response := map[string]interface{}{
			"rsp_cd":  "00000",
			"rsp_msg": "정상 처리 되었습니다.",
			"Out":     mockData,
		}

		// 응답 헤더 설정
		w.Header().Set("Content-Type", "application/json")

		// 응답 작성
		json.NewEncoder(w).Encode(response)
	}
}

// AssertFloatEqual float64 값 비교 헬퍼
func AssertFloatEqual(t *testing.T, expected, actual float64, message string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %.2f, got %.2f", message, expected, actual)
	}
}

// AssertIntEqual int64 값 비교 헬퍼
func AssertIntEqual(t *testing.T, expected, actual int64, message string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %d, got %d", message, expected, actual)
	}
}

// AssertStringEqual 문자열 값 비교 헬퍼
func AssertStringEqual(t *testing.T, expected, actual, message string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %s, got %s", message, expected, actual)
	}
}

// CreateForeignMinChartMockHandler 해외주식 분차트 조회용 모의 핸들러 생성
func CreateForeignMinChartMockHandler(t *testing.T, expectedPath string, expectedStockCode string, mockData interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 경로 확인
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 메소드 확인
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 요청 본문 파싱
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// 요청 데이터 검증
		if in, ok := req["In"].(map[string]interface{}); ok {
			if stockCode, ok := in["InputIscd1"].(string); ok {
				if stockCode != expectedStockCode {
					t.Errorf("Expected stock code %s, got %s", expectedStockCode, stockCode)
				}
			}
		}

		// 응답 생성
		response := map[string]interface{}{
			"rsp_cd":  "00000",
			"rsp_msg": "정상 처리 되었습니다.",
			"Out":     mockData,
		}

		// 응답 헤더 설정
		w.Header().Set("Content-Type", "application/json")

		// 응답 작성
		json.NewEncoder(w).Encode(response)
	}
}
// CreateForeignDayChartMockHandler 해외주식 일차트 조회용 모의 핸들러 생성
func CreateForeignDayChartMockHandler(t *testing.T, expectedPath string, expectedStockCode string, mockData interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 경로 확인
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 메소드 확인
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 요청 본문 파싱
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// 요청 데이터 검증
		if in, ok := req["In"].(map[string]interface{}); ok {
			if stockCode, ok := in["InputIscd1"].(string); ok {
				if stockCode != expectedStockCode {
					t.Errorf("Expected stock code %s, got %s", expectedStockCode, stockCode)
				}
			}
		}

		// 응답 생성
		response := map[string]interface{}{
			"rsp_cd":  "00000",
			"rsp_msg": "정상 처리 되었습니다.",
			"Out":     mockData,
		}

		// 응답 헤더 설정
		w.Header().Set("Content-Type", "application/json")

		// 응답 작성
		json.NewEncoder(w).Encode(response)
	}
}
// CreateForeignWeekChartMockHandler 해외주식 주차트 조회용 모의 핸들러 생성
func CreateForeignWeekChartMockHandler(t *testing.T, expectedPath string, expectedStockCode string, mockData interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 경로 확인
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 메소드 확인
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 요청 본문 파싱
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// 요청 데이터 검증
		if in, ok := req["In"].(map[string]interface{}); ok {
			if stockCode, ok := in["InputIscd1"].(string); ok {
				if stockCode != expectedStockCode {
					t.Errorf("Expected stock code %s, got %s", expectedStockCode, stockCode)
				}
			}
			// 기간구분코드 확인
			if periodDiv, ok := in["InputPeriodDivCode"].(string); ok {
				if periodDiv != "W" {
					t.Errorf("Expected period div code W, got %s", periodDiv)
				}
			}
		}

		// 응답 생성
		response := map[string]interface{}{
			"rsp_cd":  "00000",
			"rsp_msg": "정상 처리 되었습니다.",
			"Out":     mockData,
		}

		// 응답 헤더 설정
		w.Header().Set("Content-Type", "application/json")

		// 응답 작성
		json.NewEncoder(w).Encode(response)
	}
}

// CreateForeignMonthChartMockHandler 해외주식 월차트 조회용 모의 핸들러 생성
func CreateForeignMonthChartMockHandler(t *testing.T, expectedPath string, expectedStockCode string, mockData interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 경로 확인
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// 메소드 확인
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}

		// 요청 본문 파싱
		var req map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// 요청 데이터 검증
		if in, ok := req["In"].(map[string]interface{}); ok {
			if stockCode, ok := in["InputIscd1"].(string); ok {
				if stockCode != expectedStockCode {
					t.Errorf("Expected stock code %s, got %s", expectedStockCode, stockCode)
				}
			}
		}

		// 응답 생성
		response := map[string]interface{}{
			"rsp_cd":  "00000",
			"rsp_msg": "정상 처리 되었습니다.",
			"Out":     mockData,
		}

		// 응답 헤더 설정
		w.Header().Set("Content-Type", "application/json")

		// 응답 작성
		json.NewEncoder(w).Encode(response)
	}
}