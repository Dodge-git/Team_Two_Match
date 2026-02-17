package dto

import "time"

type SetReactionRequest struct {
	EventID      *uint64 `json:"event_id,omitempty"`
	CommentaryID *uint64 `json:"commentary_id,omitempty"`

	Type string `json:"type" binding:"required,oneof=like fire shock sad laugh"`
}

type ReactionResponse struct {
	ID uint64 `json:"id"`

	UserID uint64 `json:"user_id"`

	EventID      *uint64 `json:"event_id,omitempty"`
	CommentaryID *uint64 `json:"commentary_id,omitempty"`

	Type string `json:"type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReactionSummaryResponse struct {
	Like  int `json:"like"`
	Fire  int `json:"fire"`
	Shock int `json:"shock"`
	Sad   int `json:"sad"`
	Laugh int `json:"laugh"`

	Total int `json:"total"`
}

type UserReactionResponse struct {
	Type *string `json:"type,omitempty"`
}
