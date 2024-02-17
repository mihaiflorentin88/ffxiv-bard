package container

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/song/processor"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"fmt"
)

func GetNewEmptyUser() *user.User {
	userRepository, err := GetUserRepository()
	if err != nil {
		panic(fmt.Sprintf("Cannot instantiate the UserRepository. Reason: %s", err))
	}
	return user.NewEmptyUser(userRepository)
}

func GetMidiProcessor() contract.SongProcessorInterface {
	config := GetConfig().Song
	return processor.NewMidiProcessor(config.UnfinishedFilesPath, config.FinishedFilesPath, GetFileSystem())
}

func GetEmptySong() contract.SongInterface {
	return song.NewEmptySong(GetMidiProcessor(), GetFileSystem())
}

func GetEmptyGenre() song.Genre {
	repository, err := GetGenreRepository()
	if err != nil {
		panic(fmt.Sprintf("Cannot instantiate the GenreRepository. Reason: %s", err))
	}
	return song.NewEmptyGenre(repository)
}

func GetEmptyComment() song.Comment {
	repository, err := GetCommentRepository()
	if err != nil {
		panic(fmt.Sprintf("Cannot instantiate the GenreRepository. Reason: %s", err))
	}
	return song.NewEmptyComment(repository)
}
