package auth

import (
	"errors"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	ErrorHandler   contract.HttpErrorHandlerInterface
	Renderer       contract.HttpRenderer
	Oauth          contract.Oauth
	UserRepository contract.UserRepositoryInterface
}

func NewAuthController(errorHandler contract.HttpErrorHandlerInterface, renderer contract.HttpRenderer, oauth contract.Oauth, userRepository contract.UserRepositoryInterface) *Controller {
	return &Controller{
		ErrorHandler:   errorHandler,
		Renderer:       renderer,
		Oauth:          oauth,
		UserRepository: userRepository,
	}
}

func (a *Controller) RenderLoginPage(c *gin.Context) {
	a.Renderer.
		AddTemplate("resource/template/auth/login.html").
		AddTemplate("resource/template/auth/login_css.html").
		RemoveTemplate("resource/template/base/additional_styles.html").
		Render(c, nil)
}

func (a *Controller) LoginDiscord(c *gin.Context) {
	oauthConf := a.Oauth.Auth()
	state, err := a.Oauth.GetStateToken()
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, c)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, oauthConf.AuthCodeURL(state))
}

func (a *Controller) LoginDiscordCallback(c *gin.Context) {
	code := c.Query("code")
	retrievedState := c.Query("state")

	myState, err := a.Oauth.GetStateToken()
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, c)
		return
	}
	if retrievedState != myState {
		a.ErrorHandler.RenderTemplate(errors.New("state does not match"), c)
		return
	}
	token, err := a.Oauth.Auth().Exchange(c, code)
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, c)
		return
	}
	response, err := a.Oauth.Auth().Client(c, token).Get("https://discord.com/api/users/@me")
	if err != nil && response.StatusCode != 200 {
		a.ErrorHandler.RenderTemplate(err, c)
	}
	userDTO, err := dto.DiscordUserFromHttpResponse(response, code, token)
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, c)
	}
	loggedUser := user.FromDiscordUserDTO(userDTO, a.Oauth, a.UserRepository)
	err = loggedUser.Persist()
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, c)
	}
	defer response.Body.Close()
	c.Redirect(http.StatusPermanentRedirect, "/")
}

func (a *Controller) logout() {

}
