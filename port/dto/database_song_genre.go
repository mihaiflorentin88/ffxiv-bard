package dto

type DatabaseSongGenre struct {
	SongID  int `db:"song_id"`
	GenreID int `db:"genre_id"`
}
