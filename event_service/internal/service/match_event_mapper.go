package service

import (
	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/models"
)

func mapToMatchEventResponse(e *models.MatchEvent) *dto.MatchEventResponse {
	return &dto.MatchEventResponse{
		ID:          e.ID,
		MatchID:     e.MatchID,
		EventType:   string(e.EventType),
		Minute:      e.Minute,
		Period:      e.Period,
		TeamID:      e.TeamID,
		PlayerID:    e.PlayerID,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
	}
}
