package i18n

import (
	"strings"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func TestInit(t *testing.T) {
	// Test that Init() doesn't panic and initializes the bundle
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Init() panicked: %v", r)
		}
	}()

	Init()

	if bundle == nil {
		t.Error("Init() should initialize the bundle, but bundle is nil")
	}
}

func TestGetLocalizer(t *testing.T) {
	// Arrange
	Init() // Initialize before testing

	testCases := []struct {
		name     string
		language string
		wantNil  bool
	}{
		{
			name:     "English localizer",
			language: "en",
			wantNil:  false,
		},
		{
			name:     "Portuguese localizer",
			language: "pt",
			wantNil:  false,
		},
		{
			name:     "Unknown language",
			language: "fr",
			wantNil:  false, // Should still return a localizer, but with fallback
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			localizer, err := GetLocalizer(tc.language)

			// Assert
			if err != nil {
				t.Errorf("GetLocalizer(%s) error = %v, want nil", tc.language, err)
				return
			}
			if tc.wantNil && localizer != nil {
				t.Errorf("GetLocalizer(%s) = %v, want nil", tc.language, localizer)
			}
			if !tc.wantNil && localizer == nil {
				t.Errorf("GetLocalizer(%s) = nil, want non-nil", tc.language)
			}
		})
	}
}

// TDD: Red phase - This test should fail because we haven't implemented GetMessage yet
func TestGetMessage(t *testing.T) {
	// Arrange
	Init()

	testCases := []struct {
		name         string
		language     string
		messageID    string
		expectedText string
		expectError  bool
	}{
		{
			name:         "Get English welcome message",
			language:     "en",
			messageID:    "welcome_message",
			expectedText: "Welcome to Coup Game!",
			expectError:  false,
		},
		{
			name:         "Get Portuguese welcome message",
			language:     "pt",
			messageID:    "welcome_message",
			expectedText: "Bem-vindo ao Jogo Coup!",
			expectError:  false,
		},
		{
			name:        "Get unknown message",
			language:    "en",
			messageID:   "non_existent_message",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := GetMessage(tc.language, tc.messageID)

			// Assert
			if tc.expectError {
				if err == nil {
					t.Error("GetMessage() should return error for unknown message")
				}
			} else {
				if err != nil {
					t.Errorf("GetMessage() error = %v, wantErr false", err)
					return
				}
				if result != tc.expectedText {
					t.Errorf("GetMessage() = %v, want %v", result, tc.expectedText)
				}
			}
		})
	}
}

func TestTranslateMessage(t *testing.T) {
	// Arrange
	Init()

	testCases := []struct {
		name         string
		language     string
		messageID    string
		expectedText string
	}{
		{
			name:         "English welcome message",
			language:     "en",
			messageID:    "welcome_message",
			expectedText: "Welcome to Coup Game!",
		},
		{
			name:         "Portuguese welcome message",
			language:     "pt",
			messageID:    "welcome_message",
			expectedText: "Bem-vindo ao Jogo Coup!",
		},
		{
			name:         "English waiting message",
			language:     "en",
			messageID:    "waiting_for_players",
			expectedText: "Waiting for players to join...",
		},
		{
			name:         "Portuguese waiting message",
			language:     "pt",
			messageID:    "waiting_for_players",
			expectedText: "Aguardando jogadores entrarem...",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			localizer, err := GetLocalizer(tc.language)
			if err != nil {
				t.Errorf("GetLocalizer() error = %v, want nil", err)
				return
			}

			result, err := localizer.Localize(&i18n.LocalizeConfig{
				MessageID: tc.messageID,
			})

			// Assert
			if err != nil {
				t.Errorf("Localize() error = %v, wantErr false", err)
				return
			}
			if result != tc.expectedText {
				t.Errorf("Localize() = %v, want %v", result, tc.expectedText)
			}
		})
	}
}

func TestGetLocalizerWithoutInit(t *testing.T) {
	// Reset bundle to test behavior without Init
	originalBundle := bundle
	bundle = nil
	defer func() {
		bundle = originalBundle // Restore after test
	}()

	// Test that GetLocalizer returns an error when bundle is not initialized
	localizer, err := GetLocalizer("en")
	if err == nil {
		t.Error("GetLocalizer() should return error when bundle is not initialized")
	}
	if localizer != nil {
		t.Error("GetLocalizer() should return nil localizer when bundle is not initialized")
	}

	// Test legacy unsafe method still panics
	defer func() {
		if r := recover(); r == nil {
			t.Error("GetLocalizerUnsafe() should panic when bundle is not initialized")
		}
	}()

	GetLocalizerUnsafe("en")
}

// TDD: Add more tests for edge cases
func TestGetMessageValidation(t *testing.T) {
	// Arrange
	Init()

	testCases := []struct {
		name        string
		language    string
		messageID   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty language",
			language:    "",
			messageID:   "welcome_message",
			expectError: true,
			errorMsg:    "language cannot be empty",
		},
		{
			name:        "Empty messageID",
			language:    "en",
			messageID:   "",
			expectError: true,
			errorMsg:    "messageID cannot be empty",
		},
		{
			name:        "Valid parameters",
			language:    "en",
			messageID:   "welcome_message",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := GetMessage(tc.language, tc.messageID)

			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("GetMessage() should return error for %s", tc.name)
				}
				if tc.errorMsg != "" && err.Error() != tc.errorMsg {
					t.Errorf("GetMessage() error = %v, want %v", err.Error(), tc.errorMsg)
				}
				if result != "" {
					t.Errorf("GetMessage() result = %v, want empty string for error case", result)
				}
			} else {
				if err != nil {
					t.Errorf("GetMessage() error = %v, wantErr false", err)
				}
				if result == "" {
					t.Error("GetMessage() should return non-empty result for valid parameters")
				}
			}
		})
	}
}

