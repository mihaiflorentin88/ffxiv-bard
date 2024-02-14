package song

import (
	"errors"
	"ffxvi-bard/domain/date"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"fmt"
	"github.com/google/uuid"
	"math"
)

type Status = contract.Status
type EnsembleSize = contract.EnsembleSize

const ( // Make sure the order of these constants match the order of the Status enum in port/song.go
	Pending    = contract.Pending
	Processing = contract.Processing
	Processed  = contract.Processed
	Failed     = contract.Failed
	Deleted    = contract.Deleted
)

const ( // Make sure the order of these constants match the order of the EnsembleSize enum in port/song.go
	Solo    = contract.Solo
	Duet    = contract.Duet
	Trio    = contract.Trio
	Quartet = contract.Quartet
	Quintet = contract.Quintet
	Sextet  = contract.Sextet
	Septet  = contract.Septet
	Octet   = contract.Octet
)

type song struct {
	storageID     int
	title         string
	artist        string
	ensembleSize  EnsembleSize
	fileCode      string
	file          []byte
	rating        []Rating
	genre         []Genre
	uploader      user.User
	comments      []contract.CommentInterface
	status        Status
	statusMessage string
	songProcessor contract.SongProcessorInterface
	filesystem    contract.FileSystemInterface
	Date          date.Date
}

func NewSong(title string, artist string, ensembleSize EnsembleSize, genre []Genre, comments []contract.CommentInterface, file []byte, uploader user.User, songProcessor contract.SongProcessorInterface, filesystem contract.FileSystemInterface) (contract.SongInterface, error) {
	song := &song{}
	err := song.SetTitle(title)
	if err != nil {
		return nil, err
	}
	err = song.SetArtist(artist)
	if err != nil {
		return nil, err
	}
	song.SetEnsembleSize(ensembleSize)
	song.SetGenre(genre)
	song.SetComments(comments)
	song.SetFile(file)
	song.SetUploader(uploader)
	song.SetStatus(Pending)
	song.SetSongProcessor(songProcessor)
	song.SetFileSystem(filesystem)
	song.GenerateFileCode()
	err = song.songProcessor.WriteUnprocessedSong(song)
	return song, err
}

func (s *song) EnsembleString() string {
	ensembleStrings := [...]string{"Solo", "Duet", "Trio", "Quartet", "Quintet", "Sextet", "Septet", "Octet"}
	return ensembleStrings[s.ensembleSize]
}

func (s *song) StatusString() string {
	statusStrings := [...]string{"Pending", "Processing", "Processed", "Failed", "Deleted"}
	return statusStrings[s.status]
}

func (s *song) GenerateFileCode() {
	newUUID := uuid.New()
	s.fileCode = newUUID.String()
}

func (s *song) AddComment(c contract.CommentInterface) {
	s.comments = append(s.comments, c)
}

func (s *song) RemoveComment(c contract.CommentInterface) {
	for i, comment := range s.comments {
		if comment.GetStorageID() == c.GetStorageID() {
			s.comments = append(s.comments[:i], s.comments[i+1:]...)
			break
		}
	}
}

func (s *song) GetAverageRating() float64 {
	total := 0
	for _, rating := range s.rating {
		total += rating.GetRanking()
	}
	if total == 0 {
		return 0
	}
	average := float64(total) / float64(len(s.rating))
	return math.Round(average*100) / 100
}

func (s *song) ChangeStatus(status Status, statusMessage string) {
	s.status = status
	s.ChangeStatusMessage(statusMessage)
}

func (s *song) GetStatus() Status {
	return s.status
}

func (s *song) ChangeStatusMessage(statusMessage string) {
	s.statusMessage = statusMessage
}

func (s *song) GetStatusMessage() string {
	return s.statusMessage
}

func (s *song) ProcessSong() error {
	err := s.songProcessor.ProcessSong(s)
	if err != nil {
		msg := fmt.Sprintf("Failed to process song. Reason: %s", err.Error())
		s.ChangeStatus(Failed, msg)
		return errors.New(msg)
	}
	return nil
}

func (s *song) GetFileCode() string {
	return s.fileCode
}

func (s *song) GetFile() []byte {
	return s.file
}

func (s *song) GetTitle() string {
	return s.title
}

func (s *song) GetArtist() string {
	return s.artist
}

func (s *song) GetEnsembleSize() EnsembleSize {
	return s.ensembleSize
}

func (s *song) GetGenre() []Genre {
	return s.genre
}

func (s *song) GetUploader() user.User {
	return s.uploader
}

func (s *song) GetComments() []contract.CommentInterface {
	return s.comments
}

func (s *song) GetStorageID() int {
	return s.storageID
}

func (s *song) SetTitle(title string) error {
	if title == "" {
		return errors.New("title cannot be empty")
	}
	s.title = title
	return nil
}

func (s *song) SetArtist(artist string) error {
	if artist == "" {
		return errors.New("artist cannot be empty")
	}
	s.artist = artist
	return nil
}

func (s *song) SetEnsembleSize(ensembleSize EnsembleSize) {
	s.ensembleSize = ensembleSize
}

func (s *song) SetGenre(genre []Genre) {
	s.genre = genre
}

func (s *song) SetComments(comments []contract.CommentInterface) {
	s.comments = comments
}

func (s *song) SetFile(file []byte) {
	s.file = file
}

func (s *song) SetUploader(uploader user.User) {
	s.uploader = uploader
}

func (s *song) SetStatus(status Status) {
	s.status = status
}

func (s *song) SetSongProcessor(songProcessor contract.SongProcessorInterface) {
	s.songProcessor = songProcessor
}

func (s *song) SetFileSystem(filesystem contract.FileSystemInterface) {
	s.filesystem = filesystem
}

func (s *song) RemoveUnprocessedSong() error {
	return s.songProcessor.RemoveUnprocessedSong(s)
}

func (s *song) AddRating(rating Rating) {
	s.rating = append(s.rating, rating)
}
