package database

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
	"strings"
)

type InstrumentRepository struct {
	driver contract.DatabaseDriverInterface
}

func NewInstrumentRepository(driver contract.DatabaseDriverInterface) contract.InstrumentRepositoryInterface {
	return &InstrumentRepository{
		driver: driver,
	}
}

func (i *InstrumentRepository) FetchAll() ([]dto.DatabaseInstrument, error) {
	var instruments []dto.DatabaseInstrument
	rows, err := i.driver.FetchMany("SELECT id, name FROM instrument")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var instrument dto.DatabaseInstrument
		if err := rows.Scan(&instrument.ID, &instrument.Name); err != nil {
			return nil, err
		}
		instruments = append(instruments, instrument)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return instruments, nil
}

func (i *InstrumentRepository) FetchByIDs(instrumentIDs []int) ([]dto.DatabaseInstrument, error) {
	if len(instrumentIDs) == 0 {
		return []dto.DatabaseInstrument{}, nil // Return an empty slice if no IDs are provided
	}
	placeholder := make([]string, len(instrumentIDs))
	for i := range placeholder {
		placeholder[i] = "?"
	}
	query := fmt.Sprintf("SELECT id, name FROM instrument WHERE id IN (%s)", strings.Join(placeholder, ","))

	args := make([]interface{}, len(instrumentIDs))
	for i, id := range instrumentIDs {
		args[i] = id
	}
	rows, err := i.driver.FetchMany(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var instruments []dto.DatabaseInstrument
	for rows.Next() {
		var instrument dto.DatabaseInstrument
		if err := rows.Scan(&instrument.ID, &instrument.Name); err != nil {
			return nil, err
		}
		instruments = append(instruments, instrument)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return instruments, nil
}

func (i *InstrumentRepository) FetchBySongID(songID int) (*[]dto.DatabaseInstrument, error) {
	var instruments []dto.DatabaseInstrument
	query := `
		SELECT i.id, i.name
		FROM instrument i
		INNER JOIN song_instrument si on si.instrument_id = i.id 
		WHERE si.song_id = ?`
	rows, err := i.driver.FetchMany(query, songID)
	for rows.Next() {
		var instrument dto.DatabaseInstrument
		if err := rows.Scan(&instrument.ID, &instrument.Name); err != nil {
			return nil, err
		}
		instruments = append(instruments, instrument)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &instruments, nil
}

func diffInstruments(current, new []int) (toAdd, toRemove []int) {
	currentMap := make(map[int]bool)
	newMap := make(map[int]bool)
	for _, id := range current {
		currentMap[id] = true
	}
	for _, id := range new {
		if !currentMap[id] {
			toAdd = append(toAdd, id)
		}
		newMap[id] = true
	}
	for _, id := range current {
		if !newMap[id] {
			toRemove = append(toRemove, id)
		}
	}
	return toAdd, toRemove
}
