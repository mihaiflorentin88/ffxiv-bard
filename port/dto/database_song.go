package dto

import (
	"time"
)

type DatabaseSongDTO struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Artist        string     `json:"artist"`
	EnsembleSize  int        `json:"ensemble_size"`
	FileCode      string     `json:"file_code"`
	UploaderID    int64      `json:"uploader_id"`
	Status        int        `json:"status"`
	StatusMessage *string    `json:"status_message"` // This can be null in the database, so it is a pointer
	Checksum      string     `json:"checksum"`
	LockExpireTS  *time.Time `json:"lock_expire_ts"` // This can be null in the database, so it is a pointer
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
