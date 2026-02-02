package kafka

import (
	"encoding/json"
	"fmt"
	"jp-ru-dict/backend/internal/config"
	"jp-ru-dict/backend/pkg/logger"
	"time"

	"github.com/IBM/sarama"
)

// EventType тип события
type EventType string

const (
	EventWordCreated EventType = "WORD_CREATED"
	EventWordUpdated EventType = "WORD_UPDATED"
	EventWordDeleted EventType = "WORD_DELETED"
)

// Event структура события для Kafka
type Event struct {
	Type      EventType   `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	UserID    int         `json:"user_id"`
}

// Producer интерфейс для отправки событий
type Producer interface {
	SendEvent(eventType EventType, userID int, payload interface{}) error
	Close() error
}

// SaramaProducer реализация Producer на основе sarama
type SaramaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer создает новый Kafka producer
func NewProducer(cfg config.KafkaConfig) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{cfg.Broker}, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &SaramaProducer{
		producer: producer,
		topic:    "word-events",
	}, nil
}

func (p *SaramaProducer) SendEvent(eventType EventType, userID int, payload interface{}) error {
	event := Event{
		Type:      eventType,
		Payload:   payload,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	bytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", userID)),
		Value: sarama.ByteEncoder(bytes),
	}

	_, _, err = p.producer.SendMessage(msg)
	if err != nil {
		// NFR-010: Логируем в консоль, если Kafka недоступна
		logger.Log.Error().Err(err).Str("event", string(bytes)).Msg("Не удалось отправить событие в Kafka")
		return err
	}

	return nil
}

func (p *SaramaProducer) Close() error {
	return p.producer.Close()
}

// NoOpProducer заглушка для продюсера (используется если Kafka недоступна)
type NoOpProducer struct{}

// NewNoOpProducer создает новый NoOp продюсер
func NewNoOpProducer() Producer {
	return &NoOpProducer{}
}

// SendEvent ничего не делает
func (p *NoOpProducer) SendEvent(eventType EventType, userID int, payload interface{}) error {
	return nil
}

// Close ничего не делает
func (p *NoOpProducer) Close() error {
	return nil
}
