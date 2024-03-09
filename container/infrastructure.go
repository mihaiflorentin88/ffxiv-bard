package container

import (
	"ffxvi-bard/infrastructure/client/media"
	database "ffxvi-bard/infrastructure/database/sql"
	migration "ffxvi-bard/infrastructure/database/sql/migration"
	repository "ffxvi-bard/infrastructure/database/sql/repository"
	"ffxvi-bard/infrastructure/filesystem"
	"ffxvi-bard/infrastructure/oauth"
	"ffxvi-bard/port/contract"
	"fmt"
)

func GetDatabaseDriver() contract.DatabaseDriverInterface {
	config := GetConfig()
	driver, err := database.NewSqlDriver(&config.Database)
	if err != nil {
		panic(fmt.Sprintf("Cannot access UserRepository. Reason %s", err))
	}
	return driver
}

func GetMigrationDriver() contract.SqlMigrationDriverInterface {
	config := GetConfig()
	return migration.NewMigrationDriver(&config.Database)
}

func GetFileSystem() contract.FileSystemInterface {
	return filesystem.NewFileSystem()
}

func GetDiscordAuth() contract.Oauth {
	config := GetConfig()
	return oauth.NewDiscordOauth(&config.Discord)
}

func GetUserRepository() contract.UserRepositoryInterface {
	return repository.NewUserRepository(GetDatabaseDriver())
}

func GetGenreRepository() contract.GenreRepositoryInterface {
	return repository.NewGenreRepository(GetDatabaseDriver())
}

func GetSongRepository() contract.SongRepositoryInterface {
	return repository.NewSongRepository(GetDatabaseDriver())
}

func GetRatingRepository() contract.RatingRepositoryInterface {
	return repository.NewRatingRepository(GetDatabaseDriver())
}

func GetCommentRepository() contract.CommentRepositoryInterface {
	return repository.NewCommentRepository(GetDatabaseDriver())
}

func GetSpotifyClient() contract.MediaClientInterface {
	appConfig := GetConfig().Spotify
	return media.NewSpotifyClient(appConfig)
}

func GetInstrumentRepository() contract.InstrumentRepositoryInterface {
	return repository.NewInstrumentRepository(GetDatabaseDriver())
}
