package dto

import "time"

type DatabaseComment struct {
	ID        int       `db:"id"`
	AuthorID  int64     `db:"author_id"`
	SongID    int       `db:"song_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Status    bool      `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
