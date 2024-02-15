package song

import (
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Controller *Controller
}

func NewSongRouter(controller *Controller) contract.RouterInterface {
	return &Router{
		Controller: controller,
	}
}

func (r Router) EnableRoutes(router *gin.Engine) {
	router.GET("/song/list", r.Controller.RenderSongList)
	router.GET("/song/add", r.Controller.RenderAddNewSongForm)
	router.POST("/song", r.Controller.HandleAddNewSong)
}
