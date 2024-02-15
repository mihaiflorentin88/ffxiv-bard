package auth

import (
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Controller *Controller
}

func NewAuthRouter(controller *Controller) contract.RouterInterface {
	return &Router{
		Controller: controller,
	}
}

func (r Router) EnableRoutes(router *gin.Engine) {
	router.GET("/auth/login", r.Controller.RenderLoginPage)
	router.GET("/auth/login/discord", r.Controller.LoginDiscord)
	router.GET("/auth/login/discord/callback", r.Controller.LoginDiscordCallback)
}
