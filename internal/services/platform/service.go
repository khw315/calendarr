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

// Dispatch dispatches notification payloads using legacy/combined arguments.
func (s *Service) Dispatch(ctx context.Context, cfg *models.Config, formatted *models.DiscordPayload, slackFormatted *models.SlackPayload) error {
	return s.DispatchPayloads(ctx, cfg, models.NotificationPayloads{
		DiscordCombined: formatted,
		DiscordTV:       formatted,
		DiscordMovie:    formatted,
		SlackCombined:   slackFormatted,
		SlackTV:         slackFormatted,
		SlackMovie:      slackFormatted,
	})
}

// DispatchPayloads dispatches notifications to combined or media-segregated channels.
func (s *Service) DispatchPayloads(ctx context.Context, cfg *models.Config, payloads models.NotificationPayloads) error {
	client := &http.Client{
		Timeout: time.Duration(cfg.HTTPTimeout) * time.Second,
	}

	var errs []string
	if discordErrs := s.dispatchDiscord(ctx, client, cfg, payloads); len(discordErrs) > 0 {
		errs = append(errs, discordErrs...)
	}

	if slackErrs := s.dispatchSlack(ctx, client, cfg, payloads); len(slackErrs) > 0 {
		errs = append(errs, slackErrs...)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during dispatch: %v", errs)
	}

	return nil
}

func (s *Service) dispatchDiscord(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	if !cfg.UseDiscord {
		return nil
	}

	if cfg.HasDistinctDiscordChannels() {
		return s.dispatchDiscordSegregated(ctx, client, cfg, payloads)
	}

	if errStr := s.dispatchTarget(ctx, client, "Discord", cfg.DiscordWebhookURL, payloads.DiscordCombined); errStr != "" {
		return []string{errStr}
	}

	return nil
}

func (s *Service) dispatchDiscordSegregated(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	var errs []string
	hasAny := s.hasDiscordReleases(payloads)

	if s.shouldSendDiscordPayload(payloads.DiscordTV, hasAny) {
		if errStr := s.dispatchTarget(ctx, client, "Discord TV", cfg.GetDiscordTVWebhookURL(), payloads.DiscordTV); errStr != "" {
			errs = append(errs, errStr)
		}
	}

	if s.shouldSendDiscordPayload(payloads.DiscordMovie, hasAny) {
		if errStr := s.dispatchTarget(ctx, client, "Discord Movie", cfg.GetDiscordMovieWebhookURL(), payloads.DiscordMovie); errStr != "" {
			errs = append(errs, errStr)
		}
	}

	return errs
}

func (s *Service) shouldSendDiscordPayload(payload *models.DiscordPayload, hasAnyReleases bool) bool {
	if payload == nil {
		return false
	}
	return len(payload.Embeds) > 0 || !hasAnyReleases
}

func (s *Service) hasDiscordReleases(payloads models.NotificationPayloads) bool {
	if payloads.DiscordCombined != nil && len(payloads.DiscordCombined.Embeds) > 0 {
		return true
	}
	if payloads.DiscordTV != nil && len(payloads.DiscordTV.Embeds) > 0 {
		return true
	}
	if payloads.DiscordMovie != nil && len(payloads.DiscordMovie.Embeds) > 0 {
		return true
	}
	return false
}

func (s *Service) dispatchSlack(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	if !cfg.UseSlack {
		return nil
	}

	if cfg.HasDistinctSlackChannels() {
		return s.dispatchSlackSegregated(ctx, client, cfg, payloads)
	}

	if errStr := s.dispatchTarget(ctx, client, "Slack", cfg.SlackWebhookURL, payloads.SlackCombined); errStr != "" {
		return []string{errStr}
	}

	return nil
}

func (s *Service) dispatchSlackSegregated(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	var errs []string
	hasAny := s.hasSlackReleases(payloads)

	if s.shouldSendSlackPayload(payloads.SlackTV, hasAny) {
		if errStr := s.dispatchTarget(ctx, client, "Slack TV", cfg.GetSlackTVWebhookURL(), payloads.SlackTV); errStr != "" {
			errs = append(errs, errStr)
		}
	}

	if s.shouldSendSlackPayload(payloads.SlackMovie, hasAny) {
		if errStr := s.dispatchTarget(ctx, client, "Slack Movie", cfg.GetSlackMovieWebhookURL(), payloads.SlackMovie); errStr != "" {
			errs = append(errs, errStr)
		}
	}

	return errs
}

func (s *Service) shouldSendSlackPayload(payload *models.SlackPayload, hasAnyReleases bool) bool {
	if payload == nil {
		return false
	}
	return len(payload.Blocks) > 1 || !hasAnyReleases
}

func (s *Service) hasSlackReleases(payloads models.NotificationPayloads) bool {
	if payloads.SlackCombined != nil && len(payloads.SlackCombined.Blocks) > 1 {
		return true
	}
	if payloads.SlackTV != nil && len(payloads.SlackTV.Blocks) > 1 {
		return true
	}
	if payloads.SlackMovie != nil && len(payloads.SlackMovie.Blocks) > 1 {
		return true
	}
	return false
}

func (s *Service) dispatchTarget(ctx context.Context, client *http.Client, name, url string, payload interface{}) string {
	if url == "" || payload == nil {
		return ""
	}

	if err := s.postWebhook(ctx, client, name, url, payload); err != nil {
		log.Printf("❌ Failed to send %s webhook: %v", name, err)
		return fmt.Sprintf("%s error: %v", name, err)
	}

	log.Printf("✅ Successfully dispatched notification to %s webhook", name)
	return ""
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
