package maincontroller

import (
	"ffxvi-bard/cmd/http/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	ErrorHandler utils.ErrorHandler
	Renderer     utils.Renderer
}

func NewMainController(errorHandler utils.ErrorHandler, renderer utils.Renderer) Controller {
	return Controller{
		ErrorHandler: errorHandler,
		Renderer:     renderer,
	}
}

func (m *Controller) RenderIndex(c *gin.Context) {
	m.Renderer.AddTemplate("resource/template/main/content.gohtml")
	m.Renderer.Render(c, nil, http.StatusOK)
}
