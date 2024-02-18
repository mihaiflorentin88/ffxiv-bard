package container

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/song/processor"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
)

func GetNewEmptyUser() *user.User {
	return user.NewEmptyUser(GetUserRepository())
}

func GetMidiProcessor() contract.SongProcessorInterface {
	config := GetConfig().Song
	return processor.NewMidiProcessor(config.UnfinishedFilesPath, config.FinishedFilesPath, GetFileSystem())
}

func GetEmptySong() *song.Song {
	return song.NewEmptySong(GetMidiProcessor(), GetFileSystem())
}

func GetEmptyGenre() song.Genre {
	return song.NewEmptyGenre(GetGenreRepository())
}

func GetEmptyComment() song.Comment {
	return song.NewEmptyComment(GetCommentRepository())
}
