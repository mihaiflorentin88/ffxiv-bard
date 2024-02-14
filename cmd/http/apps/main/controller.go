package maincontroller

import (
	"embed"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"io/fs"
	"net/http"
)

type mainController struct {
	Router       *gin.Engine
	StaticFS     *embed.FS
	ErrorHandler contract.HttpError
	Renderer     contract.HttpRenderer
}

func NewMainController(router *gin.Engine, staticFs *embed.FS, errorHandler contract.HttpError, renderer contract.HttpRenderer) contract.Controller {

	return &mainController{
		Router:       router,
		StaticFS:     staticFs,
		ErrorHandler: errorHandler,
		Renderer:     renderer,
	}
}

func (m *mainController) Initialize() {
	m.enableStatic()
	m.index()
}

func (m *mainController) index() {
	m.Router.GET("/", func(c *gin.Context) {
		m.Renderer.AddTemplate("resource/template/main/content.html")
		m.Renderer.Render(c, nil)
	})
}

func (m *mainController) enableStatic() {
	cssFS, err := fs.Sub(m.StaticFS, "resource/css")
	if err != nil {
		panic("Cannot parse the css")
	}
	jsFS, err := fs.Sub(m.StaticFS, "resource/js")
	if err != nil {
		panic("Cannot parse the js")
	}
	imgFS, err := fs.Sub(m.StaticFS, "resource/img")
	if err != nil {
		panic("Cannot parse the img")
	}

	m.Router.StaticFS("/_resource/css", http.FS(cssFS))
	m.Router.StaticFS("/_resource/js", http.FS(jsFS))
	m.Router.StaticFS("/_resource/img", http.FS(imgFS))
}
