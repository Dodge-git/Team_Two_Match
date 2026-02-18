package models

import (
	"time"
)

type EventType string

const (
	EventGoal         EventType = "goal"
	EventYellowCard   EventType = "yellow_card"
	EventRedCard      EventType = "red_card"
	EventSubstitution EventType = "substitution"
	EventPenalty      EventType = "penalty"
	EventVarDecision  EventType = "var_decision"

	EventHalfStart EventType = "half_start"
	EventHalfEnd   EventType = "half_end"
	EventFullTime  EventType = "full_time"
	EventInjury    EventType = "injury"
	EventTimeout   EventType = "timeout"
)

type MatchEvent struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	MatchID   uint64    `json:"match_id" gorm:"not null;index:idx_match_time_minute,priority:1"`
	EventType EventType `gorm:"type:varchar(30);not null;index"`

	Minute int `gorm:"not null;index:idx_match_time_minute,priority:2"`
	Period *int

	TeamID   *uint64 `gorm:"index"`
	PlayerID *uint64 `gorm:"index"`

	Description string `gorm:"type:varchar(500)"`

	CreatedAt time.Time
}
