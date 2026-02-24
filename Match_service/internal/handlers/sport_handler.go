package handlers

import (
	"match_service/internal/dto"
	"match_service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SportHandler struct {
	sportService services.SportService
}

func NewSportHandler(sportService services.SportService) *SportHandler {
	return &SportHandler{sportService: sportService}
}

func (h *SportHandler) RegisterRoutes(router *gin.Engine) {
	sports := router.Group("/sports")
	{
		sports.POST("/", h.CreateSport)
		sports.GET("/", h.ListSports)
	}
}

func (h *SportHandler) CreateSport(c *gin.Context) {
	var req dto.CreateSportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sport, err := h.sportService.Create(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, &dto.SportResponse{
		ID:   sport.ID,
		Name: sport.Name,
	})
}

func (h *SportHandler) ListSports(c *gin.Context) {
	list, err := h.sportService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var dtoList []dto.SportResponse
	for _, sport := range list {
		sportResponse := dto.SportResponse{
			ID:   sport.ID,
			Name: sport.Name,
		}
		dtoList = append(dtoList, sportResponse)
	}
	c.JSON(http.StatusOK, dtoList)
}
