package contract

import "ffxvi-bard/port/dto"

type InstrumentRepositoryInterface interface {
	FetchAll() ([]dto.DatabaseInstrument, error)
	FetchByIDs(instrumentIDs []int) ([]dto.DatabaseInstrument, error)
	FetchBySongID(songID int) (*[]dto.DatabaseInstrument, error)
}
