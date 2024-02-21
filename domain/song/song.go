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
	"math"
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
	Filename       string
	File           []byte
	Checksum       string
	Rating         []Rating
	Genre          []Genre
	Uploader       *user.User
	Comments       []Comment
	status         Status
	statusMessage  string
	LockExpireTs   time.Time
	SongProcessor  contract.SongProcessorInterface
	Filesystem     contract.FileSystemInterface
	Date           date.Date
	songRepository contract.SongRepositoryInterface
	emptyRating    *Rating
	emptyComment   *Comment
	emptyGenre     *Genre
	emptyUser      *user.User
}

func NewEmptySong(songProcessor contract.SongProcessorInterface, filesystem contract.FileSystemInterface, emptyUser *user.User, emptyRating *Rating, emptyComment *Comment, emptyGenre *Genre, songRepository contract.SongRepositoryInterface) *Song {
	return &Song{
		SongProcessor:  songProcessor,
		Filesystem:     filesystem,
		Uploader:       emptyUser,
		emptyRating:    emptyRating,
		emptyComment:   emptyComment,
		emptyGenre:     emptyGenre,
		emptyUser:      emptyUser,
		songRepository: songRepository,
	}
}

func FromNewSongForm(newSongDto dto.NewSongForm, songRepository contract.SongRepositoryInterface, genreRepository contract.GenreRepositoryInterface, songProcessor contract.SongProcessorInterface) (*Song, error) {
	song := Song{
		Title:         newSongDto.Title,
		Artist:        newSongDto.Artist,
		status:        Pending,
		statusMessage: "Pending Song processing.",
	}
	song.songRepository = songRepository
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
	song.Filename = song.SongProcessor.GenerateUniqueFilename()
	err = song.SongProcessor.WriteUnprocessedSong(song.Filename, song.File)
	return &song, nil
}

func (s *Song) ToDatabaseSongDTO() dto.DatabaseSong {
	return dto.DatabaseSong{
		ID:            s.StorageID,
		Title:         s.Title,
		Artist:        s.Artist,
		EnsembleSize:  int(s.EnsembleSize),
		Filename:      s.Filename,
		UploaderID:    s.Uploader.StorageID,
		Status:        int(s.status),
		StatusMessage: s.statusMessage,
		Checksum:      s.Checksum,
		LockExpireTS:  s.LockExpireTs,
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

func (s *Song) AddComment(c Comment) {
	s.Comments = append(s.Comments, c)
}

func (s *Song) RemoveComment(c Comment) {
	for i, comment := range s.Comments {
		if comment.StorageID == c.StorageID {
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
	if s.Filename == "" {
		s.Filename = s.SongProcessor.GenerateUniqueFilename()
	}
	err := s.SongProcessor.ProcessSong(s.Filename)
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
	return s.SongProcessor.RemoveUnprocessedSong(s.Filename)
}

func (s *Song) AddRating(rating Rating) {
	s.Rating = append(s.Rating, rating)
}

func (s *Song) LoadByID(songID int) (*Song, error) {
	songDTO, err := s.songRepository.FindByID(songID)
	if err != nil {
		return s, errors.New(fmt.Sprintf("failed to find the requested song. Reason %s", err))
	}
	_, err = FromDatabaseDTO(s, &songDTO)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s *Song) IsCorrectlyInstantiated() error {
	if s.emptyUser == nil {
		return errors.New("song.Song was not correctly instantiated")
	}
	if s.emptyGenre == nil {
		return errors.New("song.Song was not correctly instantiated")
	}
	if s.emptyRating == nil {
		return errors.New("song.Song was not correctly instantiated")
	}

	if s.emptyComment == nil {
		return errors.New("song.Song was not correctly instantiated")
	}

	if s.songRepository == nil {
		return errors.New("song.Song was not correctly instantiated")
	}
	return nil
}

func FromDatabaseDTO(song *Song, databaseSong *dto.DatabaseSong) (*Song, error) {
	var err error
	if err = song.IsCorrectlyInstantiated(); err != nil {
		return song, err
	}
	song.StorageID = databaseSong.ID
	song.Title = databaseSong.Title
	song.Artist = databaseSong.Artist
	song.EnsembleSize, err = GetEnsembleSizeFromInt(databaseSong.EnsembleSize)
	if err != nil {
		return song, err
	}
	song.Filename = databaseSong.Filename
	song.Checksum = databaseSong.Checksum
	ratings, err := song.emptyRating.FetchBySongId(song.StorageID)
	if err != nil {
		return song, err
	}
	song.Rating = ratings
	song.Genre, err = song.emptyGenre.FetchBySongID(song.StorageID)
	if err != nil {
		return song, err
	}
	if song.Uploader == nil || song.Uploader.StorageID == 0 {
		song.emptyUser.StorageID = databaseSong.UploaderID
		song.Uploader = song.emptyUser
	}
	err = song.Uploader.HydrateByID()
	if err != nil {
		return song, err
	}
	comments, err := song.emptyComment.FetchBySongID(song.StorageID)
	if err != nil {
		return song, err
	}
	song.Comments = comments
	song.status = Status(databaseSong.Status)
	song.statusMessage = databaseSong.StatusMessage
	song.LockExpireTs = databaseSong.LockExpireTS
	song.Date.UpdatedAt = databaseSong.UpdatedAt
	song.Date.CreatedAt = databaseSong.CreatedAt
	return song, nil
}
