package dto

import "match_service/internal/models"

type CreateSportRequest struct {
	Name models.SportName `json:"name" binding:"required"`
}

type SportResponse struct {
	ID   uint             `json:"id"`
	Name models.SportName `json:"name"`
}
