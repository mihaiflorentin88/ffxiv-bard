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

func (s *ServiceContainer) ErrorHandler() utils.ErrorHandler {
	if s.http.errorHandler != nil {
		return *s.http.errorHandler
	}
	errHandler := utils.NewHttpErrorHandler()
	s.http.errorHandler = &errHandler
	return errHandler
}

func (s *ServiceContainer) HttpRenderer() utils.Renderer {
	if s.http.httpRenderer != nil {
		return *s.http.httpRenderer
	}
	renderer := utils.NewRenderer(s.ErrorHandler())
	s.http.httpRenderer = &renderer
	return renderer
}

func (s *ServiceContainer) GinRouter() *gin.Engine {
	if s.http.ginRouter != nil {
		return s.http.ginRouter
	}
	ginRouter := gin.Default()
	s.http.ginRouter = ginRouter
	return ginRouter
}

func (s *ServiceContainer) AuthMiddleware() middleware.AuthMiddleware {
	if s.http.authMiddleware != nil {
		return *s.http.authMiddleware
	}
	authMiddleware := middleware.NewAuthMiddleware(s.DiscordAuth(), s.EmptyUser())
	s.http.authMiddleware = &authMiddleware
	return authMiddleware
}

func (s *ServiceContainer) SongController() song.Controller {
	if s.http.songController != nil {
		return *s.http.songController
	}

	songController := song.NewSongController(s.EmptySong(), s.ErrorHandler(), s.HttpRenderer(), s.SubmitSongForm(), s.SongListingForm(), s.SongFormView(), s.SongDetailsForm(), s.SubmitSongRatingForm(), s.SubmitSongCommentForm(), s.SongEditViewForm())
	s.http.songController = &songController
	return songController
}

func (s *ServiceContainer) MainController() maincontroller.Controller {
	if s.http.mainController != nil {
		return *s.http.mainController
	}
	mainController := maincontroller.NewMainController(s.ErrorHandler(), s.HttpRenderer())
	s.http.mainController = &mainController
	return mainController
}

func (s *ServiceContainer) SongRouter() *song.Router {
	if s.http.songRouter != nil {
		return s.http.songRouter

	}
	songRouter := song.NewSongRouter(s.SongController(), s.AuthMiddleware())
	s.http.songRouter = songRouter
	return songRouter
}

func (s *ServiceContainer) AuthController() auth.Controller {
	if s.http.authController != nil {
		return *s.http.authController
	}
	authController := auth.NewAuthController(s.ErrorHandler(), s.HttpRenderer(), s.DiscordAuth(), s.UserRepository())
	s.http.authController = &authController
	return authController
}

func (s *ServiceContainer) MainRouter() *maincontroller.Router {
	if s.http.mainRouter != nil {
		return s.http.mainRouter
	}
	mainRouter := maincontroller.NewMainRouter(s.MainController())
	s.http.mainRouter = mainRouter
	return mainRouter
}

func (s *ServiceContainer) AuthRouter() *auth.Router {
	if s.http.authRouter != nil {
		return s.http.authRouter
	}
	authRouter := auth.NewAuthRouter(s.AuthController())
	s.http.authRouter = authRouter
	return authRouter
}

func (s *ServiceContainer) SubmitSongForm() form.SubmitSongForm {
	if s.http.formSongSubmit != nil {
		return *s.http.formSongSubmit
	}
	songSubmitForm := form.NewSubmitSongForm(s.SongRepository(), s.GenreRepository(), s.MidiProcessor(), s.EmptyUser(), s.EmptyGenre(), s.EmptyRating(), s.EmptyComment(), s.EmptyInstrument())
	s.http.formSongSubmit = &songSubmitForm
	return songSubmitForm
}

func (s *ServiceContainer) SongListingForm() form.SongList {
	if s.http.formSongListing != nil {
		return *s.http.formSongListing
	}
	songListingForm := form.NewSongList(s.SongRepository(), s.GenreRepository(), s.RatingRepository(), s.InstrumentRepository())
	s.http.formSongListing = &songListingForm
	return songListingForm
}

func (s *ServiceContainer) SongFormView() form.NewSongFormView {
	if s.http.formNewSongFormView != nil {
		return *s.http.formNewSongFormView
	}
	newSongFormView := form.NewAddNewSongFormView(s.GenreRepository(), s.InstrumentRepository())
	s.http.formNewSongFormView = &newSongFormView
	return newSongFormView
}

func (s *ServiceContainer) SongDetailsForm() form.SongDetails {
	if s.http.formSongDetails != nil {
		return *s.http.formSongDetails
	}
	songDetailsForm := form.NewSongDetailsForm(s.GenreRepository(), s.CommentRepository(), s.RatingRepository(), s.InstrumentRepository(), s.EmptySong())
	s.http.formSongDetails = &songDetailsForm
	return songDetailsForm
}

func (s *ServiceContainer) SubmitSongRatingForm() form.SubmitSongRatingForm {
	if s.http.formSongRating != nil {
		return *s.http.formSongRating
	}
	songRatingForm := form.NewSubmitSongRatingForm(s.RatingRepository())
	s.http.formSongRating = &songRatingForm
	return songRatingForm
}

func (s *ServiceContainer) SubmitSongCommentForm() form.SubmitCommentForm {
	if s.http.formSongCommentSubmit != nil {
		return *s.http.formSongCommentSubmit
	}
	submitCommentForm := form.NewSubmitCommentForm(s.CommentRepository())
	s.http.formSongCommentSubmit = &submitCommentForm
	return submitCommentForm
}

func (s *ServiceContainer) SongEditViewForm() form.SongEditForm {
	if s.http.formSongEditView != nil {
		return *s.http.formSongEditView
	}
	emptySong := s.EmptySong()
	songEditForm := form.NewSongEditForm(s.GenreRepository(), s.InstrumentRepository(), &emptySong)
	s.http.formSongEditView = &songEditForm
	return songEditForm
}
