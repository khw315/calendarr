package formatter

import (
	"strings"
	"testing"
	"time"

	"github.com/khw315/calendarr/internal/constants"
	"github.com/khw315/calendarr/internal/models"
)

func TestFormatterService(t *testing.T) {
	svc := NewService()
	cfg := models.DefaultConfig()
	cfg.DiscordMentionRoleID = "123456789"
	cfg.ShowTimezoneInSubheader = true
	cfg.Timezone = "Asia/Jakarta"
	loc, _ := time.LoadLocation("Asia/Jakarta")
	cfg.TimezoneLocation = loc

	now := time.Now().In(loc)

	events := []*models.Event{
		{
			Summary:    "Severance S02E01",
			StartTime:  now.Add(1 * time.Hour),
			EndTime:    now.Add(2 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		{
			Summary:    "Dune Part 2",
			StartTime:  now.Add(3 * time.Hour),
			EndTime:    now.Add(5 * time.Hour),
			SourceType: constants.EventTypeMovie,
		},
	}

	startDate := now.Add(-1 * time.Hour)
	endDate := now.Add(24 * time.Hour)

	res := svc.Format(events, cfg, startDate, endDate)
	if res == nil {
		t.Fatalf("Format returned nil result")
	}

	if res.Discord == nil || len(res.Discord.Embeds) == 0 {
		t.Errorf("Expected Discord payload with embeds")
	}

	if !strings.Contains(res.Discord.Content, "# New Releases") {
		t.Errorf("Expected '# New Releases' in Discord content when events are present, got: %s", res.Discord.Content)
	}

	if res.Slack == nil || len(res.Slack.Blocks) == 0 {
		t.Errorf("Expected Slack payload with blocks")
	}

	if res.Counts["tv_count"] != 1 || res.Counts["movie_count"] != 1 {
		t.Errorf("Unexpected event counts: %v", res.Counts)
	}

	emptyRes := svc.Format([]*models.Event{}, cfg, startDate, endDate)
	if emptyRes.Discord == nil || len(emptyRes.Discord.Embeds) != 0 {
		t.Errorf("Expected 0 Discord embeds for empty events schedule, got %d", len(emptyRes.Discord.Embeds))
	}
	if strings.Contains(emptyRes.Discord.Content, "New Releases") {
		t.Errorf("Expected 'New Releases' to be omitted when schedule is empty, got: %s", emptyRes.Discord.Content)
	}
	if !strings.HasPrefix(emptyRes.Discord.Content, "# ") {
		t.Errorf("Expected empty release message as H1 header, got: %s", emptyRes.Discord.Content)
	}
	if emptyRes.Slack == nil || len(emptyRes.Slack.Blocks) == 0 {
		t.Errorf("Expected fallback Slack blocks for empty events")
	}
	if emptyRes.Slack.Blocks[0].Text.Text == "New Releases" || emptyRes.Slack.Blocks[0].Text.Text == "" {
		t.Errorf("Expected fallback localized message instead of 'New Releases' in Slack header, got: %s", emptyRes.Slack.Blocks[0].Text.Text)
	}
}
