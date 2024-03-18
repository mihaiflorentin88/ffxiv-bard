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

func (s *ServiceContainer) DatabaseDriver() *database.SqliteDriver {
	if s.infrastructure.databaseDriver != nil {
		return s.infrastructure.databaseDriver
	}
	config := s.Config()
	driver, err := database.NewSqlDriver(&config.Database)
	if err != nil {
		panic(fmt.Sprintf("Cannot access UserRepository. Reason %s", err))
	}
	s.infrastructure.databaseDriver = driver
	return driver
}

func (s *ServiceContainer) MigrationDriver() *migration.MigrationDriver {
	if s.infrastructure.migrationDriver != nil {
		return s.infrastructure.migrationDriver
	}
	config := s.Config()
	driver := migration.NewMigrationDriver(&config.Database)
	s.infrastructure.migrationDriver = driver
	return driver
}

func (s *ServiceContainer) FileSystem() *filesystem.FileSystem {
	if s.infrastructure.filesystem != nil {
		return s.infrastructure.filesystem
	}
	_filesystem := filesystem.NewFileSystem()
	s.infrastructure.filesystem = _filesystem
	return _filesystem
}

func (s *ServiceContainer) DiscordAuth() *oauth.DiscordOauth {
	if s.infrastructure.oauth != nil {
		return s.infrastructure.oauth
	}
	config := s.Config()
	discordOauth := oauth.NewDiscordOauth(&config.Discord)
	s.infrastructure.oauth = discordOauth
	return discordOauth
}

func (s *ServiceContainer) UserRepository() *repository.UserRepository {
	if s.infrastructure.userRepository != nil {
		return s.infrastructure.userRepository
	}
	userRepository := repository.NewUserRepository(s.DatabaseDriver())
	s.infrastructure.userRepository = userRepository
	return userRepository
}

func (s *ServiceContainer) GenreRepository() *repository.GenreRepository {
	if s.infrastructure.genreRepository != nil {
		return s.infrastructure.genreRepository
	}
	genreRepository := repository.NewGenreRepository(s.DatabaseDriver())
	s.infrastructure.genreRepository = genreRepository
	return genreRepository
}

func (s *ServiceContainer) SongRepository() contract.SongRepositoryInterface {
	if s.infrastructure.songRepository != nil {
		return s.infrastructure.songRepository
	}
	songRepository := repository.NewSongRepository(s.DatabaseDriver())
	s.infrastructure.songRepository = songRepository
	return songRepository
}

func (s *ServiceContainer) RatingRepository() *repository.RatingRepository {
	if s.infrastructure.ratingRepository != nil {
		return s.infrastructure.ratingRepository
	}
	ratingRepository := repository.NewRatingRepository(s.DatabaseDriver())
	s.infrastructure.ratingRepository = ratingRepository
	return ratingRepository
}

func (s *ServiceContainer) CommentRepository() *repository.CommentRepository {
	if s.infrastructure.commentRepository != nil {
		return s.infrastructure.commentRepository
	}
	commentRepository := repository.NewCommentRepository(s.DatabaseDriver())
	s.infrastructure.commentRepository = commentRepository
	return commentRepository
}

func (s *ServiceContainer) InstrumentRepository() *repository.InstrumentRepository {
	if s.infrastructure.instrumentRepository != nil {
		return s.infrastructure.instrumentRepository
	}

	instrumentRepository := repository.NewInstrumentRepository(s.DatabaseDriver())
	s.infrastructure.instrumentRepository = instrumentRepository
	return instrumentRepository
}
