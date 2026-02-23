package handlers

import (
	"errors"
	"match_service/internal/dto"
	"match_service/internal/errs"
	"match_service/internal/models"
	"match_service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	teamService services.TeamService
}

func (h *TeamHandler) RegisterRoutes(router *gin.Engine) {
	teams := router.Group("/teams")
	{
		teams.POST("/", h.CreateTeam)
		teams.GET("/:id", h.GetTeamByID)
		teams.DELETE("/:id", h.DeleteTeam)
		teams.GET("/", h.ListTeams)
	}
}

func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	team, err := h.teamService.Create(req)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidSportID) || errors.Is(err, errs.ErrInvalidTeamName) || errors.Is(err, errs.ErrInvalidCity) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}
	team, err := h.teamService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrTeamNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, team)
}

func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team id"})
		return
	}
	_, err = h.teamService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, errs.ErrTeamNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	if err := h.teamService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "team deleted"})
}

func (h *TeamHandler) ListTeams(c *gin.Context) {
	var filter models.TeamFilter
	sportIDParam := c.Query("sport_id")

	if sportIDParam != "" {
		sportIDParam, err := strconv.ParseUint(sportIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sport id"})
			return
		}
		sportID := uint(sportIDParam)
		
			filter.SportID = &sportID
	
	}
	if pageParam := c.Query("page"); pageParam != "" {
		page, err := strconv.Atoi(pageParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
			return
		}
		filter.Page = page
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page_size"})
			return
		}
		filter.PageSize = pageSize
	}
	teams, total, err := h.teamService.List(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.TeamListResponse{
		SportID:  *filter.SportID,
		Data:     teams,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	})
}
