package utils

import (
	"strconv"
	"strings"
	"time"
)

// ParseFloat 문자열을 float64로 변환
func ParseFloat(str string) float64 {
	if str == "" {
		return 0
	}
	val, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		return 0
	}
	return val
}

// ParseInt 문자열을 int64로 변환
func ParseInt(str string) int64 {
	if str == "" {
		return 0
	}
	val, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		return 0
	}
	return val
}

// ParseDate 날짜 문자열을 time.Time으로 변환
func ParseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}

	// YYYYMMDD 형식 파싱
	if len(dateStr) == 8 {
		if t, err := time.Parse("20060102", dateStr); err == nil {
			return t
		}
	}

	// YYYY-MM-DD 형식 파싱
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t
	}

	return time.Now()
}