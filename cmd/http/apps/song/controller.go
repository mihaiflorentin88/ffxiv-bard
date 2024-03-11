package song

import (
	"errors"
	"ffxvi-bard/cmd/http/apps/song/form"
	"ffxvi-bard/cmd/http/utils"
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	Song                 song.Song
	ErrorHandler         utils.ErrorHandler
	Renderer             utils.Renderer
	addSongFormProcessor form.SubmitSongForm
	songListForm         form.SongList
	newSongFormView      form.NewSongFormView
	songDetailsForm      form.SongDetails
	submitSongRatingForm form.SubmitSongRatingForm
	submitCommentForm    form.SubmitCommentForm
	songEditForm         form.SongEditForm
}

func NewSongController(song song.Song, errorHandler utils.ErrorHandler, renderer utils.Renderer,
	addSongFormProcessor form.SubmitSongForm, songListForm form.SongList, newSongFormView form.NewSongFormView,
	songDetailsForm form.SongDetails, submitSongRatingForm form.SubmitSongRatingForm, submitCommentForm form.SubmitCommentForm,
	songEditForm form.SongEditForm) Controller {
	return Controller{
		Song:                 song,
		ErrorHandler:         errorHandler,
		Renderer:             renderer,
		addSongFormProcessor: addSongFormProcessor,
		songListForm:         songListForm,
		newSongFormView:      newSongFormView,
		songDetailsForm:      songDetailsForm,
		submitSongRatingForm: submitSongRatingForm,
		submitCommentForm:    submitCommentForm,
		songEditForm:         songEditForm,
	}
}

func (s *Controller) RenderSongList(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	title := c.Query("title")
	artist := c.Query("artist")
	audioCrafter := c.Query("audio_crafter")
	sort := c.Query("sort")
	ensembleSize, err := strconv.Atoi(c.Query("ensembleSize"))
	if err != nil {
		ensembleSize = -1
	}
	instrumentID, err := strconv.Atoi(c.Query("instrument"))
	if err != nil {
		instrumentID = -1
	}
	genre, err := strconv.Atoi(c.Query("genre"))
	if err != nil {
		genre = -1
	}
	songListForm, err := s.songListForm.Fetch(title, artist, ensembleSize, genre, audioCrafter, instrumentID, page, sort)
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
	copy := s.newSongFormView
	songForm := &copy
	newSongForm, err := songForm.Fetch()
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
	const maxFileSize = 10 << 20
	title := c.PostForm("title")
	artist := c.PostForm("artist")
	ensembleSize := c.PostForm("ensembleSize")
	source := c.PostForm("source")
	if source == "" {
		source = "N/A"
	}
	note := c.PostForm("note")
	if note == "" {
		note = "N/A"
	}
	audioCrafter := c.PostForm("crafter")
	if audioCrafter == "" {
		audioCrafter = "N/A"
	}
	genre := c.PostFormArray("genre")
	instrument := c.PostFormArray("instrument")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	if fileHeader.Size > maxFileSize {
		s.ErrorHandler.RenderTemplate(errors.New("file size exceeds the 10MB limit"), http.StatusBadRequest, c)
	}
	sessionUser, _ := c.Get("user")
	loggedUser, err := user.FromSession(sessionUser)
	if err != nil {
		err = errors.New(fmt.Sprintf("cannot load user: Reason %s. Try relogging", err))
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	copy := s.addSongFormProcessor
	songForm := &copy
	songID, err := songForm.Submit(title, artist, ensembleSize, genre, fileHeader, loggedUser, source, note, audioCrafter, instrument)
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
	copy := s.songDetailsForm
	songDetailsForm := &copy
	songDetails, err := songDetailsForm.Fetch(songID, c)
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

	copy := s.Song
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
			fmt.Sprintf("attachment; filename=%s - %s - %s - %s.mid",
				targetSong.AudioCrafter,
				targetSong.EnsembleString(),
				targetSong.Artist,
				targetSong.Title))
		_ = targetSong.IncrementDownload()
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
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the song id has to be an integer"), http.StatusBadRequest, c)
		return
	}
	rating, err := strconv.Atoi(c.PostForm("rating"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the rating has to be an integer"), http.StatusBadRequest, c)
		return
	}
	copy := s.submitSongRatingForm
	ratingForm := &copy
	err = ratingForm.Submit(loggedUser, songID, rating)
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
	songID, err := strconv.Atoi(c.Param("id"))
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
	copy := s.submitCommentForm
	commentForm := &copy
	err = commentForm.Submit(loggedUser, songID, comment)
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
	copy := s.submitCommentForm
	commentForm := &copy
	err = commentForm.SubmitUpdate(loggedUser, json.Content, json.CommentId)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
}

