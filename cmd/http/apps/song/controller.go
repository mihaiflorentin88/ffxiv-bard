package song

import (
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	ErrorHandler contract.HttpErrorHandlerInterface
	Renderer     contract.HttpRenderer
}

func NewSongController(errorHandler contract.HttpErrorHandlerInterface, renderer contract.HttpRenderer) *Controller {
	return &Controller{
		ErrorHandler: errorHandler,
		Renderer:     renderer,
	}
}

func (s *Controller) RenderSongList(c *gin.Context) {
	s.Renderer.
		AddTemplate("resource/template/song/list_songs.html").
		Render(c, nil)
}

func (s *Controller) RenderAddNewSongForm(c *gin.Context) {
	s.Renderer.
		AddTemplate("resource/template/song/add_song.html").
		AddTemplate("resource/template/song/add_songl_js.html").
		RemoveTemplate("resource/template/base/additional_js.html").
		Render(c, nil)
}

func (s *Controller) HandleAddNewSong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
