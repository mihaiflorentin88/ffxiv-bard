package form

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
)

type SongList struct {
	Songs          *[]dto.DatabaseSong
	songRepository contract.SongRepositoryInterface
	errorHandler   contract.HttpErrorHandlerInterface
}

func NewSongList(songRepository contract.SongRepositoryInterface, errorHandler contract.HttpErrorHandlerInterface) SongList {
	return SongList{
		songRepository: songRepository,
		errorHandler:   errorHandler,
	}
}

func (s *SongList) FetchData() (*SongList, error) {
	songs, err := s.songRepository.FetchAll()
	if err != nil {
		return nil, err
	}
	s.Songs = songs
	return s, nil
}
