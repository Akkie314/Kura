package models

type SessionRequest struct {
	FileName string `json:"file_name" binding:"required"`
	FileType string `json:"file_type" binding:"required"`
	FileSize int64  `json:"file_size" binding:"required"`
}
