package language

import (
	"strings"
)

// could use [github.com/biter777/countries] to get all languages (didn't have time to implement)

var SupportedLanguages = []string{
	"en", "es", "de", "fr", "it", "ru", "zh", "ja", "pt", "pl",
	"unknown",
}

const Default = "en"

func IsValidLanguage(lang string) bool {
	for _, l := range SupportedLanguages {
		if l == lang {
			return true
		}
	}
	return false
}

// ExtractLanguageCode derives the language code from the 'server_name' field
func ExtractLanguageCode(serverName string) string {
	parts := strings.Split(serverName, ".")
	if len(parts) < 3 {
		return "unknown"
	}
	langCode := parts[0]

	if IsValidLanguage(langCode) {
		return langCode
	}

	return "unknown"
}
