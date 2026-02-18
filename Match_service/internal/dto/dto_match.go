package dto

import (
	"match_service/internal/models"
	"time"
)

type CreateMatchRequest struct {
	SportID        uint      `json:"sport_id" binding:"required,gt=0"`
	HomeTeamID     uint      `json:"home_team_id" binding:"required,gt=0,nefield=AwayTeamID"`
	AwayTeamID     uint      `json:"away_team_id" binding:"required,gt=0,nefield=HomeTeamID"`
	ScheduledAt    time.Time `json:"scheduled_at" binding:"required"`
	Venue          string    `json:"venue,omitempty" binding:"omitempty,max=255"`
	TournamentName string    `json:"tournament_name,omitempty" binding:"omitempty,max=255"`
}

type StartMatchRequest struct {
	StartedAt *time.Time `json:"started_at,omitempty"`
}

type FinishMatchRequest struct {
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type UpdateScoreRequest struct {
	HomeScore int `json:"home_score" binding:"gte=0"`
	AwayScore int `json:"away_score" binding:"gte=0"`
}

type MatchResponse struct {
	ID             uint               `json:"id"`
	SportID        uint               `json:"sport_id"`
	HomeTeamID     uint               `json:"home_team_id"`
	AwayTeamID     uint               `json:"away_team_id"`
	Status         models.MatchStatus `json:"status"`
	ScheduledAt    time.Time          `json:"scheduled_at"`
	Venue          string             `json:"venue"`
	TournamentName string             `json:"tournament_name"`
}

type MatchListResponse struct {
	Data     []MatchResponse `json:"data"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Total    int64           `json:"total"`
}
