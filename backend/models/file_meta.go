package models

type FileMeta struct {
	ID           string `json:"id"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	HasThumbnail bool   `json:"has_thumbnail"`
}
