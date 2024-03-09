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

func (s *SongRepository) InsertNewSong(song dto.DatabaseSong, genreIDs []int, instrumentIDs []int) (int, error) {
	query := `
		INSERT INTO song 
		    (title, artist, ensemble_size, filename, uploader_id, status, status_message, checksum, source, note, audio_crafter, created_at, updated_at)
			  VALUES
		    (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	result, err := s.driver.Execute(query, song.Title, song.Artist, song.EnsembleSize, song.Filename, song.UploaderID, song.Status, song.StatusMessage, song.Checksum, song.Source, song.Note, song.AudioCrafter)
	if err != nil {
		return 0, fmt.Errorf("error inserting new song: %w", err)
	}

	songID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error getting last inserted song ID: %w", err)
	}
	song.ID = int(songID)

	uniqueGenreIDs := make(map[int]struct{})
	for _, genreID := range genreIDs {
		uniqueGenreIDs[genreID] = struct{}{}
	}

	for genreID := range uniqueGenreIDs {
		err = s.insertSongGenre(songID, genreID)
		if err != nil {
			return 0, fmt.Errorf("error inserting song-genre relationship: %w", err)
		}
	}

	uniqueInstrumentIDs := make(map[int]struct{})
	for _, instrumentID := range instrumentIDs {
		uniqueInstrumentIDs[instrumentID] = struct{}{}
	}

	for instrumentID := range uniqueInstrumentIDs {
		err = s.insertSongInstrument(songID, instrumentID)
		if err != nil {
			return 0, fmt.Errorf("error inserting song-instrument relationship: %w", err)
		}
	}

	return int(songID), nil
}

func (s *SongRepository) insertSongGenre(songID int64, genreID int) error {
	query := `INSERT OR IGNORE INTO song_genre (song_id, genre_id) VALUES (?, ?)`
	_, err := s.driver.Execute(query, songID, genreID)
	if err != nil {
		return fmt.Errorf("error inserting into song_genre: %w", err)
	}
	return nil
}

func (s *SongRepository) insertSongInstrument(songID int64, instrumentID int) error {
	query := `INSERT OR IGNORE INTO song_instrument (song_id, instrument_id) VALUES (?, ?)`
	_, err := s.driver.Execute(query, songID, instrumentID)
	if err != nil {
		return fmt.Errorf("error inserting into song_instrument: %w", err)
	}
	return nil
}

func (s *SongRepository) UpdateSong(song dto.DatabaseSong, newGenreIDs []int, newInstrumentIDs []int) error {
	db, err := s.driver.GetConnection()
	if err != nil {
		return err
	}
	uniqueGenreIDs := removeDuplicates(newGenreIDs)
	uniqueInstrumentIDs := removeDuplicates(newInstrumentIDs)

	// Starts the transaction

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Updates the song
	_, err = tx.Exec("UPDATE song SET title = ?, artist = ?, ensemble_size = ?, source = ?, audio_crafter = ?, note = ? WHERE id = ?",
		song.Title, song.Artist, song.EnsembleSize, song.Source, song.AudioCrafter, song.Note, song.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Updates the Genres
	var currentGenreIDs []int
	rows, err := tx.Query("SELECT genre_id FROM song_genre WHERE song_id = ?", song.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return err
		}
		currentGenreIDs = append(currentGenreIDs, id)
	}

	genresToAdd, genresToRemove := diffGenres(currentGenreIDs, uniqueGenreIDs)

	for _, genreID := range genresToRemove {
		_, err = tx.Exec("DELETE FROM song_genre WHERE song_id = ? AND genre_id = ?", song.ID, genreID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, genreID := range genresToAdd {
		_, err = tx.Exec("INSERT INTO song_genre (song_id, genre_id) VALUES (?, ?)", song.ID, genreID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Updates the Instruments
	var currentInstrumentIDs []int
	rows, err = tx.Query("SELECT instrument_id FROM song_instrument WHERE song_id = ?", song.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			tx.Rollback()
			return err
		}
		currentInstrumentIDs = append(currentInstrumentIDs, id)
	}

	instrumentsToAdd, instrumentsToRemove := diffInstruments(currentInstrumentIDs, uniqueInstrumentIDs)

	for _, instrumentID := range instrumentsToRemove {
		_, err = tx.Exec("DELETE FROM song_instrument WHERE song_id = ? AND instrument_id = ?", song.ID, instrumentID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, instrumentID := range instrumentsToAdd {
		_, err = tx.Exec("INSERT INTO song_instrument (song_id, instrument_id) VALUES (?, ?)", song.ID, instrumentID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *SongRepository) UpdateStatus(songID int, status int, message string) error {
	query := `UPDATE song SET status = ?, status_message = ? WHERE id = ?`
	_, err := s.driver.Execute(query, status, message, songID)
	if err != nil {
		return fmt.Errorf("error updating song status. Reason: %w", err)
	}
	return nil
}

func (s *SongRepository) FindByChecksum(checksum string) (dto.DatabaseSong, error) {
	var songDTO dto.DatabaseSong

	query := `SELECT id, title, artist, ensemble_size, filename, uploader_id, status, status_message, checksum, lock_expire_ts, created_at, updated_at, source, note, audio_crafter
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
		&songDTO.Source,
		&songDTO.Note,
		&songDTO.AudioCrafter,
	)
	if err != nil {
		return dto.DatabaseSong{}, fmt.Errorf("error scanning song by checksum: %w", err)
	}

	return songDTO, nil
}

