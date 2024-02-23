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

func (r *RatingRepository) FindByUserAndSong(songID int, userID int64) (int, error) {
	var rating int
	query := `
		SELECT r.rating 
		FROM rating r
		WHERE r.song_id = ? AND r.author_id = ?
	`
	result, err := r.driver.FetchOne(query, songID, userID)
	if err != nil {
		return 0, err
	}
	result.Scan(&rating)
	return rating, nil
}

func (r *RatingRepository) UpdateRating(songID int, userID int64, newRating int) error {
	query := `
		UPDATE rating 
		SET rating = ?, updated_at = CURRENT_TIMESTAMP
		WHERE song_id = ? AND author_id = ?
	`
	_, err := r.driver.Execute(query, newRating, songID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RatingRepository) InsertRating(songID int, userID int64, rating int) error {
	query := `
		INSERT INTO rating (song_id, author_id, rating) 
		VALUES (?, ?, ?)
	`
	_, err := r.driver.Execute(query, songID, userID, rating)
	if err != nil {
		return err
	}
	return nil
}
