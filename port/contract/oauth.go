package contract

import (
	"golang.org/x/oauth2"
)

type Oauth interface {
	Auth() *oauth2.Config
	GetStateToken() (string, error)
}
