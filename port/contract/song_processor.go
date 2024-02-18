package contract

type SongProcessorInterface interface {
	ProcessSong(songFilename string) error
	WriteUnprocessedSong(songFilename string, song []byte) error
	RemoveUnprocessedSong(songFilename string) error
}
