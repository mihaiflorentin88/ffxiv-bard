package song

import (
	"ffxvi-bard/cmd/http/apps/song/form"
	"ffxvi-bard/domain/song"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	Song            contract.SongInterface
	GenreRepository contract.GenreRepositoryInterface
	ErrorHandler    contract.HttpErrorHandlerInterface
	Renderer        contract.HttpRenderer
	SongRepository  contract.SongRepositoryInterface
	SongProcessor   contract.SongProcessorInterface
}

func NewSongController(song contract.SongInterface, errorHandler contract.HttpErrorHandlerInterface,
	renderer contract.HttpRenderer, genreRepository contract.GenreRepositoryInterface,
	songRepository contract.SongRepositoryInterface, songProcessor contract.SongProcessorInterface) *Controller {
	return &Controller{
		Song:            song,
		GenreRepository: genreRepository,
		ErrorHandler:    errorHandler,
		Renderer:        renderer,
		SongRepository:  songRepository,
		SongProcessor:   songProcessor,
	}
}

func (s *Controller) RenderSongList(c *gin.Context) {
	form := form.NewSongList(s.SongRepository, s.ErrorHandler)
	form.FetchData(c)
	s.Renderer.
		StartClean().
		RemoveTemplate("resource/template/base/additional_js.gohtml").
		RemoveTemplate("resource/template/base/additional_styles.gohtml").
		AddTemplate("resource/template/song/list.gohtml").
		AddTemplate("resource/template/song/list_css.gohtml").
		AddTemplate("resource/template/song/list_js.gohtml").
		Render(c, form, http.StatusOK)
}

func (s *Controller) RenderAddNewSongForm(c *gin.Context) {
	formViewData, err := form.NewAddNewSongFormView(s.Song, s.GenreRepository)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusInternalServerError, c)
		return
	}
	s.Renderer.
		StartClean().
		RemoveTemplate("resource/template/base/additional_styles.gohtml").
		AddTemplate("resource/template/song/add_song.gohtml").
		AddTemplate("resource/template/song/add_song_css.gohtml").
		Render(c, formViewData, http.StatusOK)
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
	submittedForm, err := form.NewSongFormSubmitted(title, artist, ensembleSize, genre, fileHeader, s.ErrorHandler, c)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	songDTO := dto.AddNewSong(submittedForm.Title, submittedForm.Artist, submittedForm.EnsembleSize,
		submittedForm.Genre, submittedForm.File, &submittedForm.User)
	newSong, err := song.FromNewSongDTO(songDTO, s.SongRepository, s.GenreRepository, s.SongProcessor)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
	songDatabaseDto := newSong.ToDatabaseSongDTO()
	err = s.SongRepository.InsertNewSong(songDatabaseDto, songDTO.Genre)
	if err != nil {
		s.ErrorHandler.RenderTemplate(err, http.StatusBadRequest, c)
		return
	}
}
