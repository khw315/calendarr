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
