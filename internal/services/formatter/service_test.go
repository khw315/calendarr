package formatter

import (
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

	if res.Slack == nil || len(res.Slack.Blocks) == 0 {
		t.Errorf("Expected Slack payload with blocks")
	}

	if res.Counts["tv_count"] != 1 || res.Counts["movie_count"] != 1 {
		t.Errorf("Unexpected event counts: %v", res.Counts)
	}

	// Test empty events payload
	emptyRes := svc.Format([]*models.Event{}, cfg, startDate, endDate)
	if emptyRes.Discord == nil || len(emptyRes.Discord.Embeds) == 0 {
		t.Errorf("Expected fallback Discord embed for empty events")
	}
	if emptyRes.Slack == nil || len(emptyRes.Slack.Blocks) == 0 {
		t.Errorf("Expected fallback Slack blocks for empty events")
	}
}
