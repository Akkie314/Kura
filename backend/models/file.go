package models

import "time"

type File struct {
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:id;not null" json:"id"`
	UserID       string    `gorm:"type:uuid;column:user_id; not null" json:"user_id"`
	OriginalName string    `gorm:"type:varchar(255);not null" json:"original_name"`
	MimeType     string    `gorm:"type:varchar(255);not null" json:"mime_type"`
	HasThumbnail bool      `gorm:"default:false" json:"has_thumbnail"`
	Size         int64     `gorm:"not null" json:"size"`
	FilePath     string    `gorm:"type:text;not null" json:"file_path"`
	HashSha256   string    `gorm:"type:char(64);not null" json:"hash_sha256"`
	DeleteAt     time.Time `gorm:"not null" json:"delete_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
