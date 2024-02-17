package contract

import "ffxvi-bard/port/dto"

type CommentRepositoryInterface interface {
	FindBySongId(songID int) ([]dto.DatabaseComment, error)
}
