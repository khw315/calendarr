package models

import (
	"encoding/json"
	"testing"
)

func TestPartialConfigUnmarshal(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Language = "ID"
	cfg.Timezone = "Asia/Jakarta"
	cfg.UseDiscord = true
	cfg.DiscordWebhookURL = "https://discord.com/api/webhooks/test"
	cfg.CalendarURLs = []CalendarUrl{
		{URL: "https://sonarr.test.com/ics", Type: "tv"},
	}

	// Partial update JSON modifying only DISCORD_HIDE_MENTION_INSTRUCTIONS
	partialJSON := `{"DISCORD_HIDE_MENTION_INSTRUCTIONS": true}`

	if err := json.Unmarshal([]byte(partialJSON), cfg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Verify target field was updated
	if !cfg.DiscordHideMentionInstructions {
		t.Errorf("Expected DiscordHideMentionInstructions to be true")
	}

	// Verify un-modified fields retained their values
	if cfg.Language != "ID" {
		t.Errorf("Expected Language 'ID', got %q", cfg.Language)
	}
	if cfg.Timezone != "Asia/Jakarta" {
		t.Errorf("Expected Timezone 'Asia/Jakarta', got %q", cfg.Timezone)
	}
	if !cfg.UseDiscord {
		t.Errorf("Expected UseDiscord to remain true")
	}
	if cfg.DiscordWebhookURL != "https://discord.com/api/webhooks/test" {
		t.Errorf("Expected DiscordWebhookURL to remain preserved, got %q", cfg.DiscordWebhookURL)
	}
	if len(cfg.CalendarURLs) != 1 || cfg.CalendarURLs[0].URL != "https://sonarr.test.com/ics" {
		t.Errorf("Expected CalendarURLs to be preserved, got %v", cfg.CalendarURLs)
	}
}

func TestConfigMarshalAndUnmarshalAll(t *testing.T) {
	cfg := DefaultConfig()

	data, err := json.Marshal(cfg)
	if err != nil || len(data) == 0 {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	fullJSON := `{
		"APP_LANGUAGE": "FR",
		"USE_DISCORD": true,
		"DISCORD_WEBHOOK_URL": "https://discord.com/webhook",
		"DISCORD_MENTION_ROLE_ID": "12345",
		"DISCORD_HIDE_MENTION_INSTRUCTIONS": true,
		"DISCORD_TIMESTAMP_STYLE": "R",
		"ENABLE_CUSTOM_DISCORD_FOOTER": true,
		"USE_SLACK": true,
		"SLACK_WEBHOOK_URL": "https://slack.com/webhook",
		"ENABLE_CUSTOM_SLACK_FOOTER": true,
		"PASSED_EVENT_HANDLING": "strike",
		"DEDUPLICATE_EVENTS": true,
		"SHOW_DATE_RANGE": true,
		"SHOW_TIMEZONE_IN_SUBHEADER": true,
		"TZ": "Europe/Paris",
		"USE_24_HOUR": true,
		"ADD_LEADING_ZERO": true,
		"DISPLAY_TIME": true,
		"SCHEDULE_TYPE": "cron",
		"SCHEDULE_DAY": "monday",
		"RUN_TIME": "08:00",
		"CRON_SCHEDULE": "0 8 * * *",
		"RUN_ON_STARTUP": true,
		"DEBUG": true,
		"HTTP_TIMEOUT": 30,
		"LOG_MAX_SIZE_MB": 10,
		"LOG_BACKUP_COUNT": 5,
		"CALENDAR_URLS": [
			{"url": "https://test.com/cal.ics", "type": "tv"}
		]
	}`

	var updated Config
	if err := json.Unmarshal([]byte(fullJSON), &updated); err != nil {
		t.Fatalf("Failed to unmarshal full JSON: %v", err)
	}

	if updated.Language != "FR" {
		t.Errorf("Expected FR, got %s", updated.Language)
	}
	if updated.Timezone != "Europe/Paris" {
		t.Errorf("Expected Europe/Paris, got %s", updated.Timezone)
	}
	if updated.HTTPTimeout != 30 {
		t.Errorf("Expected timeout 30, got %d", updated.HTTPTimeout)
	}
	if len(updated.CalendarURLs) != 1 {
		t.Errorf("Expected 1 calendar URL, got %d", len(updated.CalendarURLs))
	}
}
