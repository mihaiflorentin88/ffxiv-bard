package database

import (
	"ffxvi-bard/port/contract"
	"ffxvi-bard/port/dto"
	"fmt"
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
