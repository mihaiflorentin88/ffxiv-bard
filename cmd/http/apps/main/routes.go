package maincontroller

import (
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Controller *Controller
}

func NewMainRouter(controller *Controller) contract.RouterInterface {
	return &Router{
		Controller: controller,
	}
}

func (r Router) RegisterRoutes(router *gin.Engine) {
	router.GET("/", r.Controller.RenderIndex)
}
