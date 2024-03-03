package form

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"io"
	"mime/multipart"
	"strconv"
)

type SubmitSongForm struct {
	Title           string
	Artist          string
	EnsembleSize    int
	Genre           []int
	File            []byte
	User            interface{}
	songRepository  contract.SongRepositoryInterface
	genreRepository contract.GenreRepositoryInterface
	songProcessor   contract.SongProcessorInterface
	emptyUser       user.User
	emptyRating     song.Rating
	emptyGenre      song.Genre
	emptyComment    song.Comment
}

func NewSubmitSongForm(songRepository contract.SongRepositoryInterface, genreRepository contract.GenreRepositoryInterface, songProcessor contract.SongProcessorInterface, emptyUser user.User, emptyGenre song.Genre, emptyRating song.Rating, emptyComment song.Comment) SubmitSongForm {
	return SubmitSongForm{
		songRepository:  songRepository,
		genreRepository: genreRepository,
		songProcessor:   songProcessor,
		emptyUser:       emptyUser,
		emptyRating:     emptyRating,
		emptyGenre:      emptyGenre,
		emptyComment:    emptyComment,
	}
}

func (s *SubmitSongForm) Submit(title string, artist string, ensembleSize string, genre []string, fileHeader *multipart.FileHeader, user interface{}) (int, error) {
	s.Title = title
	s.Artist = artist
	s.User = user
	ensembleSizeInt, err := strconv.Atoi(ensembleSize)
	if err != nil {
		return 0, err
	}
	s.EnsembleSize = ensembleSizeInt
	file, err := fileHeader.Open()
	if err != nil {
		return 0, err
	}
	defer file.Close()
	s.File, err = io.ReadAll(file)
	if err != nil {
		return 0, err
	}
	for _, genreStr := range genre {
		genreInt, err := strconv.Atoi(genreStr)
		if err != nil {
			return 0, err
		}
		s.Genre = append(s.Genre, genreInt)
	}
	songDTO := dto.AddNewSong(s.Title, s.Artist, s.EnsembleSize,
		s.Genre, s.File, s.User)
	if err != nil {
		return 0, err
	}
	newSong, err := song.FromNewSongForm(songDTO, s.songRepository, s.genreRepository, s.songProcessor, &s.emptyUser, &s.emptyRating, &s.emptyGenre, &s.emptyComment)
	if err != nil {
		return 0, err
	}
	songDatabaseDto := newSong.ToDatabaseSongDTO()
	songID, err := s.songRepository.InsertNewSong(songDatabaseDto, songDTO.Genre)
	if err != nil {
		return 0, err
	}
	songDatabaseDto.ID = songID
	newSong, err = song.FromDatabaseDTO(newSong, &songDatabaseDto)
	if err != nil {
		return songID, err
	}
	err = newSong.ProcessSong()
	return songID, err
}
