package container

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/song/processor"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
)

type DomainContainer struct {
	emptyUser       *user.User
	emptySong       *song.Song
	emptySongRating *song.Rating
	emptyGenre      *song.Genre
	emptyInstrument *song.Instrument
	emptyComment    *song.Comment
	midiProcessor   *processor.MidiProcessor
}

func (s *ServiceContainer) GetEmptyUser() user.User {
	if s.domain.emptyUser != nil {
		return *s.domain.emptyUser
	}
	emptyUser := user.NewEmptyUser(s.GetUserRepository())
	s.domain.emptyUser = &emptyUser
	return emptyUser
}

func (s *ServiceContainer) GetMidiProcessor() contract.SongProcessorInterface {
	if s.domain.midiProcessor != nil {
		return *s.domain.midiProcessor
	}
	config := s.GetConfig().Song
	midiProcessor := processor.NewMidiProcessor(config.UnfinishedFilesPath, config.FinishedFilesPath, s.GetFileSystem())
	s.domain.midiProcessor = &midiProcessor
	return midiProcessor
}

func (s *ServiceContainer) GetEmptySong() song.Song {
	if s.domain.emptySong != nil {
		return *s.domain.emptySong
	}
	emptySong := song.NewEmptySong(s.GetMidiProcessor(), s.GetFileSystem(), s.GetEmptyUser(), s.GetEmptyRating(), s.GetEmptyComment(), s.GetEmptyGenre(), s.GetEmptyInstrument(), s.GetSongRepository())
	s.domain.emptySong = &emptySong
	return emptySong
}

func (s *ServiceContainer) GetEmptyRating() song.Rating {
	if s.domain.emptySongRating != nil {
		return *s.domain.emptySongRating
	}
	songRating := song.NewEmptyRating(s.GetRatingRepository(), s.GetEmptyUser())
	s.domain.emptySongRating = &songRating
	return songRating
}

func (s *ServiceContainer) GetEmptyGenre() song.Genre {
	if s.domain.emptyGenre != nil {
		return *s.domain.emptyGenre
	}
	genre := song.NewEmptyGenre(s.GetGenreRepository())
	s.domain.emptyGenre = &genre
	return genre
}

func (s *ServiceContainer) GetEmptyInstrument() song.Instrument {
	if s.domain.emptyInstrument != nil {
		return *s.domain.emptyInstrument
	}
	instrument := song.NewEmptyInstrument(s.GetInstrumentRepository())
	s.domain.emptyInstrument = &instrument
	return instrument
}

func (s *ServiceContainer) GetEmptyComment() song.Comment {
	if s.domain.emptyComment != nil {
		return *s.domain.emptyComment
	}
	comment := song.NewEmptyComment(s.GetCommentRepository(), s.GetEmptyUser())
	s.domain.emptyComment = &comment
	return comment
}
