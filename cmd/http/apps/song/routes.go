package song

import (
	"ffxvi-bard/cmd/http/utils/middleware"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Controller    *Controller
	JWTMiddleware middleware.AuthMiddleware
}

func NewSongRouter(controller *Controller, JWTMiddleware middleware.AuthMiddleware) contract.RouterInterface {
	return &Router{
		Controller:    controller,
		JWTMiddleware: JWTMiddleware,
	}
}

func (r Router) RegisterRoutes(router *gin.Engine) {
	router.GET("/song/list", r.Controller.RenderSongList)
	router.GET("/song/add", r.JWTMiddleware.UI(), func(c *gin.Context) {
		r.Controller.RenderAddNewSongForm(c)
	})
	router.POST("/song", r.Controller.HandleAddNewSong)
}
