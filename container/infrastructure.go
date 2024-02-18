package container

import (
	database "ffxvi-bard/infrastructure/database/sql"
	migration "ffxvi-bard/infrastructure/database/sql/migration"
	repository "ffxvi-bard/infrastructure/database/sql/repository"
	"ffxvi-bard/infrastructure/filesystem"
	"ffxvi-bard/infrastructure/oauth"
	"ffxvi-bard/port/contract"
)

func GetDatabaseDriver() contract.DatabaseDriverInterface {
	config := GetConfig()
	driver, err := database.NewSqlDriver(&config.Database)
	if err != nil {
		panic("Cannot access UserRepository.")
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
