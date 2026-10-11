package formatter

import (
	"fmt"
	"regexp"
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
	Discord       *models.DiscordPayload
	DiscordTV     *models.DiscordPayload
	DiscordMovie  *models.DiscordPayload
	Slack         *models.SlackPayload
	SlackTV       *models.SlackPayload
	SlackMovie    *models.SlackPayload
	Ntfy          *models.NtfyPayload
	NtfyTV        *models.NtfyPayload
	NtfyMovie     *models.NtfyPayload
	Telegram      *models.TelegramPayload
	TelegramTV    *models.TelegramPayload
	TelegramMovie *models.TelegramPayload
	Events        []*models.Event
	Counts        map[string]int
}

func (r *FormattedResult) ToNotificationPayloads() models.NotificationPayloads {
	return models.NotificationPayloads{
		DiscordCombined:  r.Discord,
		DiscordTV:        r.DiscordTV,
		DiscordMovie:     r.DiscordMovie,
		SlackCombined:    r.Slack,
		SlackTV:          r.SlackTV,
		SlackMovie:       r.SlackMovie,
		NtfyCombined:     r.Ntfy,
		NtfyTV:           r.NtfyTV,
		NtfyMovie:        r.NtfyMovie,
		TelegramCombined: r.Telegram,
		TelegramTV:       r.TelegramTV,
		TelegramMovie:    r.TelegramMovie,
	}
}

