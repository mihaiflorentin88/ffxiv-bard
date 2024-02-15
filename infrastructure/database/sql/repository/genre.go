package database

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
)

type GenreRepository struct {
	driver contract.DatabaseDriverInterface
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
