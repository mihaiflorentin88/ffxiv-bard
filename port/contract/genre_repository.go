package contract

import "ffxvi-bard/port/dto"

type GenreRepositoryInterface interface {
	FetchAll() ([]dto.DatabaseGenre, error)
	FetchByIDs(genreIDs []int) ([]dto.DatabaseGenre, error)
	FetchBySongID(songID int) (*[]dto.DatabaseGenre, error)
}
