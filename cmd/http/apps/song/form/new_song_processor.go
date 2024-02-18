package form

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"io"
	"mime/multipart"
	"strconv"
)

type AddSongFormProcessor struct {
	Title           string
	Artist          string
	EnsembleSize    int
	Genre           []int
	File            []byte
	User            interface{}
	songRepository  contract.SongRepositoryInterface
	genreRepository contract.GenreRepositoryInterface
	songProcessor   contract.SongProcessorInterface
}

func NewSongFormProcessor(songRepository contract.SongRepositoryInterface, genreRepository contract.GenreRepositoryInterface, songProcessor contract.SongProcessorInterface) AddSongFormProcessor {
	return AddSongFormProcessor{
		songRepository:  songRepository,
		genreRepository: genreRepository,
		songProcessor:   songProcessor,
	}
}

func (s *AddSongFormProcessor) Process(title string, artist string, ensembleSize string, genre []string, fileHeader *multipart.FileHeader, user interface{}) error {
	s.Title = title
	s.Artist = artist
	s.User = user
	ensembleSizeInt, err := strconv.Atoi(ensembleSize)
	if err != nil {
		return err
	}
	s.EnsembleSize = ensembleSizeInt
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	s.File, err = io.ReadAll(file)
	if err != nil {
		return err
	}
	for _, genreStr := range genre {
		genreInt, err := strconv.Atoi(genreStr)
		if err != nil {
			return err
		}
		s.Genre = append(s.Genre, genreInt)
	}
	songDTO := dto.AddNewSong(s.Title, s.Artist, s.EnsembleSize,
		s.Genre, s.File, s.User)
	if err != nil {
		return err
	}
	newSong, err := song.FromNewSongForm(songDTO, s.songRepository, s.genreRepository, s.songProcessor)
	if err != nil {
		return err
	}
	songDatabaseDto := newSong.ToDatabaseSongDTO()
	err = s.songRepository.InsertNewSong(songDatabaseDto, songDTO.Genre)
	if err != nil {
		return err
	}
	return nil
}
