package song

import (
	"ffxvi-bard/domain/date"
	"ffxvi-bard/port/dto"
)

type Genre struct {
	StorageID int
	Name      string
	Date      date.Date
}

func FromGenreDatabaseDTO(genre dto.DatabaseGenre) Genre {
	return Genre{
		StorageID: genre.ID,
		Name:      genre.Name,
	}

}

func FromGenresDatabaseDTO(genres []dto.DatabaseGenre) []Genre {
	var result []Genre
	for _, genre := range genres {
		result = append(result, FromGenreDatabaseDTO(genre))
	}
	return result
}
