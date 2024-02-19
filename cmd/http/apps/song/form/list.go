package form

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"sync"
)

type SongListPagination struct {
	TotalPages  int
	CurrentPage int
	NextPage    int
	PrevPage    int
}

type SongListFilter struct {
	Title              string
	Artist             string
	EnsembleSize       int
	GenreID            int
	Offset             int
	Limit              int
	EnsembleSizeString string
}

func (p *SongListPagination) PagesSequence() []int {
	pages := []int{}
	pages = append(pages, 1)
	if p.CurrentPage > 2 {
		pages = append(pages, p.CurrentPage-1)
	}
	if p.CurrentPage > 1 {
		pages = append(pages, p.CurrentPage)
	}

	if p.CurrentPage < p.TotalPages-1 {
		pages = append(pages, p.CurrentPage+1)
	}
	if p.TotalPages > 1 && !contains(pages, p.TotalPages) {
		pages = append(pages, p.TotalPages)
	}
	return pages
}

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func NewSongListPagination(totalCount, currentPage, limit int) *SongListPagination {
	totalPages := totalCount / limit
	if totalCount%limit != 0 {
		totalPages++
	}

	nextPage := currentPage + 1
	if nextPage > totalPages {
		nextPage = 0
	}

	prevPage := currentPage - 1
	if prevPage < 1 {
		prevPage = 0
	}

	return &SongListPagination{
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		NextPage:    nextPage,
		PrevPage:    prevPage,
	}
}

func NewSongListFilter(title string, artist string, ensembleSize int, genreID int, offset int, limit int) *SongListFilter {
	filter := SongListFilter{
		Title:        title,
		Artist:       artist,
		EnsembleSize: ensembleSize,
		GenreID:      genreID,
		Offset:       offset,
		Limit:        limit,
	}
	if ensembleSize != -1 {
		filter.EnsembleSizeString = song.EnsembleString(ensembleSize)
	}
	return &filter
}

type SongList struct {
	Songs            *[]dto.SongWithDetails
	Genres           *[]dto.DatabaseGenre
	EnsembleSize     map[int]string
	songRepository   contract.SongRepositoryInterface
	genreRepository  contract.GenreRepositoryInterface
	ratingRepository contract.RatingRepositoryInterface
	spotify          contract.MediaClientInterface
	Filters          *SongListFilter
	Pagination       *SongListPagination
}

func NewSongList(songRepository contract.SongRepositoryInterface, genreRepository contract.GenreRepositoryInterface, ratingRepository contract.RatingRepositoryInterface, spotify contract.MediaClientInterface) SongList {
	return SongList{
		songRepository:   songRepository,
		genreRepository:  genreRepository,
		ratingRepository: ratingRepository,
		EnsembleSize:     song.GetDetailedEnsembleString(),
		spotify:          spotify,
	}
}

func (s *SongList) addAlbumImageToSong(song *dto.SongWithDetails) {
	spotifySong, err := s.spotify.Search(song.Title, song.Artist)
	if err != nil {
		return
	}
	image, err := spotifySong.GetSmallestImage()
	if err != nil {
		return
	}
	song.ImageUrl = image.URL
}

func (s *SongList) Fetch(songTitle string, artist string, ensembleSize int, genreID int, page int) (*SongList, error) {
	var songs []dto.SongWithDetails
	var totalSongs int
	var err error
	var genres []dto.DatabaseGenre

	limit := 12
	offset := (page - 1) * limit
	s.Filters = NewSongListFilter(songTitle, artist, ensembleSize, genreID, page, limit)
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		songs, err = s.songRepository.FetchForPagination(songTitle, artist, ensembleSize, genreID, limit, offset)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		totalSongs, err = s.songRepository.FetchTotalSongsForListing(songTitle, artist, ensembleSize, genreID)
		if err != nil {
			errChan <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		genres, err = s.genreRepository.FetchAll()
		if err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return s, err
		}
	}
	s.Pagination = NewSongListPagination(totalSongs, page, limit)
	for i := range songs {
		songs[i].EnsembleSizeString = song.EnsembleString(songs[i].EnsembleSize)
	}
	s.Songs = &songs
	s.Genres = &genres
	return s, nil
}
