package container

import (
	"ffxvi-bard/cmd/http/apps/auth"
	maincontroller "ffxvi-bard/cmd/http/apps/main"
	"ffxvi-bard/cmd/http/apps/song"
	"ffxvi-bard/cmd/http/apps/song/form"
	"ffxvi-bard/cmd/http/utils"
	"ffxvi-bard/cmd/http/utils/middleware"
	"github.com/gin-gonic/gin"
)

type HttpContainer struct {
	errorHandler          *utils.ErrorHandler
	httpRenderer          *utils.Renderer
	authMiddleware        *middleware.AuthMiddleware
	songController        *song.Controller
	mainController        *maincontroller.Controller
	authController        *auth.Controller
	mainRouter            *maincontroller.Router
	songRouter            *song.Router
	authRouter            *auth.Router
	ginRouter             *gin.Engine
	formSongSubmit        *form.SubmitSongForm
	formSongListing       *form.SongList
	formNewSongFormView   *form.NewSongFormView
	formSongDetails       *form.SongDetails
	formSongEdit          *form.SongEditForm
	formSongRating        *form.SubmitSongRatingForm
	formSongCommentSubmit *form.SubmitCommentForm
	formSongEditView      *form.SongEditForm
}

func (s *ServiceContainer) GetErrorHandler() utils.ErrorHandler {
	if s.http.errorHandler != nil {
		return *s.http.errorHandler
	}
	errHandler := utils.NewHttpErrorHandler()
	s.http.errorHandler = &errHandler
	return errHandler
}

func (s *ServiceContainer) GetHttpRenderer() utils.Renderer {
	if s.http.httpRenderer != nil {
		return *s.http.httpRenderer
	}
	renderer := utils.NewRenderer(s.GetErrorHandler())
	s.http.httpRenderer = &renderer
	return renderer
}

func (s *ServiceContainer) GetGinRouter() *gin.Engine {
	if s.http.ginRouter != nil {
		return s.http.ginRouter
	}
	ginRouter := gin.Default()
	s.http.ginRouter = ginRouter
	return ginRouter
}

func (s *ServiceContainer) GetAuthMiddleware() middleware.AuthMiddleware {
	if s.http.authMiddleware != nil {
		return *s.http.authMiddleware
	}
	authMiddleware := middleware.NewAuthMiddleware(s.GetDiscordAuth(), s.GetEmptyUser())
	s.http.authMiddleware = &authMiddleware
	return authMiddleware
}

func (s *ServiceContainer) GetSongController() song.Controller {
	if s.http.songController != nil {
		return *s.http.songController
	}

	songController := song.NewSongController(s.GetEmptySong(), s.GetErrorHandler(), s.GetHttpRenderer(), s.GetNewSubmitSongForm(), s.GetNewSongListingForm(), s.GetNewSongFormView(), s.GetSongDetailsForm(), s.GetSubmitSongRatingForm(), s.GetSubmitSongCommentForm(), s.GetSongEditViewForm())
	s.http.songController = &songController
	return songController
}

func (s *ServiceContainer) GetMainController() maincontroller.Controller {
	if s.http.mainController != nil {
		return *s.http.mainController
	}
	mainController := maincontroller.NewMainController(s.GetErrorHandler(), s.GetHttpRenderer())
	s.http.mainController = &mainController
	return mainController
}

func (s *ServiceContainer) GetSongRouter() *song.Router {
	if s.http.songRouter != nil {
		return s.http.songRouter

	}
	songRouter := song.NewSongRouter(s.GetSongController(), s.GetAuthMiddleware())
	s.http.songRouter = songRouter
	return songRouter
}

func (s *ServiceContainer) getAuthController() auth.Controller {
	if s.http.authController != nil {
		return *s.http.authController
	}
	authController := auth.NewAuthController(s.GetErrorHandler(), s.GetHttpRenderer(), s.GetDiscordAuth(), s.GetUserRepository())
	s.http.authController = &authController
	return authController
}

func (s *ServiceContainer) GetMainRouter() *maincontroller.Router {
	if s.http.mainRouter != nil {
		return s.http.mainRouter
	}
	mainRouter := maincontroller.NewMainRouter(s.GetMainController())
	s.http.mainRouter = mainRouter
	return mainRouter
}

func (s *ServiceContainer) GetAuthRouter() *auth.Router {
	if s.http.authRouter != nil {
		return s.http.authRouter
	}
	authRouter := auth.NewAuthRouter(s.getAuthController())
	s.http.authRouter = authRouter
	return authRouter
}

func (s *ServiceContainer) GetNewSubmitSongForm() form.SubmitSongForm {
	if s.http.formSongSubmit != nil {
		return *s.http.formSongSubmit
	}
	songSubmitForm := form.NewSubmitSongForm(s.GetSongRepository(), s.GetGenreRepository(), s.GetMidiProcessor(), s.GetEmptyUser(), s.GetEmptyGenre(), s.GetEmptyRating(), s.GetEmptyComment(), s.GetEmptyInstrument())
	s.http.formSongSubmit = &songSubmitForm
	return songSubmitForm
}

func (s *ServiceContainer) GetNewSongListingForm() form.SongList {
	if s.http.formSongListing != nil {
		return *s.http.formSongListing
	}
	songListingForm := form.NewSongList(s.GetSongRepository(), s.GetGenreRepository(), s.GetRatingRepository(), s.GetInstrumentRepository())
	s.http.formSongListing = &songListingForm
	return songListingForm
}

func (s *ServiceContainer) GetNewSongFormView() form.NewSongFormView {
	if s.http.formNewSongFormView != nil {
		return *s.http.formNewSongFormView
	}
	newSongFormView := form.NewAddNewSongFormView(s.GetGenreRepository(), s.GetInstrumentRepository())
	s.http.formNewSongFormView = &newSongFormView
	return newSongFormView
}

func (s *ServiceContainer) GetSongDetailsForm() form.SongDetails {
	if s.http.formSongDetails != nil {
		return *s.http.formSongDetails
	}
	songDetailsForm := form.NewSongDetailsForm(s.GetGenreRepository(), s.GetCommentRepository(), s.GetRatingRepository(), s.GetInstrumentRepository(), s.GetEmptySong())
	s.http.formSongDetails = &songDetailsForm
	return songDetailsForm
}

func (s *ServiceContainer) GetSubmitSongRatingForm() form.SubmitSongRatingForm {
	if s.http.formSongRating != nil {
		return *s.http.formSongRating
	}
	songRatingForm := form.NewSubmitSongRatingForm(s.GetRatingRepository())
	s.http.formSongRating = &songRatingForm
	return songRatingForm
}

func (s *ServiceContainer) GetSubmitSongCommentForm() form.SubmitCommentForm {
	if s.http.formSongCommentSubmit != nil {
		return *s.http.formSongCommentSubmit
	}
	submitCommentForm := form.NewSubmitCommentForm(s.GetCommentRepository())
	s.http.formSongCommentSubmit = &submitCommentForm
	return submitCommentForm
}

func (s *ServiceContainer) GetSongEditViewForm() form.SongEditForm {
	if s.http.formSongEditView != nil {
		return *s.http.formSongEditView
	}
	emptySong := s.GetEmptySong()
	songEditForm := form.NewSongEditForm(s.GetGenreRepository(), s.GetInstrumentRepository(), &emptySong)
	s.http.formSongEditView = &songEditForm
	return songEditForm
}
