package contract

import "ffxvi-bard/port/dto"

type CommentRepositoryInterface interface {
	FindBySongID(songID int) ([]dto.DatabaseComment, error)
}
