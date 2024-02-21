package user

import (
	"context"
	"errors"
	"ffxvi-bard/domain/date"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"ffxvi-bard/port/helper"
	"fmt"
	"golang.org/x/oauth2"
	"time"
)

const ProviderDiscord = "discord"

type User struct {
	StorageID         int64
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
	IsAdmin           bool
	OauthToken        *oauth2.Token
	ExpiresAt         time.Time
	TokenReleasedAt   time.Time
	Oauth             contract.Oauth
	Repository        contract.UserRepositoryInterface
	Date              date.Date
}

func NewEmptyUser(repository contract.UserRepositoryInterface) *User {
	return &User{Repository: repository}
}

func (u *User) HydrateByID() error {
	if u.StorageID == 0 {
		return errors.New("user has no id assigned")
	}
	if u.Repository == nil {
		return errors.New("user.User was not properly instantiated")
	}
	userDto, err := u.Repository.FindById(u.StorageID)
	if err != nil {
		return errors.New(fmt.Sprintf("could not load userid %v. Reason: %s ", u.StorageID, err))
	}
	FromDatabaseUserDTO(userDto, u, u.Repository)
	return nil
}

func (u *User) HydrateByUsername() error {
	if u.Username == "" {
		return errors.New("user has no username assigned")
	}
	userDto, err := u.Repository.FindByUsername(u.Username)
	if err != nil {
		return errors.New(fmt.Sprintf("could not load userid. Reason: %s ", err))
	}
	FromDatabaseUserDTO(userDto, u, u.Repository)
	return nil
}

func (u *User) HydrateByEmail() error {
	if u.Username == "" {
		return errors.New("user has no email address assigned")
	}
	userDto, err := u.Repository.FindByEmail(u.Email)
	if err != nil {
		return errors.New(fmt.Sprintf("could not load userid. Reason: %s ", err))
	}
	FromDatabaseUserDTO(userDto, u, u.Repository)
	return nil
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
	if u.StorageID != 0 {
		userDto, err := u.Repository.FindById(u.StorageID)
		if err == nil || userDto != nil {
			u.StorageID = userDto.ID
			err := u.Repository.Update(u.ToDatabaseUserDTO())
			if err != nil {
				return errors.New(fmt.Sprintf("cannot update user. Reason: %s", err))
			}
			return nil
		}
	}
	if u.Username != "" {
		userDto, err := u.Repository.FindByUsername(u.Username)
		if err == nil || userDto != nil {
			u.StorageID = userDto.ID
			err := u.Repository.Update(u.ToDatabaseUserDTO())
			if err != nil {
				return errors.New(fmt.Sprintf("cannot update user. Reason: %s", err))
			}
			return nil
		}
	}

	if u.Email != "" {
		userDto, err := u.Repository.FindByEmail(u.Email)
		if err == nil || userDto != nil {
			u.StorageID = userDto.ID
			err := u.Repository.Update(u.ToDatabaseUserDTO())
			if err != nil {
				return errors.New(fmt.Sprintf("cannot update user. Reason: %s", err))
			}
			return nil
		}
	}
	err := u.Repository.Create(u.ToDatabaseUserDTO())

	if err != nil {
		return errors.New(fmt.Sprintf("cannot create new user. Reason: %s", err))
	}
	return nil
}

func (u *User) ToDatabaseUserDTO() dto.DatabaseUser {
	return dto.DatabaseUser{
		ID:                u.StorageID,
		Username:          u.Username,
		Email:             u.Email,
		Name:              &u.Name,
		AuthProvider:      &u.AuthProvider,
		ProviderID:        &u.ProviderID,
		Avatar:            &u.Avatar,
		Discriminator:     &u.Discriminator,
		AuthorizationCode: &u.AuthorizationCode,
		AccessToken:       &u.AccessToken,
		RefreshToken:      &u.RefreshToken,
		TokenType:         &u.TokenType,
		ExpiresAt:         &u.ExpiresAt,
		TokenReleasedAt:   &u.TokenReleasedAt,
		CreatedAt:         u.Date.CreatedAt,
		UpdatedAt:         u.Date.UpdatedAt,
	}
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
		Repository:        repository,
	}
}

func FromDatabaseUserDTO(userDTO *dto.DatabaseUser, user *User, userRepository contract.UserRepositoryInterface) *User {
	user.StorageID = userDTO.ID
	user.Username = userDTO.Username
	user.Email = userDTO.Email
	user.Name = helper.GetStringValue(userDTO.Name)
	user.AuthProvider = helper.GetStringValue(userDTO.AuthProvider)
	user.ProviderID = helper.GetStringValue(userDTO.ProviderID)
	user.Avatar = helper.GetStringValue(userDTO.Avatar)
	user.Discriminator = helper.GetStringValue(userDTO.Discriminator)
	user.AuthorizationCode = helper.GetStringValue(userDTO.AuthorizationCode)
	user.AccessToken = helper.GetStringValue(userDTO.AccessToken)
	user.RefreshToken = helper.GetStringValue(userDTO.RefreshToken)
	user.TokenType = helper.GetStringValue(userDTO.TokenType)
	user.ExpiresAt = helper.GetTimeValue(userDTO.ExpiresAt)
	user.TokenReleasedAt = helper.GetTimeValue(userDTO.TokenReleasedAt)
	user.IsAdmin = userDTO.IsAdmin
	user.Repository = userRepository
	return user
}

func FromSession(sessionUser interface{}) (*User, error) {
	if sessionUser != nil {
		if userObj, ok := sessionUser.(*User); ok {
			return userObj, nil
		}
	}
	return nil, errors.New("session user is not of the correct type")
}
