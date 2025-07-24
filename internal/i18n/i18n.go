package i18n

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// I18nService interface for dependency injection and testing
type I18nService interface {
	GetMessage(lang string, messageID string) (string, error)
	GetMessageWithData(lang string, messageID string, data map[string]interface{}) (string, error)
	GetLocalizer(lang string) *i18n.Localizer
	SetDefaultLanguage(lang string) error
	GetDefaultLanguage() string
}

// Bundle implementation of I18nService
type Bundle struct {
	bundle          *i18n.Bundle
	defaultLanguage language.Tag
}

// NewBundle creates a new Bundle with the specified locales path and default language
func NewBundleWithDefaults(localesPath string, defaultLang language.Tag) (*Bundle, error) {
	// Validate and sanitize the locales path
	if localesPath == "" {
		return nil, fmt.Errorf("locales path cannot be empty")
	}

	// Clean the path to prevent directory traversal attacks
	cleanPath := filepath.Clean(localesPath)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("invalid locales path: directory traversal not allowed")
	}

	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Read all JSON files from locales directory
	files, err := os.ReadDir(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read locales directory: %w", err)
	}

	loadedCount := 0
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(cleanPath, file.Name())
		if _, err := bundle.LoadMessageFile(filePath); err != nil {
			return nil, fmt.Errorf("failed to load translation file %s: %w", filePath, err)
		}
		log.Printf("Loaded translation file: %s", filePath)
		loadedCount++
	}

	if loadedCount == 0 {
		return nil, fmt.Errorf("no translation files found in directory: %s", cleanPath)
	}

	return &Bundle{
		bundle:          bundle,
		defaultLanguage: defaultLang,
	}, nil
}

// NewBundle creates a new Bundle with English as default language (backward compatibility)
func NewBundle(localesPath string) (*Bundle, error) {
	return NewBundleWithDefaults(localesPath, language.English)
}

// SetDefaultLanguage sets the default language for the bundle
func (b *Bundle) SetDefaultLanguage(lang string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return fmt.Errorf("invalid language tag %s: %w", lang, err)
	}
	b.defaultLanguage = tag
	return nil
}

// GetDefaultLanguage returns the default language of the bundle
func (b *Bundle) GetDefaultLanguage() string {
	return b.defaultLanguage.String()
}

// GetLocalizer returns a Localizer for the given language tag
func (b *Bundle) GetLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(b.bundle, lang, b.defaultLanguage.String())
}

// GetMessage returns a translated message for the given language and message ID
func (b *Bundle) GetMessage(lang string, messageID string) (string, error) {
	if lang == "" {
		return "", fmt.Errorf("language cannot be empty")
	}

	if messageID == "" {
		return "", fmt.Errorf("messageID cannot be empty")
	}

	localizer := b.GetLocalizer(lang)
	return localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

// GetMessageWithData returns a translated message with template data
func (b *Bundle) GetMessageWithData(lang string, messageID string, data map[string]interface{}) (string, error) {
	if lang == "" {
		return "", fmt.Errorf("language cannot be empty")
	}

	if messageID == "" {
		return "", fmt.Errorf("messageID cannot be empty")
	}

	localizer := b.GetLocalizer(lang)
	return localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
}

// Global instance for backward compatibility
var bundle *i18n.Bundle

// Init initializes the i18n bundle and loads translations from default path
func Init() {
	err := InitWithLocalesPath("internal/i18n/locales")
	if err != nil {
		log.Fatalf("Failed to initialize i18n: %v", err)
	}
}

// InitWithLocalesPath initializes the i18n bundle and loads translations from specified path
func InitWithLocalesPath(localesPath string) error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Read all JSON files from locales directory
	files, err := os.ReadDir(localesPath)
	if err != nil {
		return fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(localesPath, file.Name())
		if _, err := bundle.LoadMessageFile(filePath); err != nil {
			return fmt.Errorf("failed to load translation file %s: %w", filePath, err)
		}
		log.Printf("Loaded translation file: %s", filePath)
	}

	return nil
}

// GetLocalizer returns a Localizer for the given language tag (e.g., "en", "pt")
// Returns error instead of panicking for production safety
func GetLocalizer(lang string) (*i18n.Localizer, error) {
	if bundle == nil {
		return nil, fmt.Errorf("i18n bundle not initialized. Call Init() first")
	}
	return i18n.NewLocalizer(bundle, lang), nil
}

// GetLocalizerUnsafe returns a Localizer for the given language tag (legacy method that panics)
// Deprecated: Use GetLocalizer instead for production safety
func GetLocalizerUnsafe(lang string) *i18n.Localizer {
	if bundle == nil {
		panic("i18n bundle not initialized. Call Init() first")
	}
	return i18n.NewLocalizer(bundle, lang)
}

// GetMessage returns a translated message for the given language and message ID
func GetMessage(lang string, messageID string) (string, error) {
	if bundle == nil {
		return "", fmt.Errorf("i18n bundle not initialized")
	}

	if lang == "" {
		return "", fmt.Errorf("language cannot be empty")
	}

	if messageID == "" {
		return "", fmt.Errorf("messageID cannot be empty")
	}

	localizer, err := GetLocalizer(lang)
	if err != nil {
		return "", err
	}

	return localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
}

// GetMessageWithData returns a translated message with template data
func GetMessageWithData(lang string, messageID string, data map[string]interface{}) (string, error) {
	if bundle == nil {
		return "", fmt.Errorf("i18n bundle not initialized")
	}

	if lang == "" {
		return "", fmt.Errorf("language cannot be empty")
	}

	if messageID == "" {
		return "", fmt.Errorf("messageID cannot be empty")
	}

	localizer, err := GetLocalizer(lang)
	if err != nil {
		return "", err
	}

	return localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
}
