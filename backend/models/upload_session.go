package models

import "time"

type UploadSession struct {
	UploadID     string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:upload_id;not null" json:"upload_session"`
	UserID       string    `gorm:"type:uuid;column:user_id; not null" json:"user_id"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"original_name"`
	MimeType     string    `gorm:"type:varchar(255);not null" json:"mime_type"`
	TotalSize    int64     `gorm:"not null" json:"total_size"`
	StoredSize   int64     `gorm:"not null" json:"stored_size"`
	TempFilePath string    `gorm:"type:text;not null" json:"temp_file_path"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
