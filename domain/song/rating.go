package song

import (
	"errors"
	"ffxvi-bard/domain/date"
	"ffxvi-bard/domain/user"
)

type Rating struct {
	StorageId int
	Song      *Song
	Author    *user.User
	rating    int
	Date      date.Date
}

func NewSongRanking(song *Song, user *user.User, rating int) (*Rating, error) {
	if rating < 0 || rating > 10 {
		return nil, errors.New("rating must be between 0 and 10")
	}
	return &Rating{Song: song, Author: user, rating: rating}, nil
}

func (r *Rating) SetRanking(rating int) error {
	if rating < 0 || rating > 10 {
		return errors.New("rating must be between 0 and 10")
	}
	r.rating = rating
	return nil
}

func (r *Rating) GetRating() int {
	return r.rating
}
