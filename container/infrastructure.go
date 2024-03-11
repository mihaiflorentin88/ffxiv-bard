package container

import (
	database "ffxvi-bard/infrastructure/database/sql"
	migration "ffxvi-bard/infrastructure/database/sql/migration"
	repository "ffxvi-bard/infrastructure/database/sql/repository"
	"ffxvi-bard/infrastructure/filesystem"
	"ffxvi-bard/infrastructure/oauth"
	"ffxvi-bard/port/contract"
	"fmt"
)

type InfrastructureContainer struct {
	databaseDriver       *database.SqliteDriver
	migrationDriver      *migration.MigrationDriver
	filesystem           *filesystem.FileSystem
	oauth                *oauth.DiscordOauth
	userRepository       *repository.UserRepository
	genreRepository      *repository.GenreRepository
	songRepository       *repository.SongRepository
	ratingRepository     *repository.RatingRepository
	commentRepository    *repository.CommentRepository
	instrumentRepository *repository.InstrumentRepository
}

func (s *ServiceContainer) GetDatabaseDriver() *database.SqliteDriver {
	if s.infrastructure.databaseDriver != nil {
		return s.infrastructure.databaseDriver
	}
	config := s.GetConfig()
	driver, err := database.NewSqlDriver(&config.Database)
	if err != nil {
		panic(fmt.Sprintf("Cannot access UserRepository. Reason %s", err))
	}
	s.infrastructure.databaseDriver = driver
	return driver
}

func (s *ServiceContainer) GetMigrationDriver() *migration.MigrationDriver {
	if s.infrastructure.migrationDriver != nil {
		return s.infrastructure.migrationDriver
	}
	config := s.GetConfig()
	driver := migration.NewMigrationDriver(&config.Database)
	s.infrastructure.migrationDriver = driver
	return driver
}

func (s *ServiceContainer) GetFileSystem() *filesystem.FileSystem {
	if s.infrastructure.filesystem != nil {
		return s.infrastructure.filesystem
	}
	_filesystem := filesystem.NewFileSystem()
	s.infrastructure.filesystem = _filesystem
	return _filesystem
}

func (s *ServiceContainer) GetDiscordAuth() *oauth.DiscordOauth {
	if s.infrastructure.oauth != nil {
		return s.infrastructure.oauth
	}
	config := s.GetConfig()
	discordOauth := oauth.NewDiscordOauth(&config.Discord)
	s.infrastructure.oauth = discordOauth
	return discordOauth
}

func (s *ServiceContainer) GetUserRepository() *repository.UserRepository {
	if s.infrastructure.userRepository != nil {
		return s.infrastructure.userRepository
	}
	userRepository := repository.NewUserRepository(s.GetDatabaseDriver())
	s.infrastructure.userRepository = userRepository
	return userRepository
}

func (s *ServiceContainer) GetGenreRepository() *repository.GenreRepository {
	if s.infrastructure.genreRepository != nil {
		return s.infrastructure.genreRepository
	}
	genreRepository := repository.NewGenreRepository(s.GetDatabaseDriver())
	s.infrastructure.genreRepository = genreRepository
	return genreRepository
}

func (s *ServiceContainer) GetSongRepository() contract.SongRepositoryInterface {
	if s.infrastructure.songRepository != nil {
		return s.infrastructure.songRepository
	}
	songRepository := repository.NewSongRepository(s.GetDatabaseDriver())
	s.infrastructure.songRepository = songRepository
	return songRepository
}

func (s *ServiceContainer) GetRatingRepository() *repository.RatingRepository {
	if s.infrastructure.ratingRepository != nil {
		return s.infrastructure.ratingRepository
	}
	ratingRepository := repository.NewRatingRepository(s.GetDatabaseDriver())
	s.infrastructure.ratingRepository = ratingRepository
	return ratingRepository
}

func (s *ServiceContainer) GetCommentRepository() *repository.CommentRepository {
	if s.infrastructure.commentRepository != nil {
		return s.infrastructure.commentRepository
	}
	commentRepository := repository.NewCommentRepository(s.GetDatabaseDriver())
	s.infrastructure.commentRepository = commentRepository
	return commentRepository
}

func (s *ServiceContainer) GetInstrumentRepository() *repository.InstrumentRepository {
	if s.infrastructure.instrumentRepository != nil {
		return s.infrastructure.instrumentRepository
	}

	instrumentRepository := repository.NewInstrumentRepository(s.GetDatabaseDriver())
	s.infrastructure.instrumentRepository = instrumentRepository
	return instrumentRepository
}
