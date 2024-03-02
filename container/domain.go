package container

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/song/processor"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
)

func GetEmptyUser() *user.User {
	return user.NewEmptyUser(GetUserRepository())
}

func GetMidiProcessor() contract.SongProcessorInterface {
	config := GetConfig().Song
	return processor.NewMidiProcessor(config.UnfinishedFilesPath, config.FinishedFilesPath, GetFileSystem())
}

func GetEmptySong() *song.Song {
	return song.NewEmptySong(GetMidiProcessor(), GetFileSystem(), GetEmptyUser(), GetEmptyRating(), GetEmptyComment(), GetEmptyGenre(), GetSongRepository())
}

func GetEmptyRating() *song.Rating {
	return song.NewEmptyRating(GetRatingRepository(), GetEmptyUser())
}

func GetEmptyGenre() *song.Genre {
	genre := song.NewEmptyGenre(GetGenreRepository())
	return &genre
}

func GetEmptyComment() *song.Comment {
	emptyUser := GetEmptyUser()
	comment := song.NewEmptyComment(GetCommentRepository(), *emptyUser)
	return &comment
}
