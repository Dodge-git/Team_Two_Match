package handlers

import (
	"match_service/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, sportService services.SportService, teamService services.TeamService, playerService services.PlayerService, matchService services.MatchService) {

	sportHandler := NewSportHandler(sportService)
	playerHandler := NewPlayerHandler(playerService)
	teamHandler := NewTeamHandler(teamService)
	matchHandler := NewMatchHandler(matchService)

	sportHandler.RegisterRoutes(router)
	teamHandler.RegisterRoutes(router)
	playerHandler.RegisterRoutes(router)
	matchHandler.RegisterRoutes(router)
}
