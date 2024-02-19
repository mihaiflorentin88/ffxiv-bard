package form

import (
	"errors"
	"ffxvi-bard/domain/song"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
)

type NewSongFormView struct {
	EnsembleSize    map[int]string
	Genres          []dto.DatabaseGenre
	genreRepository contract.GenreRepositoryInterface
}

func NewAddNewSongFormView(genreRepository contract.GenreRepositoryInterface) NewSongFormView {
	return NewSongFormView{
		genreRepository: genreRepository,
	}
}

func (f NewSongFormView) Fetch() (NewSongFormView, error) {
	genres, err := f.genreRepository.FetchAll()
	if err != nil {
		return f, errors.New(fmt.Sprintf("failed to fetch available genres. Reason: %s", err))
	}
	f.Genres = genres
	f.EnsembleSize = song.GetDetailedEnsembleString()
	return f, nil
}
