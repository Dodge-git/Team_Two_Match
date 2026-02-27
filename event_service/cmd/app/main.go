package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/mountainman199231/event_service/internal/client/match"
	"github.com/mountainman199231/event_service/internal/config"
	"github.com/mountainman199231/event_service/internal/handler"
	"github.com/mountainman199231/event_service/internal/repository"
	"github.com/mountainman199231/event_service/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// config.SetEnv(logger)

	db, err := config.DatabaseConnect(logger)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}

	matchEventRepo := repository.NewMatchEventRepository(db)
	commentaryRepo := repository.NewCommentaryRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	reactionRepo := repository.NewReactionRepository(db)

	matchClient := match.NewClient(os.Getenv("MATCH_SERVICE_URL"))

	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		logger.Error("KAFKA_BROKER is not set")
		os.Exit(1)
	}
	brokers := []string{broker}
	producer := service.NewKafkaProducer(brokers)
	defer producer.Close()

	matchEventService := service.NewMatchEventService(
		matchEventRepo,
		commentaryRepo,
		reactionRepo,
		producer,
		matchClient,
	)

	reactionService := service.NewReactionService(
		reactionRepo,
		matchEventRepo,
		commentaryRepo,
	)

	commentService := service.NewCommentService(commentRepo)

	commentaryService := service.NewCommentaryService(
		commentaryRepo,
		db,
		matchClient,
	)

	matchEventHandler := handler.NewMatchEventHandler(matchEventService)
	reactionHandler := handler.NewReactionHandler(reactionService)
	commentHandler := handler.NewCommentHandler(commentService)
	commentaryHandler := handler.NewCommentaryHandler(commentaryService)

	router := gin.Default()

	handler.RegisterRoutes(
		router,
		reactionHandler,
		matchEventHandler,
		commentHandler,
		commentaryHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	logger.Info("server started", "port", port)

	if err := router.Run(":" + port); err != nil && err != http.ErrServerClosed {
		logger.Error("server failed", "error", err)
		os.Exit(1)
	}
}
