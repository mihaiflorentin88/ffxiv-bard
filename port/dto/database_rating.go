package dto

import "time"

type DatabaseRating struct {
	ID        int       `db:"id"`
	SongID    int       `db:"song_id"`
	AuthorID  int       `db:"author_id"`
	Rating    int       `db:"rating"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
