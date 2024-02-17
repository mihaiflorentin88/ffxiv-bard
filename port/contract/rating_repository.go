package contract

import "ffxvi-bard/port/dto"

type RatingRepositoryInterface interface {
	FindAllBySongId(songID int) ([]dto.DatabaseRating, error)
}
