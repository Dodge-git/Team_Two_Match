package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	r *gin.Engine,
	reactionHandler *ReactionHandler,
	matchEventHandler *MatchEventHandler,
	commentHandler *CommentHandler,
	commentaryHandler *CommentaryHandler,
) {
	api := r.Group("/api/v1")

	events := api.Group("/events")
	{
		events.GET("/match/:match_id", matchEventHandler.GetMatchEvents)
		events.GET("/match/:match_id/timeline", matchEventHandler.GetMatchTimeline)

		events.POST("", matchEventHandler.CreateMatchEvent)
		events.GET("/:id", matchEventHandler.GetMatchEventByID)

		events.GET("/:id/reactions/summary", reactionHandler.GetEventSummary)
		events.GET("/:id/reactions/me", reactionHandler.GetMyEventReaction)

		events.GET("/:id/comments", commentHandler.GetByEvent)
	}

	commentaries := api.Group("/commentaries")
	{
		commentaries.POST("", commentaryHandler.Create)
		commentaries.GET("/match/:match_id", commentaryHandler.GetByMatch)

		commentaries.DELETE("/:id", commentaryHandler.Delete)
		commentaries.POST("/:id/pin", commentaryHandler.Pin)

		commentaries.GET("/:id/reactions/summary", reactionHandler.GetCommentarySummary)
		commentaries.GET("/:id/reactions/me", reactionHandler.GetMyCommentaryReaction)

		commentaries.GET("/:id/comments", commentHandler.GetByCommentary)
	}

	comments := api.Group("/comments")
	{
		comments.POST("", commentHandler.CreateComment)
		comments.PATCH("/:id", commentHandler.Update)
		comments.DELETE("/:id", commentHandler.Delete)
	}

	reactions := api.Group("/reactions")
	{
		reactions.POST("", reactionHandler.SetReaction)
	}
}