func (s *Service) Format(events []*models.Event, cfg *models.Config, startDate, endDate time.Time) *FormattedResult {
	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)

	// Filter & count
	var activeEvents []*models.Event
	var tvEvents []*models.Event
	var movieEvents []*models.Event
	tvCount := 0
	movieCount := 0

	for _, ev := range events {
		if ev.IsPast(now) && cfg.PassedEventHandling == constants.PassedEventHide {
			continue
		}
		if ev.IsTV() {
			tvCount++
			tvEvents = append(tvEvents, ev)
		} else if ev.IsMovie() {
			movieCount++
			movieEvents = append(movieEvents, ev)
		}
		activeEvents = append(activeEvents, ev)
	}

	counts := map[string]int{
		"tv_count":    tvCount,
		"movie_count": movieCount,
		"total_count": len(activeEvents),
	}

	discordPayload := s.buildDiscordPayload(activeEvents, cfg, tvCount, movieCount, now)
	discordTVPayload := s.buildDiscordPayload(tvEvents, cfg, tvCount, 0, now)
	discordMoviePayload := s.buildDiscordPayload(movieEvents, cfg, 0, movieCount, now)

	slackPayload := s.buildSlackPayload(activeEvents, cfg, tvCount, movieCount, now)
	slackTVPayload := s.buildSlackPayload(tvEvents, cfg, tvCount, 0, now)
	slackMoviePayload := s.buildSlackPayload(movieEvents, cfg, 0, movieCount, now)

	ntfyPayload := s.buildNtfyPayload(activeEvents, cfg, tvCount, movieCount, now)
	ntfyTVPayload := s.buildNtfyPayload(tvEvents, cfg, tvCount, 0, now)
	ntfyMoviePayload := s.buildNtfyPayload(movieEvents, cfg, 0, movieCount, now)

	telegramPayload := s.buildTelegramPayload(activeEvents, cfg, tvCount, movieCount, now)
	telegramTVPayload := s.buildTelegramPayload(tvEvents, cfg, tvCount, 0, now)
	telegramMoviePayload := s.buildTelegramPayload(movieEvents, cfg, 0, movieCount, now)

	return &FormattedResult{
		Discord:       discordPayload,
		DiscordTV:     discordTVPayload,
		DiscordMovie:  discordMoviePayload,
		Slack:         slackPayload,
		SlackTV:       slackTVPayload,
		SlackMovie:    slackMoviePayload,
		Ntfy:          ntfyPayload,
		NtfyTV:        ntfyTVPayload,
		NtfyMovie:     ntfyMoviePayload,
		Telegram:      telegramPayload,
		TelegramTV:    telegramTVPayload,
		TelegramMovie: telegramMoviePayload,
		Events:        activeEvents,
		Counts:        counts,
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

type showGroup struct {
	showName string
	events   []*models.Event
}

func groupEventsByDay(events []*models.Event, loc *time.Location) ([]string, map[string][]*models.Event) {
	grouped := make(map[string][]*models.Event)
	var dates []string

	for _, ev := range events {
		dayKey := ev.DayKey(loc)
		if _, exists := grouped[dayKey]; !exists {
			dates = append(dates, dayKey)
		}
		grouped[dayKey] = append(grouped[dayKey], ev)
	}
	return dates, grouped
}

func groupTVShows(dayEvents []*models.Event) map[string]*showGroup {
	tvShowMap := make(map[string]*showGroup)
	for _, ev := range dayEvents {
		if ev.IsMovie() {
			continue
		}
		showName, _, _ := parseTVEventSummary(ev.Summary)
		norm := strings.ToLower(showName)
		if group, exists := tvShowMap[norm]; exists {
			group.events = append(group.events, ev)
		} else {
			tvShowMap[norm] = &showGroup{
				showName: showName,
				events:   []*models.Event{ev},
			}
		}
	}
	return tvShowMap
}

func (s *Service) formatDayDiscordLines(dayEvents []*models.Event, cfg *models.Config, now time.Time, loc *time.Location) []string {
	tvShowMap := groupTVShows(dayEvents)
	var lines []string
	renderedBulkShows := make(map[string]bool)

	for _, ev := range dayEvents {
		if ev.IsMovie() {
			lines = append(lines, s.formatMovieEventDiscord(ev, cfg, now, loc))
			continue
		}

		showName, _, _ := parseTVEventSummary(ev.Summary)
		norm := strings.ToLower(showName)
		group := tvShowMap[norm]

		if group != nil && len(group.events) > constants.BulkThresholdDiscord {
			if !renderedBulkShows[norm] {
				renderedBulkShows[norm] = true
				lines = append(lines, s.formatBulkTVEventDiscord(group.showName, group.events, cfg, now, loc))
			}
			continue
		}

		lines = append(lines, s.formatTVEventDiscord(ev, cfg, now, loc))
	}
	return lines
}

func (s *Service) buildDiscordEmbedsByDay(events []*models.Event, cfg *models.Config, now time.Time, loc *time.Location, lang string) []models.DiscordEmbed {
	dates, grouped := groupEventsByDay(events, loc)
	var embeds []models.DiscordEmbed

	for _, dayKey := range dates {
		dayEvents := grouped[dayKey]
		lines := s.formatDayDiscordLines(dayEvents, cfg, now, loc)

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

	dates, grouped := groupEventsByDay(events, loc)

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

var tvShowSeasonEpRegex = regexp.MustCompile(`(?i)\s+(?:s\d+e\d+|\d+x\d+)\b`)

func parseTVEventSummary(summary string) (showName, epNum, epTitle string) {
	parts := strings.Split(summary, " - ")
	if len(parts) >= 3 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), strings.TrimSpace(strings.Join(parts[2:], " - "))
	} else if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), ""
	}

	if loc := tvShowSeasonEpRegex.FindStringIndex(summary); loc != nil {
		show := strings.TrimSpace(summary[:loc[0]])
		ep := strings.TrimSpace(summary[loc[0]:])
		return show, ep, ""
	}

	return strings.TrimSpace(summary), "", ""
}

func (s *Service) formatTimeDiscord(startTime time.Time, cfg *models.Config, now time.Time, loc *time.Location) string {
	t := startTime.In(loc)
	timeStyle := cfg.DiscordTimestampStyle
	if timeStyle == "" && !t.Before(now) {
		timeStyle = "R"
	}

	if timeStyle != "" {
		return fmt.Sprintf(" — <t:%d:%s>", t.Unix(), timeStyle)
	} else if cfg.TimeSettings.DisplayTime {
		if cfg.TimeSettings.Use24Hour {
			return fmt.Sprintf(" — %s", t.Format(timeFormat24Hour))
		}
		return fmt.Sprintf(" — %s", t.Format(timeFormat12Hour))
	}
	return ""
}

