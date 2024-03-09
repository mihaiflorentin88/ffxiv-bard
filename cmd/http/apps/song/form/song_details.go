package form

import (
	"ffxvi-bard/domain/song"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"github.com/gin-gonic/gin"
)

const dateLayout = "Jan 2, 2006 at 3:04pm"

type Genre struct {
	ID   int
	Name string
}

type Instrument struct {
	ID   int
	Name string
}

type Comment struct {
	StorageID int
	Author    string
	Title     string
	Content   string
	Status    bool
	CreatedAt string
	UpdatedAt string
	CanEdit   bool
}

type SongDetails struct {
	ID                   int
	Title                string
	Artist               string
	Filename             string
	EnsembleSize         int
	EnsembleSizeString   string
	Genre                []Genre
	Source               string
	Note                 string
	AudioCrafter         string
	Instrument           []Instrument
	CanEdit              bool
	Comments             []Comment
	Uploader             string
	Rating               float64
	LoggedUserRating     int
	genreRepository      contract.GenreRepositoryInterface
	commentRepository    contract.CommentRepositoryInterface
	instrumentRepository contract.InstrumentRepositoryInterface
	ratingRepository     contract.RatingRepositoryInterface
	song                 song.Song
	loggedUser           user.User
}

func (s *SongDetails) GetAvailableStars() []int {
	stars := make([]int, 5)
	for i := range stars {
		stars[i] = 5 - i
	}
	return stars
}

func NewSongDetailsForm(genreRepository contract.GenreRepositoryInterface, commentRepository contract.CommentRepositoryInterface, ratingRepository contract.RatingRepositoryInterface, instrumentRepository contract.InstrumentRepositoryInterface, song *song.Song) SongDetails {
	return SongDetails{
		genreRepository:      genreRepository,
		song:                 *song,
		commentRepository:    commentRepository,
		ratingRepository:     ratingRepository,
		instrumentRepository: instrumentRepository,
	}
}

func (s *SongDetails) Fetch(songId int, c *gin.Context) (*SongDetails, error) {
	_, err := s.song.LoadByID(songId)
	if err != nil {
		return s, err
	}
	sessionUser, _ := c.Get("user")
	s.loggedUser, _ = user.FromSession(sessionUser)
	s.hydrate()
	return s, nil
}

func (s *SongDetails) hydrate() {
	var comments []Comment
	s.ID = s.song.StorageID
	s.Title = s.song.Title
	s.Artist = s.song.Artist
	s.Filename = s.song.Filename
	s.EnsembleSize = int(s.song.EnsembleSize)
	s.EnsembleSizeString = song.EnsembleString(s.EnsembleSize)
	s.Source = s.song.Source
	s.Note = s.song.Note
	s.AudioCrafter = s.song.AudioCrafter
	for _, genre := range s.song.Genre {
		renderedGenre := Genre{
			ID:   genre.StorageID,
			Name: genre.Name}
		s.Genre = append(s.Genre, renderedGenre)
	}

	for _, instrument := range s.song.Instrument {
		renderedInstrument := Instrument{
			ID:   instrument.StorageID,
			Name: instrument.Name,
		}
		s.Instrument = append(s.Instrument, renderedInstrument)
	}
	if (s.loggedUser != user.User{}) {
		if s.loggedUser.StorageID == s.song.Uploader.StorageID || s.loggedUser.IsAdmin {
			s.CanEdit = true
		}
	}
	for _, comment := range s.song.Comments {
		renderedComment := Comment{
			StorageID: comment.StorageID,
			Author:    comment.Author.Name,
			Status:    comment.Status,
			Content:   comment.Content,
		}
		if (s.loggedUser != user.User{}) {
			if s.loggedUser.StorageID == comment.Author.StorageID || s.loggedUser.IsAdmin {
				renderedComment.CanEdit = true
			}
		}
		renderedComment.CreatedAt = comment.Date.CreatedAt.Format(dateLayout)
		renderedComment.UpdatedAt = comment.Date.UpdatedAt.Format(dateLayout)
		comments = append(comments, renderedComment)
	}
	if (s.loggedUser != user.User{}) {
		s.LoggedUserRating, _ = s.ratingRepository.FindByUserAndSong(s.song.StorageID, s.loggedUser.StorageID)
	}
	s.Comments = comments
	s.Rating = s.song.GetAverageRating()
	s.Uploader = s.song.Uploader.Name
}
