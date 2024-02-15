package container

import (
	"ffxvi-bard/cmd/http/apps/auth"
	maincontroller "ffxvi-bard/cmd/http/apps/main"
	"ffxvi-bard/cmd/http/apps/song"
	"ffxvi-bard/cmd/http/utils"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"log"
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

func GetSongController() *song.Controller {
	return song.NewSongController(GetErrorHandler(), GetHttpRenderer())
}

func GetMainController() *maincontroller.Controller {
	return maincontroller.NewMainController(GetErrorHandler(), GetHttpRenderer())
}

func GetSongRouter() contract.RouterInterface {
	return song.NewSongRouter(GetSongController())
}

func getAuthController() *auth.Controller {
	r, err := GetUserRepository()
	if err != nil {
		log.Println(err)
	}
	return auth.NewAuthController(GetErrorHandler(), GetHttpRenderer(), GetDiscordAuth(), r)
}

func GetMainRouter() contract.RouterInterface {
	return maincontroller.NewMainRouter(GetMainController())
}

func GetAuthRouter() contract.RouterInterface {
	return auth.NewAuthRouter(getAuthController())
}
