package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"ffxvi-bard/config"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
	"time"
)

type DiscordOauth struct {
	config *config.DiscordConfig
	auth   *oauth2.Config
	state  string
	jwt    *jwt.Token
}

func NewDiscordOauth(config *config.DiscordConfig) *DiscordOauth {
	return &DiscordOauth{config: config}
}

func (d *DiscordOauth) Auth() *oauth2.Config {
	if d.auth != nil {
		return d.auth
	}
	d.auth = &oauth2.Config{
		RedirectURL:  d.config.RedirectURI,
		ClientID:     d.config.ClientID,
		ClientSecret: d.config.ClientSecret,
		Scopes:       d.config.Scopes,
		Endpoint:     discord.Endpoint,
	}
	return d.auth
}

func (d *DiscordOauth) GetStateToken() (string, error) {
	if d.state != "" {
		return d.state, nil
	}
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(b)
	d.state = state
	return state, nil
}

func (d *DiscordOauth) GenerateJWT(username string) *jwt.Token {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": username,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	d.jwt = token
	return token
}

func (d *DiscordOauth) EncodeJWT(token *jwt.Token) (string, error) {
	tokenString, err := token.SignedString([]byte(d.config.JwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (d *DiscordOauth) DecodeJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(d.config.JwtSecret), nil
	})

	if err != nil {
		return "", fmt.Errorf("cannot verify token signature. Reason: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string) // Make sure to use the same claim key as in GenerateJWT
		if !ok {
			return "", errors.New("user_id claim not found or wrong type")
		}
		return userID, nil
	}

	return "", errors.New("invalid or expired token")
}
