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
	song            *song.Song
}

func NewAddNewSongFormView(song *song.Song, genreRepository contract.GenreRepositoryInterface) NewSongFormView {
	return NewSongFormView{
		genreRepository: genreRepository,
		song:            song,
	}

}

func (f NewSongFormView) GetData() (NewSongFormView, error) {
	genres, err := f.genreRepository.FetchAll()
	if err != nil {
		return f, errors.New(fmt.Sprintf("failed to fetch available genres. Reason: %s", err))
	}
	f.Genres = genres
	f.EnsembleSize = f.song.GetDetailedEnsembleString()
	return f, nil
}
