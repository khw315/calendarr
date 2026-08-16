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
