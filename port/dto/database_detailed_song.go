package dto

type SongWithDetails struct {
	ID                 int     `db:"id"`
	Title              string  `db:"title"`
	Artist             string  `db:"artist"`
	EnsembleSize       int     `db:"ensemble_size"`
	AverageRating      float64 `db:"average_rating"`
	GenreName          string  `db:"genre_name"`
	UploaderName       string  `db:"uploader_name"`
	TotalComments      int     `db:"total_comments"`
	EnsembleSizeString string
	ImageUrl           string
}
