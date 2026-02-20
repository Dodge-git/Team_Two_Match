package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	SportID uint  `json:"sport_id" gorm:"not null"`
	Sport   Sport `json:"-" gorm:"foreignKey:SportID"`

	Name      string `json:"name" gorm:"not null;size:100"`
	ShortName string `json:"short_name" gorm:"not null;size:16"`
	City      string `json:"city" gorm:"not null;size:100"`

	Players []Player `json:"-"`

	HomeMatches []Match `json:"-" gorm:"foreignKey:HomeTeamID"`
	AwayMatches []Match `json:"-" gorm:"foreignKey:AwayTeamID"`
}
