package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokers []string) KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}

	return &kafkaProducer{
		writer: writer,
	}

}

func (k *kafkaProducer) PublishMatchEventCreated(msg dto.MatchEventCreatedMessage) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return k.writer.WriteMessages(ctx,
		kafka.Message{
			Topic: "match.event.created",
			Key:   []byte(strconv.FormatUint(msg.MatchID, 10)),
			Value: bytes,
			Time:  time.Now(),
		},
	)
}

func (k *kafkaProducer) PublishMatchGoal(msg dto.MatchGoalMessage) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return k.writer.WriteMessages(ctx,
		kafka.Message{
			Topic: "match.goal",
			Key:   []byte(strconv.FormatUint(msg.MatchID, 10)),
			Value: bytes,
			Time:  time.Now(),
		},
	)
}

func (k *kafkaProducer) Close() error {
	return k.writer.Close()
}
