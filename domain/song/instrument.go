package song

import (
	"ffxvi-bard/domain/date"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
)

type Instrument struct {
	StorageID            int
	Name                 string
	Date                 date.Date
	instrumentRepository contract.InstrumentRepositoryInterface
}

func NewEmptyInstrument(instrumentRepository contract.InstrumentRepositoryInterface) Instrument {
	return Instrument{
		instrumentRepository: instrumentRepository,
	}
}

func (g *Instrument) FetchBySongID(songID int) ([]Instrument, error) {
	var instruments []Instrument
	instrumentDTOs, err := g.instrumentRepository.FetchBySongID(songID)
	if err != nil {
		return instruments, err
	}
	for _, instrumentDTO := range *instrumentDTOs {
		instrument := FromInstrumentDatabaseDTO(instrumentDTO)
		instrument.instrumentRepository = g.instrumentRepository
		instruments = append(instruments, instrument)
	}
	return instruments, nil
}

func FromInstrumentDatabaseDTO(instrument dto.DatabaseInstrument) Instrument {
	return Instrument{
		StorageID: instrument.ID,
		Name:      instrument.Name,
	}
}

func FromInstrumentsDatabaseDTO(instruments []dto.DatabaseInstrument) []Instrument {
	var result []Instrument
	for _, instrument := range instruments {
		result = append(result, FromInstrumentDatabaseDTO(instrument))
	}
	return result
}
