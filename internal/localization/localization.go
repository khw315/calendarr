package localization

import (
	"embed"
	"encoding/json"
	"log"
	"math/rand"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed locales/*.json
var localesFS embed.FS

const DefaultLanguage = "EN"

var (
	translations = make(map[string]map[string]interface{})
	once         sync.Once
)

func initTranslations() {
	entries, err := localesFS.ReadDir("locales")
	if err != nil {
		log.Printf("⚠️ Failed to read embedded locales: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" || entry.Name() == "template.json" {
			continue
		}

		langCode := strings.ToUpper(strings.TrimSuffix(entry.Name(), ".json"))
		data, err := localesFS.ReadFile("locales/" + entry.Name())
		if err != nil {
			log.Printf("⚠️ Failed to read locale file %s: %v", entry.Name(), err)
			continue
		}

		var localeMap map[string]interface{}
		if err := json.Unmarshal(data, &localeMap); err != nil {
			log.Printf("⚠️ Failed to parse locale file %s: %v", entry.Name(), err)
			continue
		}

		translations[langCode] = localeMap
	}
	log.Printf("✅ Loaded %d localized languages (%v)", len(translations), SupportedLanguages())
}

func SupportedLanguages() []string {
	once.Do(initTranslations)
	var list []string
	for k := range translations {
		list = append(list, k)
	}
	return list
}

func NormalizeLanguage(lang string) string {
	once.Do(initTranslations)
	normalized := strings.ToUpper(strings.TrimSpace(lang))
	if _, ok := translations[normalized]; ok {
		return normalized
	}
	return DefaultLanguage
}

func GetText(lang, key string) string {
	once.Do(initTranslations)
	lang = NormalizeLanguage(lang)
	if data, ok := translations[lang]; ok {
		if val, exists := data[key]; exists {
			if strVal, isStr := val.(string); isStr {
				return strVal
			}
		}
	}
	// Fallback to EN
	if data, ok := translations[DefaultLanguage]; ok {
		if val, exists := data[key]; exists {
			if strVal, isStr := val.(string); isStr {
				return strVal
			}
		}
	}
	return ""
}

func GetRandomMessage(lang, messageKey string) string {
	once.Do(initTranslations)
	lang = NormalizeLanguage(lang)
	key := messageKey + "_messages"

	var messages []string
	if data, ok := translations[lang]; ok {
		if rawMsgs, exists := data[key]; exists {
			if sliceMsgs, isSlice := rawMsgs.([]interface{}); isSlice {
				for _, m := range sliceMsgs {
					if s, ok := m.(string); ok {
						messages = append(messages, s)
					}
				}
			}
		}
	}

	if len(messages) == 0 && lang != DefaultLanguage {
		if data, ok := translations[DefaultLanguage]; ok {
			if rawMsgs, exists := data[key]; exists {
				if sliceMsgs, isSlice := rawMsgs.([]interface{}); isSlice {
					for _, m := range sliceMsgs {
						if s, ok := m.(string); ok {
							messages = append(messages, s)
						}
					}
				}
			}
		}
	}

	if len(messages) == 0 {
		switch messageKey {
		case "no_new_releases":
			return "No new releases to share."
		case "no_day_content":
			return "No releases scheduled for this day."
		default:
			return ""
		}
	}

	return messages[rand.Intn(len(messages))]
}
