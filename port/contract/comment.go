package contract

type CommentInterface interface {
	Like()
	Dislike()
	GetStorageID() int
	SetStorageID(id int)
}
