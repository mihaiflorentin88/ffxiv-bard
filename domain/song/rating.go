package song

import (
	"errors"
	"ffxvi-bard/domain/date"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
)

type Rating struct {
	StorageId        int
	SongID           int
	AuthorID         int
	rating           int
	Date             date.Date
	ratingRepository contract.RatingRepositoryInterface
	emptyUser        *user.User
}

func NewSongRanting(songID int, authorID int, rating int) (*Rating, error) {
	if rating < 0 || rating > 5 {
		return nil, errors.New("rating must be between 0 and 10")
	}
	return &Rating{SongID: songID, AuthorID: authorID, rating: rating}, nil
}

func NewEmptyRating(ratingRepository contract.RatingRepositoryInterface, emptyUser *user.User) *Rating {
	return &Rating{
		ratingRepository: ratingRepository,
		emptyUser:        emptyUser,
	}
}

func (r *Rating) SetRating(rating int) error {
	if rating < 0 || rating > 10 {
		return errors.New("rating must be between 0 and 10")
	}
	r.rating = rating
	return nil
}

func (r *Rating) GetRating() int {
	return r.rating
}

func (r *Rating) FetchBySongId(songID int) ([]Rating, error) {
	if r.ratingRepository == nil {
		return nil, errors.New("song.Rating was not correctly instantiated")
	}
	var ratings []Rating
	ratingDTOs, err := r.ratingRepository.FindAllBySongId(songID)
	if err != nil {
		return nil, err
	}
	for _, ratingDTO := range ratingDTOs {
		rating, err := FromRatingDTO(ratingDTO)
		if err != nil {
			return nil, err
		}
		ratings = append(ratings, rating)
	}
	return ratings, nil
}

func FromRatingDTO(ratingDTO dto.DatabaseRating) (Rating, error) {
	rating := Rating{
		StorageId: ratingDTO.ID,
		SongID:    ratingDTO.SongID,
		AuthorID:  ratingDTO.AuthorID,
	}
	err := rating.SetRating(ratingDTO.Rating)
	if err != nil {
		return rating, err
	}
	rating.Date.CreatedAt = ratingDTO.CreatedAt
	rating.Date.UpdatedAt = ratingDTO.UpdatedAt
	return rating, nil

}
