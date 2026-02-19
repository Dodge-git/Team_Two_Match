package models

import "gorm.io/gorm"

type SportName string

const (
	SportFootball   SportName = "football"
	SportBasketball SportName = "basketball"
	SportMMA        SportName = "mma"
	SportBoxing     SportName = "boxing"
	SportVolleyball SportName = "volleyball"
)

type Sport struct {
	gorm.Model
	Name SportName `json:"name" gorm:"uniqueIndex;not null"`
}

var ValidSport = map[SportName]struct{}{
	SportBasketball: {},
	SportFootball:   {},
	SportMMA:        {},
	SportBoxing:     {},
	SportVolleyball: {},
}

type TeamFilter struct {
	SportID  *uint 
	Page     int
	PageSize int
}