func (s *SongRepository) FindByID(songID int) (dto.DatabaseSong, error) {
	var songDTO dto.DatabaseSong

	query := `SELECT id, title, artist, ensemble_size, filename, uploader_id, status, status_message, checksum, lock_expire_ts, created_at, updated_at, source, note, audio_crafter
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
		&songDTO.Source,
		&songDTO.Note,
		&songDTO.AudioCrafter,
	)
	if err != nil {
		return dto.DatabaseSong{}, fmt.Errorf("error scanning song by id: %w", err)
	}

	return songDTO, nil
}

func (s *SongRepository) FetchAll() (*[]dto.DatabaseSong, error) {
	var songs []dto.DatabaseSong
	query := `SELECT id, title, artist, ensemble_size, filename, uploader_id, status, status_message, checksum, lock_expire_ts, created_at, updated_at, source, note, audio_crafter
			  FROM song WHERE  status = 2`

	rows, err := s.driver.FetchMany(query)
	if err != nil {
		return nil, fmt.Errorf("error executing query to find song by checksum: %w", err)
	}
	for rows.Next() {
		var song dto.DatabaseSong
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.EnsembleSize, &song.Filename, &song.UploaderID, &song.Status, &song.StatusMessage, &song.Checksum, &song.LockExpireTS, &song.CreatedAt, &song.UpdatedAt, &song.Source, &song.Note, &song.AudioCrafter)
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

func (s *SongRepository) FetchForPagination(songTitle string, artist string, ensembleSize int, audioCrafter string, instrumentID int, genreID int, sort string, limit int, offset int) ([]dto.SongWithDetails, error) {
	var songs []dto.SongWithDetails
	query := `
		SELECT
			s.id,
			s.title,
			s.artist,
			s.ensemble_size,
			u.name AS uploader_name,
			s.audio_crafter as audio_crafter,
			COALESCE(GROUP_CONCAT(DISTINCT g.name), 'N/A') AS genre_name,
			COALESCE(GROUP_CONCAT(DISTINCT i.name), 'N/A') AS instrument_name,
			COALESCE(ROUND(AVG(r.rating), 2), 0) AS average_rating,
			COUNT(DISTINCT c.id) AS total_comments 
		FROM
			song s
		LEFT JOIN user u ON s.uploader_id = u.id
		LEFT JOIN song_genre sg ON s.id = sg.song_id 
		LEFT JOIN genre g ON sg.genre_id = g.id 
		LEFT JOIN song_instrument si ON s.id = si.song_id 
		LEFT JOIN instrument i ON si.instrument_id = i.id 
		LEFT JOIN rating r ON s.id = r.song_id
		LEFT JOIN comment c ON s.id = c.song_id
		WHERE s.status = 2 
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

	if audioCrafter != "" {
		conditions = append(conditions, "s.audio_crafter like ?")
		args = append(args, "%"+audioCrafter+"%")
	}
	if genreID != -1 {
		conditions = append(conditions, "g.id = ?")
		args = append(args, genreID)
	}

	if instrumentID != -1 {
		conditions = append(conditions, "i.id = ?")
		args = append(args, instrumentID)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
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
		err := rows.Scan(&song.ID, &song.Title, &song.Artist, &song.EnsembleSize, &song.UploaderName, &song.AudioCrafter, &song.GenreName, &song.InstrumentName, &song.AverageRating, &song.TotalComments)
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

func (s *SongRepository) FetchTotalSongsForListing(songTitle string, artist string, ensembleSize int, audioCrafter string, instrumentID int, genreID int) (int, error) {
	var totalCount int
	query := `
    SELECT 
        COUNT(DISTINCT s.id)
    FROM song s
    LEFT JOIN user u ON s.uploader_id = u.id
    LEFT JOIN song_genre sg ON s.id = sg.song_id
    LEFT JOIN genre g ON sg.genre_id = g.id
    LEFT JOIN song_instrument si ON s.id = si.song_id    
    LEFT JOIN instrument i ON i.id = si.instrument_id
    LEFT JOIN rating r ON s.id = r.song_id
    LEFT JOIN comment c ON s.id = c.song_id 
    WHERE s.status = 2
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

	if audioCrafter != "" {
		conditions = append(conditions, "s.audio_crafter like ?")
		args = append(args, "%"+audioCrafter+"%")
	}

	if ensembleSize != -1 {
		conditions = append(conditions, "s.ensemble_size = ?")
		args = append(args, ensembleSize)
	}

	if instrumentID != -1 {
		conditions = append(conditions, "i.id = ?")
		args = append(args, instrumentID)
	}

	if genreID != -1 {
		conditions = append(conditions, "g.id = ?")
		args = append(args, genreID)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	args = append(args)
	result, err := s.driver.FetchOne(query, args...)
	result.Scan(&totalCount)
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

func (s *SongRepository) IncrementDownloadCount(songID int) error {
	query := `UPDATE song SET download_count = download_count + 1 WHERE id = ?`
	_, err := s.driver.Execute(query, songID)
	if err != nil {
		return fmt.Errorf("error updating download count: %w", err)
	}
	return nil
}

func removeDuplicates(slice []int) []int {
	seen := make(map[int]bool)
	unique := []int{}

	for _, value := range slice {
		if _, ok := seen[value]; !ok {
			seen[value] = true
			unique = append(unique, value)
		}
	}

	return unique
}
