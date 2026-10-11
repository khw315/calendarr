package platform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/khw315/calendarr/internal/models"
)

func TestPlatformServiceDispatchSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := &Service{telegramBaseURL: ts.URL}
	cfg := models.DefaultConfig()
	cfg.UseDiscord = true
	cfg.DiscordWebhookURL = ts.URL
	cfg.UseSlack = true
	cfg.SlackWebhookURL = ts.URL
	cfg.UseNtfy = true
	cfg.NtfyServerURL = ts.URL
	cfg.NtfyTopic = "test-topic"
	cfg.UseTelegram = true
	cfg.TelegramBotToken = "test-bot-token"
	cfg.TelegramChatID = "987654321"

	payloads := models.NotificationPayloads{
		DiscordCombined:  &models.DiscordPayload{Content: "Test Discord"},
		SlackCombined:    &models.SlackPayload{Text: "Test Slack"},
		NtfyCombined:     &models.NtfyPayload{Title: "Test", Message: "Test Ntfy", HasReleases: true},
		TelegramCombined: &models.TelegramPayload{Text: "Test Telegram", HasReleases: true},
	}

	err := svc.DispatchPayloads(context.Background(), cfg, payloads)
	if err != nil {
		t.Fatalf("DispatchPayloads failed: %v", err)
	}
}

func TestPlatformServiceDispatchError(t *testing.T) {
	tsErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Error"))
	}))
	defer tsErr.Close()

	svc := &Service{telegramBaseURL: tsErr.URL}
	cfg := models.DefaultConfig()
	cfg.UseDiscord = true
	cfg.DiscordWebhookURL = tsErr.URL
	cfg.UseSlack = true
	cfg.SlackWebhookURL = tsErr.URL
	cfg.UseNtfy = true
	cfg.NtfyServerURL = tsErr.URL
	cfg.NtfyTopic = "test-topic"
	cfg.UseTelegram = true
	cfg.TelegramBotToken = "test-bot-token"
	cfg.TelegramChatID = "987654321"

	payloads := models.NotificationPayloads{
		DiscordCombined:  &models.DiscordPayload{Content: "Test Discord"},
		SlackCombined:    &models.SlackPayload{Text: "Test Slack"},
		NtfyCombined:     &models.NtfyPayload{Title: "Test", Message: "Test Ntfy", HasReleases: true},
		TelegramCombined: &models.TelegramPayload{Text: "Test Telegram", HasReleases: true},
	}

	err := svc.DispatchPayloads(context.Background(), cfg, payloads)
	if err == nil {
		t.Errorf("Expected dispatch error when webhooks return 500")
	}
}

