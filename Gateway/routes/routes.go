package routes

import (
	"github.com/Dodge-git/Team_Two_Match/gateway/auth"
	"github.com/Dodge-git/Team_Two_Match/gateway/config"
	"github.com/Dodge-git/Team_Two_Match/gateway/proxy"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, cfg *config.Config) error {

	userProxy, _ := proxy.NewReverseProxy(cfg.UserServiceURL)
	matchProxy, _ := proxy.NewReverseProxy(cfg.MatchServiceURL)
	eventProxy, _ := proxy.NewReverseProxy(cfg.EventServiceURL)

	// ================= PUBLIC =================
	router.Any("/api/auth/*path", gin.WrapH(userProxy))

	// ================= PROTECTED =================
	api := router.Group("/api")
	api.Use(auth.AuthMiddleware(cfg.JWTSecret))

	// USER
	api.Any("/users/*path", gin.WrapH(userProxy))

	// MATCH SERVICE
	api.Any("/matches/*path", gin.WrapH(matchProxy))
	api.Any("/players/*path", gin.WrapH(matchProxy))
	api.Any("/teams/*path", gin.WrapH(matchProxy))
	api.Any("/sports/*path", gin.WrapH(matchProxy))

	// EVENT SERVICE
	api.Any("/events", gin.WrapH(eventProxy))
api.Any("/events/*path", gin.WrapH(eventProxy))

api.Any("/commentaries", gin.WrapH(eventProxy))
api.Any("/commentaries/*path", gin.WrapH(eventProxy))

api.Any("/comments", gin.WrapH(eventProxy))
api.Any("/comments/*path", gin.WrapH(eventProxy))

api.Any("/reactions", gin.WrapH(eventProxy))
api.Any("/reactions/*path", gin.WrapH(eventProxy))

	return nil
}

	 
