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
	enSub := FormatSubheader("EN", 2, 3, 0)
	if enSub != "2 all-new episodes and 3 movie releases" {
		t.Errorf("Unexpected subheader for EN: %s", enSub)
	}

	// 2. Singular TV & Movie with Premiere (3 items list_join)
	enAll := FormatSubheader("EN", 1, 1, 1)
	if enAll != "1 all-new episode, 1 movie release, and 1 season premiere" {
		t.Errorf("Unexpected subheader for EN with 3 items: %s", enAll)
	}

	// 3. TV and Premiere only (matches the user's scenario)
	enTVPremiere := FormatSubheader("EN", 3, 0, 1)
	if enTVPremiere != "3 all-new episodes and 1 season premiere" {
		t.Errorf("Unexpected subheader for EN TV+Premiere: %s", enTVPremiere)
	}

	// 4. Korean localization
	koSub := FormatSubheader("KO", 3, 0, 1)
	if koSub != "3개의 최신 에피소드 및 1개의 시즌 프리미어" {
		t.Errorf("Unexpected subheader for KO TV+Premiere: %s", koSub)
	}

	// 5. Indonesian localization
	idSub := FormatSubheader("ID", 3, 0, 1)
	if idSub != "3 episode terbaru dan 1 perdana musim" {
		t.Errorf("Unexpected subheader for ID TV+Premiere: %s", idSub)
	}

	// 6. Japanese localization
	jaSub := FormatSubheader("JA", 3, 0, 1)
	if jaSub != "3件の新着エピソードと1件のシーズン初回配信" {
		t.Errorf("Unexpected subheader for JA TV+Premiere: %s", jaSub)
	}

	// 7. TV only
	if FormatSubheader("EN", 5, 0, 0) != "5 all-new episodes" {
		t.Errorf("Expected '5 all-new episodes', got: %s", FormatSubheader("EN", 5, 0, 0))
	}

	// 8. Movie only
	if FormatSubheader("EN", 0, 4, 0) != "4 movie releases" {
		t.Errorf("Expected '4 movie releases', got: %s", FormatSubheader("EN", 0, 4, 0))
	}

	// 9. Zero items -> fallback random message
	if FormatSubheader("EN", 0, 0, 0) == "" {
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
