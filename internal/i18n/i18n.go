package i18n

import (
	"encoding/json"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load English translations
	if _, err := bundle.LoadMessageFile("internal/i18n/locales/en.json"); err != nil {
		log.Fatalf("Failed to load English translations: %v", err)
	}

	// Load Portuguese translations
	if _, err := bundle.LoadMessageFile("internal/i18n/locales/pt.json"); err != nil {
		log.Fatalf("Failed to load Portuguese translations: %v", err)
	}
}

// GetLocalizer returns a Localizer for the given language tag (e.g., "en", "pt")
func GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, lang)
}
