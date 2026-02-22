package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name          string  `json:"name" gorm:"not null,size:66"`
	Email         string  `json:"email" gorm:"not null,uniqueIndex"`
	Phone         string  `json:"phone" gorm:"not null,uniqueIndex"`
	Role          string  `json:"role" gorm:"not null,size:20,default:'user',check:role IN ('user','commentator','admin','fan')"`
	PasswordHash  string  `json:"-" gorm:"not null,size:255"`
	FavoriteSport *string `json:"favorite_sport" gorm:"size:63"`
}


