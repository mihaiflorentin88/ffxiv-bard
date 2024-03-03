package auth

import (
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
		StartClean().
		AddTemplate("resource/template/auth/login.gohtml").
		AddTemplate("resource/template/auth/login_css.gohtml").
		RemoveTemplate("resource/template/base/additional_styles.gohtml").
		Render(c, nil, http.StatusOK)
}

func (a *Controller) LoginWithDiscord(c *gin.Context) {
	oauthConf := a.Oauth.Auth()
	state, err := a.Oauth.GetStateToken()
	stateJWT := a.Oauth.GenerateJWT(state)
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}
	encodedState, err := a.Oauth.EncodeJWT(stateJWT)
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}
	//c.SetCookie("lts", encodedState, 60, "/", "", false, false)
	c.Redirect(http.StatusTemporaryRedirect, oauthConf.AuthCodeURL(encodedState))
}

func (a *Controller) LoginWithDiscordCallback(c *gin.Context) {
	code := c.Query("code")
	retrievedState := c.Query("state")
	//lts, _ := c.Cookie("lts")
	_, err := a.Oauth.DecodeJWT(retrievedState)
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}

	//myState, err := a.Oauth.GetStateToken()
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}
	//if retrievedState != myState {
	//	if ltsDecoded != retrievedState {
	//		a.ErrorHandler.RenderTemplate(errors.New("state does not match"), http.StatusUnauthorized, c)
	//		return
	//	}
	//}
	token, err := a.Oauth.Auth().Exchange(c, code)
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}
	response, err := a.Oauth.Auth().Client(c, token).Get("https://discord.com/api/users/@me")
	if err != nil && response.StatusCode != 200 {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}
	userDTO, err := dto.DiscordUserFromHttpResponse(response, code, token)
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}
	loggedUser := user.FromDiscordUserDTO(userDTO, a.Oauth, a.UserRepository)
	err = loggedUser.Persist()
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}
	defer response.Body.Close()
	jwt := a.Oauth.GenerateJWT(loggedUser.Username)
	userToken, err := a.Oauth.EncodeJWT(jwt)
	if err != nil {
		a.ErrorHandler.RenderTemplate(err, http.StatusUnauthorized, c)
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "token",
		Value:    userToken,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true,
	})
	c.SetCookie("token", userToken, 72*60*60, "/", "", true, false)
	c.Redirect(http.StatusPermanentRedirect, "/")
}

func (a *Controller) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Set("user", nil)
	c.Redirect(http.StatusSeeOther, "/")
}
