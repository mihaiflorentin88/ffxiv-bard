package httpapps

import (
	"embed"
	"ffxvi-bard/cmd/http/apps/auth"
	maincontroller "ffxvi-bard/cmd/http/apps/main"
	songcontroller "ffxvi-bard/cmd/http/apps/song"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

type appManager struct {
	Song contract.Controller
	Main contract.Controller
	Auth contract.Controller
}

func NewRoot(router *gin.Engine, staticFS *embed.FS) contract.Controller {
	return &appManager{
		Song: songcontroller.NewSongController(router, staticFS, NewHttpError(), NewRenderer(NewHttpError(), staticFS)),
		Main: maincontroller.NewMainController(router, staticFS, NewHttpError(), NewRenderer(NewHttpError(), staticFS)),
		Auth: auth.NewAuth(router, staticFS, NewHttpError(), NewRenderer(NewHttpError(), staticFS)),
	}
}

func (a *appManager) Initialize() {
	a.Song.Initialize()
	a.Main.Initialize()
	a.Auth.Initialize()
}
