package song

import (
	"embed"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

type songController struct {
	Router       *gin.Engine
	StaticFS     *embed.FS
	ErrorHandler contract.HttpError
	Renderer     contract.HttpRenderer
}

func NewSongController(router *gin.Engine, staticFS *embed.FS, errorHandler contract.HttpError, renderer contract.HttpRenderer) contract.Controller {
	return &songController{
		Router:       router,
		StaticFS:     staticFS,
		ErrorHandler: errorHandler,
		Renderer:     renderer,
	}
}

func (s *songController) Initialize() {
	s.showAddNewSongForm()
	s.HandleAddNewSong()
	s.ShowSongList()
}

func (s *songController) ShowSongList() {
	s.Router.GET("/song/list", func(c *gin.Context) {
		s.Renderer.
			AddTemplate("resource/template/song/list_songs.html").
			Render(c, nil)
	})
}

func (s *songController) showAddNewSongForm() {
	s.Router.GET("/song/add", func(c *gin.Context) {
		s.Renderer.
			AddTemplate("resource/template/song/add_song.html").
			AddTemplate("resource/template/song/add_songl_js.html").
			RemoveTemplate("resource/template/base/additional_js.html").
			Render(c, nil)
	})
}

func (s *songController) HandleAddNewSong() {
	s.Router.POST("/song", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