func (s *Service) formatBulkTVEventDiscord(showName string, group []*models.Event, cfg *models.Config, now time.Time, loc *time.Location) string {
	if len(group) == 0 {
		return ""
	}
	first := group[0]
	formattedShow := fmt.Sprintf("**%s**", showName)
	timeStr := s.formatTimeDiscord(first.StartTime, cfg, now, loc)

	isPremiere := false
	allPast := true
	for _, ev := range group {
		if ev.IsPremiere() {
			isPremiere = true
		}
		if !ev.IsPast(now) {
			allPast = false
		}
	}

	line := fmt.Sprintf("%s%s", formattedShow, timeStr)
	if isPremiere {
		line += " 🎉"
	}
	if allPast && cfg.PassedEventHandling == constants.PassedEventStrike {
		line = fmt.Sprintf("~~%s~~", line)
	}
	return line
}

func (s *Service) formatTVEventDiscord(ev *models.Event, cfg *models.Config, now time.Time, loc *time.Location) string {
	showName, epNum, epTitle := parseTVEventSummary(ev.Summary)

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

	timeStr := s.formatTimeDiscord(ev.StartTime, cfg, now, loc)

	line := fmt.Sprintf("%s%s%s", formattedShow, epDetails, timeStr)
	if ev.IsPremiere() {
		line += " 🎉"
	}
	if ev.IsPast(now) && cfg.PassedEventHandling == constants.PassedEventStrike {
		line = fmt.Sprintf("~~%s~~", line)
	}
	return line
}

func (s *Service) formatMovieEventDiscord(ev *models.Event, cfg *models.Config, now time.Time, loc *time.Location) string {
	timeStr := s.formatTimeDiscord(ev.StartTime, cfg, now, loc)

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

func escapeHTML(s string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	).Replace(s)
}

func (s *Service) buildNtfyPayload(events []*models.Event, cfg *models.Config, tvCount, movieCount int, now time.Time) *models.NtfyPayload {
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
		title = "New Releases"
	}

	priority := cfg.NtfyPriority
	if priority <= 0 {
		priority = 3
	}

	if len(events) == 0 {
		return &models.NtfyPayload{
			Title:       title,
			Message:     localization.GetRandomMessage(lang, "no_new_releases"),
			Priority:    priority,
			Tags:        []string{"calendar"},
			HasReleases: false,
		}
	}

	subheader := localization.FormatSubheader(lang, tvCount, movieCount)
	if cfg.ShowTimezoneInSubheader && cfg.Timezone != "" {
		subheader += fmt.Sprintf(" (%s)", cfg.Timezone)
	}

	tags := []string{"calendar"}
	if tvCount > 0 {
		tags = append(tags, "tv")
	}
	if movieCount > 0 {
		tags = append(tags, "film")
	}

	for _, ev := range events {
		if ev.IsPremiere() {
			tags = append(tags, "tada")
			break
		}
	}

	dates, grouped := groupEventsByDay(events, loc)
	var sections []string
	sections = append(sections, fmt.Sprintf("**%s**", subheader))

	for _, dayKey := range dates {
		dayEvents := grouped[dayKey]
		lines := s.formatDayNtfyLines(dayEvents, cfg, loc)
		dateHeader := localization.FormatDateHeader(dayEvents[0].StartTime.In(loc), lang)
		sections = append(sections, fmt.Sprintf("### 🗓️ %s\n%s", dateHeader, strings.Join(lines, "\n")))
	}

	return &models.NtfyPayload{
		Title:       title,
		Message:     strings.Join(sections, "\n\n"),
		Priority:    priority,
		Tags:        tags,
		HasReleases: true,
	}
}

