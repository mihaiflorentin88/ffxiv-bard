package dto

type NewSongForm struct {
	Title        string
	Artist       string
	EnsembleSize int
	Source       string
	Note         string
	AudioCrafter string
	Genre        []int
	Instrument   []int
	File         []byte
	User         interface{} // this is type user.User()
}

func AddNewSong(title string, artist string, ensembleSize int, genre []int, file []byte, user interface{}, source string, note string, audioCrafter string, instrument []int) NewSongForm {
	return NewSongForm{
		Title:        title,
		Artist:       artist,
		EnsembleSize: ensembleSize,
		Genre:        genre,
		File:         file,
		User:         user,
		Source:       source,
		Note:         note,
		AudioCrafter: audioCrafter,
		Instrument:   instrument,
	}
}
