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
	"time"
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

var EnsembleSizes []EnsembleSize = []EnsembleSize{
	contract.Solo,    // 0
	contract.Duet,    // 1
	contract.Trio,    // 2
	contract.Quartet, // 3
	contract.Quintet, // 4
	contract.Sextet,  // 5
	contract.Septet,  // 6
	contract.Octet,   // 7
}

func GetEnsembleSize(size int) (contract.EnsembleSize, error) {
	if size >= 0 && size < len(EnsembleSizes) {
		return EnsembleSizes[size], nil
	}
	return 0, errors.New("%v is not a valid EnsembleSize")
}

func EnsembleString(i int) string {
	ensembleStrings := [...]string{"Solo", "Duet", "Trio", "Quartet", "Quintet", "Sextet", "Septet", "Octet"}
	return ensembleStrings[i]
}

type song struct {
	storageID      int
	title          string
	artist         string
	ensembleSize   EnsembleSize
	fileCode       string
	file           []byte
	checksum       string
	rating         []Rating
	genre          []Genre
	uploader       *user.User
	comments       []contract.CommentInterface
	status         Status
	statusMessage  string
	lockTs         time.Time
	songProcessor  contract.SongProcessorInterface
	filesystem     contract.FileSystemInterface
	Date           date.Date
	songRepository contract.SongRepositoryInterface
}

func NewEmptySong(songProcessor contract.SongProcessorInterface, filesystem contract.FileSystemInterface) contract.SongInterface {
	return &song{
		songProcessor: songProcessor,
		filesystem:    filesystem,
	}
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
	song.SetUploader(&uploader)
	song.SetStatus(Pending)
	song.SetSongProcessor(songProcessor)
	song.SetFileSystem(filesystem)
	song.GenerateFileCode()
	err = song.songProcessor.WriteUnprocessedSong(song)
	return song, err
}

func FromNewSongDTO(newSongDto dto.NewSong, songRepository contract.SongRepositoryInterface, genreRepository contract.GenreRepositoryInterface, songProcessor contract.SongProcessorInterface) (contract.SongInterface, error) {
	song := song{
		title:         newSongDto.Title,
		artist:        newSongDto.Artist,
		status:        Pending,
		statusMessage: "Pending song processing.",
	}
	song.songRepository = songRepository
	song.file = newSongDto.File
	song.ComputeChecksum()
	duplicatedEntry, err := songRepository.FindByChecksum(song.checksum)
	if duplicatedEntry.ID != 0 {
		return &song, errors.New(fmt.Sprintf("song already exists under id `%v`", duplicatedEntry.ID))
	}
	ensSize, err := GetEnsembleSize(newSongDto.EnsembleSize)
	if err != nil {
		return &song, err
	}
	song.ensembleSize = ensSize
	genresDTO, err := genreRepository.FetchByIDs(newSongDto.Genre)
	if err != nil {
		return &song, errors.New(fmt.Sprintf("one of the genres might not be valid. Error %s", err))
	}
	if userObj, ok := newSongDto.User.(*user.User); ok {
		song.uploader = userObj
	} else {
		return &song, errors.New("song uploader is not of the correct type")
	}
	song.genre = FromGenresDatabaseDTO(genresDTO)
	song.songProcessor = songProcessor
	song.GenerateFileCode()
	err = song.songProcessor.WriteUnprocessedSong(&song)
	return &song, nil
}

func (s *song) ToDatabaseSongDTO() dto.DatabaseSongDTO {
	return dto.DatabaseSongDTO{
		ID:            s.storageID,
		Title:         s.title,
		Artist:        s.artist,
		EnsembleSize:  int(s.ensembleSize),
		FileCode:      s.fileCode,
		UploaderID:    s.uploader.StorageID,
		Status:        int(s.status),
		StatusMessage: &s.statusMessage,
		Checksum:      s.checksum,
		LockExpireTS:  &s.lockTs,
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
}

func (s *song) EnsembleString() string {
	ensembleStrings := [...]string{"Solo", "Duet", "Trio", "Quartet", "Quintet", "Sextet", "Septet", "Octet"}
	return ensembleStrings[s.ensembleSize]
}

func (s *song) ComputeChecksum() {
	if s.file == nil {
		errors.New("no files to compute checksum")
	}
	hash := sha256.New()
	hash.Write(s.file)
	hashBytes := hash.Sum(nil)
	s.checksum = base64.StdEncoding.EncodeToString(hashBytes)
}

func (a *song) GetDetailedEnsembleString() map[int]string {
	var result = make(map[int]string)
	for i := range 8 { // don't mind the ide. this is syntax added in golang 1.22. the ide just didn had time to catch up with it.
		result[i] = EnsembleString(i)
	}
	return result
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

func (s *song) GetUploader() *user.User {
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

func (s *song) SetUploader(uploader *user.User) {
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

func FromSubmittedForm() {

}
