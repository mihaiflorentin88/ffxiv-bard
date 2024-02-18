package dto

type NewSongForm struct {
	Title        string
	Artist       string
	EnsembleSize int
	Genre        []int
	File         []byte
	User         interface{} // this is type user.User()
}

func AddNewSong(title string, artist string, ensembleSize int, genre []int, file []byte, user interface{}) NewSongForm {
	return NewSongForm{
		Title:        title,
		Artist:       artist,
		EnsembleSize: ensembleSize,
		Genre:        genre,
		File:         file,
		User:         user,
	}
}
