package dto

type DatabaseSongInstrument struct {
	SongID       int `db:"song_id"`
	InstrumentID int `db:"instrument_id"`
}
