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
	songDetailsForm      form.SongDetails
}

func NewSongController(song *song.Song, errorHandler contract.HttpErrorHandlerInterface, renderer contract.HttpRenderer,
	addSongFormProcessor form.AddSongFormProcessor, songListForm form.SongList, newSongFormView form.NewSongFormView,
	songDetailsForm form.SongDetails) *Controller {
	return &Controller{
		Song:                 song,
		ErrorHandler:         errorHandler,
		Renderer:             renderer,
		addSongFormProcessor: addSongFormProcessor,
		songListForm:         songListForm,
		newSongFormView:      newSongFormView,
		songDetailsForm:      songDetailsForm,
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
	songID, err := s.addSongFormProcessor.Process(title, artist, ensembleSize, genre, fileHeader, loggedUser)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/song/%v", songID))
}

func (s *Controller) SongDetails(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	songDetails, err := s.songDetailsForm.Fetch(songID, c)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	s.Renderer.
		StartClean().
		RemoveTemplate("resource/template/base/additional_styles.gohtml").
		RemoveTemplate("resource/template/base/additional_js.gohtml").
		AddTemplate("resource/template/song/song_details.gohtml").
		AddTemplate("resource/template/song/song_details_css.gohtml").
		AddTemplate("resource/template/song/song_details_js.gohtml").
		Render(c, songDetails, http.StatusOK)
}

func (s *Controller) DownloadSong(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	c.Header("Content-Type", "audio/midi")
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	disposition := c.DefaultQuery("disposition", "inline")

	copy := *s.Song
	targetSong := &copy
	_, err = targetSong.LoadByID(songID)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	if targetSong.GetStatus() != song.Processed {
		s.ErrorHandler.RenderTemplate(errors.New("cannot download file. the song is still being processed"), http.StatusBadRequest, c)
		return
	}
	filepath := targetSong.GetFilePath()
	if disposition == "attachment" {
		c.Header("Content-Disposition",
			fmt.Sprintf("attachment; filename=[%s]%s_%s.mid",
				targetSong.EnsembleString(),
				targetSong.Artist,
				targetSong.Title))
	}
	c.File(filepath)
}
