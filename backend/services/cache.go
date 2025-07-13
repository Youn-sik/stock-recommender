package services

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-recommender/backend/config"
	"stock-recommender/backend/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	client *redis.Client
	ctx    context.Context
}

func NewCacheService(cfg *config.Config) *CacheService {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: "", // no password
		DB:       0,  // default DB
	})

	return &CacheService{
		client: rdb,
		ctx:    context.Background(),
	}
}

// StockPrice 캐싱
func (c *CacheService) SetStockPrice(symbol string, price *models.StockPrice) error {
	key := fmt.Sprintf("stock:price:%s", symbol)
	data, err := json.Marshal(price)
	if err != nil {
		return err
	}
	
	return c.client.Set(c.ctx, key, data, time.Minute*5).Err()
}

func (c *CacheService) GetStockPrice(symbol string) (*models.StockPrice, error) {
	key := fmt.Sprintf("stock:price:%s", symbol)
	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var price models.StockPrice
	err = json.Unmarshal([]byte(data), &price)
	return &price, err
}

// 기술지표 캐싱
func (c *CacheService) SetIndicators(symbol string, indicators map[string]float64) error {
	key := fmt.Sprintf("indicators:%s", symbol)
	
	// Redis HSET으로 저장
	fields := make(map[string]interface{})
	for k, v := range indicators {
		fields[k] = v
	}
	
	err := c.client.HMSet(c.ctx, key, fields).Err()
	if err != nil {
		return err
	}
	
	// TTL 설정
	return c.client.Expire(c.ctx, key, time.Minute*10).Err()
}

func (c *CacheService) GetIndicators(symbol string) (map[string]float64, error) {
	key := fmt.Sprintf("indicators:%s", symbol)
	
	result, err := c.client.HGetAll(c.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	indicators := make(map[string]float64)
	for k, v := range result {
		var value float64
		if err := json.Unmarshal([]byte(v), &value); err == nil {
			indicators[k] = value
		}
	}
	
	return indicators, nil
}

// 매매 신호 캐싱
func (c *CacheService) SetSignals(symbol string, signals []models.TradingSignal) error {
	key := fmt.Sprintf("signals:%s", symbol)
	data, err := json.Marshal(signals)
	if err != nil {
		return err
	}
	
	return c.client.Set(c.ctx, key, data, time.Minute*15).Err()
}

func (c *CacheService) GetSignals(symbol string) ([]models.TradingSignal, error) {
	key := fmt.Sprintf("signals:%s", symbol)
	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var signals []models.TradingSignal
	err = json.Unmarshal([]byte(data), &signals)
	return signals, err
}

// 종목 목록 캐싱
func (c *CacheService) SetStocks(market string, stocks []models.Stock) error {
	key := fmt.Sprintf("stocks:%s", market)
	data, err := json.Marshal(stocks)
	if err != nil {
		return err
	}
	
	return c.client.Set(c.ctx, key, data, time.Hour).Err()
}

func (c *CacheService) GetStocks(market string) ([]models.Stock, error) {
	key := fmt.Sprintf("stocks:%s", market)
	data, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	
	var stocks []models.Stock
	err = json.Unmarshal([]byte(data), &stocks)
	return stocks, err
}

// 캐시 무효화
func (c *CacheService) InvalidateStock(symbol string) error {
	pattern := fmt.Sprintf("*:%s", symbol)
	keys, err := c.client.Keys(c.ctx, pattern).Result()
	if err != nil {
		return err
	}
	
	if len(keys) > 0 {
		return c.client.Del(c.ctx, keys...).Err()
	}
	
	return nil
}

// 헬스 체크
func (c *CacheService) Ping() error {
	return c.client.Ping(c.ctx).Err()
}