func TestGetMessageWithoutInit(t *testing.T) {
	// Reset bundle to test behavior without Init
	originalBundle := bundle
	bundle = nil
	defer func() {
		bundle = originalBundle // Restore after test
	}()

	// Act
	result, err := GetMessage("en", "welcome_message")

	// Assert
	if err == nil {
		t.Error("GetMessage() should return error when bundle is not initialized")
	}
	if result != "" {
		t.Errorf("GetMessage() result = %v, want empty string when bundle not initialized", result)
	}
	if err.Error() != "i18n bundle not initialized" {
		t.Errorf("GetMessage() error = %v, want 'i18n bundle not initialized'", err.Error())
	}
}

// TDD: Test for dynamic language loading
func TestLoadLanguagesDynamically(t *testing.T) {
	// Create a temporary bundle for testing
	originalBundle := bundle
	defer func() {
		bundle = originalBundle // Restore after test
	}()

	// Act
	err := InitWithLocalesPath("internal/i18n/locales")

	// Assert
	if err != nil {
		t.Errorf("InitWithLocalesPath() error = %v, want nil", err)
	}
	if bundle == nil {
		t.Error("InitWithLocalesPath() should initialize bundle")
	}

	// Test that we can get messages from loaded languages
	result, err := GetMessage("en", "welcome_message")
	if err != nil {
		t.Errorf("GetMessage() after dynamic loading error = %v, want nil", err)
	}
	if result == "" {
		t.Error("GetMessage() should return non-empty result after dynamic loading")
	}
}

// TDD: Test for message with variables/pluralization
func TestGetMessageWithVariables(t *testing.T) {
	// Arrange
	Init()

	testCases := []struct {
		name      string
		language  string
		messageID string
		data      map[string]interface{}
		expected  string
	}{
		{
			name:      "English message with player count",
			language:  "en",
			messageID: "players_connected",
			data:      map[string]interface{}{"Count": 3},
			expected:  "3 players connected",
		},
		{
			name:      "Portuguese message with player count",
			language:  "pt",
			messageID: "players_connected",
			data:      map[string]interface{}{"Count": 2},
			expected:  "2 jogadores conectados",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := GetMessageWithData(tc.language, tc.messageID, tc.data)

			// Assert
			if err != nil {
				t.Errorf("GetMessageWithData() error = %v, want nil", err)
			}
			if result != tc.expected {
				t.Errorf("GetMessageWithData() = %v, want %v", result, tc.expected)
			}
		})
	}
}

// TDD: Test new Bundle interface with dependency injection
func TestBundleInterface(t *testing.T) {
	// Test creating bundle with custom default language
	bundle, err := NewBundleWithDefaults("internal/i18n/locales", language.Portuguese)
	if err != nil {
		t.Errorf("NewBundleWithDefaults() error = %v, want nil", err)
	}

	// Test default language
	defaultLang := bundle.GetDefaultLanguage()
	if defaultLang != "pt" {
		t.Errorf("GetDefaultLanguage() = %v, want 'pt'", defaultLang)
	}

	// Test setting new default language
	err = bundle.SetDefaultLanguage("en")
	if err != nil {
		t.Errorf("SetDefaultLanguage() error = %v, want nil", err)
	}

	newDefault := bundle.GetDefaultLanguage()
	if newDefault != "en" {
		t.Errorf("GetDefaultLanguage() after set = %v, want 'en'", newDefault)
	}
}

// TDD: Test path validation and security
func TestPathValidation(t *testing.T) {
	testCases := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty path",
			path:        "",
			expectError: true,
			errorMsg:    "locales path cannot be empty",
		},
		{
			name:        "Directory traversal attack",
			path:        "../../../etc/passwd",
			expectError: true,
			errorMsg:    "invalid locales path: directory traversal not allowed",
		},
		{
			name:        "Non-existent directory",
			path:        "non/existent/path",
			expectError: true,
		},
		{
			name:        "Valid path",
			path:        "internal/i18n/locales",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			bundle, err := NewBundle(tc.path)

			// Assert
			if tc.expectError {
				if err == nil {
					t.Errorf("NewBundle(%s) should return error", tc.path)
				}
				if tc.errorMsg != "" && !strings.Contains(err.Error(), tc.errorMsg) {
					t.Errorf("NewBundle(%s) error = %v, want to contain %v", tc.path, err.Error(), tc.errorMsg)
				}
				if bundle != nil {
					t.Errorf("NewBundle(%s) should return nil bundle on error", tc.path)
				}
			} else {
				if err != nil {
					t.Errorf("NewBundle(%s) error = %v, want nil", tc.path, err)
				}
				if bundle == nil {
					t.Errorf("NewBundle(%s) should return non-nil bundle", tc.path)
				}
			}
		})
	}
}
