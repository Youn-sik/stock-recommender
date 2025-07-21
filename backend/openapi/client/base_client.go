package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"stock-recommender/backend/config"
	"stock-recommender/backend/openapi/errors"
	"stock-recommender/backend/openapi/logger"
	"stock-recommender/backend/openapi/models"
	"stock-recommender/backend/openapi/utils"
)

type DBSecClient struct {
	baseURL           string
	appKey            string
	appSecret         string
	accessToken       string
	httpClient        *http.Client
	rateLimiter       chan struct{}
	tokenGenerateTime time.Time
	logger            logger.Logger
}

// 인증 토큰 응답 구조체
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func NewDBSecClient(cfg *config.Config) *DBSecClient {
	// Rate limiter: 초당 20요청으로 제한
	rateLimiter := make(chan struct{}, 20)
	go func() {
		for {
			time.Sleep(50 * time.Millisecond) // 20 requests per second
			select {
			case rateLimiter <- struct{}{}:
			default:
			}
		}
	}()

	client := &DBSecClient{
		baseURL:     "https://openapi.dbsec.co.kr:8443",
		appKey:      cfg.API.DBSecAppKey,
		appSecret:   cfg.API.DBSecAppSecret,
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		rateLimiter: rateLimiter,
		logger:      logger.GetDefaultLogger().With(logger.Field{Key: "component", Value: "dbsec_client"}),
	}

	// 시작시 토큰 발급
	if client.appKey != "" && client.appSecret != "" {
		err := client.authenticate()
		if err != nil {
			client.logger.Warn("Failed to authenticate with DBSec API during initialization", logger.Field{Key: "error", Value: err})
		}
	}

	return client
}

// OAuth 인증 토큰 발급
func (c *DBSecClient) authenticate() error {
	c.logger.Debug("Starting authentication process")
	
	authURL := c.baseURL + "/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("appkey", c.appKey)
	data.Set("appsecretkey", c.appSecret)
	data.Set("scope", "oob")

	req, err := http.NewRequest("POST", authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.NewNetworkError("failed to create auth request", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.NewNetworkError("auth request failed", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewNetworkError("failed to read auth response", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Authentication failed", fmt.Errorf("status: %d", resp.StatusCode),
			logger.Field{Key: "status_code", Value: resp.StatusCode},
			logger.Field{Key: "response_body", Value: string(body)})
		return errors.NewAuthError("authentication failed", fmt.Errorf("status %d: %s", resp.StatusCode, string(body)))
	}

	var tokenResp TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return errors.NewParseError("failed to parse token response", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenGenerateTime = time.Now()
	
	c.logger.Info("Successfully authenticated with DBSec API",
		logger.Field{Key: "token_type", Value: tokenResp.TokenType},
		logger.Field{Key: "scope", Value: tokenResp.Scope},
		logger.Field{Key: "expires_in", Value: tokenResp.ExpiresIn})

	return nil
}

// API 호출을 위한 공통 함수
func (c *DBSecClient) makeRequest(method, path string, queryParams map[string]string, body interface{}) ([]byte, error) {
	return c.MakeRequestWithHeaders(method, path, queryParams, body, nil)
}

// MakeRequestWithHeaders 추가 헤더를 포함한 API 호출
func (c *DBSecClient) MakeRequestWithHeaders(method, path string, queryParams map[string]string, body interface{}, additionalHeaders map[string]string) ([]byte, error) {
	return c.MakeRequestWithResponse(method, path, queryParams, body, additionalHeaders)
}

// MakeRequestWithResponse 응답 헤더를 포함한 API 호출
func (c *DBSecClient) MakeRequestWithResponse(method, path string, queryParams map[string]string, body interface{}, additionalHeaders map[string]string) ([]byte, error) {
	// Rate limiting
	<-c.rateLimiter

	// 토큰이 없으면 인증 시도
	if c.accessToken == "" {
		if err := c.authenticate(); err != nil {
			return nil, fmt.Errorf("authentication failed: %w", err)
		}
	}

	// URL 구성
	fullURL := c.baseURL + path
	if queryParams != nil && len(queryParams) > 0 {
		params := url.Values{}
		for k, v := range queryParams {
			params.Set(k, v)
		}
		fullURL += "?" + params.Encode()
	}

	// 요청 본문 준비
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonData)
	}

	// HTTP 요청 생성
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 헤더 설정
	c.setCommonHeaders(req, path, queryParams)

	// 추가 헤더 설정
	for key, value := range additionalHeaders {
		req.Header.Set(key, value)
	}

	// 요청 실행
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// 응답 읽기
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 상태 코드 확인
	if resp.StatusCode != http.StatusOK {
		// 토큰 만료 등의 경우 재인증 시도
		if resp.StatusCode == http.StatusUnauthorized {
			c.logger.Info("Token expired, attempting re-authentication")
			if err := c.authenticate(); err == nil {
				c.logger.Debug("Re-authentication successful, retrying request")
				// 재인증 성공시 요청 재시도
				return c.MakeRequestWithResponse(method, path, queryParams, body, additionalHeaders)
			} else {
				c.logger.Error("Re-authentication failed", err)
			}
		}
		
		c.logger.Error("API request failed", fmt.Errorf("status: %d", resp.StatusCode),
			logger.Field{Key: "method", Value: method},
			logger.Field{Key: "path", Value: path},
			logger.Field{Key: "status_code", Value: resp.StatusCode},
			logger.Field{Key: "response_body", Value: string(respBody)})
		
		return nil, errors.NewNetworkError("API request failed", fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody)))
	}

	return respBody, nil
}

