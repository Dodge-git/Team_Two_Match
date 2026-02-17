package models

import "time"

type Commentary struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	MatchID uint64 `json:"match_id" gorm:"not null;index:idx_commentary_match_minute,priority:1"`

	Minute int `json:"minute" gorm:"not null;index:idx_commentary_match_minute,priority:2"`

	Text string `json:"text" gorm:"type:text;not null"`

	IsPinned bool `json:"is_pinned" gorm:"default:false;index"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
