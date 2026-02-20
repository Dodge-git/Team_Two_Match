package dto

import "match_service/internal/models"

type CreateSportRequest struct {
	Name models.SportName `json:"name" binding:"required"`
}

type SportResponse struct {
	Name models.SportName `json:"name"`
}

