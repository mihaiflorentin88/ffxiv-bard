package database

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
)

type CommentRepository struct {
	driver contract.DatabaseDriverInterface
}

func NewCommentRepository(driver contract.DatabaseDriverInterface) contract.CommentRepositoryInterface {
	return &CommentRepository{
		driver: driver,
	}
}

func (c CommentRepository) FindBySongID(songID int) ([]dto.DatabaseComment, error) {
	var comments []dto.DatabaseComment
	query := `
		SELECT * 
		FROM comment c 
		WHERE c.song_id = ?
	`
	result, err := c.driver.FetchMany(query, songID)
	if err != nil {
		return comments, err
	}
	for result.Next() {
		var comment dto.DatabaseComment
		err := result.Scan(&comment.ID, &comment.AuthorID, &comment.SongID, &comment.Title, &comment.Content, &comment.Status, &comment.CreatedAt, &comment.UpdatedAt)
		if err != nil {
			return comments, fmt.Errorf("error scanning comment for song id `%v`: %w", songID, err)
		}
		comments = append(comments, comment)
	}
	if err = result.Err(); err != nil {
		return comments, err
	}
	return comments, nil
}