func (s *Controller) EditSongView(c *gin.Context) {
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the song id has to be an integer"), http.StatusBadRequest, c)
		return
	}
	sessionUser, exists := c.Get("user")
	if !exists {
		s.ErrorHandler.RenderTemplate(errors.New("you must be logged in order to edit"), http.StatusBadRequest, c)
		return
	}
	loggedUser, err := user.FromSession(sessionUser)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	copy := s.songEditForm
	editForm := &copy

	songEditForm, err := editForm.Fetch(songID, loggedUser)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	s.Renderer.
		StartClean().
		RemoveTemplate("resource/template/base/additional_styles.gohtml").
		RemoveTemplate("resource/template/base/additional_js.gohtml").
		AddTemplate("resource/template/song/edit_song.gohtml").
		AddTemplate("resource/template/song/edit_song_css.gohtml").
		AddTemplate("resource/template/song/edit_song_js.gohtml").
		Render(c, songEditForm, http.StatusOK)
}

func (s *Controller) SubmitEditSong(c *gin.Context) {
	sessionUser, exists := c.Get("user")
	if !exists {
		s.ErrorHandler.RenderTemplate(errors.New("you must be logged in order to edit"), http.StatusBadRequest, c)
		return
	}
	loggedUser, err := user.FromSession(sessionUser)
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the song id has to be an integer"), http.StatusBadRequest, c)
		return
	}
	title := c.PostForm("title")
	artist := c.PostForm("artist")
	source := c.PostForm("source")
	note := c.PostForm("note")
	audioCrafter := c.PostForm("crafter")
	ensembleSize, err := strconv.Atoi(c.PostForm("ensembleSize"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the ensembleSize has to be an integer"), http.StatusBadRequest, c)
		return
	}
	var genreIDs []int
	genreStrings := c.PostFormArray("genre")
	for _, genreStringID := range genreStrings {
		genreID, err := strconv.Atoi(genreStringID)
		if err != nil {
			s.ErrorHandler.RenderTemplate(errors.New("the ensembleSize has to be an integer"), http.StatusBadRequest, c)
			return
		}
		genreIDs = append(genreIDs, genreID)
	}

	var instrumentIDs []int
	instrumentStrings := c.PostFormArray("instrument")
	for _, instrumentStringID := range instrumentStrings {
		instrumentID, err := strconv.Atoi(instrumentStringID)
		if err != nil {
			s.ErrorHandler.RenderTemplate(errors.New("the instrument id has to be an integer"), http.StatusBadRequest, c)
			return
		}
		instrumentIDs = append(instrumentIDs, instrumentID)
	}
	copy := s.songEditForm
	editForm := &copy

	err = editForm.HandleSubmittedForm(songID, title, artist, ensembleSize, genreIDs, loggedUser, source, note, audioCrafter, instrumentIDs)
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New(fmt.Sprintf("failed to update song. Reason: %s", err)), http.StatusBadRequest, c)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/song/%v", songID))
}

func (s *Controller) DeleteSong(c *gin.Context) {
	sessionUser, exists := c.Get("user")
	if !exists {
		s.ErrorHandler.RenderTemplate(errors.New("you must be logged in order to edit"), http.StatusBadRequest, c)
		return
	}
	loggedUser, err := user.FromSession(sessionUser)
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.ErrorHandler.RenderTemplate(errors.New("the song id has to be an integer"), http.StatusBadRequest, c)
		return
	}
	if !loggedUser.IsAdmin {
		s.ErrorHandler.RenderTemplate(errors.New("cannot perform action. Reason: permission"), http.StatusBadRequest, c)
		return
	}
	copy := s.Song
	_song := &copy
	currentSong, err := _song.LoadByID(songID)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	err = currentSong.Inactivate()
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
}
