package service

import "github.com/mountainman199231/event_service/internal/dto"

type NoOpKafkaProducer struct{}

func NewNoOpKafkaProducer() KafkaProducer {
	return &NoOpKafkaProducer{}
}

func (n *NoOpKafkaProducer) PublishMatchEventCreated(msg dto.MatchEventCreatedMessage) error {
	return nil
}

func (n *NoOpKafkaProducer) PublishMatchGoal(msg dto.MatchGoalMessage) error {
	return nil
}
