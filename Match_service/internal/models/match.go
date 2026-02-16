package models

import (
	"time"

	"gorm.io/gorm"
)

type MatchStatus string

const (
	MatchStatusScheduled MatchStatus = "scheduled"
	MatchStatusLive      MatchStatus = "live"
	MatchStatusFinished  MatchStatus = "finished"
	MatchStatusCancelled MatchStatus = "cancelled"
)

type Match struct {
	gorm.Model
	SportID uint  `json:"sport_id" gorm:"not null"`
	Sport   Sport `json:"-"`

	HomeTeamID uint `json:"home_team_id" gorm:"not null"`
	HomeTeam   Team `json:"-"`

	AwayTeamID uint `json:"away_team_id" gorm:"not null"`
	AwayTeam   Team `json:"-"`

	Status MatchStatus `json:"status" gorm:"type:varchar(16);not null;default:'scheduled'; check:status IN ('scheduled','live','finished','cancelled')"`

	ScheduledAt time.Time  `json:"scheduled_at" gorm:"not null"`
	StartedAt   *time.Time `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at"`

	Venue          string `json:"venue" gorm:"size:255"`
	TournamentName string `json:"tournament_name" gorm:"size:255"`

	HomeMatches []Match `json:"-" gorm:"foreignKey:HomeTeamID"`
	AwayMatches []Match `json:"-" gorm:"foreignKey:AwayTeamID"`
}
