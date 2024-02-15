package contract

import "ffxvi-bard/port/dto"

type GenreRepositoryInterface interface {
	FetchAll() ([]dto.DatabaseGenre, error)
}