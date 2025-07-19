package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel 로그 레벨 타입
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String LogLevel의 문자열 표현
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Field 로그 필드 구조체
type Field struct {
	Key   string
	Value interface{}
}

// Logger 인터페이스
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, err error, fields ...Field)
	With(fields ...Field) Logger
}

// DefaultLogger 기본 로거 구현
type DefaultLogger struct {
	level  LogLevel
	fields []Field
	logger *log.Logger
}

// NewDefaultLogger 새로운 기본 로거 생성
func NewDefaultLogger(level LogLevel) *DefaultLogger {
	return &DefaultLogger{
		level:  level,
		fields: make([]Field, 0),
		logger: log.New(os.Stdout, "", 0),
	}
}

// Debug 디버그 로그 출력
func (l *DefaultLogger) Debug(msg string, fields ...Field) {
	if l.level <= DEBUG {
		l.log(DEBUG, msg, nil, fields...)
	}
}

// Info 정보 로그 출력
func (l *DefaultLogger) Info(msg string, fields ...Field) {
	if l.level <= INFO {
		l.log(INFO, msg, nil, fields...)
	}
}

// Warn 경고 로그 출력
func (l *DefaultLogger) Warn(msg string, fields ...Field) {
	if l.level <= WARN {
		l.log(WARN, msg, nil, fields...)
	}
}

// Error 에러 로그 출력
func (l *DefaultLogger) Error(msg string, err error, fields ...Field) {
	if l.level <= ERROR {
		l.log(ERROR, msg, err, fields...)
	}
}

// With 필드를 추가한 새로운 로거 반환
func (l *DefaultLogger) With(fields ...Field) Logger {
	newFields := make([]Field, len(l.fields)+len(fields))
	copy(newFields, l.fields)
	copy(newFields[len(l.fields):], fields)
	
	return &DefaultLogger{
		level:  l.level,
		fields: newFields,
		logger: l.logger,
	}
}

// log 실제 로그 출력
func (l *DefaultLogger) log(level LogLevel, msg string, err error, fields ...Field) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	
	// 모든 필드 결합
	allFields := make([]Field, len(l.fields)+len(fields))
	copy(allFields, l.fields)
	copy(allFields[len(l.fields):], fields)
	
	// 필드 문자열 생성
	fieldsStr := ""
	if len(allFields) > 0 {
		fieldsStr = " "
		for i, field := range allFields {
			if i > 0 {
				fieldsStr += " "
			}
			fieldsStr += fmt.Sprintf("%s=%v", field.Key, field.Value)
		}
	}
	
	// 에러 정보 추가
	errStr := ""
	if err != nil {
		errStr = fmt.Sprintf(" error=%v", err)
	}
	
	l.logger.Printf("[%s] %s %s%s%s", timestamp, level.String(), msg, fieldsStr, errStr)
}

// 전역 로거 인스턴스
var defaultLogger Logger

func init() {
	defaultLogger = NewDefaultLogger(INFO)
}

// SetDefaultLogger 기본 로거 설정
func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

// GetDefaultLogger 기본 로거 반환
func GetDefaultLogger() Logger {
	return defaultLogger
}

// 편의 함수들
func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, err error, fields ...Field) {
	defaultLogger.Error(msg, err, fields...)
}

func With(fields ...Field) Logger {
	return defaultLogger.With(fields...)
}