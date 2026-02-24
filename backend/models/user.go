package models

import "time"

type User struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid();column:id; not null" json:"id"`
	DiscordID string    `gorm:"unique;column:discord_id;type:varchar(20);not null" json:"discord_id"`
	Username  string    `gorm:"type:varchar(32);not null" json:"username"`
	Role      string    `gorm:"type:text;default:'user'" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
