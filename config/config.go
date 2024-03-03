package config

import (
	"embed"
	"errors"
	"ffxvi-bard/infrastructure/filesystem"
	"ffxvi-bard/port/contract"
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"path/filepath"
)

//go:embed config.toml
var configFS embed.FS

type Config struct {
	Database DatabaseConfig `toml:"database"`
	Song     SongConfig     `toml:"song"`
	Discord  DiscordConfig  `toml:"discord"`
	Spotify  SpotifyConfig  `toml:"spotify"`
}

type DatabaseConfig struct {
	Database string `toml:"database"`
	Path     string `toml:"path"`
}

func EnsureFolder(path string, fs contract.FileSystemInterface) {
	err := fs.EnsureDir(path)
	if err != nil {
		log.Panicf("Cannot create new folder on path %s", path)
	}
}

type SongConfig struct {
	UnfinishedFilesPath string `toml:"unfinished_files_path"`
	FinishedFilesPath   string `toml:"finished_files_path"`
	MaxFileSize         int    `toml:"max_file_size"`
	filesystem          contract.FileSystemInterface
}

type DiscordConfig struct {
	ClientID     string   `toml:"client_id"`
	ClientSecret string   `toml:"client_secret"`
	RedirectURI  string   `toml:"redirect_uri"`
	Scopes       []string `toml:"scopes"`
	JwtSecret    string   `toml:"jwt_secret"`
}

type SpotifyConfig struct {
	Url         string `toml:"url"`
	AccountsUrl string `toml:"accounts_url"`
	Secret      string `toml:"secret"`
	ClientID    string `toml:"client_id"`
}

func (s *SongConfig) ensureFolders() {
	err := s.filesystem.EnsureDir(s.UnfinishedFilesPath)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error creating unfinished files directory: %s", err.Error())))
	}

	err = s.filesystem.EnsureDir(s.FinishedFilesPath)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error creating finished files directory: %s", err.Error())))
	}
}

func NewConfig() (*Config, error) {
	fs := filesystem.FileSystem{}
	var config Config
	content, err := configFS.ReadFile("config.toml")
	if err != nil {
		return nil, err
	}
	if _, err := toml.Decode(string(content), &config); err != nil {
		return nil, err
	}
	config.Song.filesystem = fs
	dbPath := filepath.Dir(config.Database.Path)
	EnsureFolder(dbPath, fs)
	config.Song.ensureFolders()
	return &config, err
}
