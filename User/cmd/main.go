package main

import (
	"User/internal/config"
	"User/internal/repository"
	"User/internal/services"
	"User/internal/transport"
	"github.com/gin-gonic/gin"
)

func main() {

	db := config.SetUpDatabaseConnection()

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := transport.NewUserHandler(userService)

	r := gin.Default()

// AUTH
auth := r.Group("/auth")
{
	auth.POST("/register", userHandler.Register)
	auth.POST("/login", userHandler.Login)
}

// USERS
users := r.Group("/users")
{
	users.GET("/:id", userHandler.GetProfile)
	users.PUT("/:id", userHandler.UpdateUser)
	users.DELETE("/:id", userHandler.DeleteUser)
}

r.Run(":8000")
}