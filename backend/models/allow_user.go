package models

type AllowUser struct {
	DiscordID string `gorm:"primaryKey;column:discord_id;type:varchar(20);not null" json:"discord_id"`
	Memo      string `gorm:"type:text" json:"memo"`
}
