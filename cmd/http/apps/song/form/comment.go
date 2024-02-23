package form

import (
	"errors"
	"ffxvi-bard/domain/user"
	"ffxvi-bard/port/contract"
	"fmt"
)

type SubmitCommentForm struct {
	commentRepository contract.CommentRepositoryInterface
}

func NewSubmitCommentForm(commentRepository contract.CommentRepositoryInterface) SubmitCommentForm {
	return SubmitCommentForm{
		commentRepository: commentRepository,
	}
}

func (f SubmitCommentForm) Submit(loggedUser *user.User, songID int, comment string) error {
	_, err := f.commentRepository.InsertComment(loggedUser.StorageID, songID, comment)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot submit comment. Reason %s", err))
	}
	return nil
}
