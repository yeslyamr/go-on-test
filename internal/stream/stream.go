package stream

import (
	"encoding/json"
	"github.com/r3labs/sse/v2"
	"goon/internal/language"
	"goon/internal/models"
	"goon/internal/storage"
	"log"
	"time"
)

type Handler struct {
	client  *sse.Client
	storage *storage.Storage
}

func NewHandler(storage *storage.Storage) *Handler {
	client := sse.NewClient("https://stream.wikimedia.org/v2/stream/recentchange")
	return &Handler{
		client:  client,
		storage: storage,
	}
}

func (h *Handler) Start() {
	err := h.client.Subscribe("message", func(msg *sse.Event) {
		var change models.Change
		if err := json.Unmarshal(msg.Data, &change); err != nil {
			log.Printf("Error parsing change: %v", err)
			return
		}

		langCode := language.ExtractLanguageCode(change.ServerName)
		change.LanguageCode = langCode

		// update the recent change for the language
		_ = h.storage.SetRecentChangeByLanguage(change, langCode)

		// increment ChangeCount for the language and date
		date := time.Unix(change.Timestamp, 0)
		_ = h.storage.IncrementDailyChangeCount(langCode, date)
	})

	if err != nil {
		log.Printf("Error starting event subscription: %v", err)
	}
}
