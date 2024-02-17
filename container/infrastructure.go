package container

import (
	database "ffxvi-bard/infrastructure/database/sql"
	migration "ffxvi-bard/infrastructure/database/sql/migration"
	repository "ffxvi-bard/infrastructure/database/sql/repository"
	"ffxvi-bard/infrastructure/filesystem"
	"ffxvi-bard/infrastructure/oauth"
	"ffxvi-bard/port/contract"
)

func GetDatabaseDriver() (contract.DatabaseDriverInterface, error) {
	config := GetConfig()
	return database.NewSqlDriver(&config.Database)
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

func GetUserRepository() (contract.UserRepositoryInterface, error) {
	driver, err := GetDatabaseDriver()
	if err != nil {
		return nil, err
	}
	return repository.NewUserRepository(driver), nil
}

func GetGenreRepository() (contract.GenreRepositoryInterface, error) {
	driver, err := GetDatabaseDriver()
	if err != nil {
		return nil, err
	}
	return repository.NewGenreRepository(driver), nil
}

func GetSongRepository() (contract.SongRepositoryInterface, error) {
	driver, err := GetDatabaseDriver()
	if err != nil {
		return nil, err
	}
	return repository.NewSongRepository(driver), nil
}

func GetRatingRepository() (contract.RatingRepositoryInterface, error) {
	driver, err := GetDatabaseDriver()
	if err != nil {
		return nil, err
	}
	return repository.NewRatingRepository(driver), nil
}

func GetCommentRepository() (contract.CommentRepositoryInterface, error) {
	driver, err := GetDatabaseDriver()
	if err != nil {
		return nil, err
	}
	return repository.NewCommentRepository(driver), nil
}
