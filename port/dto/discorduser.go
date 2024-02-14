package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type DiscordUser struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Name              string `json:"global_name"`
	Discriminator     string `json:"discriminator"`
	Email             string `json:"email,omitempty"`
	Avatar            string `json:"avatar"`
	AuthorizationCode string
	AccessToken       *oauth2.Token
	TokenReleasedAt   time.Time
}

func DiscordUserFromHttpResponse(response *http.Response, authorizationCode string, accessToken *oauth2.Token) (DiscordUser, error) {
	var user DiscordUser
	err := json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return user, errors.New(fmt.Sprintf("failed to decode discord user data. Reason: %s", err))
	}
	user.AuthorizationCode = authorizationCode
	user.AccessToken = accessToken
	user.TokenReleasedAt = time.Now()
	return user, nil
}
