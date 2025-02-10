package bot

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"goon/internal/language"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

// handleSetLangCommand processes the !setLang command
func (b *Bot) handleSetLangCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !setLang [language_code]")
		return
	}

	langCode := strings.ToLower(args[1])
	if !language.IsValidLanguage(langCode) {
		supported := strings.Join(language.SupportedLanguages, ", ")
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unsupported language code. Supported languages: %s", supported))
		return
	}

	guildID := m.GuildID
	userID := m.Author.ID

	// Set the language preference
	err := b.storage.SetLanguagePreference(guildID, userID, langCode)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to set language preference. Please try again later.")
		return
	}

	// response message
	if guildID != "" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Server language preference set to `%s`.", langCode))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Your language preference set to `%s`.", langCode))
	}
}

// handleRecentCommand processes the !recent command
func (b *Bot) handleRecentCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var langCode string
	var err error

	// check if language is provided
	if len(args) >= 2 {
		langCode = strings.ToLower(args[1])
		if !language.IsValidLanguage(langCode) {
			supported := strings.Join(language.SupportedLanguages, ", ")
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unsupported language code. Supported languages: %s", supported))
			return
		}
	} else {
		guildID := m.GuildID
		userID := m.Author.ID

		langCode, err = b.storage.GetLanguagePreference(guildID, userID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Failed to retrieve language preference.")
			return
		}
	}

	// get recent change for the language
	change, err := b.storage.GetRecentChangeByLanguage(langCode)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No recent change found for this language: `%s`.", langCode))
		return
	} else if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve recent change. Please try again later.")
		return
	}

	// response
	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Recent Wikipedia Changes (%s)", strings.ToUpper(langCode)),
		Color: 0x00ff00,
	}

	t := time.Unix(change.Timestamp, 0).Format(time.RFC1123)
	fieldValue := fmt.Sprintf("**User:** %s\n**Time:** %s\n**Server URL:**(%s)", change.User, t, change.URL)
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   change.Title,
		Value:  fieldValue,
		Inline: false,
	})

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// handleStatsCommand processes the !stats command
func (b *Bot) handleStatsCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !stats [yyyy-mm-dd]")
		return
	}

	dateStr := args[1]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid date format. Please use `yyyy-mm-dd`.")
		return
	}

	guildID := m.GuildID
	userID := m.Author.ID
	langCode, err := b.storage.GetLanguagePreference(guildID, userID)
	if err != nil {
		log.Printf("Error retrieving language preference: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve language preference.")
		return
	}

	count, err := b.storage.GetDailyChangeCount(langCode, date)
	if err != nil {
		log.Printf("Error querying daily stats: %v", err)
		s.ChannelMessageSend(m.ChannelID, "Failed to retrieve daily stats. Please try again later.")
		return
	}

	response := fmt.Sprintf("📊 **Wikipedia Changes on %s (%s):**\nTotal Changes: %d", dateStr, strings.ToUpper(langCode), count)
	s.ChannelMessageSend(m.ChannelID, response)
}
