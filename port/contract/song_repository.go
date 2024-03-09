package contract

import "ffxvi-bard/port/dto"

type SongRepositoryInterface interface {
	InsertNewSong(song dto.DatabaseSong, genreIDs []int, instrumentIDs []int) (int, error)
	UpdateSong(song dto.DatabaseSong, newGenreIDs []int, newInstrumentIDs []int) error
	UpdateStatus(songID int, status int, message string) error
	FindByChecksum(checksum string) (dto.DatabaseSong, error)
	FindByID(songID int) (dto.DatabaseSong, error)
	FetchAll() (*[]dto.DatabaseSong, error)
	FetchForPagination(songTitle string, artist string, ensembleSize int, audioCrafter string, instrumentID int, genreID int, sort string, limit int, offset int) ([]dto.SongWithDetails, error)
	FetchTotalSongsForListing(songTitle string, artist string, ensembleSize int, audioCrafter string, instrumentID int, genreID int) (int, error)
	IncrementDownloadCount(songID int) error
}
