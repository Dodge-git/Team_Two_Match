package models

import "time"

type ReactionType string

const (
	ReactionLike  ReactionType = "like"
	ReactionFire  ReactionType = "fire"
	ReactionShock ReactionType = "shock"
	ReactionSad   ReactionType = "sad"
	ReactionLaugh ReactionType = "laugh"
)

type Reaction struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	UserID uint64 `json:"user_id" gorm:"not null;uniqueIndex:idx_user_event;uniqueIndex:idx_user_commentary"`

	EventID      *uint64 `json:"event_id,omitempty" gorm:"uniqueIndex:idx_user_event"`
	CommentaryID *uint64 `json:"commentary_id,omitempty" gorm:"uniqueIndex:idx_user_commentary"`

	Type ReactionType `json:"type" gorm:"type:varchar(20);not null"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
