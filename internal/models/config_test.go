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
