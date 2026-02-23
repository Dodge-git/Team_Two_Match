package handlers

import (
	"errors"
	"match_service/internal/dto"
	"match_service/internal/errs"
	"match_service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PlayerHandler struct {
	playerService services.PlayerService
}

func NewPlayerHandler(playerService services.PlayerService) *PlayerHandler {
	return &PlayerHandler{playerService: playerService}
}
func (h *PlayerHandler) RegisterRoutes(router *gin.Engine) {
	players := router.Group("/players")
	{
		players.POST("/", h.CreatePlayer)
		players.PATCH("/:id", h.UpdatePlayer)
		players.DELETE("/:id", h.DeletePlayer)
		players.GET("/:id", h.GetPlayerByID)
		players.GET("/team/:team_id", h.ListPlayersByTeam)
	}
}

func (h *PlayerHandler) CreatePlayer(c *gin.Context) {
	var req dto.CreatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player, err := h.playerService.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, player)
}

func (h *PlayerHandler) UpdatePlayer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	var req dto.UpdatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player, err := h.playerService.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) DeletePlayer(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	if err := h.playerService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "player deleted "})
}

func (h *PlayerHandler) GetPlayerByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	player, err := h.playerService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrPlayerNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) ListPlayersByTeam(c *gin.Context) {
	teamIDParam := c.Param("team_id")
	teamID, err := strconv.ParseUint(teamIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}
	players, err := h.playerService.List(uint(teamID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}
