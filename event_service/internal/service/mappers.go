package service

import (
	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/models"
	"math"
)

func mapToMatchEventResponse(e *models.MatchEvent) *dto.MatchEventResponse {
	return &dto.MatchEventResponse{
		ID:          e.ID,
		MatchID:     e.MatchID,
		EventType:   string(e.EventType),
		Minute:      e.Minute,
		Period:      e.Period,
		TeamID:      e.TeamID,
		PlayerID:    e.PlayerID,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
	}
}

func mapCommentToDTO(c *models.Comment) *dto.CommentResponse {
	return &dto.CommentResponse{
		ID:           c.ID,
		UserID:       c.UserID,
		EventID:      c.EventID,
		CommentaryID: c.CommentaryID,
		Text:         c.Text,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		IsEdited:     !c.CreatedAt.Equal(c.UpdatedAt),
	}
}

func buildListResponse(comments []*models.Comment, total int64, page, pageSize int) *dto.CommentListResponse {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	items := make([]dto.CommentResponse, 0, len(comments))

	for _, c := range comments {
		items = append(items, *mapCommentToDTO(c))
	}

	return &dto.CommentListResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

func mapReactionToDTO(r *models.Reaction) *dto.ReactionResponse {
	return &dto.ReactionResponse{
		ID:           r.ID,
		UserID:       r.UserID,
		EventID:      r.EventID,
		CommentaryID: r.CommentaryID,
		Type:         string(r.Type),
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func buildReactionSummary(data map[string]int64) *dto.ReactionSummaryResponse {

	like := int(data["like"])
	fire := int(data["fire"])
	shock := int(data["shock"])
	sad := int(data["sad"])
	laugh := int(data["laugh"])

	return &dto.ReactionSummaryResponse{
		Like:  like,
		Fire:  fire,
		Shock: shock,
		Sad:   sad,
		Laugh: laugh,
		Total: like + fire + shock + sad + laugh,
	}

}
