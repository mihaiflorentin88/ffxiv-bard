package song

import (
	"ffxvi-bard/domain/date"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
)

type Comment struct {
	StorageID         int
	Author            user.User
	Content           string
	Status            bool
	Date              date.Date
	commentRepository contract.CommentRepositoryInterface
	emptyUser         user.User
}

func NewComment(content string, author user.User, status bool) *Comment {
	return &Comment{
		Content: content,
		Author:  author,
		Status:  status,
	}
}

func NewEmptyComment(commentRepository contract.CommentRepositoryInterface, emptyUser user.User) Comment {
	return Comment{
		commentRepository: commentRepository,
		emptyUser:         emptyUser,
	}
}

func (c *Comment) FetchBySongID(songID int) ([]Comment, error) {
	var comments []Comment
	commentDTOs, err := c.commentRepository.FindBySongID(songID)
	if err != nil {
		return comments, err
	}
	for _, commentDTO := range commentDTOs {
		comment, err := FromCommentDTO(commentDTO, c.emptyUser)
		if err != nil {
			return comments, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func FromCommentDTO(commentDTO dto.DatabaseComment, emptyUser user.User) (Comment, error) {
	comment := Comment{
		StorageID: commentDTO.ID,
		Content:   commentDTO.Content,
		Status:    commentDTO.Status,
	}
	comment.Date.CreatedAt = commentDTO.CreatedAt
	comment.Date.UpdatedAt = commentDTO.UpdatedAt
	emptyUser.StorageID = commentDTO.AuthorID
	err := emptyUser.HydrateByID()
	if err != nil {
		return comment, err
	}
	comment.Author = emptyUser
	return comment, err
}
