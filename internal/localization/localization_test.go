package localization

import (
	"testing"
	"time"
)

func TestSupportedLanguagesAndNormalize(t *testing.T) {
	langs := SupportedLanguages()
	if len(langs) == 0 {
		t.Errorf("Expected supported languages to be loaded")
	}

	langList := GetLanguageList()
	if len(langList) == 0 {
		t.Errorf("Expected GetLanguageList to return items")
	}
	foundEN := false
	for _, l := range langList {
		if l.Code == "EN" {
			foundEN = true
			if l.Name == "" {
				t.Errorf("Expected EN language name to be non-empty")
			}
		}
	}
	if !foundEN {
		t.Errorf("Expected EN in GetLanguageList")
	}

	if NormalizeLanguage("en") != "EN" {
		t.Errorf("Expected EN")
	}

	if NormalizeLanguage("INVALID_LANG") != DefaultLanguage {
		t.Errorf("Expected default language %s for invalid lang", DefaultLanguage)
	}
}

func TestGetTextAndMessages(t *testing.T) {
	if GetText("EN", "header_text") == "" {
		t.Errorf("Expected header_text in EN to be non-empty")
	}

	if GetText("INVALID_LANG", "header_text") == "" {
		t.Errorf("Expected fallback header_text to be non-empty")
	}

	if GetRandomMessage("EN", "no_new_releases") == "" {
		t.Errorf("Expected random message to be non-empty")
	}

	if GetRandomMessage("ID", "no_new_releases") == "" {
		t.Errorf("Expected random message in ID to be non-empty")
	}

	if fallbackMessage("no_new_releases") != "No new releases to share." {
		t.Errorf("Unexpected fallback message")
	}

	if fallbackMessage("no_day_content") != "No releases scheduled for this day." {
		t.Errorf("Unexpected fallback message")
	}

	if fallbackMessage("unknown_key") != "" {
		t.Errorf("Expected empty fallback for unknown key")
	}

	if getRandomItem([]string{"item1"}) != "item1" {
		t.Errorf("Expected item1")
	}

	if getRandomItem([]string{}) != "" {
		t.Errorf("Expected empty string for empty slice")
	}
}

func TestFormatDateHeader(t *testing.T) {
	testTime := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)

	if FormatDateHeader(testTime, "EN") == "" {
		t.Errorf("Expected non-empty date header for EN")
	}

	if FormatDateHeader(testTime, "ID") == "" {
		t.Errorf("Expected non-empty date header for ID")
	}

	if FormatDateHeader(testTime, "FR") == "" {
		t.Errorf("Expected non-empty date header for FR")
	}

	if FormatDateHeader(testTime, "INVALID") == "" {
		t.Errorf("Expected non-empty date header for fallback language")
	}
}

func TestFormatSubheader(t *testing.T) {
	// 1. Both TV and Movie (plural)
	if FormatSubheader("EN", 2, 3) == "" {
		t.Errorf("Expected non-empty subheader")
	}

	// 2. Singular TV & Movie
	if FormatSubheader("EN", 1, 1) == "" {
		t.Errorf("Expected non-empty subheader for singular items")
	}

	// 3. TV only
	if FormatSubheader("EN", 5, 0) == "" {
		t.Errorf("Expected non-empty subheader for TV only")
	}

	// 4. Movie only
	if FormatSubheader("EN", 0, 4) == "" {
		t.Errorf("Expected non-empty subheader for Movie only")
	}

	// 5. Zero items -> fallback random message
	if FormatSubheader("EN", 0, 0) == "" {
		t.Errorf("Expected fallback message when count is 0")
	}
}

func TestHelpers(t *testing.T) {
	// extractString
	if extractString(nil, "key") != "" {
		t.Errorf("Expected empty string for nil map")
	}
	if extractString(map[string]interface{}{"num": 123}, "num") != "" {
		t.Errorf("Expected empty string for non-string val")
	}

	// extractMessages
	if extractMessages(nil, "key") != nil {
		t.Errorf("Expected nil for nil map")
	}
	if extractMessages(map[string]interface{}{"key_messages": "not a slice"}, "key") != nil {
		t.Errorf("Expected nil for non-slice value")
	}

	// getSubheaderLabel
	if getSubheaderLabel(nil, "tv", 1) != "tv" {
		t.Errorf("Expected 'tv' fallback for nil map")
	}
	if getSubheaderLabel(map[string]interface{}{"subheader_labels": "invalid"}, "tv", 1) != "tv" {
		t.Errorf("Expected 'tv' fallback for invalid subheader_labels")
	}

	// extractNestedString
	if extractNestedString(nil, "section", "key") != "" {
		t.Errorf("Expected empty string for nil map")
	}
	if extractNestedString(map[string]interface{}{"section": "invalid"}, "section", "key") != "" {
		t.Errorf("Expected empty string for non-map section")
	}
}
