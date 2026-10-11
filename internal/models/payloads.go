package models

// NtfyPayload represents an outgoing message to NTFY.
type NtfyPayload struct {
	Title       string   `json:"title"`
	Message     string   `json:"message"`
	Priority    int      `json:"priority"`
	Tags        []string `json:"tags"`
	HasReleases bool     `json:"-"`
}

// TelegramPayload represents an outgoing message to Telegram Bot API sendMessage.
type TelegramPayload struct {
	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	HasReleases           bool   `json:"-"`
}

// NotificationPayloads aggregates formatted notification payloads for combined and dedicated media channels.
type NotificationPayloads struct {
	DiscordCombined  *DiscordPayload
	DiscordTV        *DiscordPayload
	DiscordMovie     *DiscordPayload
	SlackCombined    *SlackPayload
	SlackTV          *SlackPayload
	SlackMovie       *SlackPayload
	NtfyCombined     *NtfyPayload
	NtfyTV           *NtfyPayload
	NtfyMovie        *NtfyPayload
	TelegramCombined *TelegramPayload
	TelegramTV       *TelegramPayload
	TelegramMovie    *TelegramPayload
}

// TestWebhookRequest represents the parameters for testing a notification webhook/platform endpoint.
type TestWebhookRequest struct {
	Platform   string `json:"platform"`
	WebhookURL string `json:"webhook_url,omitempty"`
	ServerURL  string `json:"server_url,omitempty"`
	Topic      string `json:"topic,omitempty"`
	Token      string `json:"token,omitempty"`
	Priority   string `json:"priority,omitempty"`
	ChatID     string `json:"chat_id,omitempty"`
}
