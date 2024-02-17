package dto

import "time"

type DatabaseComment struct {
	ID        int       `db:"id"`
	AuthorID  int       `db:"author_id"`
	SongID    int       `db:"song_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Likes     int       `db:"likes"`
	Status    bool      `db:"status"`
	Dislikes  int       `db:"dislikes"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
