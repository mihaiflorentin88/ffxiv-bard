package song

import (
	"ffxvi-bard/domain/date"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
)

type Genre struct {
	StorageID       int
	Name            string
	Date            date.Date
	genreRepository contract.GenreRepositoryInterface
}

func NewEmptyGenre(genreRepository contract.GenreRepositoryInterface) Genre {
	return Genre{
		genreRepository: genreRepository,
	}
}

func (g *Genre) FetchBySongID(songID int) ([]Genre, error) {
	var genres []Genre
	genreDTOs, err := g.genreRepository.FetchBySongID(songID)
	if err != nil {
		return genres, err
	}
	for _, genreDTO := range *genreDTOs {
		genre := FromGenreDatabaseDTO(genreDTO)
		genre.genreRepository = g.genreRepository
		genres = append(genres, genre)
	}
	return genres, nil
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
