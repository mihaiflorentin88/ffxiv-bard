package contract

import "ffxvi-bard/port/dto"

type SongRepositoryInterface interface {
	InsertNewSong(dto dto.DatabaseSong, genresIDs []int) error
	FindByChecksum(checksum string) (dto.DatabaseSong, error)
	FetchAll() (*[]dto.DatabaseSong, error)
	FetchForPagination(songTitle string, artist string, ensembleSize int, genreID int, sort string, limit int, offset int) ([]dto.SongWithDetails, error)
	FetchTotalSongsForListing(songTitle string, artist string, ensembleSize int, genreID int) (int, error)
}
