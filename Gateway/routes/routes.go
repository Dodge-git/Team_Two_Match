package routes

import (
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

	// Тут маршруты Публиные
	public := router.Group("/")
	 {
        public.Any("/")
	 }




	 
}