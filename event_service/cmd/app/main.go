package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/mountainman199231/event_service/internal/client/match"
	"github.com/mountainman199231/event_service/internal/config"
	"github.com/mountainman199231/event_service/internal/handler"
	"github.com/mountainman199231/event_service/internal/repository"
	"github.com/mountainman199231/event_service/internal/service"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		logger.Error("REDIS_ADDR is not set")
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		PoolSize:     20,
		MinIdleConns: 5,
	})
	defer rdb.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Не удалось подключиться к redis:", err)
	}

	logger.Info("Redis подключен:", pong)

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

	commentService := service.NewCommentService(commentRepo, rdb)

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
