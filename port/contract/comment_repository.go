package contract

import "ffxvi-bard/port/dto"

type CommentRepositoryInterface interface {
	FindBySongID(songID int) ([]dto.DatabaseComment, error)
	FindByID(commentID int) (dto.DatabaseComment, error)
	InsertComment(authorID int64, songID int, content string) (int64, error)
	UpdateComment(commentID int, content string) error
}
