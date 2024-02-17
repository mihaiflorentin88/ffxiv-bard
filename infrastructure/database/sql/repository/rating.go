package database

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
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
		result.Scan(&rating.ID)
		result.Scan(&rating.SongID)
		result.Scan(&rating.AuthorID)
		result.Scan(&rating.Rating)
		result.Scan(&rating.CreatedAt)
		result.Scan(&rating.UpdatedAt)
		ratings = append(ratings, rating)
	}
	if err = result.Err(); err != nil {
		return ratings, err
	}
	return ratings, nil
}
