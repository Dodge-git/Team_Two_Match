package main

import (
	"log"

	"github.com/Dodge-git/Team_Two_Match/gateway/config"
	"github.com/Dodge-git/Team_Two_Match/gateway/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	
	cfg := config.Load()

	r := gin.Default()

	if err := routes.RegisterRoutes(r,cfg); err != nil{
		log.Fatal(err)
	}

	r.Run(":8080")
}