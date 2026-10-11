package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/khw315/calendarr/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Dispatch(ctx context.Context, cfg *models.Config, formatted *models.DiscordPayload, slackFormatted *models.SlackPayload) error {
	client := &http.Client{
		Timeout: time.Duration(cfg.HTTPTimeout) * time.Second,
	}

	var errs []string

	if cfg.UseDiscord && cfg.DiscordWebhookURL != "" {
		if err := s.postWebhook(ctx, client, "Discord", cfg.DiscordWebhookURL, formatted); err != nil {
			log.Printf("❌ Failed to send Discord webhook: %v", err)
			errs = append(errs, fmt.Sprintf("Discord error: %v", err))
		} else {
			log.Println("✅ Successfully dispatched notification to Discord webhook")
		}
	}

	if cfg.UseSlack && cfg.SlackWebhookURL != "" {
		if err := s.postWebhook(ctx, client, "Slack", cfg.SlackWebhookURL, slackFormatted); err != nil {
			log.Printf("❌ Failed to send Slack webhook: %v", err)
			errs = append(errs, fmt.Sprintf("Slack error: %v", err))
		} else {
			log.Println("✅ Successfully dispatched notification to Slack webhook")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during dispatch: %v", errs)
	}

	return nil
}

func (s *Service) postWebhook(ctx context.Context, client *http.Client, platformName, webhookURL string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal %s payload: %w", platformName, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create %s request: %w", platformName, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("%s HTTP POST failed: %w", platformName, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s returned status %d: %s", platformName, resp.StatusCode, string(body))
	}

	return nil
}

// SendTest sends an immediate verification test notification to the specified platform webhook.
func (s *Service) SendTest(ctx context.Context, platformName, webhookURL string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	normPlatform := strings.ToLower(strings.TrimSpace(platformName))
	switch normPlatform {
	case "discord":
		payload := s.buildDiscordTestPayload()
		return s.postWebhook(ctx, client, "Discord", webhookURL, payload)
	case "slack":
		payload := s.buildSlackTestPayload()
		return s.postWebhook(ctx, client, "Slack", webhookURL, payload)
	default:
		return fmt.Errorf("unsupported platform for webhook test: %s", platformName)
	}
}

func (s *Service) buildDiscordTestPayload() *models.DiscordPayload {
	return &models.DiscordPayload{
		Embeds: []models.DiscordEmbed{
			{
				Title:       "🔔 Calendarr Webhook Test",
				Description: "This is a test notification from **Calendarr**.\nYour Discord webhook is configured and communicating properly!",
				Color:       0x5865F2, // Discord Blurple
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
				Footer: &models.DiscordFooter{
					Text: "Calendarr • Webhook Verification",
				},
			},
		},
	}
}

func (s *Service) buildSlackTestPayload() *models.SlackPayload {
	return &models.SlackPayload{
		Text: "🔔 Calendarr Webhook Test: Your Slack integration is active!",
		Blocks: []models.SlackBlock{
			{
				Type: "header",
				Text: &models.SlackText{
					Type: "plain_text",
					Text: "🔔 Calendarr Webhook Test",
				},
			},
			{
				Type: "section",
				Text: &models.SlackText{
					Type: "mrkdwn",
					Text: "This is a test notification from *Calendarr*.\nYour Slack webhook integration is active and working properly!",
				},
			},
		},
	}
}
