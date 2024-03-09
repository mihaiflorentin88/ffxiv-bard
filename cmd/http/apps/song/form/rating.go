package form

import (
	"errors"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"fmt"
)

type SubmitSongRatingForm struct {
	ratingRepository contract.RatingRepositoryInterface
}

func NewSubmitSongRatingForm(ratingRepository contract.RatingRepositoryInterface) SubmitSongRatingForm {
	return SubmitSongRatingForm{
		ratingRepository: ratingRepository,
	}
}

func (p *SubmitSongRatingForm) Submit(loggedUser user.User, songID int, rating int) error {
	existingRating, _ := p.ratingRepository.FindByUserAndSong(songID, loggedUser.StorageID)
	if existingRating == 0 {
		err := p.ratingRepository.InsertRating(songID, loggedUser.StorageID, rating)
		if err != nil {
			return errors.New(fmt.Sprintf("cannot add new rating. Reason: %s", err))
		}
		return nil
	}
	err := p.ratingRepository.UpdateRating(songID, loggedUser.StorageID, rating)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot update rating. Reason: %s", err))
	}
	return nil
}
