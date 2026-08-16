package platform

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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
		if err := s.sendDiscord(ctx, client, cfg.DiscordWebhookURL, formatted); err != nil {
			log.Printf("❌ Failed to send Discord webhook: %v", err)
			errs = append(errs, fmt.Sprintf("Discord error: %v", err))
		} else {
			log.Println("✅ Successfully dispatched notification to Discord webhook")
		}
	}

	if cfg.UseSlack && cfg.SlackWebhookURL != "" {
		if err := s.sendSlack(ctx, client, cfg.SlackWebhookURL, slackFormatted); err != nil {
			log.Printf("❌ Failed to send Slack webhook: %v", err)
			errs = append(errs, fmt.Sprintf("Slack error: %v", err))
		} else {
			log.Println("✅ Successfully dispatched notification to Slack webhook")
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors during dispatch: %s", fmt.Sprintf("%v", errs))
	}

	return nil
}

func (s *Service) sendDiscord(ctx context.Context, client *http.Client, webhookURL string, payload *models.DiscordPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create Discord request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Discord HTTP POST failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Discord returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Service) sendSlack(ctx context.Context, client *http.Client, webhookURL string, payload *models.SlackPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create Slack request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Slack HTTP POST failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Slack returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
