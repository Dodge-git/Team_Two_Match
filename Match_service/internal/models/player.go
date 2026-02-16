package models

import "gorm.io/gorm"

type Player struct {
	gorm.Model
	TeamID uint `json:"team_id" gorm:"not null;unique"`
	Team   Team `json:"-"`

	Name     string `json:"name" gorm:"not null;size:100"`
	Number   uint   `json:"number" gorm:"not null;unique"`
	Position string `json:"position" gorm:"not null;size:50"`
}
