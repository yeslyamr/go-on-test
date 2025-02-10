package models

import "time"

type DailyStat struct {
	LanguageCode string `gorm:"uniqueIndex"`
	ChangeDate   time.Time
	ChangeCount  int64
}
