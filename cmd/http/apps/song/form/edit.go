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
	ID                 int
	Title              string
	Artist             string
	Filename           string
	EnsembleSize       song.EnsembleSize
	EnsembleSizeString string
	AllEnsembleSizes   map[int]string
	Genre              []Genre
	AllGenres          []dto.DatabaseGenre
	CanDelete          bool
	genreRepository    contract.GenreRepositoryInterface
	song               *song.Song
	loggedUser         *user.User
}

func NewSongEditForm(genreRepository contract.GenreRepositoryInterface, song *song.Song) SongEditForm {
	return SongEditForm{
		song:            song,
		genreRepository: genreRepository,
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

func (s *SongEditForm) resetForm() {
	s.ID = 0
	s.Title = ""
	s.Artist = ""
	s.Filename = ""
	s.EnsembleSize = 0
	s.CanDelete = false
	s.EnsembleSizeString = ""
	s.AllEnsembleSizes = make(map[int]string)
	s.Genre = []Genre{}
	s.AllGenres = []dto.DatabaseGenre{}
	s.loggedUser = &user.User{}
}

func (s *SongEditForm) Fetch(songID int, loggedUser *user.User) (SongEditForm, error) {
	s.resetForm()
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
	s.CanDelete = loggedUser.IsAdmin
	s.AllGenres = genres
	s.ID = s.song.StorageID
	s.Title = s.song.Title
	s.Artist = s.song.Artist
	s.EnsembleSize = s.song.EnsembleSize
	s.EnsembleSizeString = s.song.EnsembleString()
	for _, genre := range s.song.Genre {
		s.Genre = append(s.Genre, Genre{ID: genre.StorageID, Name: genre.Name})
	}
	return *s, nil
}

func (s *SongEditForm) HandleSubmittedForm(songID int, title string, artist string, ensembleSize int, genreIDs []int, loggedUser *user.User) error {
	s.resetForm()
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
	if err != nil {
		return err
	}

	s.song.Title = title
	s.song.Artist = artist
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
	err = s.song.Update()
	return err
}
