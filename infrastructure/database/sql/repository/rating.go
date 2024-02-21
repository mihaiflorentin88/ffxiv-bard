package database

import (
	"errors"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
)

type RatingRepository struct {
	driver contract.DatabaseDriverInterface
}

func NewRatingRepository(driver contract.DatabaseDriverInterface) contract.RatingRepositoryInterface {
	return &RatingRepository{
		driver: driver,
	}
}

func (r *RatingRepository) FindAllBySongId(songID int) ([]dto.DatabaseRating, error) {
	var ratings []dto.DatabaseRating

	query := `
		SELECT * 
		FROM rating r
		WHERE r.song_id = ?
`
	result, err := r.driver.FetchMany(query, &songID)
	if err != nil {
		return ratings, err
	}
	for result.Next() {
		var rating dto.DatabaseRating
		err := result.Scan(&rating.ID, &rating.SongID, &rating.AuthorID, &rating.Rating, &rating.CreatedAt, &rating.UpdatedAt)
		if err != nil {
			return ratings, errors.New(fmt.Sprintf("cannot retrieve rating for song id `%v`. Reason %s", songID, err))
		}
		ratings = append(ratings, rating)
	}
	if err = result.Err(); err != nil {
		return ratings, err
	}
	return ratings, nil
}
