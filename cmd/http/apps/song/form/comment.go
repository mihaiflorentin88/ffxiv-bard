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

func (f SubmitCommentForm) SubmitUpdate(loggedUser *user.User, comment string, commentID int) error {
	oldComment, err := f.commentRepository.FindByID(commentID)
	if err != nil {
		return err
	}
	if !loggedUser.IsAdmin && loggedUser.StorageID != oldComment.AuthorID {
		return errors.New("you do not have permissions to edit this comment")
	}

	err = f.commentRepository.UpdateComment(commentID, comment)
	if err != nil {
		return errors.New(fmt.Sprintf("cannot submit comment. Reason %s", err))
	}
	return nil
}
