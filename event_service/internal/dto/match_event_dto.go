package dto

import (
	"time"
)

type TimelineItemType string

const (
	TimelineTypeEvent      TimelineItemType = "event"
	TimelineTypeCommentary TimelineItemType = "commentary"
)

type CreateMatchEventRequest struct {
	MatchID   uint64 `json:"match_id" binding:"required"`
	EventType string `json:"event_type" binding:"required,oneof=goal yellow_card red_card substitution penalty var_decision half_start half_end full_time injury timeout"`
	Minute    int    `json:"minute" binding:"required,min=0"`
	Period    *int   `json:"period,omitempty"`

	TeamID   *uint64 `json:"team_id,omitempty"`
	PlayerID *uint64 `json:"player_id,omitempty"`

	Description string `json:"description,omitempty"`
}

type UpdateMatchEventRequest struct {
	Minute      *int    `json:"minute,omitempty"`
	Period      *int    `json:"period,omitempty"`
	Description *string `json:"description,omitempty"`
}

type MatchEventResponse struct {
	ID        uint64 `json:"id"`
	MatchID   uint64 `json:"match_id"`
	EventType string `json:"event_type"`

	Minute int  `json:"minute"`
	Period *int `json:"period,omitempty"`

	TeamID   *uint64 `json:"team_id,omitempty"`
	PlayerID *uint64 `json:"player_id,omitempty"`

	Description string `json:"description,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

type TimelineEventDTO struct {
	ID        uint64           `json:"id"`
	Type      TimelineItemType `json:"type"`
	EventType string           `json:"event_type"`

	Minute int `json:"minute"`

	TeamID   *uint64 `json:"team_id,omitempty"`
	PlayerID *uint64 `json:"player_id,omitempty"`

	Description string `json:"description,omitempty"`

	Reactions map[string]int `json:"reactions"`
}

type MatchEventCreatedMessage struct {
	MatchID     uint64 `json:"match_id"`
	EventType   string `json:"event_type"`
	Minute      int    `json:"minute"`
	Description string `json:"description"`
}

type MatchGoalMessage struct {
	MatchID  uint64  `json:"match_id"`
	TeamID   uint64  `json:"team_id"`
	PlayerID *uint64 `json:"player_id,omitempty"`
	Minute   int     `json:"minute"`
}
