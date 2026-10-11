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

	var errs []string

	if cfg.HasDistinctDiscordChannels() {
		hasAnyReleases := s.hasDiscordReleases(payloads)

		tvURL := cfg.GetDiscordTVWebhookURL()
		if tvURL != "" && payloads.DiscordTV != nil {
			if len(payloads.DiscordTV.Embeds) > 0 || !hasAnyReleases {
				if err := s.postWebhook(ctx, client, "Discord (TV)", tvURL, payloads.DiscordTV); err != nil {
					log.Printf("❌ Failed to send Discord TV webhook: %v", err)
					errs = append(errs, fmt.Sprintf("Discord TV error: %v", err))
				} else {
					log.Println("✅ Successfully dispatched notification to Discord TV webhook")
				}
			}
		}

		movieURL := cfg.GetDiscordMovieWebhookURL()
		if movieURL != "" && payloads.DiscordMovie != nil {
			if len(payloads.DiscordMovie.Embeds) > 0 || !hasAnyReleases {
				if err := s.postWebhook(ctx, client, "Discord (Movie)", movieURL, payloads.DiscordMovie); err != nil {
					log.Printf("❌ Failed to send Discord Movie webhook: %v", err)
					errs = append(errs, fmt.Sprintf("Discord Movie error: %v", err))
				} else {
					log.Println("✅ Successfully dispatched notification to Discord Movie webhook")
				}
			}
		}
		return errs
	}

	if cfg.DiscordWebhookURL != "" && payloads.DiscordCombined != nil {
		if err := s.postWebhook(ctx, client, "Discord", cfg.DiscordWebhookURL, payloads.DiscordCombined); err != nil {
			log.Printf("❌ Failed to send Discord webhook: %v", err)
			errs = append(errs, fmt.Sprintf("Discord error: %v", err))
		} else {
			log.Println("✅ Successfully dispatched notification to Discord webhook")
		}
	}

	return errs
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

func (s *Service) dispatchSlack(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	if !cfg.UseSlack {
		return nil
	}

	var errs []string

	if cfg.HasDistinctSlackChannels() {
		hasAnyReleases := s.hasSlackReleases(payloads)

		tvURL := cfg.GetSlackTVWebhookURL()
		if tvURL != "" && payloads.SlackTV != nil {
			if len(payloads.SlackTV.Blocks) > 1 || !hasAnyReleases {
				if err := s.postWebhook(ctx, client, "Slack (TV)", tvURL, payloads.SlackTV); err != nil {
					log.Printf("❌ Failed to send Slack TV webhook: %v", err)
					errs = append(errs, fmt.Sprintf("Slack TV error: %v", err))
				} else {
					log.Println("✅ Successfully dispatched notification to Slack TV webhook")
				}
			}
		}

		movieURL := cfg.GetSlackMovieWebhookURL()
		if movieURL != "" && payloads.SlackMovie != nil {
			if len(payloads.SlackMovie.Blocks) > 1 || !hasAnyReleases {
				if err := s.postWebhook(ctx, client, "Slack (Movie)", movieURL, payloads.SlackMovie); err != nil {
					log.Printf("❌ Failed to send Slack Movie webhook: %v", err)
					errs = append(errs, fmt.Sprintf("Slack Movie error: %v", err))
				} else {
					log.Println("✅ Successfully dispatched notification to Slack Movie webhook")
				}
			}
		}
		return errs
	}

	if cfg.SlackWebhookURL != "" && payloads.SlackCombined != nil {
		if err := s.postWebhook(ctx, client, "Slack", cfg.SlackWebhookURL, payloads.SlackCombined); err != nil {
			log.Printf("❌ Failed to send Slack webhook: %v", err)
			errs = append(errs, fmt.Sprintf("Slack error: %v", err))
		} else {
			log.Println("✅ Successfully dispatched notification to Slack webhook")
		}
	}

	return errs
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
