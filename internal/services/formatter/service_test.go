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

func TestBulkAndPremiereDiscordFormatting(t *testing.T) {
	svc := NewService()
	cfg := models.DefaultConfig()
	loc := time.UTC
	cfg.TimezoneLocation = loc
	cfg.DiscordTimestampStyle = "R"

	refTime := time.Date(2026, time.August, 7, 12, 0, 0, 0, loc)

	events := []*models.Event{
		// Show with 12 episodes (> 2 episodes) -> bulk grouped
		{
			Summary:    "Our Sticky Love - 1x01 - Episode 1",
			StartTime:  refTime,
			EndTime:    refTime.Add(1 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		{
			Summary:    "Our Sticky Love - 1x02 - Episode 2",
			StartTime:  refTime,
			EndTime:    refTime.Add(1 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		{
			Summary:    "Our Sticky Love - 1x03 - Episode 3",
			StartTime:  refTime,
			EndTime:    refTime.Add(1 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		// Show with 1 episode (non-premiere)
		{
			Summary:    "A Bona Fide Killer - 1x03 - Episode 3",
			StartTime:  refTime.Add(2 * time.Hour),
			EndTime:    refTime.Add(3 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		// Show with 1 episode (premiere)
		{
			Summary:    "Flex x Cop - 2x01 - Episode 1",
			StartTime:  refTime.Add(4 * time.Hour),
			EndTime:    refTime.Add(5 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		// Show with exactly 2 episodes (<= 2 episodes) -> individual lines
		{
			Summary:    "Double Trouble - 1x01 - Episode 1",
			StartTime:  refTime.Add(6 * time.Hour),
			EndTime:    refTime.Add(7 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		{
			Summary:    "Double Trouble - 1x02 - Episode 2",
			StartTime:  refTime.Add(7 * time.Hour),
			EndTime:    refTime.Add(8 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
	}

	startDate := refTime.Add(-24 * time.Hour)
	endDate := refTime.Add(24 * time.Hour)

	res := svc.Format(events, cfg, startDate, endDate)
	if res.Discord == nil || len(res.Discord.Embeds) == 0 {
		t.Fatalf("Expected discord embeds")
	}

	desc := res.Discord.Embeds[0].Description
	lines := strings.Split(desc, "\n")

	// Expected lines:
	// 1. **Our Sticky Love** — <t:...:R> 🎉
	// 2. **A Bona Fide Killer** - 1x03 - *Episode 3* — <t:...:R>
	// 3. **Flex x Cop** - 2x01 - *Episode 1* — <t:...:R> 🎉
	// 4. **Double Trouble** - 1x01 - *Episode 1* — <t:...:R> 🎉
	// 5. **Double Trouble** - 1x02 - *Episode 2* — <t:...:R>

	if len(lines) != 5 {
		t.Fatalf("Expected 5 lines in description, got %d:\n%s", len(lines), desc)
	}

	// 1. Bulk Our Sticky Love
	if !strings.HasPrefix(lines[0], "**Our Sticky Love** — ") {
		t.Errorf("Line 1 should be bulk format for Our Sticky Love, got: %s", lines[0])
	}
	if !strings.HasSuffix(lines[0], "🎉") {
		t.Errorf("Line 1 should end with premiere emoji 🎉, got: %s", lines[0])
	}
	if strings.Contains(lines[0], "1x01") || strings.Contains(lines[0], "Episode 1") {
		t.Errorf("Line 1 should NOT contain episode info, got: %s", lines[0])
	}

	// 2. A Bona Fide Killer
	if !strings.Contains(lines[1], "**A Bona Fide Killer** - 1x03 - *Episode 3*") {
		t.Errorf("Line 2 unexpected: %s", lines[1])
	}
	if strings.Contains(lines[1], "🎉") {
		t.Errorf("Line 2 should NOT have premiere emoji, got: %s", lines[1])
	}

	// 3. Flex x Cop (Season Premiere)
	if !strings.Contains(lines[2], "**Flex x Cop** - 2x01 - *Episode 1*") {
		t.Errorf("Line 3 unexpected: %s", lines[2])
	}
	if !strings.HasSuffix(lines[2], "🎉") {
		t.Errorf("Line 3 should have premiere emoji 🎉, got: %s", lines[2])
	}

	// 4 & 5. Double Trouble (2 episodes <= 2, not bulk)
	if !strings.Contains(lines[3], "**Double Trouble** - 1x01 - *Episode 1*") {
		t.Errorf("Line 4 should be individual episode 1: %s", lines[3])
	}
	if !strings.Contains(lines[4], "**Double Trouble** - 1x02 - *Episode 2*") {
		t.Errorf("Line 5 should be individual episode 2: %s", lines[4])
	}
}
