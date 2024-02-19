package form

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"ffxvi-bard/port/helper"
	"sync"
)

const (
	SortTitleAsc   = "title_asc"
	SortTitleDesc  = "title_desc"
	SortAddedAsc   = "added_asc"
	SortAddedDesc  = "added_desc"
	SortArtistAsc  = "artist_asc"
	SortArtistDesc = "artist_desc"
	SortRatingAsc  = "rating_low"
	SortRatingDesc = "rating_high"
)

func getSortCriteriaTitle(option string) string {
	var titles = map[string]string{
		SortTitleAsc:   "Title (A - Z)",
		SortTitleDesc:  "Title (Z - A)",
		SortAddedAsc:   "Date Added (Oldest First)",
		SortAddedDesc:  "Date Added (Newest First)",
		SortArtistAsc:  "Artist (A - Z)",
		SortArtistDesc: "Artist (Z - A)",
		SortRatingAsc:  "Rating (Lowest First)",
		SortRatingDesc: "Rating (Highest First)",
	}
	return titles[option]
}

func getSortOptions() map[string]string {
	sortOptions := make(map[string]string)
	for _, option := range []string{SortTitleAsc, SortTitleDesc, SortAddedAsc, SortAddedDesc, SortArtistAsc, SortArtistDesc, SortRatingAsc, SortRatingDesc} {
		sortOptions[option] = getSortCriteriaTitle(option)
	}
	return sortOptions
}

type SongListPagination struct {
	TotalPages  int
	CurrentPage int
	NextPage    int
	PrevPage    int
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
	if p.TotalPages > 1 && !helper.Contains(pages, p.TotalPages) {
		pages = append(pages, p.TotalPages)
	}
	return pages
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

type SongListFilter struct {
	Title              string
	Artist             string
	EnsembleSize       int
	GenreID            int
	Offset             int
	Limit              int
	Sort               string
	EnsembleSizeString string
}

func NewSongListFilter(title string, artist string, ensembleSize int, genreID int, offset int, limit int, sort string) *SongListFilter {
	filter := SongListFilter{
		Title:        title,
		Artist:       artist,
		EnsembleSize: ensembleSize,
		GenreID:      genreID,
		Offset:       offset,
		Limit:        limit,
		Sort:         sort,
	}
	if ensembleSize != -1 {
		filter.EnsembleSizeString = song.EnsembleString(ensembleSize)
	}
	return &filter
}

type SongList struct {
	Songs            *[]dto.SongWithDetails
	Genres           *[]dto.DatabaseGenre
	SortOptions      map[string]string
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

func (s *SongList) Fetch(songTitle string, artist string, ensembleSize int, genreID int, page int, sort string) (*SongList, error) {
	var songs []dto.SongWithDetails
	var totalSongs int
	var err error
	var genres []dto.DatabaseGenre

	limit := 15
	offset := (page - 1) * limit
	s.Filters = NewSongListFilter(songTitle, artist, ensembleSize, genreID, page, limit, sort)
	s.SortOptions = getSortOptions()
	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()
		songs, err = s.songRepository.FetchForPagination(songTitle, artist, ensembleSize, genreID, sort, limit, offset)
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
