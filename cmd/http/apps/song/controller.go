package song

import (
	"ffxvi-bard/cmd/http/apps/song/form"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	Song            contract.SongInterface
	GenreRepository contract.GenreRepositoryInterface
	ErrorHandler    contract.HttpErrorHandlerInterface
	Renderer        contract.HttpRenderer
}

func NewSongController(song contract.SongInterface, errorHandler contract.HttpErrorHandlerInterface,
	renderer contract.HttpRenderer, genreRepository contract.GenreRepositoryInterface) *Controller {
	return &Controller{
		Song:            song,
		GenreRepository: genreRepository,
		ErrorHandler:    errorHandler,
		Renderer:        renderer,
	}
}

func (s *Controller) RenderSongList(c *gin.Context) {
	s.Renderer.
		AddTemplate("resource/template/song/list_songs.gohtml").
		Render(c, nil, http.StatusOK)
}

func (s *Controller) RenderAddNewSongForm(c *gin.Context) {
	formViewData, err := form.NewAddNewSongFormView(s.Song, s.GenreRepository)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusInternalServerError, c)
	}
	s.Renderer.
		AddTemplate("resource/template/song/add_song.gohtml").
		AddTemplate("resource/template/song/add_songl_js.gohtml").
		RemoveTemplate("resource/template/base/additional_js.gohtml").
		Render(c, formViewData, http.StatusOK)
}

func (s *Controller) HandleAddNewSong(c *gin.Context) {
	title := c.PostForm("title")
	artist := c.PostForm("artist")
	ensembleSize := c.PostForm("ensembleSize")
	genre := c.PostFormArray("genre")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	_ = title
	_ = artist
	_ = ensembleSize
	_ = genre
	_ = file
}
