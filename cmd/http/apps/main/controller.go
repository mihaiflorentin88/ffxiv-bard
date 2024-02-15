package maincontroller

import (
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	ErrorHandler contract.HttpErrorHandlerInterface
	Renderer     contract.HttpRenderer
}

func NewMainController(errorHandler contract.HttpErrorHandlerInterface, renderer contract.HttpRenderer) *Controller {
	return &Controller{
		ErrorHandler: errorHandler,
		Renderer:     renderer,
	}
}

func (m *Controller) RenderIndex(c *gin.Context) {
	m.Renderer.AddTemplate("resource/template/main/content.gohtml")
	m.Renderer.Render(c, nil, http.StatusOK)
}
