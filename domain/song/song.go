package song

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"ffxvi-bard/domain/date"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
	"github.com/google/uuid"
	"math"
	"strings"
	"time"
)

type Status int
type EnsembleSize int

const ( // If you change anything here. You will also need to change the StatusString() function
	Pending Status = iota
	Processing
	Processed
	Failed
	Deleted
)

const ( // If you change anything here. You will also need to change the EnsembleString() function
	Solo EnsembleSize = iota
	Duet
	Trio
	Quartet
	Quintet
	Sextet
	Septet
	Octet
)

var EnsembleSizes []EnsembleSize = []EnsembleSize{
	Solo,    // 0
	Duet,    // 1
	Trio,    // 2
	Quartet, // 3
	Quintet, // 4
	Sextet,  // 5
	Septet,  // 6
	Octet,   // 7
}

func GetEnsembleSizeFromInt(size int) (EnsembleSize, error) {
	if size >= 0 && size < len(EnsembleSizes) {
		return EnsembleSizes[size], nil
	}
	return 0, errors.New("%v is not a valid EnsembleSize")
}

func GetDetailedEnsembleString() map[int]string {
	var result = make(map[int]string)
	for i := range 8 { // don't mind the ide. this is syntax added in golang 1.22. the ide just didn had time to catch up with it.
		result[i] = EnsembleString(i)
	}
	return result
}

func EnsembleString(i int) string {
	ensembleStrings := [...]string{"Solo", "Duet", "Trio", "Quartet", "Quintet", "Sextet", "Septet", "Octet"}
	return ensembleStrings[i]
}

type Song struct {
	StorageID      int
	Title          string
	Artist         string
	EnsembleSize   EnsembleSize
	FileCode       string
	File           []byte
	Checksum       string
	Rating         []Rating
	Genre          []Genre
	Uploader       *user.User
	Comments       []contract.CommentInterface
	status         Status
	statusMessage  string
	LockTs         time.Time
	SongProcessor  contract.SongProcessorInterface
	Filesystem     contract.FileSystemInterface
	Date           date.Date
	SongRepository contract.SongRepositoryInterface
}

func NewEmptySong(songProcessor contract.SongProcessorInterface, filesystem contract.FileSystemInterface) *Song {
	return &Song{
		SongProcessor: songProcessor,
		Filesystem:    filesystem,
	}
}

func FromNewSongForm(newSongDto dto.NewSongForm, songRepository contract.SongRepositoryInterface, genreRepository contract.GenreRepositoryInterface, songProcessor contract.SongProcessorInterface) (*Song, error) {
	song := Song{
		Title:         newSongDto.Title,
		Artist:        newSongDto.Artist,
		status:        Pending,
		statusMessage: "Pending Song processing.",
	}
	song.SongRepository = songRepository
	song.File = newSongDto.File
	err := song.ComputeChecksum()
	if err != nil {
		return &song, err
	}
	duplicatedEntry, err := songRepository.FindByChecksum(song.Checksum)
	if duplicatedEntry.ID != 0 {
		return &song, errors.New(fmt.Sprintf("Song already exists under id `%v`", duplicatedEntry.ID))
	}
	ensSize, err := GetEnsembleSizeFromInt(newSongDto.EnsembleSize)
	if err != nil {
		return &song, err
	}
	song.EnsembleSize = ensSize
	genresDTO, err := genreRepository.FetchByIDs(newSongDto.Genre)
	if err != nil {
		return &song, errors.New(fmt.Sprintf("one of the genres might not be valid. Error %s", err))
	}
	song.Uploader, err = user.FromSession(newSongDto.User)
	if err != nil {
		return &song, errors.New("Song uploader is not of the correct type")
	}
	song.Genre = FromGenresDatabaseDTO(genresDTO)
	song.SongProcessor = songProcessor
	song.GenerateFileCode()
	err = song.SongProcessor.WriteUnprocessedSong(song.FileCode, song.File)
	return &song, nil
}

func (s *Song) ToDatabaseSongDTO() dto.DatabaseSong {
	return dto.DatabaseSong{
		ID:            s.StorageID,
		Title:         strings.ToLower(s.Title),
		Artist:        strings.ToLower(s.Artist),
		EnsembleSize:  int(s.EnsembleSize),
		FileCode:      s.FileCode,
		UploaderID:    s.Uploader.StorageID,
		Status:        int(s.status),
		StatusMessage: &s.statusMessage,
		Checksum:      s.Checksum,
		LockExpireTS:  &s.LockTs,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
}

func (s *Song) EnsembleString() string {
	ensembleStrings := [...]string{"Solo", "Duet", "Trio", "Quartet", "Quintet", "Sextet", "Septet", "Octet"}
	return ensembleStrings[s.EnsembleSize]
}

func (s *Song) ComputeChecksum() error {
	if s.File == nil {
		return errors.New("no files to compute checksum")
	}
	hash := sha256.New()
	hash.Write(s.File)
	hashBytes := hash.Sum(nil)
	s.Checksum = base64.StdEncoding.EncodeToString(hashBytes)
	return nil
}

func (s *Song) StatusString() string {
	statusStrings := [...]string{"Pending", "Processing", "Processed", "Failed", "Deleted"}
	return statusStrings[s.status]
}

func (s *Song) GenerateFileCode() {
	newUUID := uuid.New()
	s.FileCode = newUUID.String()
}

func (s *Song) AddComment(c contract.CommentInterface) {
	s.Comments = append(s.Comments, c)
}

func (s *Song) RemoveComment(c contract.CommentInterface) {
	for i, comment := range s.Comments {
		if comment.GetStorageID() == c.GetStorageID() {
			s.Comments = append(s.Comments[:i], s.Comments[i+1:]...)
			break
		}
	}
}

func (s *Song) GetAverageRating() float64 {
	total := 0
	for _, rating := range s.Rating {
		total += rating.rating
	}
	if total == 0 {
		return 0
	}
	average := float64(total) / float64(len(s.Rating))
	return math.Round(average*100) / 100
}

func (s *Song) ChangeStatus(status Status, statusMessage string) {
	s.status = status
	s.ChangeStatusMessage(statusMessage)
}

func (s *Song) GetStatus() Status {
	return s.status
}

func (s *Song) ChangeStatusMessage(statusMessage string) {
	s.statusMessage = statusMessage
}

func (s *Song) GetStatusMessage() string {
	return s.statusMessage
}

func (s *Song) ProcessSong() error {
	s.GenerateFileCode()
	err := s.SongProcessor.ProcessSong(s.FileCode)
	if err != nil {
		msg := fmt.Sprintf("Failed to process Song. Reason: %s", err.Error())
		s.ChangeStatus(Failed, msg)
		return errors.New(msg)
	}
	return nil
}

func (s *Song) GetStorageID() int {
	return s.StorageID
}

func (s *Song) SetTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	s.Title = title
	return nil
}

func (s *Song) SetArtist(artist string) error {
	if artist == "" {
		return errors.New("artist cannot be empty")
	}
	s.Artist = artist
	return nil
}

func (s *Song) RemoveUnprocessedSong() error {
	s.GenerateFileCode()
	return s.SongProcessor.RemoveUnprocessedSong(s.FileCode)
}

func (s *Song) AddRating(rating Rating) {
	s.Rating = append(s.Rating, rating)
}
