package localization

import (
	"crypto/rand"
	"embed"
	"encoding/binary"
	"encoding/json"
	"log"
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

	var keys []string
	for k := range translations {
		keys = append(keys, k)
	}
	log.Printf("✅ Loaded %d localized languages (%v)", len(translations), keys)
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
	if strVal := extractString(translations[lang], key); strVal != "" {
		return strVal
	}
	return extractString(translations[DefaultLanguage], key)
}

func extractString(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}
	if val, exists := data[key]; exists {
		if strVal, isStr := val.(string); isStr {
			return strVal
		}
	}
	return ""
}

func GetRandomMessage(lang, messageKey string) string {
	once.Do(initTranslations)
	lang = NormalizeLanguage(lang)

	messages := extractMessages(translations[lang], messageKey)
	if len(messages) == 0 && lang != DefaultLanguage {
		messages = extractMessages(translations[DefaultLanguage], messageKey)
	}

	if len(messages) == 0 {
		return fallbackMessage(messageKey)
	}

	return getRandomItem(messages)
}

func extractMessages(data map[string]interface{}, messageKey string) []string {
	if data == nil {
		return nil
	}
	key := messageKey + "_messages"
	rawMsgs, exists := data[key]
	if !exists {
		return nil
	}

	sliceMsgs, isSlice := rawMsgs.([]interface{})
	if !isSlice {
		return nil
	}

	var messages []string
	for _, m := range sliceMsgs {
		if s, ok := m.(string); ok {
			messages = append(messages, s)
		}
	}
	return messages
}

func fallbackMessage(messageKey string) string {
	switch messageKey {
	case "no_new_releases":
		return "No new releases to share."
	case "no_day_content":
		return "No releases scheduled for this day."
	default:
		return ""
	}
}

func getRandomItem(items []string) string {
	if len(items) == 0 {
		return ""
	}
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return items[0]
	}
	val := binary.LittleEndian.Uint64(b[:])
	idx := int(val % uint64(len(items)))
	return items[idx]
}
