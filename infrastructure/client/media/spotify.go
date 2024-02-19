package media

import (
	"encoding/base64"
	"encoding/json"
	"ffxvi-bard/config"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const tokenRoute = "/token"

type spotifyToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

func (s *spotifyToken) IsValid() bool {
	if time.Now().UnixMilli() > s.ExpiresAt.UnixMilli() {
		return false
	}
	return true
}

type spotifyAuth struct {
	config *config.SpotifyConfig
	token  *spotifyToken
	client *http.Client
}

func (s *spotifyAuth) getTokenUrl() string {
	return fmt.Sprintf("%s%s", s.config.AccountsUrl, tokenRoute)
}

func (s *spotifyAuth) auth() error {
	var tokenResp spotifyToken
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", s.getTokenUrl(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	authValue := s.config.ClientID + ":" + s.config.Secret
	basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(authValue))
	req.Header.Add("Authorization", basicAuth)
	tokenResp.CreatedAt = time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return err
	}
	tokenResp.ExpiresAt = tokenResp.CreatedAt.Add(time.Duration(tokenResp.ExpiresIn-10) * time.Second)
	s.client = client
	s.token = &tokenResp
	return nil
}

func (s *spotifyAuth) GetClient() (*http.Client, error) {
	if s.client == nil {
		err := s.auth()
		if err != nil {
			return nil, err
		}
	}
	if !s.token.IsValid() {
		err := s.auth()
		if err != nil {
			return nil, err
		}
	}
	return s.client, nil
}

func (s *spotifyAuth) GetToken() (*spotifyToken, error) {
	if s.token == nil {
		err := s.auth()
		if err != nil {
			return nil, err
		}
	}
	if !s.token.IsValid() {
		err := s.auth()
		if err != nil {
			return nil, err
		}
	}
	return s.token, nil
}

type SpotifyClient struct {
	config config.SpotifyConfig
	auth   spotifyAuth
}

func NewSpotifyClient(config config.SpotifyConfig) contract.MediaClientInterface {
	return &SpotifyClient{
		config: config,
		auth:   spotifyAuth{config: &config},
	}
}

func (s *SpotifyClient) getSearchURL(track string, artist string) string {
	route := "/search?q="
	query := "remaster%20"
	if track != "" {
		query += fmt.Sprintf("track:%s", track)
	}
	if artist != "" {
		query += fmt.Sprintf("artist:%s", artist)
	}
	query += "&type=album"
	//encodedQuery := url.QueryEscape(query)
	route += query
	return fmt.Sprintf("%s%s", s.config.Url, route)
}

func (s *SpotifyClient) Search(track string, artist string) (dto.MediaResponse, error) {
	var response dto.MediaResponse
	searchUrl := s.getSearchURL(track, artist)
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		return response, err
	}
	token, err := s.auth.GetToken()
	if err != nil {
		return response, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	resp, err := s.auth.client.Do(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	response.Hydrate(body)
	return response, nil
}
