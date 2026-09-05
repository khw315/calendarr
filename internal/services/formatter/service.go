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

const (
	timeFormat12Hour = "03:04 PM"
	timeFormat24Hour = "15:04"
)

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

	discordPayload := s.buildDiscordPayload(activeEvents, cfg, tvCount, movieCount, now)
	slackPayload := s.buildSlackPayload(activeEvents, cfg, tvCount, movieCount, now)

	return &FormattedResult{
		Discord: discordPayload,
		Slack:   slackPayload,
		Events:  activeEvents,
		Counts:  counts,
	}
}

func (s *Service) buildDiscordPayload(events []*models.Event, cfg *models.Config, tvCount, movieCount int, now time.Time) *models.DiscordPayload {
	lang := cfg.Language
	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	_, content := s.buildDiscordHeaderContent(lang, tvCount, movieCount, cfg)

	if len(events) == 0 {
		return &models.DiscordPayload{
			Content: content,
			Embeds:  nil,
		}
	}

	embeds := s.buildDiscordEmbedsByDay(events, cfg, now, loc, lang)

	return &models.DiscordPayload{
		Content: content,
		Embeds:  embeds,
	}
}

func (s *Service) buildDiscordHeaderContent(lang string, tvCount, movieCount int, cfg *models.Config) (string, string) {
	title := localization.GetText(lang, "header_text")
	if title == "" {
		title = localization.GetText(lang, "title")
	}
	if title == "" {
		title = "Rilis Baru"
	}

	var contentParts []string
	if tvCount == 0 && movieCount == 0 {
		emptyMessage := localization.GetRandomMessage(lang, "no_new_releases")
		contentParts = append(contentParts, fmt.Sprintf("# %s", emptyMessage))
	} else {
		contentParts = append(contentParts, fmt.Sprintf("# %s", title))
		subheader := localization.FormatSubheader(lang, tvCount, movieCount)
		if subheader != "" {
			contentParts = append(contentParts, fmt.Sprintf("## %s", subheader))
		}
	}
	if cfg.ShowTimezoneInSubheader && cfg.Timezone != "" {
		contentParts = append(contentParts, fmt.Sprintf("(%s)", cfg.Timezone))
	}
	if cfg.DiscordMentionRoleID != "" {
		mentionStr := fmt.Sprintf("<@&%s>", cfg.DiscordMentionRoleID)
		if !cfg.DiscordHideMentionInstructions {
			instruction := localization.GetText(lang, "mention_instruction")
			if instruction != "" {
				mentionStr += fmt.Sprintf("\n*%s*", instruction)
			}
		}
		contentParts = append(contentParts, mentionStr)
	}

	return title, strings.Join(contentParts, "\n\n")
}

func (s *Service) buildDiscordEmbedsByDay(events []*models.Event, cfg *models.Config, now time.Time, loc *time.Location, lang string) []models.DiscordEmbed {
	grouped := make(map[string][]*models.Event)
	var dates []string

	for _, ev := range events {
		dayKey := ev.DayKey(loc)
		if _, exists := grouped[dayKey]; !exists {
			dates = append(dates, dayKey)
		}
		grouped[dayKey] = append(grouped[dayKey], ev)
	}

	var embeds []models.DiscordEmbed
	for _, dayKey := range dates {
		dayEvents := grouped[dayKey]
		var lines []string

		for _, ev := range dayEvents {
			line := s.formatEventLineDiscord(ev, cfg, now, loc)
			lines = append(lines, line)
		}

		dayStartTime := dayEvents[0].StartTime.In(loc)
		dateHeader := localization.FormatDateHeader(dayStartTime, lang)
		embedColor := constants.DayColorsDiscord[dayStartTime.Weekday()]
		if embedColor == 0 {
			embedColor = constants.DiscordColors["blue"]
		}

		embeds = append(embeds, models.DiscordEmbed{
			Title:       dateHeader,
			Description: strings.Join(lines, "\n"),
			Color:       embedColor,
		})
	}
	return embeds
}

