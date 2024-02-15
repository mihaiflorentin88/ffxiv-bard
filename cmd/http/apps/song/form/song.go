package form

import (
	"errors"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
)

type NewSongFormView struct {
	EnsembleSize map[int]string
	Genres       []dto.DatabaseGenre
}

func NewAddNewSongFormView(song contract.SongInterface, genreRepository contract.GenreRepositoryInterface) (NewSongFormView, error) {
	form := NewSongFormView{}
	genres, err := genreRepository.FetchAll()
	if err != nil {
		return form, errors.New(fmt.Sprintf("failed to fetch available genres. Reason: %s", err))
	}
	form.Genres = genres
	form.EnsembleSize = song.GetDetailedEnsembleString()
	return form, nil
}
