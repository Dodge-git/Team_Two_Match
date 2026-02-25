package ports

import (
	"context"
	"match_service/internal/models"
)

type MatchService interface {
	UpdateScoreFromKafka(ctx context.Context, goalEvent models.GoalEvent) error
}
