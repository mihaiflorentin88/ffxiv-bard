package database

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
	"strings"
)

type GenreRepository struct {
	driver contract.DatabaseDriverInterface
}

func diffGenres(current, new []int) (toAdd, toRemove []int) {
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

func NewGenreRepository(driver contract.DatabaseDriverInterface) contract.GenreRepositoryInterface {
	return &GenreRepository{
		driver: driver,
	}
}

func (g *GenreRepository) FetchAll() ([]dto.DatabaseGenre, error) {
	var genres []dto.DatabaseGenre
	rows, err := g.driver.FetchMany("SELECT id, name FROM genre")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var genre dto.DatabaseGenre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return genres, nil
}

func (g *GenreRepository) FetchByIDs(genreIDs []int) ([]dto.DatabaseGenre, error) {
	if len(genreIDs) == 0 {
		return []dto.DatabaseGenre{}, nil // Return an empty slice if no IDs are provided
	}
	placeholder := make([]string, len(genreIDs))
	for i := range placeholder {
		placeholder[i] = "?"
	}
	query := fmt.Sprintf("SELECT id, name FROM genre WHERE id IN (%s)", strings.Join(placeholder, ","))

	args := make([]interface{}, len(genreIDs))
	for i, id := range genreIDs {
		args[i] = id
	}
	rows, err := g.driver.FetchMany(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []dto.DatabaseGenre
	for rows.Next() {
		var genre dto.DatabaseGenre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return genres, nil
}

func (g *GenreRepository) FetchBySongID(songID int) (*[]dto.DatabaseGenre, error) {
	var genres []dto.DatabaseGenre
	query := `
		SELECT g.id, g.name
		FROM genre g
		INNER JOIN song_genre sg on sg.genre_id = g.id 
		WHERE sg.song_id = ?`
	rows, err := g.driver.FetchMany(query, songID)
	for rows.Next() {
		var genre dto.DatabaseGenre
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &genres, nil
}
