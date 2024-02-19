package middleware

import (
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const BearerSchema = "Bearer "

type AuthMiddleware struct {
	oauth contract.Oauth
	user  *user.User
}

func NewJwtMiddleware(oauth contract.Oauth, user *user.User) AuthMiddleware {
	return AuthMiddleware{
		oauth: oauth,
		user:  user,
	}
}

func (a *AuthMiddleware) API() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, BearerSchema) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or malformed"})
			return
		}
		tokenString := authHeader[len(BearerSchema):]
		username, err := a.oauth.DecodeJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired JWT token"})
			return
		}
		a.user.Username = username
		err = a.user.HydrateByUsername()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "JWT is not associated with any user."})
			return
		}
		c.Set("user", a.user)
		c.Next()
	}
}

func (a *AuthMiddleware) UI() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			return
		}
		username, err := a.oauth.DecodeJWT(tokenString)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			return
		}
		sessionUsername, exists := c.Get("username")
		if !exists && sessionUsername != username {
			a.user.Username = username
			err = a.user.HydrateByUsername()
			if err != nil {
				c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
			}
			c.Set("user", a.user)
			c.Set("username", a.user.Username)
		}
		c.Next()
	}
}

func (a *AuthMiddleware) GetLoggedUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			return
		}
		username, err := a.oauth.DecodeJWT(tokenString)
		if err != nil {
			return
		}
		a.user.Username = username
		err = a.user.HydrateByUsername()
		c.Set("user", a.user)
		c.Set("username", a.user.Username)
		c.Next()
	}
}