func (s *Service) formatDayNtfyLines(dayEvents []*models.Event, cfg *models.Config, loc *time.Location) []string {
	tvShowMap := groupTVShows(dayEvents)
	var lines []string
	renderedBulkShows := make(map[string]bool)

	for _, ev := range dayEvents {
		if ev.IsMovie() {
			lines = append(lines, s.formatMovieEventNtfy(ev, cfg, loc))
			continue
		}

		showName, _, _ := parseTVEventSummary(ev.Summary)
		norm := strings.ToLower(showName)
		group := tvShowMap[norm]

		if group != nil && len(group.events) > constants.BulkThresholdDiscord {
			if !renderedBulkShows[norm] {
				renderedBulkShows[norm] = true
				lines = append(lines, s.formatBulkTVEventNtfy(group.showName, group.events, cfg, loc))
			}
			continue
		}

		lines = append(lines, s.formatTVEventNtfy(ev, cfg, loc))
	}
	return lines
}

func (s *Service) formatBulkTVEventNtfy(showName string, group []*models.Event, cfg *models.Config, loc *time.Location) string {
	if len(group) == 0 {
		return ""
	}
	first := group[0]
	timeStr := s.formatEventTimeSimple(first.StartTime, cfg, loc)
	isPremiere := false
	for _, ev := range group {
		if ev.IsPremiere() {
			isPremiere = true
			break
		}
	}

	line := fmt.Sprintf("- **%s** (%d episodes)", showName, len(group))
	if timeStr != "" {
		line += fmt.Sprintf(" `%s`", timeStr)
	}
	if isPremiere {
		line += " 🎉"
	}
	return line
}

func (s *Service) formatTVEventNtfy(ev *models.Event, cfg *models.Config, loc *time.Location) string {
	showName, epNum, epTitle := parseTVEventSummary(ev.Summary)
	timeStr := s.formatEventTimeSimple(ev.StartTime, cfg, loc)

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

	line := fmt.Sprintf("- **%s**%s", showName, epDetails)
	if timeStr != "" {
		line += fmt.Sprintf(" `%s`", timeStr)
	}
	if ev.IsPremiere() {
		line += " 🎉"
	}
	return line
}

func (s *Service) formatMovieEventNtfy(ev *models.Event, cfg *models.Config, loc *time.Location) string {
	timeStr := s.formatEventTimeSimple(ev.StartTime, cfg, loc)
	line := fmt.Sprintf("- 🎬 **%s**", ev.Summary)
	if timeStr != "" {
		line += fmt.Sprintf(" `%s`", timeStr)
	}
	return line
}

func (s *Service) buildTelegramPayload(events []*models.Event, cfg *models.Config, tvCount, movieCount int, now time.Time) *models.TelegramPayload {
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
		title = "New Releases"
	}

	if len(events) == 0 {
		text := fmt.Sprintf("<b>📅 %s</b>\n\n<i>%s</i>", escapeHTML(title), escapeHTML(localization.GetRandomMessage(lang, "no_new_releases")))
		return &models.TelegramPayload{
			Text:                  text,
			ParseMode:             "HTML",
			DisableNotification:   cfg.TelegramSilentNotification,
			DisableWebPagePreview: true,
			HasReleases:           false,
		}
	}

	subheader := localization.FormatSubheader(lang, tvCount, movieCount)
	if cfg.ShowTimezoneInSubheader && cfg.Timezone != "" {
		subheader += fmt.Sprintf(" (%s)", cfg.Timezone)
	}

	var sections []string
	sections = append(sections, fmt.Sprintf("<b>📅 %s</b>\n<i>%s</i>", escapeHTML(title), escapeHTML(subheader)))

	dates, grouped := groupEventsByDay(events, loc)
	for _, dayKey := range dates {
		dayEvents := grouped[dayKey]
		lines := s.formatDayTelegramLines(dayEvents, cfg, loc)
		dateHeader := localization.FormatDateHeader(dayEvents[0].StartTime.In(loc), lang)
		sections = append(sections, fmt.Sprintf("<b>🗓️ %s</b>\n%s", escapeHTML(dateHeader), strings.Join(lines, "\n")))
	}

	return &models.TelegramPayload{
		Text:                  strings.Join(sections, "\n\n"),
		ParseMode:             "HTML",
		DisableNotification:   cfg.TelegramSilentNotification,
		DisableWebPagePreview: true,
		HasReleases:           true,
	}
}

