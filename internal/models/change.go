package models

import (
	"gorm.io/gorm"
)

type Change struct {
	gorm.Model
	LanguageCode string `gorm:"uniqueIndex"`
	Title        string `json:"title"`
	User         string `json:"user"`
	Timestamp    int64  `json:"timestamp"`
	URL          string `json:"title_url"`
	ServerName   string `json:"server_name"`
}
