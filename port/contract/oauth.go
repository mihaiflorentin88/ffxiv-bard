package contract

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

type Oauth interface {
	Auth() *oauth2.Config
	GetStateToken() (string, error)
	GenerateJWT(username string) *jwt.Token
	EncodeJWT(token *jwt.Token) (string, error)
	DecodeJWT(tokenString string) (string, error)
}
