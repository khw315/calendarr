package platform

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khw315/calendarr/internal/models"
)

func TestPlatformServiceDispatchSuccess(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	svc := NewService()
	cfg := models.DefaultConfig()
	cfg.UseDiscord = true
	cfg.DiscordWebhookURL = ts.URL
	cfg.UseSlack = true
	cfg.SlackWebhookURL = ts.URL

	discordPayload := &models.DiscordPayload{Content: "Test Discord"}
	slackPayload := &models.SlackPayload{Text: "Test Slack"}

	err := svc.Dispatch(context.Background(), cfg, discordPayload, slackPayload)
	if err != nil {
		t.Fatalf("Dispatch failed: %v", err)
	}
}

func TestPlatformServiceDispatchError(t *testing.T) {
	tsErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Internal Error"))
	}))
	defer tsErr.Close()

	svc := NewService()
	cfg := models.DefaultConfig()
	cfg.UseDiscord = true
	cfg.DiscordWebhookURL = tsErr.URL
	cfg.UseSlack = true
	cfg.SlackWebhookURL = tsErr.URL

	discordPayload := &models.DiscordPayload{Content: "Test Discord"}
	slackPayload := &models.SlackPayload{Text: "Test Slack"}

	err := svc.Dispatch(context.Background(), cfg, discordPayload, slackPayload)
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

	svc := NewService()
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
}
