package song

import (
	"ffxvi-bard/domain/date"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
)

type Comment struct {
	storageID         int
	Author            user.User
	Title             string
	Content           string
	Status            bool
	Date              date.Date
	commentRepository contract.CommentRepositoryInterface
}

func NewComment(title string, content string, author user.User, status bool) *Comment {
	return &Comment{
		Title:   title,
		Content: content,
		Author:  author,
		Status:  status,
	}
}

func NewEmptyComment(commentRepository contract.CommentRepositoryInterface) Comment {
	return Comment{
		commentRepository: commentRepository,
	}
}

func (c *Comment) GetStorageID() int {
	return c.storageID
}

func (c *Comment) SetStorageID(id int) {
	c.storageID = id
}
