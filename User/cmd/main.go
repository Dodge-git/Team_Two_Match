package main

import (
	"User/internal/config"
	"User/internal/repository"
	"User/internal/services"
	"User/internal/transport"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	db := config.SetUpDatabaseConnection()

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := transport.NewUserHandler(userService)

	r := gin.Default()

	api := r.Group("/api")
	{
		// AUTH
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
		}

		// USERS
		users := api.Group("/users")
		{
			users.GET("/:id", userHandler.GetProfile)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	log.Println("User service running on :8000")
	r.Run(":8000")
}