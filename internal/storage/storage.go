package storage

import (
	"errors"
	"fmt"
	"goon/internal/language"
	"goon/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

type Storage struct {
	DB *gorm.DB
}

// New creates a new Storage instance
func New() *Storage {
	return &Storage{}
}

// Initialize initializes the SQLite database and performs migrations
func (s *Storage) Initialize() {
	var err error
	s.DB, err = gorm.Open(sqlite.Open("pref.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = s.DB.AutoMigrate(&models.Preference{}, &models.Change{}, &models.DailyStat{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// LANGUAGE PREFERENCE

func (s *Storage) GetLanguagePreference(guildID, userID string) (string, error) {
	var pref models.Preference

	// check user preference first
	result := s.DB.Where("type = ? AND discord_id = ?", "user", userID).First(&pref)
	if result.Error == nil {
		return pref.Language, nil
	}

	// check server preference
	if guildID != "" {
		result = s.DB.Where("type = ? AND discord_id = ?", "server", guildID).First(&pref)
		if result.Error == nil {
			return pref.Language, nil
		}
	}

	return language.Default, nil
}

func (s *Storage) SetLanguagePreference(guildID, userID, lang string) error {
	var pref models.Preference
	var configType, discordID string

	if guildID != "" {
		// server preference
		configType = "server"
		discordID = guildID
	} else {
		// user preference
		configType = "user"
		discordID = userID
	}

	// push preference
	result := s.DB.Where("type = ? AND discord_id = ?", configType, discordID).First(&pref)
	if result.Error != nil {
		// create new pref
		pref = models.Preference{
			Type:      configType,
			DiscordID: discordID,
			Language:  lang,
		}
		if err := s.DB.Create(&pref).Error; err != nil {
			return err
		}
	} else {
		// update existing pref
		pref.Language = lang
		if err := s.DB.Save(&pref).Error; err != nil {
			return err
		}
	}
	return nil
}

// RECENT CHANGES

func (s *Storage) GetRecentChangeByLanguage(lang string) (models.Change, error) {
	var change models.Change
	result := s.DB.Where("language_code = ?", lang).First(&change)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Change{}, gorm.ErrRecordNotFound
		}
		return models.Change{}, fmt.Errorf("database error: %w", result.Error)
	}

	return change, nil
}

func (s *Storage) SetRecentChangeByLanguage(change models.Change, lang string) error {
	if result := s.DB.Where("language_code = ?", lang).Assign(change).FirstOrCreate(&change); result.Error != nil {
		return result.Error
	}
	return nil
}

// STATS

func (s *Storage) GetDailyChangeCount(lang string, date time.Time) (int64, error) {
	var count int64

	// Truncate the date to the day
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	result := s.DB.Model(&models.DailyStat{}).
		Where("language_code = ? AND change_date = ?", lang, date).
		Pluck("change_count", &count)

	return count, result.Error
}

func (s *Storage) IncrementDailyChangeCount(lang string, date time.Time) error {
	// Truncate the date to the day
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	result := s.DB.Model(&models.DailyStat{}).
		Where("language_code = ? AND change_date = ?", lang, date).
		Update("change_count", gorm.Expr("change_count + 1"))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		newStat := models.DailyStat{LanguageCode: lang, ChangeDate: date, ChangeCount: 1}
		if err := s.DB.Create(&newStat).Error; err != nil {
			return err
		}
	}

	return nil
}
