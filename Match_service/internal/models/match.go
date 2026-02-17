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
	Sport   Sport `json:"-" gorm:"foreignKey:SportID"`

	HomeTeamID uint `json:"home_team_id" gorm:"not null"`
	HomeTeam   Team `json:"-" gorm:"foreignKey:HomeTeamID;check:home_team_id <> away_team_id"`

	AwayTeamID uint `json:"away_team_id" gorm:"not null"`
	AwayTeam   Team `json:"-" gorm:"foreignKey:AwayTeamID"`

	Status MatchStatus `json:"status" gorm:"type:varchar(16);not null;default:'scheduled'; check:status IN ('scheduled','live','finished','cancelled')"`

	ScheduledAt time.Time  `json:"scheduled_at" gorm:"not null"`
	StartedAt   *time.Time `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at"`

	HomeScore int `json:"home_score" gorm:"not null;default:0"`
	AwayScore int `json:"away_score" gorm:"not null;default:0"`

	Venue          string `json:"venue" gorm:"size:255"`
	TournamentName string `json:"tournament_name" gorm:"size:255"`
}
