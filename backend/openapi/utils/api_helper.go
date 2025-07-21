package utils

import (
	"encoding/json"
	"fmt"
)

// APIResponse 공통 API 응답 인터페이스
type APIResponse interface {
	GetResponseCode() string
	GetResponseMessage() string
}

// BaseAPIResponse 기본 API 응답 구조
type BaseAPIResponse struct {
	RspCd  string `json:"rsp_cd"`
	RspMsg string `json:"rsp_msg"`
}

// GetResponseCode 응답 코드 반환
func (r *BaseAPIResponse) GetResponseCode() string {
	return r.RspCd
}

// GetResponseMessage 응답 메시지 반환
func (r *BaseAPIResponse) GetResponseMessage() string {
	return r.RspMsg
}

// ParseAPIResponse API 응답을 파싱하고 검증
func ParseAPIResponse(respBody []byte, response APIResponse) error {
	if err := json.Unmarshal(respBody, response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if response.GetResponseCode() != "00000" {
		return fmt.Errorf("API error %s: %s", response.GetResponseCode(), response.GetResponseMessage())
	}

	return nil
}

// PaginationHelper 페이지네이션 처리 헬퍼
type PaginationHelper struct {
	ContKey string
}

// NewPaginationHelper 새로운 페이지네이션 헬퍼 생성
func NewPaginationHelper() *PaginationHelper {
	return &PaginationHelper{}
}

// GetContYn 연속거래 여부 반환
func (p *PaginationHelper) GetContYn() string {
	if p.ContKey == "" {
		return "N"
	}
	return "Y"
}

// SetNextKey 다음 연속키 설정
func (p *PaginationHelper) SetNextKey(contKey string) {
	p.ContKey = contKey
}

// HasNext 다음 페이지 존재 여부
func (p *PaginationHelper) HasNext() bool {
	return p.ContKey != "" && p.ContKey != "N"
}

// IsSuccessResponse 성공 응답인지 확인
func IsSuccessResponse(responseCode string) bool {
	return responseCode == "00000"
}