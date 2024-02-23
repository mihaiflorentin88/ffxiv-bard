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
		err := result.Scan(&comment.ID, &comment.AuthorID, &comment.SongID, &comment.Content, &comment.Status, &comment.CreatedAt, &comment.UpdatedAt)
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

func (c CommentRepository) InsertComment(authorID int64, songID int, content string) (int64, error) {
	query := `
		INSERT INTO comment (author_id, song_id, content)
		VALUES (?, ?, ?)
	`
	result, err := c.driver.Execute(query, authorID, songID, content)
	if err != nil {
		return 0, err
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return commentID, nil
}
