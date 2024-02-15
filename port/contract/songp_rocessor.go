package contract

type SongProcessorInterface interface {
	ProcessSong(currentSong SongInterface) error
	IsCorrectFormat() bool
	WriteUnprocessedSong(currentSong SongInterface) error
	RemoveUnprocessedSong(currentSong SongInterface) error
}
