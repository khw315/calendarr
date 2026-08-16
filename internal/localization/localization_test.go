package localization

import (
	"testing"
	"time"
)

func TestLocalization(t *testing.T) {
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

	normEN := NormalizeLanguage("en")
	if normEN != "EN" {
		t.Errorf("Expected EN, got %s", normEN)
	}

	normInvalid := NormalizeLanguage("INVALID_LANG")
	if normInvalid != DefaultLanguage {
		t.Errorf("Expected default language %s for invalid lang, got %s", DefaultLanguage, normInvalid)
	}

	headerEN := GetText("EN", "header_text")
	if headerEN == "" {
		t.Errorf("Expected header_text in EN to be non-empty")
	}

	headerFallback := GetText("INVALID_LANG", "header_text")
	if headerFallback == "" {
		t.Errorf("Expected fallback header_text to be non-empty")
	}

	msgEN := GetRandomMessage("EN", "no_new_releases")
	if msgEN == "" {
		t.Errorf("Expected random message to be non-empty")
	}

	msgFallback := GetRandomMessage("ID", "no_new_releases")
	if msgFallback == "" {
		t.Errorf("Expected random message in ID to be non-empty")
	}

	fbMsg1 := fallbackMessage("no_new_releases")
	if fbMsg1 != "No new releases to share." {
		t.Errorf("Unexpected fallback message: %s", fbMsg1)
	}

	fbMsg2 := fallbackMessage("no_day_content")
	if fbMsg2 != "No releases scheduled for this day." {
		t.Errorf("Unexpected fallback message: %s", fbMsg2)
	}

	fbMsg3 := fallbackMessage("unknown_key")
	if fbMsg3 != "" {
		t.Errorf("Expected empty fallback for unknown key, got %s", fbMsg3)
	}

	item := getRandomItem([]string{"item1"})
	if item != "item1" {
		t.Errorf("Expected item1, got %s", item)
	}

	emptyItem := getRandomItem([]string{})
	if emptyItem != "" {
		t.Errorf("Expected empty string for empty slice, got %s", emptyItem)
	}
}

func TestFormatDateHeader(t *testing.T) {
	testTime := time.Date(2026, 8, 16, 12, 0, 0, 0, time.UTC)

	headerEN := FormatDateHeader(testTime, "EN")
	if headerEN == "" {
		t.Errorf("Expected non-empty date header for EN")
	}

	headerID := FormatDateHeader(testTime, "ID")
	if headerID == "" {
		t.Errorf("Expected non-empty date header for ID")
	}

	headerFR := FormatDateHeader(testTime, "FR")
	if headerFR == "" {
		t.Errorf("Expected non-empty date header for FR")
	}

	headerInvalid := FormatDateHeader(testTime, "INVALID")
	if headerInvalid == "" {
		t.Errorf("Expected non-empty date header for fallback language")
	}
}

func TestFormatSubheader(t *testing.T) {
	// 1. Both TV and Movie (plural)
	subBoth := FormatSubheader("EN", 2, 3)
	if subBoth == "" {
		t.Errorf("Expected non-empty subheader")
	}

	// 2. Singular TV & Movie
	subSingular := FormatSubheader("EN", 1, 1)
	if subSingular == "" {
		t.Errorf("Expected non-empty subheader for singular items")
	}

	// 3. TV only
	subTV := FormatSubheader("EN", 5, 0)
	if subTV == "" {
		t.Errorf("Expected non-empty subheader for TV only")
	}

	// 4. Movie only
	subMovie := FormatSubheader("EN", 0, 4)
	if subMovie == "" {
		t.Errorf("Expected non-empty subheader for Movie only")
	}

	// 5. Zero items -> fallback random message
	subZero := FormatSubheader("EN", 0, 0)
	if subZero == "" {
		t.Errorf("Expected fallback message when count is 0")
	}
}

func TestHelpers(t *testing.T) {
	// extractString
	if str := extractString(nil, "key"); str != "" {
		t.Errorf("Expected empty string for nil map")
	}
	if str := extractString(map[string]interface{}{"num": 123}, "num"); str != "" {
		t.Errorf("Expected empty string for non-string val")
	}

	// extractMessages
	if msgs := extractMessages(nil, "key"); msgs != nil {
		t.Errorf("Expected nil for nil map")
	}
	if msgs := extractMessages(map[string]interface{}{"key_messages": "not a slice"}, "key"); msgs != nil {
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
	if val := extractNestedString(nil, "section", "key"); val != "" {
		t.Errorf("Expected empty string for nil map")
	}
	if val := extractNestedString(map[string]interface{}{"section": "invalid"}, "section", "key"); val != "" {
		t.Errorf("Expected empty string for non-map section")
	}
}