func (s *Service) formatDayTelegramLines(dayEvents []*models.Event, cfg *models.Config, loc *time.Location) []string {
	tvShowMap := groupTVShows(dayEvents)
	var lines []string
	renderedBulkShows := make(map[string]bool)

	for _, ev := range dayEvents {
		if ev.IsMovie() {
			lines = append(lines, s.formatMovieEventTelegram(ev, cfg, loc))
			continue
		}

		showName, _, _ := parseTVEventSummary(ev.Summary)
		norm := strings.ToLower(showName)
		group := tvShowMap[norm]

		if group != nil && len(group.events) > constants.BulkThresholdDiscord {
			if !renderedBulkShows[norm] {
				renderedBulkShows[norm] = true
				lines = append(lines, s.formatBulkTVEventTelegram(group.showName, group.events, cfg, loc))
			}
			continue
		}

		lines = append(lines, s.formatTVEventTelegram(ev, cfg, loc))
	}
	return lines
}

func (s *Service) formatBulkTVEventTelegram(showName string, group []*models.Event, cfg *models.Config, loc *time.Location) string {
	if len(group) == 0 {
		return ""
	}
	first := group[0]
	timeStr := s.formatEventTimeSimple(first.StartTime, cfg, loc)
	isPremiere := false
	for _, ev := range group {
		if ev.IsPremiere() {
			isPremiere = true
			break
		}
	}

	line := fmt.Sprintf("• <b>%s</b> (%d episodes)", escapeHTML(showName), len(group))
	if timeStr != "" {
		line += fmt.Sprintf(" <code>(%s)</code>", escapeHTML(timeStr))
	}
	if isPremiere {
		line += " 🎉"
	}
	return line
}

func (s *Service) formatTVEventTelegram(ev *models.Event, cfg *models.Config, loc *time.Location) string {
	showName, epNum, epTitle := parseTVEventSummary(ev.Summary)
	timeStr := s.formatEventTimeSimple(ev.StartTime, cfg, loc)

	epDetails := ""
	if epTitle != "" {
		if epNum != "" {
			epDetails = fmt.Sprintf(" - %s - <i>%s</i>", escapeHTML(epNum), escapeHTML(epTitle))
		} else {
			epDetails = fmt.Sprintf(" - <i>%s</i>", escapeHTML(epTitle))
		}
	} else if epNum != "" {
		epDetails = fmt.Sprintf(" - %s", escapeHTML(epNum))
	}

	line := fmt.Sprintf("• <b>%s</b>%s", escapeHTML(showName), epDetails)
	if timeStr != "" {
		line += fmt.Sprintf(" <code>(%s)</code>", escapeHTML(timeStr))
	}
	if ev.IsPremiere() {
		line += " 🎉"
	}
	return line
}

func (s *Service) formatMovieEventTelegram(ev *models.Event, cfg *models.Config, loc *time.Location) string {
	timeStr := s.formatEventTimeSimple(ev.StartTime, cfg, loc)
	line := fmt.Sprintf("• 🎬 <b>%s</b>", escapeHTML(ev.Summary))
	if timeStr != "" {
		line += fmt.Sprintf(" <code>(%s)</code>", escapeHTML(timeStr))
	}
	return line
}

func (s *Service) formatEventTimeSimple(startTime time.Time, cfg *models.Config, loc *time.Location) string {
	if !cfg.TimeSettings.DisplayTime {
		return ""
	}
	t := startTime.In(loc)
	if cfg.TimeSettings.Use24Hour {
		return t.Format(timeFormat24Hour)
	}
	return t.Format(timeFormat12Hour)
}
