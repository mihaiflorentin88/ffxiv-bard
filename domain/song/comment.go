package song

import (
	"ffxvi-bard/domain/date"
	"ffxvi-bard/domain/user"
)

type Comment struct {
	storageID int
	Author    user.User
	Title     string
	Content   string
	Likes     int
	Dislikes  int
	Date      date.Date
}

func NewComment(title string, content string, author user.User, likes int, dislikes int) *Comment {
	return &Comment{
		Title:    title,
		Content:  content,
		Author:   author,
		Likes:    likes,
		Dislikes: dislikes,
	}
}

func (c *Comment) Like() {
	c.Likes++
}

func (c *Comment) Dislike() {
	c.Dislikes++
}

func (c *Comment) GetStorageID() int {
	return c.storageID
}

func (c *Comment) SetStorageID(id int) {
	c.storageID = id
}
