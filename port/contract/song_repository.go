package contract

import "ffxvi-bard/port/dto"

type SongRepositoryInterface interface {
	InsertNewSong(dto dto.DatabaseSong, genresIDs []int) error
	FindByChecksum(checksum string) (dto.DatabaseSong, error)
	FetchAll() (*[]dto.DatabaseSong, error)
}
