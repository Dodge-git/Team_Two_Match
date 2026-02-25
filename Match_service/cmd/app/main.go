package main

import (
	"context"
	"log"
	"match_service/internal/config"
	"match_service/internal/handlers"
	"match_service/internal/kafka"
	"match_service/internal/models"
	"match_service/internal/repository"
	"match_service/internal/services"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

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

	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	producer := kafka.NewProducer(brokers)

	matchService := services.NewMatchService(matchRepo, sportRepo, teamRepo, producer)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewConsumer(
		brokers,
		"match.goal",
		"match-service-group",
		matchService,
	)
	go consumer.Start(ctx)

	router := gin.Default()

	handlers.RegisterRoutes(router, sportService, teamService, playerService, matchService)

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down...")

	cancel()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	log.Println("Match Service stopped")

}
