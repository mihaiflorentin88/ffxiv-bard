package contract

import "ffxvi-bard/port/dto"

type SongRepositoryInterface interface {
	InsertNewSong(dto dto.DatabaseSongDTO, genresIDs []int) error
	FindByChecksum(checksum string) (dto.DatabaseSongDTO, error)
}
