package localization

import (
	"testing"
)

func TestLocalization(t *testing.T) {
	langs := SupportedLanguages()
	if len(langs) == 0 {
		t.Errorf("Expected supported languages to be loaded")
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
