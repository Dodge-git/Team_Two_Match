package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/service"
)

type MatchEventHandler struct {
	service service.MatchEventService
}

func NewMatchEventHandler(service service.MatchEventService) *MatchEventHandler {
	return &MatchEventHandler{service: service}
}

func (h *MatchEventHandler) CreateMatchEvent(c *gin.Context) {
	var req dto.CreateMatchEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resp, err := h.service.CreateMatchEvent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *MatchEventHandler) GetMatchEventByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	resp, err := h.service.GetMatchEventByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not found",
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MatchEventHandler) GetMatchEvents(c *gin.Context) {
	idParam := c.Param("match_id")
	matchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	resp, err := h.service.GetMatchEvents(c.Request.Context(), matchID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *MatchEventHandler) GetMatchTimeline(c *gin.Context) {
	idParam := c.Param("match_id")
	matchID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	resp, err := h.service.GetMatchTimeline(c.Request.Context(), matchID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
