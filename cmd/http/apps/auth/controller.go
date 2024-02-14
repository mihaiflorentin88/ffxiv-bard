package auth

import (
	"embed"
	"errors"
	"ffxvi-bard/config"
	oauth2 "ffxvi-bard/infrastructure/oauth"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type authController struct {
	Router       *gin.Engine
	StaticFS     *embed.FS
	ErrorHandler contract.HttpError
	Renderer     contract.HttpRenderer
	Oauth        contract.Oauth
}

func NewAuth(router *gin.Engine, staticFs *embed.FS, errorHandler contract.HttpError,
	renderer contract.HttpRenderer) contract.Controller {
	appConfig, err := config.NewConfig()
	if err != nil {
		panic("Cannot load application configuration")
	}
	oauth := oauth2.NewDiscordOauth(&appConfig.Discord)
	return &authController{
		Router:       router,
		StaticFS:     staticFs,
		ErrorHandler: errorHandler,
		Renderer:     renderer,
		Oauth:        oauth,
	}
}

func (a *authController) Initialize() {
	a.loginDiscord()
	a.loginDiscordCallback()
	a.showLogin()
	a.logout()
}

func (a *authController) showLogin() {
	a.Router.GET("/auth/login", func(c *gin.Context) {
		a.Renderer.
			AddTemplate("resource/template/auth/login.html").
			AddTemplate("resource/template/auth/login_css.html").
			RemoveTemplate("resource/template/base/additional_styles.html").
			Render(c, nil)
	})
}

func (a *authController) loginDiscord() {
	a.Router.GET("/auth/login/discord", func(c *gin.Context) {
		oauthConf := a.Oauth.Auth()
		state, err := a.Oauth.GetStateToken()
		if err != nil {
			a.ErrorHandler.RenderTemplate(err, c, a.StaticFS)
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, oauthConf.AuthCodeURL(state))
	})
}

func (a *authController) loginDiscordCallback() {
	a.Router.GET("/auth/login/discord/callback", func(c *gin.Context) {
		code := c.Query("code")
		retrievedState := c.Query("state")

		myState, err := a.Oauth.GetStateToken()
		if err != nil {
			a.ErrorHandler.RenderTemplate(err, c, a.StaticFS)
			return
		}
		if retrievedState != myState {
			a.ErrorHandler.RenderTemplate(errors.New("state does not match"), c, a.StaticFS)
			return
		}
		token, err := a.Oauth.Auth().Exchange(c, code)
		if err != nil {
			a.ErrorHandler.RenderTemplate(err, c, a.StaticFS)
			return
		}
		response, err := a.Oauth.Auth().Client(c, token).Get("https://discord.com/api/users/@me")
		if err != nil && response.StatusCode != 200 {
			a.ErrorHandler.RenderTemplate(err, c, a.StaticFS)
		}
		userDTO, err := dto.DiscordUserFromHttpResponse(response, code, token)
		if err != nil {
			a.ErrorHandler.RenderTemplate(err, c, a.StaticFS)
		}
		//loggedUser := user.FromDiscordUserDTO(userDTO, a.Oauth)
		userDTO.Name = "Test"
		defer response.Body.Close()
	})
}

func (a *authController) logout() {

}
