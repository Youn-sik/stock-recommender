package services

import (
	"encoding/json"
	"fmt"
	"log"
	"stock-recommender/backend/config"
	"time"

	"github.com/streadway/amqp"
)

type QueueService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// 메시지 타입
type Message struct {
	Type      string      `json:"type"`
	Symbol    string      `json:"symbol"`
	Market    string      `json:"market"`
	Data      interface{} `json:"data"`
	Timestamp string      `json:"timestamp"`
}

// 메시지 타입 상수
const (
	MessageTypePriceUpdate      = "price_update"
	MessageTypeIndicatorRequest = "indicator_request"
	MessageTypeIndicatorResult  = "indicator_result"
	MessageTypeAIRequest        = "ai_request"
	MessageTypeAIResponse       = "ai_response"
	MessageTypeSignalGenerated  = "signal_generated"
	MessageTypeNewsUpdate       = "news_update"
)

func NewQueueService(cfg *config.Config) (*QueueService, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)

	conn, err := amqp.Dial(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	qs := &QueueService{
		conn:    conn,
		channel: ch,
	}

	// Exchange와 Queue 설정
	if err := qs.setupExchangesAndQueues(); err != nil {
		return nil, fmt.Errorf("failed to setup exchanges and queues: %w", err)
	}

	return qs, nil
}

func (qs *QueueService) setupExchangesAndQueues() error {
	// Exchange 선언
	exchanges := []string{
		"stock.data",
		"trading.signals",
		"news.analysis",
	}

	for _, exchange := range exchanges {
		err := qs.channel.ExchangeDeclare(
			exchange, // name
			"topic",  // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare exchange %s: %w", exchange, err)
		}
	}

	// Queue 선언 및 바인딩
	queues := map[string]string{
		"price.updates":         "stock.data",
		"indicator.calculation": "stock.data",
		"ai.requests":          "stock.data",
		"signal.generation":    "trading.signals",
		"signal.notifications": "trading.signals",
		"news.crawling":        "news.analysis",
		"sentiment.analysis":   "news.analysis",
	}

	for queueName, exchange := range queues {
		_, err := qs.channel.QueueDeclare(
			queueName, // name
			true,      // durable
			false,     // delete when unused
			false,     // exclusive
			false,     // no-wait
			nil,       // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}

		// 큐를 Exchange에 바인딩
		err = qs.channel.QueueBind(
			queueName,        // queue name
			queueName,        // routing key
			exchange,         // exchange
			false,            // no-wait
			nil,              // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueName, exchange, err)
		}
	}

	return nil
}

// 메시지 발행
func (qs *QueueService) Publish(exchange, routingKey string, message Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = qs.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   message.Timestamp,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to %s/%s: %s", exchange, routingKey, message.Type)
	return nil
}

// 메시지 구독
func (qs *QueueService) Subscribe(queueName string, handler func(Message) error) error {
	msgs, err := qs.channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			var message Message
			if err := json.Unmarshal(d.Body, &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				d.Nack(false, false) // 메시지 거부
				continue
			}

			if err := handler(message); err != nil {
				log.Printf("Failed to handle message: %v", err)
				d.Nack(false, true) // 메시지 거부 후 재큐잉
			} else {
				d.Ack(false) // 메시지 확인
			}
		}
	}()

	log.Printf("Started consuming from queue: %s", queueName)
	return nil
}

// 편의 메서드들
func (qs *QueueService) PublishPriceUpdate(symbol, market string, data interface{}) error {
	message := Message{
		Type:      MessageTypePriceUpdate,
		Symbol:    symbol,
		Market:    market,
		Data:      data,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}
	return qs.Publish("stock.data", "price.updates", message)
}

func (qs *QueueService) PublishIndicatorRequest(symbol, market string) error {
	message := Message{
		Type:      MessageTypeIndicatorRequest,
		Symbol:    symbol,
		Market:    market,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}
	return qs.Publish("stock.data", "indicator.calculation", message)
}

func (qs *QueueService) PublishAIRequest(symbol, market string, indicators interface{}) error {
	message := Message{
		Type:      MessageTypeAIRequest,
		Symbol:    symbol,
		Market:    market,
		Data:      indicators,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}
	return qs.Publish("stock.data", "ai.requests", message)
}

func (qs *QueueService) PublishSignal(symbol, market string, signal interface{}) error {
	message := Message{
		Type:      MessageTypeSignalGenerated,
		Symbol:    symbol,
		Market:    market,
		Data:      signal,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
	}
	return qs.Publish("trading.signals", "signal.generation", message)
}

// 리소스 정리
func (qs *QueueService) Close() error {
	if qs.channel != nil {
		qs.channel.Close()
	}
	if qs.conn != nil {
		return qs.conn.Close()
	}
	return nil
}

// 헬스 체크
func (qs *QueueService) HealthCheck() error {
	if qs.conn == nil || qs.conn.IsClosed() {
		return fmt.Errorf("connection is closed")
	}
	return nil
}