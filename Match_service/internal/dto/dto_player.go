package dto

type CreatePlayerRequest struct {
	TeamID   uint   `json:"team_id" binding:"required,gte=1"`
	Name     string `json:"name" binding:"required,min=1,max=100"`
	Number   uint   `json:"number" binding:"required,gte=1"`
	Position string `json:"position" binding:"required,max=30"`
}

type UpdatePlayerRequest struct {
	Name     *string `json:"name" binding:"omitempty,min=1,max=100"`
	Number   *uint   `json:"number" binding:"omitempty,gte=1"`
	Position *string `json:"position" binding:"omitempty"`
}
