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

	if norm := NormalizeLanguage("en"); norm != "EN" {
		t.Errorf("Expected EN, got %s", norm)
	}

	if norm := NormalizeLanguage("INVALID_LANG"); norm != DefaultLanguage {
		t.Errorf("Expected default language %s for invalid lang, got %s", DefaultLanguage, norm)
	}
}

func TestGetTextAndMessages(t *testing.T) {
	if header := GetText("EN", "header_text"); header == "" {
		t.Errorf("Expected header_text in EN to be non-empty")
	}

	if header := GetText("INVALID_LANG", "header_text"); header == "" {
		t.Errorf("Expected fallback header_text to be non-empty")
	}

	if msg := GetRandomMessage("EN", "no_new_releases"); msg == "" {
		t.Errorf("Expected random message to be non-empty")
	}

	if msg := GetRandomMessage("ID", "no_new_releases"); msg == "" {
		t.Errorf("Expected random message in ID to be non-empty")
	}

	if fb := fallbackMessage("no_new_releases"); fb != "No new releases to share." {
		t.Errorf("Unexpected fallback message: %s", fb)
	}

	if fb := fallbackMessage("no_day_content"); fb != "No releases scheduled for this day." {
		t.Errorf("Unexpected fallback message: %s", fb)
	}

	if fb := fallbackMessage("unknown_key"); fb != "" {
		t.Errorf("Expected empty fallback for unknown key, got %s", fb)
	}

	if item := getRandomItem([]string{"item1"}); item != "item1" {
		t.Errorf("Expected item1, got %s", item)
	}

	if emptyItem := getRandomItem([]string{}); emptyItem != "" {
		t.Errorf("Expected empty string for empty slice, got %s", emptyItem)
	}
}

func TestFormatDateHeader(t *testing.T) {
	testTime := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)

	if header := FormatDateHeader(testTime, "EN"); header == "" {
		t.Errorf("Expected non-empty date header for EN")
	}

	if header := FormatDateHeader(testTime, "ID"); header == "" {
		t.Errorf("Expected non-empty date header for ID")
	}

	if header := FormatDateHeader(testTime, "FR"); header == "" {
		t.Errorf("Expected non-empty date header for FR")
	}

	if header := FormatDateHeader(testTime, "INVALID"); header == "" {
		t.Errorf("Expected non-empty date header for fallback language")
	}
}

func TestFormatSubheader(t *testing.T) {
	// 1. Both TV and Movie (plural)
	if sub := FormatSubheader("EN", 2, 3); sub == "" {
		t.Errorf("Expected non-empty subheader")
	}

	// 2. Singular TV & Movie
	if sub := FormatSubheader("EN", 1, 1); sub == "" {
		t.Errorf("Expected non-empty subheader for singular items")
	}

	// 3. TV only
	if sub := FormatSubheader("EN", 5, 0); sub == "" {
		t.Errorf("Expected non-empty subheader for TV only")
	}

	// 4. Movie only
	if sub := FormatSubheader("EN", 0, 4); sub == "" {
		t.Errorf("Expected non-empty subheader for Movie only")
	}

	// 5. Zero items -> fallback random message
	if sub := FormatSubheader("EN", 0, 0); sub == "" {
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
	if lbl := getSubheaderLabel(nil, "tv", 1); lbl != "tv" {
		t.Errorf("Expected 'tv' fallback for nil map, got %s", lbl)
	}
	if lbl := getSubheaderLabel(map[string]interface{}{"subheader_labels": "invalid"}, "tv", 1); lbl != "tv" {
		t.Errorf("Expected 'tv' fallback for invalid subheader_labels, got %s", lbl)
	}

	// extractNestedString
	if extractNestedString(nil, "section", "key") != "" {
		t.Errorf("Expected empty string for nil map")
	}
	if extractNestedString(map[string]interface{}{"section": "invalid"}, "section", "key") != "" {
		t.Errorf("Expected empty string for non-map section")
	}
}
