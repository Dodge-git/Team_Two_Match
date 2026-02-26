package routes

import (
	"github.com/Dodge-git/Team_Two_Match/gateway/auth"
	"github.com/Dodge-git/Team_Two_Match/gateway/config"
	"github.com/Dodge-git/Team_Two_Match/gateway/proxy"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine,cfg *config.Config)error{

	userProxy , err := proxy.NewReverseProxy(cfg.UserServiceURL)
	if err != nil{
		return err
	}
	matchProxy, err := proxy.NewReverseProxy(cfg.MatchServiceURL)
	if err != nil{
		return err
	}
	eventProxy, err := proxy.NewReverseProxy(cfg.EventServiceURL)
	if err != nil{
		return err
	}

	// ================= PUBLIC ROUTES =================

	router.Any("/api/auth/*path", gin.WrapH(userProxy))

	// ================= PROTECTED ROUTES =================

	api := router.Group("/api")
	api.Use(auth.AuthMiddleware(cfg.JWTSecret))

	// ===================================================
	// ================= USER SERVICE ====================
	// ===================================================

	users := api.Group("/users")
	{
		users.GET("/:id", gin.WrapH(userProxy))

		// Удаление профиля (сервис сам проверяет ownership)
		users.DELETE("/:id",
			auth.RBACMiddleware(auth.DeleteProfile),
			gin.WrapH(userProxy),
		)
	}

	// ===================================================
	// ================= MATCH SERVICE ===================
	// ===================================================

	// ---------- MATCHES ----------
	matches := api.Group("/matches")
	{
		// Чтение
		matches.GET("/", gin.WrapH(matchProxy))
		matches.GET("/:id", gin.WrapH(matchProxy))
		matches.GET("/active", gin.WrapH(matchProxy))

		// Создание матча
		matches.POST("/",
			auth.RBACMiddleware(auth.CreateMatch),
			gin.WrapH(matchProxy),
		)

		// Управление статусом
		matches.PATCH("/start/:id",
			auth.RBACMiddleware(auth.ManageMatch),
			gin.WrapH(matchProxy),
		)

		matches.PATCH("/finish/:id",
			auth.RBACMiddleware(auth.ManageMatch),
			gin.WrapH(matchProxy),
		)

		matches.PATCH("/cancel/:id",
			auth.RBACMiddleware(auth.ManageMatch),
			gin.WrapH(matchProxy),
		)

		// Гол
		matches.POST("/goal",
			auth.RBACMiddleware(auth.AdGoal),
			gin.WrapH(matchProxy),
		)

		// Удаление (только admin)
		matches.DELETE("/:id",
			auth.RBACMiddleware(auth.DeleteMatch),
			gin.WrapH(matchProxy),
		)
	}

	// ---------- PLAYERS ----------
	players := api.Group("/players")
	{
		players.GET("/:id", gin.WrapH(matchProxy))
		players.GET("/team/:team_id", gin.WrapH(matchProxy))

		players.POST("/",
			auth.RBACMiddleware(auth.AdPlayer),
			gin.WrapH(matchProxy),
		)

		players.PATCH("/:id",
			auth.RBACMiddleware(auth.AdPlayer),
			gin.WrapH(matchProxy),
		)

		players.DELETE("/:id",
			auth.RBACMiddleware(auth.AdPlayer),
			gin.WrapH(matchProxy),
		)
	}

	// ---------- TEAMS ----------
	teams := api.Group("/teams")
	{
		teams.GET("/", gin.WrapH(matchProxy))
		teams.GET("/:id", gin.WrapH(matchProxy))

		teams.POST("/",
			auth.RBACMiddleware(auth.AdTeam),
			gin.WrapH(matchProxy),
		)

		teams.DELETE("/:id",
			auth.RBACMiddleware(auth.AdTeam),
			gin.WrapH(matchProxy),
		)
	}

	// ---------- SPORTS ----------
	sports := api.Group("/sports")
	{
		sports.GET("/", gin.WrapH(matchProxy))

		sports.POST("/",
			auth.RBACMiddleware(auth.AdSport),
			gin.WrapH(matchProxy),
		)
	}

	// ===================================================
	// ================= EVENT SERVICE ===================
	// ===================================================

	events := api.Group("/events")
	{
		events.GET("/*path", gin.WrapH(eventProxy))

		events.POST("",
			auth.RBACMiddleware(auth.AdEvent),
			gin.WrapH(eventProxy),
		)

		events.DELETE("/:id",
			auth.RBACMiddleware(auth.DeleteEvent),
			gin.WrapH(eventProxy),
		)
	}

	commentaries := api.Group("/commentaries")
	{
		commentaries.GET("/*path", gin.WrapH(eventProxy))

		commentaries.POST("",
			auth.RBACMiddleware(auth.AdCommentary),
			gin.WrapH(eventProxy),
		)

		commentaries.DELETE("/:id",
			auth.RBACMiddleware(auth.DeleteCommentary),
			gin.WrapH(eventProxy),
		)
	}

	comments := api.Group("/comments")
	{
		comments.POST("",
			auth.RBACMiddleware(auth.AdComment),
			gin.WrapH(eventProxy),
		)

		comments.PATCH("/:id",
			auth.RBACMiddleware(auth.AdComment),
			gin.WrapH(eventProxy),
		)

		comments.DELETE("/:id",
			auth.RBACMiddleware(auth.DeleteComment),
			gin.WrapH(eventProxy),
		)
	}

	reactions := api.Group("/reactions")
	{
		reactions.POST("",
			gin.WrapH(eventProxy),
		)
	}

	return nil
}



	 
