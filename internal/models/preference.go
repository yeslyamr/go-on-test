package models

import "gorm.io/gorm"

type Preference struct {
	gorm.Model
	Type      string `gorm:"index"` // user or server
	DiscordID string `gorm:"index"` // UserID or ServerID
	Language  string
}
