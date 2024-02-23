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
	addSongFormProcessor form.SubmitSongForm
	songListForm         form.SongList
	newSongFormView      form.NewSongFormView
	songDetailsForm      form.SongDetails
	submitSongRatingForm form.SubmitSongRatingForm
	submitCommentForm    form.SubmitCommentForm
}

func NewSongController(song *song.Song, errorHandler contract.HttpErrorHandlerInterface, renderer contract.HttpRenderer,
	addSongFormProcessor form.SubmitSongForm, songListForm form.SongList, newSongFormView form.NewSongFormView,
	songDetailsForm form.SongDetails, submitSongRatingForm form.SubmitSongRatingForm, submitCommentForm form.SubmitCommentForm) *Controller {
	return &Controller{
		Song:                 song,
		ErrorHandler:         errorHandler,
		Renderer:             renderer,
		addSongFormProcessor: addSongFormProcessor,
		songListForm:         songListForm,
		newSongFormView:      newSongFormView,
		songDetailsForm:      songDetailsForm,
		submitSongRatingForm: submitSongRatingForm,
		submitCommentForm:    submitCommentForm,
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
	songID, err := s.addSongFormProcessor.Submit(title, artist, ensembleSize, genre, fileHeader, loggedUser)
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

func (s *Controller) SubmitSongRating(c *gin.Context) {
	sessionUser, exists := c.Get("user")
	if !exists {
		s.ErrorHandler.RenderTemplate(errors.New("you must be logged in order to rate"), http.StatusBadRequest, c)
		return
	}
	loggedUser, err := user.FromSession(sessionUser)
	songID, err := strconv.Atoi(c.Param("songID"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the song id has to be an integer"), http.StatusBadRequest, c)
		return
	}
	rating, err := strconv.Atoi(c.PostForm("rating"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the rating has to be an integer"), http.StatusBadRequest, c)
		return
	}
	err = s.submitSongRatingForm.Submit(loggedUser, songID, rating)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/song/%v", songID))
}

func (s *Controller) SubmitSongComment(c *gin.Context) {
	sessionUser, exists := c.Get("user")
	if !exists {
		s.ErrorHandler.RenderTemplate(errors.New("you must be logged in order to rate"), http.StatusBadRequest, c)
		return
	}
	loggedUser, err := user.FromSession(sessionUser)
	songID, err := strconv.Atoi(c.Param("songID"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the song id has to be an integer"), http.StatusBadRequest, c)
		return
	}
	comment := c.PostForm("comment")
	if comment == "" {
		s.ErrorHandler.RenderTemplate(errors.New("cannot submit an empty comment"), http.StatusBadRequest, c)
		return
	}
	if len(comment) > 500 {
		s.ErrorHandler.RenderTemplate(errors.New("your comment cannot have more then 500 characters"), http.StatusBadRequest, c)
		return
	}
	err = s.submitCommentForm.Submit(loggedUser, songID, comment)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/song/%v", songID))
}

func (s *Controller) SubmitSongCommentUpdate(c *gin.Context) {
	var json struct {
		CommentId int    `json:"commentId"`
		Content   string `json:"content"`
	}
	sessionUser, exists := c.Get("user")
	if !exists {
		s.ErrorHandler.RenderTemplate(errors.New("you must be logged in order to rate"), http.StatusBadRequest, c)
		return
	}
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loggedUser, err := user.FromSession(sessionUser)
	//songID, err := strconv.Atoi(c.Param("songID"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the song id has to be an integer"), http.StatusBadRequest, c)
		return
	}
	if err != nil && json.CommentId != 0 {
		s.ErrorHandler.RenderTemplate(errors.New("the comment id has to be an integer"), http.StatusBadRequest, c)
		return
	}
	if len(json.Content) > 500 {
		s.ErrorHandler.RenderTemplate(errors.New("your comment cannot have more then 500 characters"), http.StatusBadRequest, c)
		return
	}
	err = s.submitCommentForm.SubmitUpdate(loggedUser, json.Content, json.CommentId)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	//c.Redirect(http.StatusFound, fmt.Sprintf("/song/%v", songID))
}