func TestPlatformServiceSendTest(t *testing.T) {
	tsSuccess := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer tsSuccess.Close()

	tsErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"message": "Invalid webhook"}`))
	}))
	defer tsErr.Close()

	svc := &Service{telegramBaseURL: tsSuccess.URL}
	ctx := context.Background()

	t.Run("Discord Success", func(t *testing.T) {
		err := svc.SendTest(ctx, "discord", tsSuccess.URL)
		if err != nil {
			t.Errorf("Expected successful Discord test send, got: %v", err)
		}
	})

	t.Run("Slack Success", func(t *testing.T) {
		err := svc.SendTest(ctx, "slack", tsSuccess.URL)
		if err != nil {
			t.Errorf("Expected successful Slack test send, got: %v", err)
		}
	})

	t.Run("Ntfy Success with options", func(t *testing.T) {
		err := svc.SendTestWithOptions(ctx, models.TestWebhookRequest{
			Platform:  "ntfy",
			ServerURL: tsSuccess.URL,
			Topic:     "alerts",
		})
		if err != nil {
			t.Errorf("Expected successful Ntfy test send, got: %v", err)
		}
	})

	t.Run("Ntfy Success via WebhookURL fallback", func(t *testing.T) {
		err := svc.SendTestWithOptions(ctx, models.TestWebhookRequest{
			Platform:   "ntfy",
			WebhookURL: tsSuccess.URL + "/my-topic",
		})
		if err != nil {
			t.Errorf("Expected successful Ntfy test send with WebhookURL, got: %v", err)
		}
	})

	t.Run("Telegram Success", func(t *testing.T) {
		err := svc.SendTestWithOptions(ctx, models.TestWebhookRequest{
			Platform: "telegram",
			Token:    "12345:ABCDE",
			ChatID:   "998877",
		})
		if err != nil {
			t.Errorf("Expected successful Telegram test send, got: %v", err)
		}
	})

	t.Run("Ntfy Missing Topic", func(t *testing.T) {
		err := svc.SendTestWithOptions(ctx, models.TestWebhookRequest{
			Platform: "ntfy",
		})
		if err == nil {
			t.Errorf("Expected error for missing Ntfy topic, got nil")
		}
	})

	t.Run("Telegram Missing Token or ChatID", func(t *testing.T) {
		err := svc.SendTestWithOptions(ctx, models.TestWebhookRequest{
			Platform: "telegram",
			Token:    "token-only",
		})
		if err == nil {
			t.Errorf("Expected error for missing Telegram ChatID, got nil")
		}
	})

	t.Run("Unsupported Platform", func(t *testing.T) {
		err := svc.SendTest(ctx, "unsupported_platform", tsSuccess.URL)
		if err == nil {
			t.Errorf("Expected error for unsupported platform, got nil")
		}
	})

	t.Run("Upstream Error", func(t *testing.T) {
		err := svc.SendTest(ctx, "discord", tsErr.URL)
		if err == nil {
			t.Errorf("Expected error on 400 Bad Request response from upstream, got nil")
		}
	})
}

func TestPlatformServiceDispatchPayloads(t *testing.T) {
	var tvHitCount, movieHitCount, combinedHitCount int
	tsTV := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tvHitCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer tsTV.Close()

	tsMovie := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		movieHitCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer tsMovie.Close()

	tsCombined := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		combinedHitCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer tsCombined.Close()

	svc := NewService()
	ctx := context.Background()

	t.Run("Distinct Discord Channels with both TV and Movie payloads", func(t *testing.T) {
		tvHitCount, movieHitCount, combinedHitCount = 0, 0, 0
		cfg := models.DefaultConfig()
		cfg.UseDiscord = true
		cfg.DiscordWebhookURL = tsCombined.URL
		cfg.DiscordTVWebhookURL = tsTV.URL
		cfg.DiscordMovieWebhookURL = tsMovie.URL

		payloads := models.NotificationPayloads{
			DiscordCombined: &models.DiscordPayload{Content: "Combined"},
			DiscordTV:       &models.DiscordPayload{Content: "TV", Embeds: []models.DiscordEmbed{{Title: "TV Show"}}},
			DiscordMovie:    &models.DiscordPayload{Content: "Movie", Embeds: []models.DiscordEmbed{{Title: "Movie Title"}}},
		}

		err := svc.DispatchPayloads(ctx, cfg, payloads)
		if err != nil {
			t.Fatalf("DispatchPayloads failed: %v", err)
		}
		if tvHitCount != 1 {
			t.Errorf("Expected tvHitCount 1, got %d", tvHitCount)
		}
		if movieHitCount != 1 {
			t.Errorf("Expected movieHitCount 1, got %d", movieHitCount)
		}
		if combinedHitCount != 0 {
			t.Errorf("Expected combinedHitCount 0, got %d", combinedHitCount)
		}
	})

	t.Run("Distinct Discord Channels with only TV events", func(t *testing.T) {
		tvHitCount, movieHitCount, combinedHitCount = 0, 0, 0
		cfg := models.DefaultConfig()
		cfg.UseDiscord = true
		cfg.DiscordWebhookURL = tsCombined.URL
		cfg.DiscordTVWebhookURL = tsTV.URL
		cfg.DiscordMovieWebhookURL = tsMovie.URL

		payloads := models.NotificationPayloads{
			DiscordCombined: &models.DiscordPayload{Content: "Combined"},
			DiscordTV:       &models.DiscordPayload{Content: "TV", Embeds: []models.DiscordEmbed{{Title: "TV Show"}}},
			DiscordMovie:    &models.DiscordPayload{Content: "Movie", Embeds: []models.DiscordEmbed{}},
		}

		err := svc.DispatchPayloads(ctx, cfg, payloads)
		if err != nil {
			t.Fatalf("DispatchPayloads failed: %v", err)
		}
		if tvHitCount != 1 {
			t.Errorf("Expected tvHitCount 1, got %d", tvHitCount)
		}
		if movieHitCount != 0 {
			t.Errorf("Expected movieHitCount 0, got %d", movieHitCount)
		}
		if combinedHitCount != 0 {
			t.Errorf("Expected combinedHitCount 0, got %d", combinedHitCount)
		}
	})

	t.Run("Distinct Slack Channels with both TV and Movie payloads", func(t *testing.T) {
		tvHitCount, movieHitCount, combinedHitCount = 0, 0, 0
		cfg := models.DefaultConfig()
		cfg.UseSlack = true
		cfg.SlackWebhookURL = tsCombined.URL
		cfg.SlackTVWebhookURL = tsTV.URL
		cfg.SlackMovieWebhookURL = tsMovie.URL

		payloads := models.NotificationPayloads{
			SlackCombined: &models.SlackPayload{Text: "Combined"},
			SlackTV:       &models.SlackPayload{Text: "TV", Blocks: []models.SlackBlock{{Type: "section"}}},
			SlackMovie:    &models.SlackPayload{Text: "Movie", Blocks: []models.SlackBlock{{Type: "section"}}},
		}

		err := svc.DispatchPayloads(ctx, cfg, payloads)
		if err != nil {
			t.Fatalf("DispatchPayloads failed: %v", err)
		}
		if tvHitCount != 1 {
			t.Errorf("Expected tvHitCount 1, got %d", tvHitCount)
		}
		if movieHitCount != 1 {
			t.Errorf("Expected movieHitCount 1, got %d", movieHitCount)
		}
		if combinedHitCount != 0 {
			t.Errorf("Expected combinedHitCount 0, got %d", combinedHitCount)
		}
	})

	t.Run("Fallback to default webhook when distinct channels are not configured", func(t *testing.T) {
		tvHitCount, movieHitCount, combinedHitCount = 0, 0, 0
		cfg := models.DefaultConfig()
		cfg.UseDiscord = true
		cfg.DiscordWebhookURL = tsCombined.URL
		cfg.UseSlack = true
		cfg.SlackWebhookURL = tsCombined.URL

		payloads := models.NotificationPayloads{
			DiscordCombined: &models.DiscordPayload{Content: "Combined Discord"},
			SlackCombined:   &models.SlackPayload{Text: "Combined Slack"},
		}

		err := svc.DispatchPayloads(ctx, cfg, payloads)
		if err != nil {
			t.Fatalf("DispatchPayloads failed: %v", err)
		}
		if combinedHitCount != 2 {
			t.Errorf("Expected combinedHitCount 2 (1 Discord + 1 Slack), got %d", combinedHitCount)
		}
	})

	t.Run("Segregated dispatch error paths", func(t *testing.T) {
		ts500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts500.Close()

		cfg := models.DefaultConfig()
		cfg.UseDiscord = true
		cfg.DiscordTVWebhookURL = ts500.URL
		cfg.DiscordMovieWebhookURL = ts500.URL
		cfg.UseSlack = true
		cfg.SlackTVWebhookURL = ts500.URL
		cfg.SlackMovieWebhookURL = ts500.URL

		payloads := models.NotificationPayloads{
			DiscordTV:    &models.DiscordPayload{Content: "TV", Embeds: []models.DiscordEmbed{{Title: "TV"}}},
			DiscordMovie: &models.DiscordPayload{Content: "Movie", Embeds: []models.DiscordEmbed{{Title: "Movie"}}},
			SlackTV:      &models.SlackPayload{Text: "TV", Blocks: []models.SlackBlock{{Type: "section"}, {Type: "section"}}},
			SlackMovie:   &models.SlackPayload{Text: "Movie", Blocks: []models.SlackBlock{{Type: "section"}, {Type: "section"}}},
		}

		err := svc.DispatchPayloads(ctx, cfg, payloads)
		if err == nil {
			t.Errorf("Expected errors during segregated dispatch when endpoints return 500")
		}
	})

	t.Run("Distinct NTFY Topics with both TV and Movie payloads", func(t *testing.T) {
		tvHitCount, movieHitCount, combinedHitCount = 0, 0, 0
		tsNtfy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case strings.Contains(r.URL.Path, "tv-topic"):
				tvHitCount++
			case strings.Contains(r.URL.Path, "movie-topic"):
				movieHitCount++
			default:
				combinedHitCount++
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer tsNtfy.Close()

		cfg := models.DefaultConfig()
		cfg.UseNtfy = true
		cfg.NtfyServerURL = tsNtfy.URL
		cfg.NtfyTopic = "combined-topic"
		cfg.NtfyTVTopic = "tv-topic"
		cfg.NtfyMovieTopic = "movie-topic"

		payloads := models.NotificationPayloads{
			NtfyCombined: &models.NtfyPayload{Title: "Combined", Message: "Combined", HasReleases: true},
			NtfyTV:       &models.NtfyPayload{Title: "TV", Message: "TV Show", HasReleases: true},
			NtfyMovie:    &models.NtfyPayload{Title: "Movie", Message: "Movie Title", HasReleases: true},
		}

		err := svc.DispatchPayloads(ctx, cfg, payloads)
		if err != nil {
			t.Fatalf("DispatchPayloads failed: %v", err)
		}
		if tvHitCount != 1 {
			t.Errorf("Expected tvHitCount 1, got %d", tvHitCount)
		}
		if movieHitCount != 1 {
			t.Errorf("Expected movieHitCount 1, got %d", movieHitCount)
		}
		if combinedHitCount != 0 {
			t.Errorf("Expected combinedHitCount 0, got %d", combinedHitCount)
		}
	})

	t.Run("Distinct Telegram Channels with both TV and Movie payloads", func(t *testing.T) {
		tvHitCount, movieHitCount, combinedHitCount = 0, 0, 0
		tsTg := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body struct {
				ChatID string `json:"chat_id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			switch body.ChatID {
			case "tv-chat":
				tvHitCount++
			case "movie-chat":
				movieHitCount++
			default:
				combinedHitCount++
			}
			w.WriteHeader(http.StatusOK)
		}))
		defer tsTg.Close()

		svcTg := &Service{telegramBaseURL: tsTg.URL}
		cfg := models.DefaultConfig()
		cfg.UseTelegram = true
		cfg.TelegramBotToken = "bot-token"
		cfg.TelegramChatID = "default-chat"
		cfg.TelegramTVChatID = "tv-chat"
		cfg.TelegramMovieChatID = "movie-chat"

		payloads := models.NotificationPayloads{
			TelegramCombined: &models.TelegramPayload{Text: "Combined", HasReleases: true},
			TelegramTV:       &models.TelegramPayload{Text: "TV Show", HasReleases: true},
			TelegramMovie:    &models.TelegramPayload{Text: "Movie Title", HasReleases: true},
		}

		err := svcTg.DispatchPayloads(ctx, cfg, payloads)
		if err != nil {
			t.Fatalf("DispatchPayloads failed: %v", err)
		}
		if tvHitCount != 1 {
			t.Errorf("Expected tvHitCount 1, got %d", tvHitCount)
		}
		if movieHitCount != 1 {
			t.Errorf("Expected movieHitCount 1, got %d", movieHitCount)
		}
		if combinedHitCount != 0 {
			t.Errorf("Expected combinedHitCount 0, got %d", combinedHitCount)
		}
	})

	t.Run("Helper functions nil and empty edge cases", func(t *testing.T) {
		if svc.shouldSendDiscordPayload(nil, true) {
			t.Errorf("Expected shouldSendDiscordPayload to return false for nil payload")
		}
		if svc.shouldSendSlackPayload(nil, true) {
			t.Errorf("Expected shouldSendSlackPayload to return false for nil payload")
		}
		if svc.shouldSendNtfyPayload(nil, true) {
			t.Errorf("Expected shouldSendNtfyPayload to return false for nil payload")
		}
		if svc.shouldSendTelegramPayload(nil, true) {
			t.Errorf("Expected shouldSendTelegramPayload to return false for nil payload")
		}
		if got := svc.dispatchTarget(ctx, http.DefaultClient, "Test", "", nil); got != "" {
			t.Errorf("Expected empty string for empty url and nil payload, got %q", got)
		}
		if got := svc.dispatchTargetNtfy(ctx, http.DefaultClient, "Test", "", models.DefaultConfig(), nil); got != "" {
			t.Errorf("Expected empty string for empty topic and nil payload, got %q", got)
		}
		if got := svc.dispatchTargetTelegram(ctx, http.DefaultClient, "Test", "", models.DefaultConfig(), nil); got != "" {
			t.Errorf("Expected empty string for empty chat ID and nil payload, got %q", got)
		}
		emptyPayloads := models.NotificationPayloads{}
		if svc.hasDiscordReleases(emptyPayloads) {
			t.Errorf("Expected hasDiscordReleases to return false for empty payloads")
		}
		if svc.hasSlackReleases(emptyPayloads) {
			t.Errorf("Expected hasSlackReleases to return false for empty payloads")
		}
		if svc.hasNtfyReleases(emptyPayloads) {
			t.Errorf("Expected hasNtfyReleases to return false for empty payloads")
		}
		if svc.hasTelegramReleases(emptyPayloads) {
			t.Errorf("Expected hasTelegramReleases to return false for empty payloads")
		}

		serverURL, topic := svc.resolveNtfyTarget(models.TestWebhookRequest{WebhookURL: "https://custom.ntfy.sh/demo-topic"})
		if serverURL != "https://custom.ntfy.sh" || topic != "demo-topic" {
			t.Errorf("Unexpected resolveNtfyTarget result: %s, %s", serverURL, topic)
		}
	})
}
