package form

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SongListForm struct {
	Songs          *[]dto.DatabaseSong
	songRepository contract.SongRepositoryInterface
	errorHandler   contract.HttpErrorHandlerInterface
}

func NewSongList(songRepository contract.SongRepositoryInterface, errorHandler contract.HttpErrorHandlerInterface) *SongListForm {
	return &SongListForm{
		songRepository: songRepository,
		errorHandler:   errorHandler,
	}
}

func (s *SongListForm) FetchData(c *gin.Context) (*SongListForm, error) {
	songs, err := s.songRepository.FetchAll()
	if err != nil {
		s.errorHandler.RenderTemplate(err, http.StatusInternalServerError, c)
		return nil, err
	}
	s.Songs = songs
	return s, nil
}
