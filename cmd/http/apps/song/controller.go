package song

import (
	"errors"
	"ffxvi-bard/cmd/http/apps/song/form"
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	Song                 *song.Song
	ErrorHandler         contract.HttpErrorHandlerInterface
	Renderer             contract.HttpRenderer
	addSongFormProcessor form.AddSongFormProcessor
	songListForm         form.SongList
	newSongFormView      form.NewSongFormView
}

func NewSongController(song *song.Song, errorHandler contract.HttpErrorHandlerInterface, renderer contract.HttpRenderer,
	addSongFormProcessor form.AddSongFormProcessor, songListForm form.SongList, newSongFormView form.NewSongFormView) *Controller {
	return &Controller{
		Song:                 song,
		ErrorHandler:         errorHandler,
		Renderer:             renderer,
		addSongFormProcessor: addSongFormProcessor,
		songListForm:         songListForm,
		newSongFormView:      newSongFormView,
	}
}

func (s *Controller) RenderSongList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	title := c.Query("title")
	artist := c.Query("artist")
	sort := c.Query("sort")
	ensembleSize, err := strconv.Atoi(c.Query("ensembleSize"))
	if err != nil {
		ensembleSize = -1
	}
	genre, err := strconv.Atoi(c.Query("genre"))
	if err != nil {
		genre = -1
	}
	songListForm, err := s.songListForm.Fetch(title, artist, ensembleSize, genre, page, sort)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusInternalServerError, c)
		return
	}
	s.Renderer.
		StartClean().
		RemoveTemplate("resource/template/base/additional_js.gohtml").
		RemoveTemplate("resource/template/base/additional_styles.gohtml").
		AddTemplate("resource/template/song/list.gohtml").
		AddTemplate("resource/template/song/list_css.gohtml").
		AddTemplate("resource/template/song/list_js.gohtml").
		Render(c, songListForm, http.StatusOK)
}

func (s *Controller) RenderAddNewSongForm(c *gin.Context) {
	newSongForm, err := s.newSongFormView.Fetch()
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusInternalServerError, c)
		return
	}
	s.Renderer.
		StartClean().
		RemoveTemplate("resource/template/base/additional_styles.gohtml").
		AddTemplate("resource/template/song/add_song.gohtml").
		AddTemplate("resource/template/song/add_song_css.gohtml").
		Render(c, newSongForm, http.StatusOK)
}

func (s *Controller) HandleAddNewSong(c *gin.Context) {
	title := c.PostForm("title")
	artist := c.PostForm("artist")
	ensembleSize := c.PostForm("ensembleSize")
	genre := c.PostFormArray("genre")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	sessionUser, _ := c.Get("user")
	loggedUser, err := user.FromSession(sessionUser)
	if err != nil {
		err = errors.New(fmt.Sprintf("cannot load user: Reason %s. Try relogging", err))
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	err = s.addSongFormProcessor.Process(title, artist, ensembleSize, genre, fileHeader, loggedUser)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	c.Redirect(http.StatusFound, "/song/add")
}
