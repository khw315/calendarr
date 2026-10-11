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

	if res.DiscordTV == nil || len(res.DiscordTV.Embeds) == 0 {
		t.Errorf("Expected DiscordTV payload with embeds")
	}
	if res.DiscordMovie == nil || len(res.DiscordMovie.Embeds) == 0 {
		t.Errorf("Expected DiscordMovie payload with embeds")
	}
	if res.SlackTV == nil || len(res.SlackTV.Blocks) == 0 {
		t.Errorf("Expected SlackTV payload with blocks")
	}
	if res.SlackMovie == nil || len(res.SlackMovie.Blocks) == 0 {
		t.Errorf("Expected SlackMovie payload with blocks")
	}

	if res.Ntfy == nil || res.Ntfy.Message == "" {
		t.Errorf("Expected Ntfy payload with message")
	}
	if res.NtfyTV == nil || res.NtfyTV.Message == "" {
		t.Errorf("Expected NtfyTV payload with message")
	}
	if res.NtfyMovie == nil || res.NtfyMovie.Message == "" {
		t.Errorf("Expected NtfyMovie payload with message")
	}

	if res.Telegram == nil || res.Telegram.Text == "" {
		t.Errorf("Expected Telegram payload with text")
	}
	if res.TelegramTV == nil || res.TelegramTV.Text == "" {
		t.Errorf("Expected TelegramTV payload with text")
	}
	if res.TelegramMovie == nil || res.TelegramMovie.Text == "" {
		t.Errorf("Expected TelegramMovie payload with text")
	}

	notifPayloads := res.ToNotificationPayloads()
	if notifPayloads.DiscordCombined == nil || notifPayloads.DiscordTV == nil || notifPayloads.DiscordMovie == nil {
		t.Errorf("Expected non-nil Discord payloads in ToNotificationPayloads")
	}
	if notifPayloads.SlackCombined == nil || notifPayloads.SlackTV == nil || notifPayloads.SlackMovie == nil {
		t.Errorf("Expected non-nil Slack payloads in ToNotificationPayloads")
	}
	if notifPayloads.NtfyCombined == nil || notifPayloads.NtfyTV == nil || notifPayloads.NtfyMovie == nil {
		t.Errorf("Expected non-nil Ntfy payloads in ToNotificationPayloads")
	}
	if notifPayloads.TelegramCombined == nil || notifPayloads.TelegramTV == nil || notifPayloads.TelegramMovie == nil {
		t.Errorf("Expected non-nil Telegram payloads in ToNotificationPayloads")
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
	if emptyRes.Ntfy == nil || emptyRes.Ntfy.Message == "" {
		t.Errorf("Expected fallback Ntfy message for empty events")
	}
	if emptyRes.Telegram == nil || emptyRes.Telegram.Text == "" {
		t.Errorf("Expected fallback Telegram text for empty events")
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

func TestBulkAndPremiereNtfyFormatting(t *testing.T) {
	svc := NewService()
	cfg := models.DefaultConfig()
	loc := time.UTC
	cfg.TimezoneLocation = loc
	cfg.TimeSettings.DisplayTime = false

	refTime := time.Date(2026, time.August, 7, 12, 0, 0, 0, loc)

	events := []*models.Event{
		// Show with 3 episodes -> bulk grouped
		{
			Summary:    "Our Sticky Love - 1x01 - Episode 1",
			StartTime:  refTime,
			SourceType: constants.EventTypeTV,
		},
		{
			Summary:    "Our Sticky Love - 1x02 - Episode 2",
			StartTime:  refTime,
			SourceType: constants.EventTypeTV,
		},
		{
			Summary:    "Our Sticky Love - 1x03 - Episode 3",
			StartTime:  refTime,
			SourceType: constants.EventTypeTV,
		},
		// Season premiere
		{
			Summary:    "Flex x Cop - 2x01 - Episode 1",
			StartTime:  refTime.Add(1 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		// Movie
		{
			Summary:    "Dune: Part Two",
			StartTime:  refTime.Add(2 * time.Hour),
			SourceType: constants.EventTypeMovie,
		},
	}

	res := svc.Format(events, cfg, refTime.Add(-24*time.Hour), refTime.Add(24*time.Hour))
	if res.Ntfy == nil || res.Ntfy.Message == "" {
		t.Fatalf("Expected non-empty Ntfy message")
	}

	msg := res.Ntfy.Message
	if !strings.Contains(msg, "- **Our Sticky Love** (3 episodes) 🎉") {
		t.Errorf("Expected bulk Ntfy line with premiere emoji for Our Sticky Love, got:\n%s", msg)
	}
	if !strings.Contains(msg, "- **Flex x Cop** - 2x01 - *Episode 1* 🎉") {
		t.Errorf("Expected TV premiere line for Flex x Cop, got:\n%s", msg)
	}
	if !strings.Contains(msg, "- 🎬 **Dune: Part Two**") {
		t.Errorf("Expected movie line for Dune: Part Two, got:\n%s", msg)
	}
}

func TestBulkAndPremiereTelegramFormatting(t *testing.T) {
	svc := NewService()
	cfg := models.DefaultConfig()
	loc := time.UTC
	cfg.TimezoneLocation = loc
	cfg.TimeSettings.DisplayTime = false

	refTime := time.Date(2026, time.August, 7, 12, 0, 0, 0, loc)

	events := []*models.Event{
		// Show with 3 episodes and special characters -> bulk grouped & escaped
		{
			Summary:    "Tom & Jerry <Special> - 1x01 - Ep 1",
			StartTime:  refTime,
			SourceType: constants.EventTypeTV,
		},
		{
			Summary:    "Tom & Jerry <Special> - 1x02 - Ep 2",
			StartTime:  refTime,
			SourceType: constants.EventTypeTV,
		},
		{
			Summary:    "Tom & Jerry <Special> - 1x03 - Ep 3",
			StartTime:  refTime,
			SourceType: constants.EventTypeTV,
		},
		// Season premiere
		{
			Summary:    "Flex x Cop - 2x01 - Episode 1",
			StartTime:  refTime.Add(1 * time.Hour),
			SourceType: constants.EventTypeTV,
		},
		// Movie with ampersand
		{
			Summary:    "Fast & Furious",
			StartTime:  refTime.Add(2 * time.Hour),
			SourceType: constants.EventTypeMovie,
		},
	}

	res := svc.Format(events, cfg, refTime.Add(-24*time.Hour), refTime.Add(24*time.Hour))
	if res.Telegram == nil || res.Telegram.Text == "" {
		t.Fatalf("Expected non-empty Telegram text")
	}

	text := res.Telegram.Text
	if !strings.Contains(text, "• <b>Tom &amp; Jerry &lt;Special&gt;</b> (3 episodes) 🎉") {
		t.Errorf("Expected bulk Telegram line with HTML escaped characters, got:\n%s", text)
	}
	if !strings.Contains(text, "• <b>Flex x Cop</b> - 2x01 - <i>Episode 1</i> 🎉") {
		t.Errorf("Expected single premiere Telegram line, got:\n%s", text)
	}
	if !strings.Contains(text, "• 🎬 <b>Fast &amp; Furious</b>") {
		t.Errorf("Expected movie line with escaped ampersand, got:\n%s", text)
	}
	if res.Telegram.ParseMode != "HTML" {
		t.Errorf("Expected ParseMode HTML, got: %s", res.Telegram.ParseMode)
	}
}
