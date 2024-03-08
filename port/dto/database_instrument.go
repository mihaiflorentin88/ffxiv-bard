package dto

type DatabaseInstrument struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
