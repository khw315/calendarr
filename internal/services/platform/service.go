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

type Service struct {
	telegramBaseURL string
}

func NewService() *Service {
	return &Service{
		telegramBaseURL: "https://api.telegram.org",
	}
}

func (s *Service) getTelegramBaseURL() string {
	if s.telegramBaseURL != "" {
		return strings.TrimRight(s.telegramBaseURL, "/")
	}
	return "https://api.telegram.org"
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

// DispatchPayloads dispatches notifications to combined or media-segregated channels across all enabled platforms.
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

	if ntfyErrs := s.dispatchNtfy(ctx, client, cfg, payloads); len(ntfyErrs) > 0 {
		errs = append(errs, ntfyErrs...)
	}

	if telegramErrs := s.dispatchTelegram(ctx, client, cfg, payloads); len(telegramErrs) > 0 {
		errs = append(errs, telegramErrs...)
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

func (s *Service) dispatchNtfy(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	if !cfg.UseNtfy {
		return nil
	}

	if cfg.HasDistinctNtfyTopics() {
		return s.dispatchNtfySegregated(ctx, client, cfg, payloads)
	}

	if errStr := s.dispatchTargetNtfy(ctx, client, "NTFY", cfg.NtfyTopic, cfg, payloads.NtfyCombined); errStr != "" {
		return []string{errStr}
	}

	return nil
}

func (s *Service) dispatchNtfySegregated(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	var errs []string
	hasAny := s.hasNtfyReleases(payloads)

	if s.shouldSendNtfyPayload(payloads.NtfyTV, hasAny) {
		if errStr := s.dispatchTargetNtfy(ctx, client, "NTFY TV", cfg.GetNtfyTVTopic(), cfg, payloads.NtfyTV); errStr != "" {
			errs = append(errs, errStr)
		}
	}

	if s.shouldSendNtfyPayload(payloads.NtfyMovie, hasAny) {
		if errStr := s.dispatchTargetNtfy(ctx, client, "NTFY Movie", cfg.GetNtfyMovieTopic(), cfg, payloads.NtfyMovie); errStr != "" {
			errs = append(errs, errStr)
		}
	}

	return errs
}

func (s *Service) shouldSendNtfyPayload(payload *models.NtfyPayload, hasAnyReleases bool) bool {
	if payload == nil {
		return false
	}
	return payload.HasReleases || !hasAnyReleases
}

func (s *Service) hasNtfyReleases(payloads models.NotificationPayloads) bool {
	if payloads.NtfyCombined != nil && payloads.NtfyCombined.HasReleases {
		return true
	}
	if payloads.NtfyTV != nil && payloads.NtfyTV.HasReleases {
		return true
	}
	if payloads.NtfyMovie != nil && payloads.NtfyMovie.HasReleases {
		return true
	}
	return false
}

func (s *Service) dispatchTargetNtfy(ctx context.Context, client *http.Client, name, topic string, cfg *models.Config, payload *models.NtfyPayload) string {
	if topic == "" || payload == nil {
		return ""
	}

	serverURL := cfg.NtfyServerURL
	if serverURL == "" {
		serverURL = "https://ntfy.sh"
	}

	priorityStr := ""
	if cfg.NtfyPriority > 0 {
		priorityStr = fmt.Sprintf("%d", cfg.NtfyPriority)
	}

	if err := s.postNtfy(ctx, client, serverURL, topic, cfg.NtfyAccessToken, priorityStr, payload); err != nil {
		log.Printf("❌ Failed to send %s notification: %v", name, err)
		return fmt.Sprintf("%s error: %v", name, err)
	}

	log.Printf("✅ Successfully dispatched notification to %s topic", name)
	return ""
}

func (s *Service) dispatchTelegram(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	if !cfg.UseTelegram {
		return nil
	}

	if cfg.HasDistinctTelegramChannels() {
		return s.dispatchTelegramSegregated(ctx, client, cfg, payloads)
	}

	if errStr := s.dispatchTargetTelegram(ctx, client, "Telegram", cfg.TelegramChatID, cfg, payloads.TelegramCombined); errStr != "" {
		return []string{errStr}
	}

	return nil
}

func (s *Service) dispatchTelegramSegregated(ctx context.Context, client *http.Client, cfg *models.Config, payloads models.NotificationPayloads) []string {
	var errs []string
	hasAny := s.hasTelegramReleases(payloads)

	if s.shouldSendTelegramPayload(payloads.TelegramTV, hasAny) {
		if errStr := s.dispatchTargetTelegram(ctx, client, "Telegram TV", cfg.GetTelegramTVChatID(), cfg, payloads.TelegramTV); errStr != "" {
			errs = append(errs, errStr)
		}
	}

	if s.shouldSendTelegramPayload(payloads.TelegramMovie, hasAny) {
		if errStr := s.dispatchTargetTelegram(ctx, client, "Telegram Movie", cfg.GetTelegramMovieChatID(), cfg, payloads.TelegramMovie); errStr != "" {
			errs = append(errs, errStr)
		}
	}

	return errs
}

func (s *Service) shouldSendTelegramPayload(payload *models.TelegramPayload, hasAnyReleases bool) bool {
	if payload == nil {
		return false
	}
	return payload.HasReleases || !hasAnyReleases
}

func (s *Service) hasTelegramReleases(payloads models.NotificationPayloads) bool {
	if payloads.TelegramCombined != nil && payloads.TelegramCombined.HasReleases {
		return true
	}
	if payloads.TelegramTV != nil && payloads.TelegramTV.HasReleases {
		return true
	}
	if payloads.TelegramMovie != nil && payloads.TelegramMovie.HasReleases {
		return true
	}
	return false
}

func (s *Service) dispatchTargetTelegram(ctx context.Context, client *http.Client, name, chatID string, cfg *models.Config, payload *models.TelegramPayload) string {
	if chatID == "" || cfg.TelegramBotToken == "" || payload == nil {
		return ""
	}

	if err := s.postTelegram(ctx, client, cfg.TelegramBotToken, chatID, payload); err != nil {
		log.Printf("❌ Failed to send %s notification: %v", name, err)
		return fmt.Sprintf("%s error: %v", name, err)
	}

	log.Printf("✅ Successfully dispatched notification to %s chat", name)
	return ""
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

func (s *Service) postNtfy(ctx context.Context, client *http.Client, serverURL, topic, token, priority string, payload *models.NtfyPayload) error {
	cleanServer := strings.TrimRight(strings.TrimSpace(serverURL), "/")
	if cleanServer == "" {
		cleanServer = "https://ntfy.sh"
	}
	cleanTopic := strings.TrimLeft(strings.TrimSpace(topic), "/")
	targetURL := fmt.Sprintf("%s/%s", cleanServer, cleanTopic)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, strings.NewReader(payload.Message))
	if err != nil {
		return fmt.Errorf("failed to create NTFY request: %w", err)
	}

	req.Header.Set("Title", payload.Title)
	req.Header.Set("Markdown", "yes")

	prio := priority
	if prio == "" && payload.Priority > 0 {
		prio = fmt.Sprintf("%d", payload.Priority)
	}
	if prio != "" {
		req.Header.Set("Priority", prio)
	}

	if len(payload.Tags) > 0 {
		req.Header.Set("Tags", strings.Join(payload.Tags, ","))
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("NTFY HTTP POST failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("NTFY returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Service) postTelegram(ctx context.Context, client *http.Client, botToken, chatID string, payload *models.TelegramPayload) error {
	targetURL := fmt.Sprintf("%s/bot%s/sendMessage", s.getTelegramBaseURL(), botToken)

	type telegramReq struct {
		ChatID                string `json:"chat_id"`
		Text                  string `json:"text"`
		ParseMode             string `json:"parse_mode"`
		DisableNotification   bool   `json:"disable_notification,omitempty"`
		DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	}

	body := telegramReq{
		ChatID:                chatID,
		Text:                  payload.Text,
		ParseMode:             payload.ParseMode,
		DisableNotification:   payload.DisableNotification,
		DisableWebPagePreview: payload.DisableWebPagePreview,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal Telegram request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create Telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Telegram HTTP POST failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Telegram returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// SendTest sends an immediate verification test notification to the specified platform webhook.
func (s *Service) SendTest(ctx context.Context, platformName, webhookURL string) error {
	return s.SendTestWithOptions(ctx, models.TestWebhookRequest{
		Platform:   platformName,
		WebhookURL: webhookURL,
	})
}

// SendTestWithOptions sends an immediate verification test notification using provided options.
func (s *Service) SendTestWithOptions(ctx context.Context, req models.TestWebhookRequest) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	normPlatform := strings.ToLower(strings.TrimSpace(req.Platform))
	switch normPlatform {
	case "discord":
		payload := s.buildDiscordTestPayload()
		return s.postWebhook(ctx, client, "Discord", req.WebhookURL, payload)
	case "slack":
		payload := s.buildSlackTestPayload()
		return s.postWebhook(ctx, client, "Slack", req.WebhookURL, payload)
	case "ntfy":
		payload := s.buildNtfyTestPayload()
		serverURL, topic := s.resolveNtfyTarget(req)
		if topic == "" {
			return fmt.Errorf("NTFY topic cannot be empty")
		}
		return s.postNtfy(ctx, client, serverURL, topic, req.Token, req.Priority, payload)
	case "telegram":
		payload := s.buildTelegramTestPayload()
		if req.Token == "" || req.ChatID == "" {
			return fmt.Errorf("both Telegram bot token and chat ID are required")
		}
		return s.postTelegram(ctx, client, req.Token, req.ChatID, payload)
	default:
		return fmt.Errorf("unsupported platform for webhook test: %s", req.Platform)
	}
}

func (s *Service) resolveNtfyTarget(req models.TestWebhookRequest) (string, string) {
	serverURL := strings.TrimSpace(req.ServerURL)
	topic := strings.TrimSpace(req.Topic)

	if topic == "" && req.WebhookURL != "" {
		clean := strings.TrimRight(req.WebhookURL, "/")
		lastSlash := strings.LastIndex(clean, "/")
		if lastSlash != -1 && lastSlash > 8 {
			return clean[:lastSlash], clean[lastSlash+1:]
		}
		topic = clean
	}

	if serverURL == "" {
		serverURL = "https://ntfy.sh"
	}
	return serverURL, topic
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

func (s *Service) buildNtfyTestPayload() *models.NtfyPayload {
	return &models.NtfyPayload{
		Title:       "🔔 Calendarr NTFY Test",
		Message:     "This is a test notification from **Calendarr**.\nYour NTFY integration is configured and communicating properly!",
		Priority:    3,
		Tags:        []string{"bell", "white_check_mark"},
		HasReleases: true,
	}
}

func (s *Service) buildTelegramTestPayload() *models.TelegramPayload {
	return &models.TelegramPayload{
		Text:                  "<b>🔔 Calendarr Telegram Test</b>\n\nThis is a test notification from <b>Calendarr</b>.\nYour Telegram bot integration is active and working properly!",
		ParseMode:             "HTML",
		DisableNotification:   false,
		DisableWebPagePreview: true,
		HasReleases:           true,
	}
}
