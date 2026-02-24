package kafka

import (
	"context"
	"encoding/json"
	"log"
	"match_service/internal/models"
	"match_service/internal/services"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader       *kafka.Reader
	matchService services.MatchService
}

func NewConsumer(brokers []string, topic string, groupID string, matchService services.MatchService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	return &Consumer{
		reader:       reader,
		matchService: matchService,
	}

}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started...")
	defer c.reader.Close()

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("Kafka consumer stopped")
				break
			}
			log.Printf("Kafka read error: %v", err)
			continue
		}
		var event models.GoalEvent

		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("error unmarshalling message:", err)
			_ = c.reader.CommitMessages(ctx, msg)
			continue
		}
		if err := c.matchService.IncrementGoal(event); err != nil {
			log.Println("failed to increment goal:", err)
			continue
		}
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("failed to commit message: %v", err)
		}
	}
}
