package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"ffxvi-bard/config"
	"ffxvi-bard/port/contract"
	"github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

type discordOauth struct {
	config *config.DiscordConfig
	auth   *oauth2.Config
	state  string
}

func NewDiscordOauth(config *config.DiscordConfig) contract.Oauth {
	return &discordOauth{config: config}
}

func (d *discordOauth) Auth() *oauth2.Config {
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

func (d *discordOauth) GetStateToken() (string, error) {
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
