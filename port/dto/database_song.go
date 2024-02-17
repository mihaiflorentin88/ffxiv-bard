package dto

import (
	"time"
)

type DatabaseSongDTO struct {
	ID            int        `json:"id" db:"id"`
	Title         string     `json:"title" db:"title"`
	Artist        string     `json:"artist" db:"artist"`
	EnsembleSize  int        `json:"ensemble_size" db:"ensemble_size"`
	FileCode      string     `json:"file_code" db:"file_code"`
	UploaderID    int64      `json:"uploader_id" db:"uploader_id"`
	Status        int        `json:"status" db:"status"`
	StatusMessage *string    `json:"status_message" db:"status_message"`
	Checksum      string     `json:"checksum" db:"checksum"`
	LockExpireTS  *time.Time `json:"lock_expire_ts" db:"lock_expire_ts"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}
