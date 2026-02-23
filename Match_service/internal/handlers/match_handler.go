package handlers

import (
	"errors"
	"match_service/internal/dto"
	"match_service/internal/errs"
	"match_service/internal/models"
	services "match_service/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MatchHandler struct {
	matchService services.MatchService
}

func NewMatchHandler(matchService services.MatchService) *MatchHandler {
	return &MatchHandler{matchService: matchService}
}

func (h *MatchHandler) RegisterRoutes(router *gin.Engine) {
	matches := router.Group("/matches")
	{
		matches.POST("/", h.CreateMatch)
		matches.GET("/", h.ListMatches)
		matches.GET("/:id", h.GetMatchByID)
		matches.DELETE("/:id", h.DeleteMatch)
		matches.PATCH("/finish/:id", h.FinishMatch)
		matches.PATCH("/start/:id", h.StartMatch)
		matches.PATCH("/cancel/:id", h.CancelMatch)
		matches.POST("/goal", h.GoalEvent)
		matches.GET("/active", h.GetActiveMatches)
	}
}

func (h *MatchHandler) CreateMatch(c *gin.Context) {
	var req dto.CreateMatchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	match, err := h.matchService.CreateMatch(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, match)
}

func (h *MatchHandler) GetMatchByID(c *gin.Context) {
	idParam := c.Param("id")
	matchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	match, err := h.matchService.GetMatchByID(uint(matchID))
	if err != nil {
		if errors.Is(err, errs.ErrMatchNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, match)
}

func (h *MatchHandler) DeleteMatch(c *gin.Context) {
	idParam := c.Param("id")
	matchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	err = h.matchService.DeleteMatch(uint(matchID))
	if err != nil {
		if errors.Is(err, errs.ErrMatchNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match deleted"})
}

func (h *MatchHandler) ListMatches(c *gin.Context) {
	var filter models.MatchFilter
	if sportIDParam := c.Query("sport_id"); sportIDParam != "" {
		sportID, err := strconv.ParseUint(sportIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sport id"})
			return
		}
		sportIDUint := uint(sportID)
		filter.SportID = &sportIDUint
	}
	if statusStr := c.Query("status"); statusStr != "" {
		status := models.MatchStatus(statusStr)
		filter.Status = &status
	}

	if dateFromSrt := c.Query("date_from"); dateFromSrt != "" {
		dateFromSrt, err := time.Parse("2006-01-02", dateFromSrt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_from format"})
			return
		}
		filter.DateFrom = &dateFromSrt
	}
	if dateToStr := c.Query("date_to"); dateToStr != "" {
		dateToStr, err := time.Parse("2006-01-02", dateToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date_to format"})
			return
		}
		filter.DateTo = &dateToStr
	}
	if pageStr := c.Query("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
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

	resp, err := h.matchService.ListMatches(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MatchHandler) StartMatch(c *gin.Context) {
	idParam := c.Param("id")
	matchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	if err = h.matchService.StartMatch(uint(matchID)); err != nil {
		if errors.Is(err, errs.ErrMatchNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, errs.ErrInvalidMatchStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match started"})
}
func (h *MatchHandler) FinishMatch(c *gin.Context) {
	idParam := c.Param("id")
	matchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	if err = h.matchService.FinishMatch(uint(matchID)); err != nil {
		if errors.Is(err, errs.ErrMatchNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, errs.ErrInvalidMatchStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match finished"})
}

func (h *MatchHandler) CancelMatch(c *gin.Context) {
	idParam := c.Param("id")
	matchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}

	if err = h.matchService.CancelMatch(uint(matchID)); err != nil {
		if errors.Is(err, errs.ErrMatchNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, errs.ErrInvalidMatchStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match canceled"})
}

func (h *MatchHandler) GoalEvent(c *gin.Context) {



}

func (h *MatchHandler) GetActiveMatches(c *gin.Context) {

	matches, err := h.matchService.GetActive()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, matches)
}
