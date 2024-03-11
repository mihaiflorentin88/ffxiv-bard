package auth

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	Controller Controller
}

func NewAuthRouter(controller Controller) *Router {
	return &Router{
		Controller: controller,
	}
}

func (r Router) RegisterRoutes(router *gin.Engine) {
	router.GET("/auth/login", r.Controller.RenderLoginPage)
	router.GET("/auth/logout", r.Controller.Logout)
	router.GET("/auth/login/discord", r.Controller.LoginWithDiscord)
	router.GET("/auth/login/discord/callback", r.Controller.LoginWithDiscordCallback)
}