// 공통 헤더 설정
func (c *DBSecClient) setCommonHeaders(req *http.Request, path string, queryParams map[string]string) {
	// 기본 헤더
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("appkey", c.appKey)
	req.Header.Set("appsecret", c.appSecret)

	// 트랜잭션 ID 설정 (API 별로 다름)
	trId := c.getTransactionId(path)
	req.Header.Set("tr_id", trId)

	// 고객 타입 (기본값)
	req.Header.Set("custtype", "P")

	// 해시키 생성 (POST 요청의 경우)
	if req.Method == "POST" && req.Body != nil {
		hashKey := c.generateHashKey(queryParams)
		req.Header.Set("hashkey", hashKey)
	}
}

// API 경로에 따른 트랜잭션 ID 반환
func (c *DBSecClient) getTransactionId(path string) string {
	switch path {
	case models.PathDomesticStockPrice:
		return "FHKST01010100"
	case models.PathDomesticStockAsking:
		return "FHKST01010200"
	case models.PathDomesticStockDaily:
		return "FHKST03010100"
	case models.PathDomesticStockList:
		return "CTPF1002R"
	case models.PathDomesticStockTicker:
		return models.TrIdStockTicker
	case models.PathDomesticStockCurrentPrice:
		return models.TrIdStockCurrentPrice
	case models.PathForeignStockTicker:
		return models.TrIdForeignStockTicker
	case models.PathForeignStockCurrentPrice:
		return models.TrIdForeignStockCurrentPrice
	case models.PathForeignStockMinChart:
		return models.TrIdForeignStockMinChart
	case models.PathForeignStockDayChart:
		return models.TrIdForeignStockDayChart
	case models.PathForeignStockWeekChart:
		return models.TrIdForeignStockWeekChart
	case models.PathForeignStockMonthChart:
		return models.TrIdForeignStockMonthChart
	case models.PathForeignStockPrice:
		return "HHDFS00000300"
	case models.PathForeignStockDaily:
		return "HHDFS76240000"
	case models.PathIndexPrice:
		return "FHPUP02100000"
	default:
		return "FHKST01010100" // 기본값
	}
}

// 해시키 생성 (POST 요청용)
func (c *DBSecClient) generateHashKey(params map[string]string) string {
	// 파라미터를 정렬된 문자열로 변환
	var paramPairs []string
	for k, v := range params {
		paramPairs = append(paramPairs, fmt.Sprintf("%s=%s", k, v))
	}

	// 파라미터 문자열 생성
	paramString := strings.Join(paramPairs, "&")

	// HMAC-SHA256으로 해시 생성
	h := hmac.New(sha256.New, []byte(c.appSecret))
	h.Write([]byte(paramString))
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return hash
}

// 헬스체크
func (c *DBSecClient) HealthCheck() error {
	if c.appKey == "" || c.appSecret == "" {
		return fmt.Errorf("API credentials not configured")
	}

	if c.accessToken == "" {
		return c.authenticate()
	}

	if time.Since(c.tokenGenerateTime) > time.Duration(23)*time.Hour {
		return c.authenticate()
	}

	return nil
}

// 유틸리티 함수들 (레거시 지원을 위해 유지, 새 코드는 utils 패키지 사용 권장)
func (c *DBSecClient) parseFloat(s string) float64 {
	return utils.ParseFloat(s)
}

func (c *DBSecClient) parseInt(s string) int64 {
	return utils.ParseInt(s)
}

func (c *DBSecClient) parseDate(dateStr string) time.Time {
	return utils.ParseDate(dateStr)
}

// API 키 유효성 검사
func (c *DBSecClient) HasValidCredentials() bool {
	return c.appKey != "" && c.appSecret != ""
}

// 토큰 재발급
func (c *DBSecClient) RefreshToken() error {
	c.accessToken = ""
	return c.authenticate()
}
