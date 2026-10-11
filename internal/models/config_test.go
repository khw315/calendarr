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

func TestConfigWebhookChannels(t *testing.T) {
	t.Run("Discord distinct channels configured", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.DiscordWebhookURL = "https://discord.com/default"
		cfg.DiscordTVWebhookURL = "https://discord.com/tv"
		cfg.DiscordMovieWebhookURL = "https://discord.com/movie"

		if !cfg.HasDistinctDiscordChannels() {
			t.Errorf("Expected HasDistinctDiscordChannels to be true")
		}
		if got := cfg.GetDiscordTVWebhookURL(); got != "https://discord.com/tv" {
			t.Errorf("Expected TV webhook URL, got %q", got)
		}
		if got := cfg.GetDiscordMovieWebhookURL(); got != "https://discord.com/movie" {
			t.Errorf("Expected Movie webhook URL, got %q", got)
		}
	})

	t.Run("Discord fallback to default webhook", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.DiscordWebhookURL = "https://discord.com/default"

		if cfg.HasDistinctDiscordChannels() {
			t.Errorf("Expected HasDistinctDiscordChannels to be false")
		}
		if got := cfg.GetDiscordTVWebhookURL(); got != "https://discord.com/default" {
			t.Errorf("Expected fallback to default webhook URL, got %q", got)
		}
		if got := cfg.GetDiscordMovieWebhookURL(); got != "https://discord.com/default" {
			t.Errorf("Expected fallback to default webhook URL, got %q", got)
		}
	})

	t.Run("Slack distinct channels configured", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.SlackWebhookURL = "https://slack.com/default"
		cfg.SlackTVWebhookURL = "https://slack.com/tv"
		cfg.SlackMovieWebhookURL = "https://slack.com/movie"

		if !cfg.HasDistinctSlackChannels() {
			t.Errorf("Expected HasDistinctSlackChannels to be true")
		}
		if got := cfg.GetSlackTVWebhookURL(); got != "https://slack.com/tv" {
			t.Errorf("Expected TV webhook URL, got %q", got)
		}
		if got := cfg.GetSlackMovieWebhookURL(); got != "https://slack.com/movie" {
			t.Errorf("Expected Movie webhook URL, got %q", got)
		}
	})

	t.Run("Slack fallback to default webhook", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.SlackWebhookURL = "https://slack.com/default"

		if cfg.HasDistinctSlackChannels() {
			t.Errorf("Expected HasDistinctSlackChannels to be false")
		}
		if got := cfg.GetSlackTVWebhookURL(); got != "https://slack.com/default" {
			t.Errorf("Expected fallback to default webhook URL, got %q", got)
		}
		if got := cfg.GetSlackMovieWebhookURL(); got != "https://slack.com/default" {
			t.Errorf("Expected fallback to default webhook URL, got %q", got)
		}
	})

	t.Run("JSON unmarshal with distinct webhook URLs", func(t *testing.T) {
		jsonStr := `{
			"DISCORD_TV_WEBHOOK_URL": "https://discord.com/custom_tv",
			"DISCORD_MOVIE_WEBHOOK_URL": "https://discord.com/custom_movie",
			"SLACK_TV_WEBHOOK_URL": "https://slack.com/custom_tv",
			"SLACK_MOVIE_WEBHOOK_URL": "https://slack.com/custom_movie",
			"USE_NTFY": true,
			"NTFY_SERVER_URL": "https://ntfy.example.com",
			"NTFY_TOPIC": "calendarr-default",
			"NTFY_TV_TOPIC": "calendarr-tv",
			"NTFY_MOVIE_TOPIC": "calendarr-movies",
			"NTFY_ACCESS_TOKEN": "tk_12345",
			"NTFY_PRIORITY": 4,
			"USE_TELEGRAM": true,
			"TELEGRAM_BOT_TOKEN": "123456:ABC-DEF",
			"TELEGRAM_CHAT_ID": "-1001234567890",
			"TELEGRAM_TV_CHAT_ID": "-1001234567891",
			"TELEGRAM_MOVIE_CHAT_ID": "-1001234567892",
			"TELEGRAM_SILENT_NOTIFICATION": true
		}`

		cfg := DefaultConfig()
		if err := json.Unmarshal([]byte(jsonStr), cfg); err != nil {
			t.Fatalf("Failed to unmarshal distinct URLs: %v", err)
		}

		if cfg.DiscordTVWebhookURL != "https://discord.com/custom_tv" {
			t.Errorf("Unexpected DiscordTVWebhookURL: %q", cfg.DiscordTVWebhookURL)
		}
		if cfg.DiscordMovieWebhookURL != "https://discord.com/custom_movie" {
			t.Errorf("Unexpected DiscordMovieWebhookURL: %q", cfg.DiscordMovieWebhookURL)
		}
		if cfg.SlackTVWebhookURL != "https://slack.com/custom_tv" {
			t.Errorf("Unexpected SlackTVWebhookURL: %q", cfg.SlackTVWebhookURL)
		}
		if cfg.SlackMovieWebhookURL != "https://slack.com/custom_movie" {
			t.Errorf("Unexpected SlackMovieWebhookURL: %q", cfg.SlackMovieWebhookURL)
		}
		if !cfg.UseNtfy || cfg.NtfyServerURL != "https://ntfy.example.com" || cfg.NtfyTopic != "calendarr-default" || cfg.NtfyTVTopic != "calendarr-tv" || cfg.NtfyMovieTopic != "calendarr-movies" || cfg.NtfyAccessToken != "tk_12345" || cfg.NtfyPriority != 4 {
			t.Errorf("Unexpected NTFY config values: %+v", cfg)
		}
		if !cfg.UseTelegram || cfg.TelegramBotToken != "123456:ABC-DEF" || cfg.TelegramChatID != "-1001234567890" || cfg.TelegramTVChatID != "-1001234567891" || cfg.TelegramMovieChatID != "-1001234567892" || !cfg.TelegramSilentNotification {
			t.Errorf("Unexpected Telegram config values: %+v", cfg)
		}

		data, err := json.Marshal(cfg)
		if err != nil {
			t.Fatalf("Failed to marshal config: %v", err)
		}
		if len(data) == 0 {
			t.Errorf("Expected non-empty JSON data")
		}
	})

	t.Run("NTFY distinct topics and fallback", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.NtfyTopic = "default-topic"

		if cfg.HasDistinctNtfyTopics() {
			t.Errorf("Expected HasDistinctNtfyTopics to be false initially")
		}
		if got := cfg.GetNtfyTVTopic(); got != "default-topic" {
			t.Errorf("Expected fallback TV topic, got %q", got)
		}
		if got := cfg.GetNtfyMovieTopic(); got != "default-topic" {
			t.Errorf("Expected fallback Movie topic, got %q", got)
		}

		cfg.NtfyTVTopic = "tv-topic"
		cfg.NtfyMovieTopic = "movie-topic"
		if !cfg.HasDistinctNtfyTopics() {
			t.Errorf("Expected HasDistinctNtfyTopics to be true")
		}
		if got := cfg.GetNtfyTVTopic(); got != "tv-topic" {
			t.Errorf("Expected TV topic, got %q", got)
		}
		if got := cfg.GetNtfyMovieTopic(); got != "movie-topic" {
			t.Errorf("Expected Movie topic, got %q", got)
		}
	})

	t.Run("Telegram distinct channels and fallback", func(t *testing.T) {
		cfg := DefaultConfig()
		cfg.TelegramChatID = "-1000"

		if cfg.HasDistinctTelegramChannels() {
			t.Errorf("Expected HasDistinctTelegramChannels to be false initially")
		}
		if got := cfg.GetTelegramTVChatID(); got != "-1000" {
			t.Errorf("Expected fallback TV chat ID, got %q", got)
		}
		if got := cfg.GetTelegramMovieChatID(); got != "-1000" {
			t.Errorf("Expected fallback Movie chat ID, got %q", got)
		}

		cfg.TelegramTVChatID = "-2000"
		cfg.TelegramMovieChatID = "-3000"
		if !cfg.HasDistinctTelegramChannels() {
			t.Errorf("Expected HasDistinctTelegramChannels to be true")
		}
		if got := cfg.GetTelegramTVChatID(); got != "-2000" {
			t.Errorf("Expected TV chat ID, got %q", got)
		}
		if got := cfg.GetTelegramMovieChatID(); got != "-3000" {
			t.Errorf("Expected Movie chat ID, got %q", got)
		}
	})
}
