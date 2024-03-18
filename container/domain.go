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

func (s *ServiceContainer) EmptyUser() user.User {
	if s.domain.emptyUser != nil {
		return *s.domain.emptyUser
	}
	emptyUser := user.NewEmptyUser(s.UserRepository())
	s.domain.emptyUser = &emptyUser
	return emptyUser
}

func (s *ServiceContainer) MidiProcessor() contract.SongProcessorInterface {
	if s.domain.midiProcessor != nil {
		return *s.domain.midiProcessor
	}
	config := s.Config().Song
	midiProcessor := processor.NewMidiProcessor(config.UnfinishedFilesPath, config.FinishedFilesPath, s.FileSystem())
	s.domain.midiProcessor = &midiProcessor
	return midiProcessor
}

func (s *ServiceContainer) EmptySong() song.Song {
	if s.domain.emptySong != nil {
		return *s.domain.emptySong
	}
	emptySong := song.NewEmptySong(s.MidiProcessor(), s.FileSystem(), s.EmptyUser(), s.EmptyRating(), s.EmptyComment(), s.EmptyGenre(), s.EmptyInstrument(), s.SongRepository())
	s.domain.emptySong = &emptySong
	return emptySong
}

func (s *ServiceContainer) EmptyRating() song.Rating {
	if s.domain.emptySongRating != nil {
		return *s.domain.emptySongRating
	}
	songRating := song.NewEmptyRating(s.RatingRepository(), s.EmptyUser())
	s.domain.emptySongRating = &songRating
	return songRating
}

func (s *ServiceContainer) EmptyGenre() song.Genre {
	if s.domain.emptyGenre != nil {
		return *s.domain.emptyGenre
	}
	genre := song.NewEmptyGenre(s.GenreRepository())
	s.domain.emptyGenre = &genre
	return genre
}

func (s *ServiceContainer) EmptyInstrument() song.Instrument {
	if s.domain.emptyInstrument != nil {
		return *s.domain.emptyInstrument
	}
	instrument := song.NewEmptyInstrument(s.InstrumentRepository())
	s.domain.emptyInstrument = &instrument
	return instrument
}

func (s *ServiceContainer) EmptyComment() song.Comment {
	if s.domain.emptyComment != nil {
		return *s.domain.emptyComment
	}
	comment := song.NewEmptyComment(s.CommentRepository(), s.EmptyUser())
	s.domain.emptyComment = &comment
	return comment
}
