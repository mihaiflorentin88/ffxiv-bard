package container

import (
	"ffxvi-bard/cmd/http/apps/auth"
	maincontroller "ffxvi-bard/cmd/http/apps/main"
	"ffxvi-bard/cmd/http/apps/song"
	"ffxvi-bard/cmd/http/apps/song/form"
	"ffxvi-bard/cmd/http/utils"
	"ffxvi-bard/cmd/http/utils/middleware"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

func GetErrorHandler() contract.HttpErrorHandlerInterface {
	return utils.NewHttpErrorHandler()
}

func GetHttpRenderer() contract.HttpRenderer {
	return utils.NewRenderer(GetErrorHandler())
}

func GetGinRouter() *gin.Engine {
	return gin.Default()
}

func GetAuthMiddleware() middleware.AuthMiddleware {
	return middleware.NewJwtMiddleware(GetDiscordAuth(), GetNewEmptyUser())
}

func GetSongController() *song.Controller {
	return song.NewSongController(GetEmptySong(), GetErrorHandler(), GetHttpRenderer(), GetAddSongFormProcessor(), GetNewSongListingForm(), GetNewSongFormView(), GetSongDetailsForm(), GetSubmitSongRatingForm(), GetSubmitSongCommentForm())
}

func GetMainController() *maincontroller.Controller {
	return maincontroller.NewMainController(GetErrorHandler(), GetHttpRenderer())
}

func GetSongRouter() contract.RouterInterface {
	return song.NewSongRouter(GetSongController(), GetAuthMiddleware())
}

func getAuthController() *auth.Controller {
	return auth.NewAuthController(GetErrorHandler(), GetHttpRenderer(), GetDiscordAuth(), GetUserRepository())
}

func GetMainRouter() contract.RouterInterface {
	return maincontroller.NewMainRouter(GetMainController())
}

func GetAuthRouter() contract.RouterInterface {
	return auth.NewAuthRouter(getAuthController())
}

func GetAddSongFormProcessor() form.AddSongFormProcessor {
	return form.NewSongFormProcessor(GetSongRepository(), GetGenreRepository(), GetMidiProcessor())
}

func GetNewSongListingForm() form.SongList {
	return form.NewSongList(GetSongRepository(), GetGenreRepository(), GetRatingRepository(), GetSpotifyClient())
}

func GetNewSongFormView() form.NewSongFormView {
	return form.NewAddNewSongFormView(GetGenreRepository())
}

func GetSongDetailsForm() form.SongDetails {
	commentRepository := GetCommentRepository()
	return form.NewSongDetailsForm(GetGenreRepository(), &commentRepository, GetRatingRepository(), GetEmptySong())
}

func GetSubmitSongRatingForm() form.SubmitSongRatingForm {
	return form.NewSubmitSongRatingForm(GetRatingRepository())
}

func GetSubmitSongCommentForm() form.SubmitCommentForm {
	return form.NewSubmitCommentForm(GetCommentRepository())
}
