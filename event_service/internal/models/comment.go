package models

import "time"

type Comment struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	UserID uint64 `json:"user_id" gorm:"not null;index"`

	EventID      *uint64 `json:"event_id,omitempty" gorm:"index"`
	CommentaryID *uint64 `json:"commentary_id,omitempty" gorm:"index"`

	Text string `json:"text" gorm:"type:text;not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
