package dto

import "time"

type CreateCommentaryRequest struct {
	MatchID uint64 `json:"match_id" binding:"required"`
	Minute  int    `json:"minute" binding:"required,min=0"`

	Text string `json:"text" binding:"required,min=1,max=2000"`
}

type UpdateCommentaryRequest struct {
	Text *string `json:"text,omitempty" binding:"omitempty,min=1,max=2000"`
}

type CommentaryResponse struct {
	ID      uint64 `json:"id"`
	MatchID uint64 `json:"match_id"`

	Minute int    `json:"minute"`
	Text   string `json:"text"`

	IsPinned bool `json:"is_pinned"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
