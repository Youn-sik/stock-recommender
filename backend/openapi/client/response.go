package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// APIResponse API 응답 구조체
type APIResponse struct {
	Body    []byte
	Headers http.Header
}

// MakeRequestWithFullResponse 응답 헤더를 포함한 API 호출
func (c *DBSecClient) MakeRequestWithFullResponse(method, path string, queryParams map[string]string, body interface{}, additionalHeaders map[string]string) (*APIResponse, error) {
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
			fmt.Println("Token expired, re-authenticating...")
			if err := c.authenticate(); err == nil {
				// 재인증 성공시 요청 재시도
				return c.MakeRequestWithFullResponse(method, path, queryParams, body, additionalHeaders)
			}
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	return &APIResponse{
		Body:    respBody,
		Headers: resp.Header,
	}, nil
}