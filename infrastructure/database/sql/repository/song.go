package database

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
	"strings"
)

type SongRepository struct {
	driver contract.DatabaseDriverInterface
}

func NewSongRepository(driver contract.DatabaseDriverInterface) contract.SongRepositoryInterface {
	return &SongRepository{
		driver: driver,
	}
}

func (s *SongRepository) InsertNewSong(song dto.DatabaseSong, genreIDs []int) error {
	query := `
		INSERT INTO song 
		    (title, artist, ensemble_size, file_code, uploader_id, status, status_message, checksum, created_at, updated_at)
			  VALUES
		    (?, ?, ?, ?, ?, ?, ?, ?,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := s.driver.Execute(query, song.Title, song.Artist, song.EnsembleSize, song.FileCode, song.UploaderID, song.Status, song.StatusMessage, song.Checksum)
	if err != nil {
		return fmt.Errorf("error inserting new song: %w", err)
	}

	songID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last inserted song ID: %w", err)
	}

	for _, genreID := range genreIDs {
		err = s.insertSongGenre(songID, genreID)
		if err != nil {
			return fmt.Errorf("error inserting song-genre relationship: %w", err)
		}
	}

	return nil
}

func (s *SongRepository) insertSongGenre(songID int64, genreID int) error {
	query := `INSERT INTO song_genre (song_id, genre_id) VALUES (?, ?)`
	_, err := s.driver.Execute(query, songID, genreID)
	if err != nil {
		return fmt.Errorf("error inserting into song_genre: %w", err)
	}
	return nil
}

func (s *SongRepository) FindByChecksum(checksum string) (dto.DatabaseSong, error) {
	var songDTO dto.DatabaseSong

	query := `SELECT id, title, artist, ensemble_size, file_code, uploader_id, status, status_message, checksum, lock_expire_ts, created_at, updated_at
			  FROM song
			  WHERE checksum = ?`

	row, err := s.driver.FetchOne(query, checksum)
	if err != nil {
		return dto.DatabaseSong{}, fmt.Errorf("error executing query to find song by checksum: %w", err)
	}
	err = row.Scan(
		&songDTO.ID,
		&songDTO.Title,
		&songDTO.Artist,
		&songDTO.EnsembleSize,
		&songDTO.FileCode,
		&songDTO.UploaderID,
		&songDTO.Status,
		&songDTO.StatusMessage,
		&songDTO.Checksum,
		&songDTO.LockExpireTS,
		&songDTO.CreatedAt,
		&songDTO.UpdatedAt,
	)
	if err != nil {
		return dto.DatabaseSong{}, fmt.Errorf("error scanning song by checksum: %w", err)
	}

	return songDTO, nil
}

func (s *SongRepository) FetchAll() (*[]dto.DatabaseSong, error) {
	var songs []dto.DatabaseSong
	query := `SELECT id, title, artist, ensemble_size, file_code, uploader_id, status, status_message, checksum, lock_expire_ts, created_at, updated_at
			  FROM song`

	rows, err := s.driver.FetchMany(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query to find song by checksum: %w", err)
	}
	for rows.Next() {
		var song dto.DatabaseSong
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.EnsembleSize, &song.FileCode, &song.UploaderID, &song.Status, &song.StatusMessage, &song.Checksum, &song.LockExpireTS, &song.CreatedAt, &song.UpdatedAt)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return &songs, nil
}

func (s *SongRepository) FetchForPagination(songTitle string, artist string, ensembleSize int, genreID int, limit int, offset int) ([]dto.SongWithDetails, int, error) {
	var songs []dto.SongWithDetails
	var totalCount int

	baseQuery := `
    SELECT s.id, s.title, s.artist, s.ensemble_size, g.name AS genre_name, u.name AS uploader_name, 
           COALESCE(ROUND(AVG(r.rating),2), 0) AS average_rating, COALESCE(COUNT(c.id), 0) AS total_comments
    FROM song s
    LEFT JOIN user u ON s.uploader_id = u.id
    LEFT JOIN song_genre sg ON s.id = sg.song_id
    LEFT JOIN genre g ON sg.genre_id = g.id
    LEFT JOIN rating r ON s.id = r.song_id
    LEFT JOIN comment c ON s.id = c.song_id
    `

	var conditions []string
	var args []interface{}

	if songTitle != "" {
		conditions = append(conditions, "s.title LIKE ?")
		args = append(args, strings.ToLower("%"+songTitle+"%"))
	}
	if artist != "" {
		conditions = append(conditions, "s.artist LIKE ?")
		args = append(args, strings.ToLower("%"+artist+"%"))
	}
	if ensembleSize != -1 {
		conditions = append(conditions, "s.ensemble_size = ?")
		args = append(args, ensembleSize)
	}
	if genreID != -1 {
		conditions = append(conditions, "g.id = ?")
		args = append(args, genreID)
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += `
    GROUP BY s.id, s.title, s.artist, s.ensemble_size, g.name, u.name
	ORDER BY s.artist, s.created_at DESC
    LIMIT ? OFFSET ?
    `
	args = append(args, limit, offset)

	// Execute the query
	rows, err := s.driver.FetchMany(baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var song dto.SongWithDetails
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.EnsembleSize, &song.GenreName, &song.UploaderName, &song.AverageRating, &song.TotalComments)
		if err != nil {
			continue
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := `
    SELECT COUNT(DISTINCT s.id)
    FROM song s
    LEFT JOIN user u ON s.uploader_id = u.id
    LEFT JOIN song_genre sg ON s.id = sg.song_id
    LEFT JOIN genre g ON sg.genre_id = g.id
    LEFT JOIN rating r ON s.id = r.song_id
    LEFT JOIN comment c ON s.id = c.song_id
    `

	// Append the same conditions to the count query
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	result, err := s.driver.FetchOne(countQuery, args[:len(args)-2]...)
	result.Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}
	return songs, totalCount, nil
}
