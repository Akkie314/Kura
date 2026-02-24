package models

type CompleteRequest struct {
	UserID string `json:"user_id" binding:"required"`
}