func (s *Service) buildSlackPayload(events []*models.Event, cfg *models.Config, tvCount, movieCount int, now time.Time) *models.SlackPayload {
	lang := cfg.Language
	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

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
			line := s.formatEventLineSlack(ev, cfg, now, loc)
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

func (s *Service) formatEventLineDiscord(ev *models.Event, cfg *models.Config, now time.Time, loc *time.Location) string {
	if ev.IsMovie() {
		return s.formatMovieEventDiscord(ev, cfg, now, loc)
	}
	return s.formatTVEventDiscord(ev, cfg, now, loc)
}

func (s *Service) formatTVEventDiscord(ev *models.Event, cfg *models.Config, now time.Time, loc *time.Location) string {
	summary := ev.Summary
	showName := summary
	epNum := ""
	epTitle := ""

	parts := strings.Split(summary, " - ")
	if len(parts) >= 3 {
		showName = parts[0]
		epNum = parts[1]
		epTitle = strings.Join(parts[2:], " - ")
	} else if len(parts) == 2 {
		showName = parts[0]
		epNum = parts[1]
	}

	formattedShow := fmt.Sprintf("**%s**", showName)
	epDetails := ""
	if epTitle != "" {
		if epNum != "" {
			epDetails = fmt.Sprintf(" - %s - *%s*", epNum, epTitle)
		} else {
			epDetails = fmt.Sprintf(" - *%s*", epTitle)
		}
	} else if epNum != "" {
		epDetails = fmt.Sprintf(" - %s", epNum)
	}

	timeStr := ""
	t := ev.StartTime.In(loc)
	if !t.Before(now) {
		timeStr = fmt.Sprintf(" — <t:%d:R>", t.Unix())
	} else if cfg.TimeSettings.DisplayTime {
		if cfg.TimeSettings.Use24Hour {
			timeStr = fmt.Sprintf(" — %s", t.Format(timeFormat24Hour))
		} else {
			timeStr = fmt.Sprintf(" — %s", t.Format(timeFormat12Hour))
		}
	}

	line := fmt.Sprintf("%s%s%s", formattedShow, epDetails, timeStr)
	if ev.IsPast(now) && cfg.PassedEventHandling == constants.PassedEventStrike {
		line = fmt.Sprintf("~~%s~~", line)
	}
	return line
}

func (s *Service) formatMovieEventDiscord(ev *models.Event, cfg *models.Config, now time.Time, loc *time.Location) string {
	t := ev.StartTime.In(loc)
	timeStr := ""
	if !t.Before(now) {
		timeStr = fmt.Sprintf(" — <t:%d:R>", t.Unix())
	} else if cfg.TimeSettings.DisplayTime {
		if cfg.TimeSettings.Use24Hour {
			timeStr = fmt.Sprintf(" — %s", t.Format(timeFormat24Hour))
		} else {
			timeStr = fmt.Sprintf(" — %s", t.Format(timeFormat12Hour))
		}
	}

	line := fmt.Sprintf("🎬  **%s**%s", ev.Summary, timeStr)
	if ev.IsPast(now) && cfg.PassedEventHandling == constants.PassedEventStrike {
		line = fmt.Sprintf("~~%s~~", line)
	}
	return line
}

func (s *Service) formatEventLineSlack(ev *models.Event, cfg *models.Config, now time.Time, loc *time.Location) string {
	timePrefix := ""
	if cfg.TimeSettings.DisplayTime {
		t := ev.StartTime.In(loc)
		if cfg.TimeSettings.Use24Hour {
			timePrefix = fmt.Sprintf("`%s` ", t.Format(timeFormat24Hour))
		} else {
			timePrefix = fmt.Sprintf("`%s` ", t.Format(timeFormat12Hour))
		}
	}

	icon := "📺"
	if ev.IsMovie() {
		icon = "🎬"
	}

	title := ev.Summary
	if ev.IsPast(now) && cfg.PassedEventHandling == constants.PassedEventStrike {
		title = fmt.Sprintf("~%s~", title)
	} else {
		title = fmt.Sprintf("*%s*", title)
	}

	return fmt.Sprintf("%s%s %s", timePrefix, icon, title)
}
