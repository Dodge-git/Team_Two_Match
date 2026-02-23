package main

import (
	"log"
	"match_service/internal/config"
	"match_service/internal/handlers"
	"match_service/internal/models"
	"match_service/internal/repository"
	"match_service/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	db := config.SetUpDatabaseConnection()

	if err := db.AutoMigrate(&models.Sport{}, &models.Player{}, &models.Team{}, &models.Match{}); err != nil {
		log.Fatalf("не удалось выполнить миграции: %v", err)
	}

	sportRepo := repository.NewSportRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	matchRepo := repository.NewMatchRepository(db)

	sportService := services.NewSportService(sportRepo)
	teamService := services.NewTeamService(teamRepo, sportRepo)
	playerService := services.NewPlayerService(playerRepo)
	matchService := services.NewMatchService(matchRepo, sportRepo, teamRepo)

	router := gin.Default()

	handlers.RegisterRoutes(router, sportService, teamService, playerService, matchService)

	if err := router.Run(); err != nil {
		log.Fatalf("не удалось запустить HTTP-сервер: %v", err)
	}

}
