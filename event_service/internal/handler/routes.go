package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine,
	reactionHandler *ReactionHandler,
	matchEventHandler *MatchEventHandler,
) {
	api := r.Group("/api")

	events := api.Group("/events")
	{
		events.POST("", matchEventHandler.CreateMatchEvent)
		events.GET("/:id", matchEventHandler.GetMatchEventByID)
		events.GET("/match/:match_id", matchEventHandler.GetMatchEvents)
		events.GET("/match/:match_id/timeline", matchEventHandler.GetMatchTimeline)
	}

	reactions := api.Group("/reactions")
	{
		reactions.POST("", reactionHandler.SetReaction)
	}

	events.GET("/:id/reactions/summary", reactionHandler.GetEventSummary)
	events.GET("/:id/reactions/me", reactionHandler.GetMyEventReaction)

	commentaries := api.Group("/commentaries")
	{
		commentaries.GET("/:id/reactions/summary", reactionHandler.GetCommentarySummary)
		commentaries.GET("/:id/reactions/me", reactionHandler.GetMyCommentaryReaction)
	}
}
