package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, handler *ReactionHandler) {
	api := r.Group("/api")

	reactions := api.Group("/reactions")
	{
		reactions.POST("", handler.SetReaction)
	}

	events := api.Group("/events")
	{
		events.GET("/:id/reactions/summary", handler.GetEventSummary)
		events.GET("/:id/my-reaction", handler.GetMyEventReaction)
	}

	commentary := api.Group("/commentary")
	{
		commentary.GET("/:id/reactions/summary", handler.GetCommentarySummary)
		commentary.GET("/:id/my-reaction", handler.GetMyCommentaryReaction)
	}
}
