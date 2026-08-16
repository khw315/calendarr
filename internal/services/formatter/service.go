package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/khw315/calendarr/internal/constants"
	"github.com/khw315/calendarr/internal/localization"
	"github.com/khw315/calendarr/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type FormattedResult struct {
	Discord *models.DiscordPayload
	Slack   *models.SlackPayload
	Events  []*models.Event
	Counts  map[string]int
}

func (s *Service) Format(events []*models.Event, cfg *models.Config, startDate, endDate time.Time) *FormattedResult {
	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)

	// Filter & count
	var activeEvents []*models.Event
	tvCount := 0
	movieCount := 0

	for _, ev := range events {
		if ev.IsPast(now) && cfg.PassedEventHandling == constants.PassedEventHide {
			continue
		}
		if ev.IsTV() {
			tvCount++
		} else if ev.IsMovie() {
			movieCount++
		}
		activeEvents = append(activeEvents, ev)
	}

	counts := map[string]int{
		"tv_count":    tvCount,
		"movie_count": movieCount,
		"total_count": len(activeEvents),
	}

	discordPayload := s.buildDiscordPayload(activeEvents, cfg, startDate, endDate, tvCount, movieCount, now)
	slackPayload := s.buildSlackPayload(activeEvents, cfg, startDate, endDate, tvCount, movieCount, now)

	return &FormattedResult{
		Discord: discordPayload,
		Slack:   slackPayload,
		Events:  activeEvents,
		Counts:  counts,
	}
}

func (s *Service) buildDiscordPayload(events []*models.Event, cfg *models.Config, startDate, endDate time.Time, tvCount, movieCount int, now time.Time) *models.DiscordPayload {
	lang := cfg.Language

	title := localization.GetText(lang, "header_text")
	if title == "" {
		title = localization.GetText(lang, "title")
	}
	if title == "" {
		title = "Rilis Baru"
	}

	// Role mention header
	content := ""
	if cfg.DiscordMentionRoleID != "" {
		content = fmt.Sprintf("<@&%s>", cfg.DiscordMentionRoleID)
	}

	if len(events) == 0 {
		return &models.DiscordPayload{
			Content: content,
			Embeds: []models.DiscordEmbed{
				{
					Title:       title,
					Description: localization.GetRandomMessage(lang, "no_new_releases"),
					Color:       constants.DiscordColors["blue"],
				},
			},
		}
	}

	// Subheader
	subheader := localization.FormatSubheader(lang, tvCount, movieCount)
	if cfg.ShowTimezoneInSubheader && cfg.Timezone != "" {
		subheader += fmt.Sprintf(" (%s)", cfg.Timezone)
	}

	embed := models.DiscordEmbed{
		Title:       title,
		Description: subheader,
		Color:       constants.DiscordColors["blue"],
	}

	// Group by date
	grouped := make(map[string][]*models.Event)
	var dates []string

	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	for _, ev := range events {
		dayKey := ev.DayKey(loc)
		if _, exists := grouped[dayKey]; !exists {
			dates = append(dates, dayKey)
		}
		grouped[dayKey] = append(grouped[dayKey], ev)
	}

	for _, dayKey := range dates {
		dayEvents := grouped[dayKey]
		var lines []string

		for _, ev := range dayEvents {
			line := s.formatEventLine(ev, cfg, true, now)
			lines = append(lines, line)
		}

		dateHeader := localization.FormatDateHeader(dayEvents[0].StartTime.In(loc), lang)

		embed.Fields = append(embed.Fields, models.DiscordField{
			Name:   fmt.Sprintf("📅 %s", dateHeader),
			Value:  strings.Join(lines, "\n"),
			Inline: false,
		})
	}

	return &models.DiscordPayload{
		Content: content,
		Embeds:  []models.DiscordEmbed{embed},
	}
}

func (s *Service) buildSlackPayload(events []*models.Event, cfg *models.Config, startDate, endDate time.Time, tvCount, movieCount int, now time.Time) *models.SlackPayload {
	lang := cfg.Language

	title := localization.GetText(lang, "header_text")
	if title == "" {
		title = localization.GetText(lang, "title")
	}
	if title == "" {
		title = "Rilis Baru"
	}

	if len(events) == 0 {
		return &models.SlackPayload{
			Blocks: []models.SlackBlock{
				{
					Type: "header",
					Text: &models.SlackText{
						Type: "plain_text",
						Text: title,
					},
				},
				{
					Type: "section",
					Text: &models.SlackText{
						Type: "mrkdwn",
						Text: localization.GetRandomMessage(lang, "no_new_releases"),
					},
				},
			},
		}
	}

	subheader := localization.FormatSubheader(lang, tvCount, movieCount)
	if cfg.ShowTimezoneInSubheader && cfg.Timezone != "" {
		subheader += fmt.Sprintf(" (%s)", cfg.Timezone)
	}

	blocks := []models.SlackBlock{
		{
			Type: "header",
			Text: &models.SlackText{
				Type: "plain_text",
				Text: title,
			},
		},
		{
			Type: "section",
			Text: &models.SlackText{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*%s*", subheader),
			},
		},
		{
			Type: "divider",
		},
	}

	// Group by day
	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	grouped := make(map[string][]*models.Event)
	var dates []string

	for _, ev := range events {
		dayKey := ev.DayKey(loc)
		if _, exists := grouped[dayKey]; !exists {
			dates = append(dates, dayKey)
		}
		grouped[dayKey] = append(grouped[dayKey], ev)
	}

	for _, dayKey := range dates {
		dayEvents := grouped[dayKey]
		var lines []string

		for _, ev := range dayEvents {
			line := s.formatEventLine(ev, cfg, false, now)
			lines = append(lines, line)
		}

		dateHeader := localization.FormatDateHeader(dayEvents[0].StartTime.In(loc), lang)

		blocks = append(blocks, models.SlackBlock{
			Type: "section",
			Text: &models.SlackText{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*📅 %s*\n%s", dateHeader, strings.Join(lines, "\n")),
			},
		})
	}

	return &models.SlackPayload{
		Blocks: blocks,
	}
}

func (s *Service) formatEventLine(ev *models.Event, cfg *models.Config, isDiscord bool, now time.Time) string {
	timePrefix := s.formatTimePrefix(ev, cfg)
	icon := "📺"
	if ev.IsMovie() {
		icon = "🎬"
	}

	title := s.formatTitleMarkup(ev.Summary, ev.IsPast(now), cfg.PassedEventHandling, isDiscord)
	return fmt.Sprintf("%s%s %s", timePrefix, icon, title)
}

func (s *Service) formatTimePrefix(ev *models.Event, cfg *models.Config) string {
	if !cfg.TimeSettings.DisplayTime {
		return ""
	}
	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}
	t := ev.StartTime.In(loc)
	if cfg.TimeSettings.Use24Hour {
		return fmt.Sprintf("`%s` ", t.Format("15:04"))
	}
	return fmt.Sprintf("`%s` ", t.Format("03:04 PM"))
}

func (s *Service) formatTitleMarkup(summary string, isPast bool, handling string, isDiscord bool) string {
	if isPast && handling == constants.PassedEventStrike {
		if isDiscord {
			return fmt.Sprintf("~~%s~~", summary)
		}
		return fmt.Sprintf("~%s~", summary)
	}
	if isDiscord {
		return fmt.Sprintf("**%s**", summary)
	}
	return fmt.Sprintf("*%s*", summary)
}
