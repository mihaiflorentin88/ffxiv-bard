package contract

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

type SongInterface interface {
	EnsembleString() string
	GetDetailedEnsembleString() map[int]string
	StatusString() string
	GenerateFileCode()
	AddComment(c CommentInterface)
	RemoveComment(c CommentInterface)
	GetStatus() Status
	ChangeStatus(status Status, statusMessage string)
	ChangeStatusMessage(statusMessage string)
	GetStatusMessage() string
	ProcessSong() error
	GetAverageRating() float64
	GetFileCode() string
	GetFile() []byte
}
