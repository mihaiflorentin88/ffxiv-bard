package form

type Genre struct {
	ID   int
	Name string
}

type Song struct {
	ID           int
	Title        string
	Artist       string
	EnsembleSize int
	Genre        []Genre
}
