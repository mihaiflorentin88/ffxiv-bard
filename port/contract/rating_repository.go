package contract

import "ffxvi-bard/port/dto"

type RatingRepositoryInterface interface {
	FindAllBySongId(songID int) ([]dto.DatabaseRating, error)
	FindByUserAndSong(songID int, userID int64) (int, error)
	UpdateRating(songID int, userID int64, newRating int) error
	InsertRating(songID int, userID int64, rating int) error
}
