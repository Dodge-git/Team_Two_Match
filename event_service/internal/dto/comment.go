package dto

import "time"

type CreateCommentRequest struct {
	EventID      *uint64 `json:"event_id,omitempty"`
	CommentaryID *uint64 `json:"commentary_id,omitempty"`

	Text string `json:"text" binding:"required,min=1,max=2000"`
}

type UpdateCommentRequest struct {
	Text *string `json:"text,omitempty" binding:"omitempty,min=1,max=2000"`
}

type CommentResponse struct {
	ID     uint64 `json:"id"`
	UserID uint64 `json:"user_id"`

	EventID      *uint64 `json:"event_id,omitempty"`
	CommentaryID *uint64 `json:"commentary_id,omitempty"`

	Text string `json:"text"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsEdited bool `json:"is_edited"`
}

type CommentListResponse struct {
	Items      []CommentResponse `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}
