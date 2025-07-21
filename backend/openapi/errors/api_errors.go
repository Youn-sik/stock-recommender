package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode API 에러 코드 타입
type ErrorCode string

const (
	// 인증 관련 에러
	ErrCodeAuthFailed     ErrorCode = "AUTH_FAILED"
	ErrCodeTokenExpired   ErrorCode = "TOKEN_EXPIRED"
	ErrCodeInvalidKey     ErrorCode = "INVALID_KEY"
	
	// 네트워크 관련 에러
	ErrCodeTimeout        ErrorCode = "TIMEOUT"
	ErrCodeNetworkError   ErrorCode = "NETWORK_ERROR"
	ErrCodeRateLimit      ErrorCode = "RATE_LIMIT"
	
	// 데이터 관련 에러
	ErrCodeInvalidData    ErrorCode = "INVALID_DATA"
	ErrCodeParseError     ErrorCode = "PARSE_ERROR"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	
	// 시스템 관련 에러
	ErrCodeServerError    ErrorCode = "SERVER_ERROR"
	ErrCodeUnknown        ErrorCode = "UNKNOWN"
)

// APIError API 에러 구조체
type APIError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	StatusCode int       `json:"status_code,omitempty"`
	Cause      error     `json:"-"`
}

// Error error 인터페이스 구현
func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 원본 에러 반환
func (e *APIError) Unwrap() error {
	return e.Cause
}

// NewAPIError 새로운 API 에러 생성
func NewAPIError(code ErrorCode, message string, cause error) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewAuthError 인증 에러 생성
func NewAuthError(message string, cause error) *APIError {
	return &APIError{
		Code:       ErrCodeAuthFailed,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		Cause:      cause,
	}
}

// NewNetworkError 네트워크 에러 생성
func NewNetworkError(message string, cause error) *APIError {
	return &APIError{
		Code:       ErrCodeNetworkError,
		Message:    message,
		StatusCode: http.StatusServiceUnavailable,
		Cause:      cause,
	}
}

// NewParseError 파싱 에러 생성
func NewParseError(message string, cause error) *APIError {
	return &APIError{
		Code:       ErrCodeParseError,
		Message:    message,
		StatusCode: http.StatusUnprocessableEntity,
		Cause:      cause,
	}
}

// NewRateLimitError 레이트 리미트 에러 생성
func NewRateLimitError(message string) *APIError {
	return &APIError{
		Code:       ErrCodeRateLimit,
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

// NewValidationError 유효성 검증 에러 생성
func NewValidationError(message string, cause error) *APIError {
	return &APIError{
		Code:       ErrCodeValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Cause:      cause,
	}
}

// IsRetryableError 재시도 가능한 에러인지 확인
func IsRetryableError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		switch apiErr.Code {
		case ErrCodeTimeout, ErrCodeNetworkError, ErrCodeServerError:
			return true
		}
	}
	return false
}

// IsAuthError 인증 에러인지 확인
func IsAuthError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Code == ErrCodeAuthFailed || apiErr.Code == ErrCodeTokenExpired
	}
	return false
}