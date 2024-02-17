package database

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
)

type CommentRepository struct {
	driver contract.DatabaseDriverInterface
}

func NewCommentRepository(driver contract.DatabaseDriverInterface) contract.CommentRepositoryInterface {
	return &CommentRepository{
		driver: driver,
	}
}

func (c CommentRepository) FindBySongId(songID int) ([]dto.DatabaseComment, error) {
	var comments []dto.DatabaseComment
	query := `
		SELECT * 
		FROM comment c 
		WHERE c.song_id = ?
	`
	result, err := c.driver.FetchMany(query, &songID)
	if err != nil {
		return comments, err
	}
	for result.Next() {
		var comment dto.DatabaseComment
		result.Scan(&comment.ID)
		result.Scan(&comment.SongID)
		result.Scan(&comment.AuthorID)
		result.Scan(&comment.Status)
		result.Scan(&comment.Title)
		result.Scan(&comment.Content)
		result.Scan(&comment.CreatedAt)
		result.Scan(&comment.UpdatedAt)
		comments = append(comments, comment)
	}
	if err = result.Err(); err != nil {
		return comments, err
	}
	return comments, nil
}
