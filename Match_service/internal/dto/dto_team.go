package dto

import (
	"match_service/internal/models"
)

type CreateTeamRequest struct {
	SportID   uint   `json:"sport_id" binding:"required,gt=0"`
	Name      string `json:"name" binding:"required,min=2,max=100"`
	ShortName string `json:"short_name" binding:"required,min=2,max=16"`
	City      string `json:"city" binding:"required,min=2,max=100"`
}
type TeamListResponse struct {
	Data       []models.Team `json:"data"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
}

