package song

import (
	"errors"
	"ffxvi-bard/domain/date"
	"ffxvi-bard/domain/user"
)

type Rating struct {
	storageId int
	Song      contract.SongInterface
	Author    user.User
	rating    int
	Date      date.Date
}

func NewSongRanking(song contract.SongInterface, user user.User, rating int) (*Rating, error) {
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

func (r *Rating) GetRanking() int {
	return r.rating
}

func (r *Rating) GetStorageID() int {
	return r.storageId
}

func (r *Rating) SetStorageID(id int) {
	r.storageId = id
}

func (r *Rating) GetSong() contract.SongInterface {
	return r.Song
}

func (r *Rating) GetAuthor() user.User {
	return r.Author
}

func (r *Rating) SetSong(song contract.SongInterface) {
	r.Song = song
}

func (r *Rating) SetAuthor(author user.User) {
	r.Author = author
}

func (r *Rating) Like() {
	r.rating++
}

func (r *Rating) Dislike() {
	r.rating--
}
