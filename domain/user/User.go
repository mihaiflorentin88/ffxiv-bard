package user

import (
	"context"
	"errors"
	"ffxvi-bard/domain/date"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"golang.org/x/oauth2"
	"time"
)

const ProviderDiscord = "discord"

type User struct {
	StorageID         int
	Username          string
	Email             string
	Name              string
	AuthProvider      string
	ProviderID        string
	Avatar            string
	Discriminator     string
	AuthorizationCode string
	AccessToken       string
	RefreshToken      string
	TokenType         string
	OauthToken        *oauth2.Token
	ExpiresAt         time.Time
	TokenReleasedAt   time.Time
	Oauth             contract.Oauth
	UserRepository    contract.UserRepositoryInterface
	Date              date.Date
}

func (u *User) HasTokenExpired() bool {
	if u.OauthToken != nil {
		return u.OauthToken.Valid()
	}
	if u.ExpiresAt.Unix()-10 < time.Now().Unix() {
		return false
	}
	return true
}

func (u *User) RenewToken() error {
	if u.Oauth == nil || u.RefreshToken == "" {
		return errors.New("OAuth service or refresh token not available")
	}

	newToken, err := u.Oauth.Auth().TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: u.RefreshToken,
	}).Token()
	if err != nil {
		return err
	}
	u.AccessToken = newToken.AccessToken
	u.RefreshToken = newToken.RefreshToken
	u.ExpiresAt = newToken.Expiry
	u.TokenType = newToken.TokenType
	u.OauthToken = newToken

	return u.Persist()
}

func (u *User) Persist() error {
	return nil
}

func (u *User) GetName() string {
	return u.Name
}

func FromDiscordUserDTO(discordUser dto.DiscordUser, oauth contract.Oauth, repository contract.UserRepositoryInterface) User {
	return User{
		Username:          discordUser.Username,
		Email:             discordUser.Email,
		Name:              discordUser.Name,
		AuthProvider:      ProviderDiscord,
		ProviderID:        discordUser.ID,
		Avatar:            discordUser.Avatar,
		Discriminator:     discordUser.Discriminator,
		AuthorizationCode: discordUser.AuthorizationCode,
		AccessToken:       discordUser.AccessToken.AccessToken,
		RefreshToken:      discordUser.AccessToken.RefreshToken,
		TokenType:         discordUser.AccessToken.TokenType,
		ExpiresAt:         discordUser.AccessToken.Expiry,
		OauthToken:        discordUser.AccessToken,
		Oauth:             oauth,
		TokenReleasedAt:   discordUser.TokenReleasedAt,
	}
}
