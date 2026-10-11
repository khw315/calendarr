package models

// NotificationPayloads aggregates formatted notification payloads for combined and dedicated media channels.
type NotificationPayloads struct {
	DiscordCombined *DiscordPayload
	DiscordTV       *DiscordPayload
	DiscordMovie    *DiscordPayload
	SlackCombined   *SlackPayload
	SlackTV         *SlackPayload
	SlackMovie      *SlackPayload
}
