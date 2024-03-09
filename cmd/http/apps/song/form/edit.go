package form

import (
	"errors"
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
)

type SongEditForm struct {
	ID                   int
	Title                string
	Artist               string
	Source               string
	Note                 string
	AudioCrafter         string
	Filename             string
	EnsembleSize         song.EnsembleSize
	EnsembleSizeString   string
	AllEnsembleSizes     map[int]string
	Genre                []Genre
	Instrument           []Instrument
	AllGenres            []dto.DatabaseGenre
	AllInstruments       []dto.DatabaseInstrument
	CanDelete            bool
	genreRepository      contract.GenreRepositoryInterface
	instrumentRepository contract.InstrumentRepositoryInterface
	song                 *song.Song
	loggedUser           user.User
}

func NewSongEditForm(genreRepository contract.GenreRepositoryInterface, instrumentRepository contract.InstrumentRepositoryInterface, song *song.Song) SongEditForm {
	return SongEditForm{
		song:                 song,
		genreRepository:      genreRepository,
		instrumentRepository: instrumentRepository,
	}
}

func (s *SongEditForm) ContainsGenre(id int) bool {
	for _, item := range s.Genre {
		if item.ID == id {
			return true
		}
	}
	return false
}

func (s *SongEditForm) Fetch(songID int, loggedUser user.User) (SongEditForm, error) {
	//s.resetForm()
	_, err := s.song.LoadByID(songID)
	if err != nil {
		return *s, err
	}

	if s.song.Uploader.StorageID != loggedUser.StorageID {
		if !loggedUser.IsAdmin {
			return *s, errors.New("you do not have permissions to edit this song")
		}
	}

	s.loggedUser = loggedUser
	s.AllEnsembleSizes = song.GetDetailedEnsembleString()
	genres, err := s.genreRepository.FetchAll()
	if err != nil {
		return *s, errors.New(fmt.Sprintf("failed to fetch available genres. Reason: %s", err))
	}
	instruments, err := s.instrumentRepository.FetchAll()
	if err != nil {
		return *s, errors.New(fmt.Sprintf("failed to fetch available instruments. Reason: %s", err))
	}
	s.CanDelete = loggedUser.IsAdmin
	s.AllGenres = genres
	s.AllInstruments = instruments
	s.Source = s.song.Source
	s.Note = s.song.Note
	s.AudioCrafter = s.song.AudioCrafter
	s.ID = s.song.StorageID
	s.Title = s.song.Title
	s.Artist = s.song.Artist
	s.EnsembleSize = s.song.EnsembleSize
	s.EnsembleSizeString = s.song.EnsembleString()
	for _, genre := range s.song.Genre {
		s.Genre = append(s.Genre, Genre{ID: genre.StorageID, Name: genre.Name})
	}
	for _, instrument := range s.song.Instrument {
		s.Instrument = append(s.Instrument, Instrument{ID: instrument.StorageID, Name: instrument.Name})
	}
	return *s, nil
}

func (s *SongEditForm) HandleSubmittedForm(songID int, title string, artist string, ensembleSize int, genreIDs []int, loggedUser user.User, source string, note string, audioCrafter string, instrumentIDs []int) error {
	//s.resetForm()
	_, err := s.song.LoadByID(songID)
	if err != nil {
		return err
	}

	if s.song.Uploader.StorageID != loggedUser.StorageID {
		if loggedUser.IsAdmin {
		} else {
			return errors.New("you do not have permissions to edit this song")
		}
	}
	s.loggedUser = loggedUser

	genresDTOs, err := s.genreRepository.FetchAll()
	instrumentDTOs, err := s.instrumentRepository.FetchAll()
	if err != nil {
		return err
	}

	s.song.Title = title
	s.song.Artist = artist
	s.song.Source = source
	s.song.Note = note
	s.song.AudioCrafter = audioCrafter
	s.song.EnsembleSize, err = song.GetEnsembleSizeFromInt(ensembleSize)
	if err != nil {
		return err
	}

	s.song.Genre = []song.Genre{}
	for _, genreID := range genreIDs {
		for _, genreDTO := range genresDTOs {
			if genreID == genreDTO.ID {
				s.song.Genre = append(s.song.Genre, song.Genre{
					StorageID: genreDTO.ID,
					Name:      genreDTO.Name,
				})
			}
		}
	}
	s.song.Instrument = []song.Instrument{}
	for _, instrumentID := range instrumentIDs {
		for _, instrumentDTO := range instrumentDTOs {
			if instrumentID == instrumentDTO.ID {
				s.song.Instrument = append(s.song.Instrument, song.Instrument{
					StorageID: instrumentDTO.ID,
					Name:      instrumentDTO.Name,
				})
			}
		}
	}
	err = s.song.Update()
	return err
}
