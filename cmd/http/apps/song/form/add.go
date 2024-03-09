package form

import (
	"errors"
	"ffxvi-bard/domain/song"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
	"strings"
	"sync"
)

type NewSongFormView struct {
	EnsembleSize         map[int]string
	Genres               []dto.DatabaseGenre
	Instruments          []dto.DatabaseInstrument
	genreRepository      contract.GenreRepositoryInterface
	instrumentRepository contract.InstrumentRepositoryInterface
}

func NewAddNewSongFormView(genreRepository contract.GenreRepositoryInterface, instrumentRepository contract.InstrumentRepositoryInterface) NewSongFormView {
	return NewSongFormView{
		genreRepository:      genreRepository,
		instrumentRepository: instrumentRepository,
	}
}

func (f NewSongFormView) Fetch() (NewSongFormView, error) {
	var genres []dto.DatabaseGenre
	var err error
	var instruments []dto.DatabaseInstrument

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		genres, err = f.genreRepository.FetchAll()
		if err != nil {
			errChan <- errors.New(fmt.Sprintf("failed to fetch available genres. Reason: %s", err))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		instruments, err = f.instrumentRepository.FetchAll()
		if err != nil {
			errChan <- errors.New(fmt.Sprintf("failed to fetch available instruments. Reason: %s", err))
		}
	}()
	if err != nil {
		return f, errors.New(fmt.Sprintf("failed to fetch available instruments. Reason: %s", err))
	}

	wg.Wait()

	close(errChan)
	var _errors []string
	for err := range errChan {
		if err != nil {
			_errors = append(_errors, err.Error())
		}
		if len(_errors) > 0 {
			return f, errors.New(strings.Join(_errors, "; "))
		}
	}

	f.Genres = genres
	f.Instruments = instruments
	f.EnsembleSize = song.GetDetailedEnsembleString()
	return f, nil
}
