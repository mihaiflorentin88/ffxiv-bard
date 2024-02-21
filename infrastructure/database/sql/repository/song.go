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

func (s *SongRepository) InsertNewSong(song dto.DatabaseSong, genreIDs []int) (int, error) {
	query := `
		INSERT INTO song 
		    (title, artist, ensemble_size, filename, uploader_id, status, status_message, checksum, created_at, updated_at)
			  VALUES
		    (?, ?, ?, ?, ?, ?, ?, ?,  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := s.driver.Execute(query, song.Title, song.Artist, song.EnsembleSize, song.Filename, song.UploaderID, song.Status, song.StatusMessage, song.Checksum)
	if err != nil {
		return 0, fmt.Errorf("error inserting new song: %w", err)
	}

	songID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted song ID: %w", err)
	}
	song.ID = int(songID)

	for _, genreID := range genreIDs {
		err = s.insertSongGenre(songID, genreID)
		if err != nil {
			return 0, fmt.Errorf("error inserting song-genre relationship: %w", err)
		}
	}

	return int(songID), nil
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

	query := `SELECT id, title, artist, ensemble_size, filename, uploader_id, status, status_message, checksum, lock_expire_ts, created_at, updated_at
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
		&songDTO.Filename,
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

func (s *SongRepository) FindByID(songID int) (dto.DatabaseSong, error) {
	var songDTO dto.DatabaseSong

	query := `SELECT id, title, artist, ensemble_size, filename, uploader_id, status, status_message, checksum, lock_expire_ts, created_at, updated_at
			  FROM song
			  WHERE id = ?`

	row, err := s.driver.FetchOne(query, songID)
	if err != nil {
		return songDTO, fmt.Errorf("error executing query to find song by id: %w", err)
	}
	err = row.Scan(
		&songDTO.ID,
		&songDTO.Title,
		&songDTO.Artist,
		&songDTO.EnsembleSize,
		&songDTO.Filename,
		&songDTO.UploaderID,
		&songDTO.Status,
		&songDTO.StatusMessage,
		&songDTO.Checksum,
		&songDTO.LockExpireTS,
		&songDTO.CreatedAt,
		&songDTO.UpdatedAt,
	)
	if err != nil {
		return dto.DatabaseSong{}, fmt.Errorf("error scanning song by id: %w", err)
	}

	return songDTO, nil
}

func (s *SongRepository) FetchAll() (*[]dto.DatabaseSong, error) {
	var songs []dto.DatabaseSong
	query := `SELECT id, title, artist, ensemble_size, filename, uploader_id, status, status_message, checksum, lock_expire_ts, created_at, updated_at
			  FROM song`

	rows, err := s.driver.FetchMany(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query to find song by checksum: %w", err)
	}
	for rows.Next() {
		var song dto.DatabaseSong
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.EnsembleSize, &song.Filename, &song.UploaderID, &song.Status, &song.StatusMessage, &song.Checksum, &song.LockExpireTS, &song.CreatedAt, &song.UpdatedAt)
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

func (s *SongRepository) FetchForPagination(songTitle string, artist string, ensembleSize int, genreID int, sort string, limit int, offset int) ([]dto.SongWithDetails, error) {
	var songs []dto.SongWithDetails

	query := `
			SELECT
				s.id,
				s.title,
				s.artist,
				s.ensemble_size,
				u.name AS uploader_name,
						(SELECT GROUP_CONCAT(g.name)
						 FROM genre g
						 LEFT JOIN song_genre sg ON g.id = sg.genre_id
						 WHERE s.id = sg.song_id) as genre_name,
				COALESCE(
						(SELECT ROUND(AVG(r.rating), 2)
						 FROM rating r
						 WHERE r.song_id = s.id),
						0) AS average_rating,
				COALESCE(
						(SELECT COUNT(c.id)
						 FROM comment c
						 WHERE c.song_id = s.id),
						0) AS total_comments
			FROM
				song s
			LEFT JOIN user u ON s.uploader_id = u.id
    `

	var conditions []string
	var args []interface{}

	if songTitle != "" {
		conditions = append(conditions, "s.title LIKE ?")
		args = append(args, "%"+songTitle+"%")
	}
	if artist != "" {
		conditions = append(conditions, "s.artist LIKE ?")
		args = append(args, "%"+artist+"%")
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
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += "GROUP BY s.id "

	switch sort {
	case "title_asc":
		query += "ORDER BY s.title ASC "
	case "title_desc":
		query += "ORDER BY s.title DESC "
	case "added_asc":
		query += "ORDER BY s.created_at ASC "
	case "added_desc":
		query += "ORDER BY s.created_at DESC "
	case "artist_asc":
		query += "ORDER BY s.artist ASC "
	case "artist_desc":
		query += "ORDER BY s.artist DESC "
	case "rating_high":
		query += "ORDER BY average_rating DESC "
	case "rating_low":
		query += "ORDER BY average_rating ASC "
	default:
		query += "ORDER BY s.artist, s.created_at DESC "
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := s.driver.FetchMany(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var song dto.SongWithDetails
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.EnsembleSize, &song.UploaderName, &song.GenreName, &song.AverageRating, &song.TotalComments)
		if err != nil {
			continue
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
}

func (s *SongRepository) FetchTotalSongsForListing(songTitle string, artist string, ensembleSize int, genreID int) (int, error) {
	var totalCount int
	query := `
    SELECT 
        COUNT(DISTINCT s.id)
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
		args = append(args, "%"+songTitle+"%")
	}
	if artist != "" {
		conditions = append(conditions, "s.artist LIKE ?")
		args = append(args, "%"+artist+"%")
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
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args)
	result, err := s.driver.FetchOne(query, args...)
	result.Scan(&totalCount)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}
