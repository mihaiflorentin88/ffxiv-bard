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
	Source          string
	Note            string
	AudioCrafter    string
	Instrument      []int
	songRepository  contract.SongRepositoryInterface
	genreRepository contract.GenreRepositoryInterface
	songProcessor   contract.SongProcessorInterface
	emptyUser       user.User
	emptyRating     song.Rating
	emptyGenre      song.Genre
	emptyComment    song.Comment
	emptyInstrument song.Instrument
}

func NewSubmitSongForm(songRepository contract.SongRepositoryInterface, genreRepository contract.GenreRepositoryInterface, songProcessor contract.SongProcessorInterface, emptyUser user.User, emptyGenre song.Genre, emptyRating song.Rating, emptyComment song.Comment, emptyInstrument song.Instrument) SubmitSongForm {
	return SubmitSongForm{
		songRepository:  songRepository,
		genreRepository: genreRepository,
		songProcessor:   songProcessor,
		emptyUser:       emptyUser,
		emptyRating:     emptyRating,
		emptyGenre:      emptyGenre,
		emptyComment:    emptyComment,
		emptyInstrument: emptyInstrument,
	}
}

func (s *SubmitSongForm) Submit(title string, artist string, ensembleSize string, genre []string,
	fileHeader *multipart.FileHeader, user interface{}, source string, note string,
	audioCrafter string, instrument []string) (int, error) {
	s.Title = title
	s.Artist = artist
	s.User = user
	s.Source = source
	s.Note = note
	s.AudioCrafter = audioCrafter
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
	for _, instrumentStr := range instrument {
		instrumentInt, err := strconv.Atoi(instrumentStr)
		if err != nil {
			return 0, err
		}
		s.Instrument = append(s.Instrument, instrumentInt)
	}
	songDTO := dto.AddNewSong(s.Title, s.Artist, s.EnsembleSize, s.Genre, s.File, s.User, s.Source, s.Note, s.AudioCrafter, s.Instrument)
	if err != nil {
		return 0, err
	}
	newSong, err := song.FromNewSongForm(songDTO, s.songRepository, s.genreRepository, s.songProcessor, &s.emptyUser, &s.emptyRating, &s.emptyGenre, &s.emptyComment, &s.emptyInstrument)
	if err != nil {
		return 0, err
	}
	songDatabaseDto := newSong.ToDatabaseSongDTO()
	songID, err := s.songRepository.InsertNewSong(songDatabaseDto, songDTO.Genre, songDTO.Instrument)
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
