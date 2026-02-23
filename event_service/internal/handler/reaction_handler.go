package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mountainman199231/event_service/internal/dto"
	"github.com/mountainman199231/event_service/internal/service"
)

type ReactionHandler struct {
	service service.ReactionService
}

func NewReactionHandler(service service.ReactionService) *ReactionHandler {
	return &ReactionHandler{service: service}
}

func getUserID(c *gin.Context) (uint64, bool) {
	/*
		val, exists := c.Get("user_id")
		if !exists {
			return 0, false
		}

		switch v := val.(type) {
		case uint64:
			return v, true
		case float64:
			return uint64(v), true
		default:
			return 0, false
		}
	*/
	return 1, true
}

func (h *ReactionHandler) SetReaction(c *gin.Context) {
	var req dto.SetReactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	response, err := h.service.SetReaction(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if response == nil {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ReactionHandler) GetEventSummary(c *gin.Context) {
	idParam := c.Param("id")
	eventID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid event id",
		})
		return
	}

	resp, err := h.service.GetEventSummary(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ReactionHandler) GetCommentarySummary(c *gin.Context) {
	idParam := c.Param("id")
	commentaryID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid commentary id",
		})
		return
	}

	resp, err := h.service.GetCommentarySummary(c.Request.Context(), commentaryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ReactionHandler) GetMyEventReaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	eventID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid event id",
		})
		return
	}

	resp, err := h.service.GetUserEventReaction(c.Request.Context(), userID, eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ReactionHandler) GetMyCommentaryReaction(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	idParam := c.Param("id")
	commentaryID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid commentary id",
		})
		return
	}

	resp, err := h.service.GetUserCommentaryReaction(c.Request.Context(), userID, commentaryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
