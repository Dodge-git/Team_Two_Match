package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	TeamID uint `json:"team_id" gorm:"not null;uniqueIndex:idx_team_number"`
	Team   Team `json:"-"`

	Name     string `json:"name" gorm:"not null;size:100"`
	Number   uint   `json:"number" gorm:"not null;uniqueIndex:idx_team_number"`
	Position string `json:"position" gorm:"not null;size:50"`
}
