package config

import (
	"os"
)

type Config struct {
	Port     string
	Database DatabaseConfig
	Redis    RedisConfig
	RabbitMQ RabbitMQConfig
	API      APIConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host string
	Port string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type APIConfig struct {
	DBSecAPIKey    string
	DBSecAppKey    string
	DBSecAppSecret string
	AIServiceURL   string
}

func Load() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "stockuser"),
			Password: getEnv("DB_PASSWORD", "stockpass"),
			Name:     getEnv("DB_NAME", "stockdb"),
		},
		Redis: RedisConfig{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnv("RABBITMQ_PORT", "5672"),
			User:     getEnv("RABBITMQ_USER", "stockmq"),
			Password: getEnv("RABBITMQ_PASS", "stockmqpass"),
		},
		API: APIConfig{
			DBSecAPIKey:    getEnv("DBSEC_APP_KEY", ""),
			DBSecAppKey:    getEnv("DBSEC_APP_KEY", ""),
			DBSecAppSecret: getEnv("DBSEC_APP_SECRET", ""),
			AIServiceURL:   getEnv("AI_SERVICE_URL", "http://localhost:8001"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
