package constants

import (
	"testing"
)

func TestConstantsValues(t *testing.T) {
	if EventTypeTV != "tv" {
		t.Errorf("Expected event type tv, got %s", EventTypeTV)
	}
	if EventTypeMovie != "movie" {
		t.Errorf("Expected event type movie, got %s", EventTypeMovie)
	}
	if PlatformDiscord != "discord" || PlatformSlack != "slack" {
		t.Errorf("Unexpected platform constants")
	}
	if PassedEventDisplay != "DISPLAY" || PassedEventHide != "HIDE" || PassedEventStrike != "STRIKE" {
		t.Errorf("Unexpected passed event handling constants")
	}

	// Test Color Maps
	if blue, ok := DiscordColors["blue"]; !ok || blue == 0 {
		t.Errorf("Expected valid blue color for Discord")
	}
	if red, ok := SlackColors["red"]; !ok || red == "" {
		t.Errorf("Expected valid red color for Slack")
	}
}
