package models

import (
	"encoding/json"
	"time"

	"github.com/khw315/calendarr/internal/constants"
)

type CalendarUrl struct {
	URL  string `json:"url"`
	Type string `json:"type"` // "tv" or "movie"
}

type TimeSettings struct {
	Use24Hour      bool `json:"use_24_hour"`
	AddLeadingZero bool `json:"add_leading_zero"`
	DisplayTime    bool `json:"display_time"`
}

type ScheduleSettings struct {
	ScheduleType string `json:"schedule_type"` // "DAILY", "WEEKLY", "CRON"
	RunTime      string `json:"run_time"`      // "09:00"
	ScheduleDay  string `json:"schedule_day"`  // "1" (Monday) to "7" (Sunday)
	CronSchedule string `json:"cron_schedule,omitempty"`
	RunOnStartup bool   `json:"run_on_startup"`
}

type LoggingSettings struct {
	LogDir      string `json:"log_dir"`
	LogFile     string `json:"log_file"`
	BackupCount int    `json:"backup_count"`
	MaxSizeMB   int    `json:"max_size_mb"`
	DebugMode   bool   `json:"debug_mode"`
}

type Config struct {
	DiscordWebhookURL               string           `json:"discord_webhook_url,omitempty"`
	SlackWebhookURL                 string           `json:"slack_webhook_url,omitempty"`
	UseDiscord                      bool             `json:"use_discord"`
	UseSlack                        bool             `json:"use_slack"`
	ShowDateRange                   bool             `json:"show_date_range"`
	ShowTimezoneInSubheader         bool             `json:"show_timezone_in_subheader"`
	DeduplicateEvents               bool             `json:"deduplicate_events"`
	DiscordMentionRoleID            string           `json:"discord_mention_role_id,omitempty"`
	DiscordHideMentionInstructions bool             `json:"discord_hide_mention_instructions"`
	DiscordTimestampStyle           string           `json:"discord_timestamp_style,omitempty"`
	CalendarURLs                    []CalendarUrl    `json:"calendar_urls"`
	PassedEventHandling             string           `json:"passed_event_handling"`
	TimeSettings                    TimeSettings     `json:"time_settings"`
	ScheduleSettings                ScheduleSettings `json:"schedule_settings"`
	LoggingSettings                 LoggingSettings  `json:"logging_settings"`
	Timezone                        string           `json:"timezone"`
	HTTPTimeout                     int              `json:"http_timeout"`
	EnableCustomDiscordFooter       bool             `json:"enable_custom_discord_footer"`
	EnableCustomSlackFooter         bool             `json:"enable_custom_slack_footer"`
	Language                        string           `json:"language"`

	// Derived fields
	TimezoneLocation *time.Location `json:"-"`
}

func DefaultConfig() *Config {
	return &Config{
		UseDiscord:                     constants.DefaultUseDiscord,
		UseSlack:                       constants.DefaultUseSlack,
		ShowDateRange:                  constants.DefaultShowDateRange,
		ShowTimezoneInSubheader:        constants.DefaultShowTimezoneInSubheader,
		DeduplicateEvents:              constants.DefaultDeduplicateEvents,
		DiscordHideMentionInstructions: constants.DefaultDiscordHideMentionInstructions,
		CalendarURLs:                   make([]CalendarUrl, 0),
		PassedEventHandling:            constants.DefaultPassedEventHandling,
		TimeSettings: TimeSettings{
			Use24Hour:      constants.DefaultUse24Hour,
			AddLeadingZero: constants.DefaultAddLeadingZero,
			DisplayTime:    constants.DefaultDisplayTime,
		},
		ScheduleSettings: ScheduleSettings{
			ScheduleType: constants.DefaultScheduleType,
			RunTime:      constants.DefaultRunTime,
			ScheduleDay:  constants.DefaultScheduleDay,
			RunOnStartup: constants.DefaultRunOnStartup,
		},
		LoggingSettings: LoggingSettings{
			LogDir:      constants.DefaultLogDir,
			LogFile:     constants.DefaultLogFile,
			BackupCount: constants.DefaultLogBackupCount,
			MaxSizeMB:   constants.DefaultLogMaxSizeMB,
			DebugMode:   constants.DefaultDebugMode,
		},
		Timezone:                  "UTC",
		HTTPTimeout:               constants.DefaultHTTPTimeout,
		EnableCustomDiscordFooter: constants.DefaultEnableCustomDiscordFooter,
		EnableCustomSlackFooter:   constants.DefaultEnableCustomSlackFooter,
		Language:                  "EN",
	}
}

func (c *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(c, "", "  ")
}

func (c *Config) MarshalJSON() ([]byte, error) {
	calURLs := c.CalendarURLs
	if calURLs == nil {
		calURLs = make([]CalendarUrl, 0)
	}

	m := map[string]interface{}{
		// FE Uppercase Keys for React Settings.tsx compatibility
		"APP_LANGUAGE":                       c.Language,
		"USE_DISCORD":                        c.UseDiscord,
		"DISCORD_WEBHOOK_URL":                c.DiscordWebhookURL,
		"DISCORD_MENTION_ROLE_ID":            c.DiscordMentionRoleID,
		"DISCORD_HIDE_MENTION_INSTRUCTIONS": c.DiscordHideMentionInstructions,
		"DISCORD_TIMESTAMP_STYLE":           c.DiscordTimestampStyle,
		"ENABLE_CUSTOM_DISCORD_FOOTER":       c.EnableCustomDiscordFooter,
		"USE_SLACK":                          c.UseSlack,
		"SLACK_WEBHOOK_URL":                  c.SlackWebhookURL,
		"ENABLE_CUSTOM_SLACK_FOOTER":         c.EnableCustomSlackFooter,
		"CALENDAR_URLS":                      calURLs,
		"PASSED_EVENT_HANDLING":             c.PassedEventHandling,
		"DEDUPLICATE_EVENTS":                 c.DeduplicateEvents,
		"USE_24_HOUR":                        c.TimeSettings.Use24Hour,
		"ADD_LEADING_ZERO":                   c.TimeSettings.AddLeadingZero,
		"DISPLAY_TIME":                       c.TimeSettings.DisplayTime,
		"SHOW_DATE_RANGE":                    c.ShowDateRange,
		"SHOW_TIMEZONE_IN_SUBHEADER":         c.ShowTimezoneInSubheader,
		"TZ":                                 c.Timezone,
		"SCHEDULE_TYPE":                      c.ScheduleSettings.ScheduleType,
		"SCHEDULE_DAY":                       c.ScheduleSettings.ScheduleDay,
		"RUN_TIME":                           c.ScheduleSettings.RunTime,
		"CRON_SCHEDULE":                      c.ScheduleSettings.CronSchedule,
		"RUN_ON_STARTUP":                     c.ScheduleSettings.RunOnStartup,
		"DEBUG":                              c.LoggingSettings.DebugMode,
		"HTTP_TIMEOUT":                       c.HTTPTimeout,
		"LOG_MAX_SIZE_MB":                    c.LoggingSettings.MaxSizeMB,
		"LOG_BACKUP_COUNT":                   c.LoggingSettings.BackupCount,

		// Standard snake_case & struct keys
		"language":                          c.Language,
		"use_discord":                       c.UseDiscord,
		"discord_webhook_url":               c.DiscordWebhookURL,
		"discord_mention_role_id":           c.DiscordMentionRoleID,
		"discord_hide_mention_instructions": c.DiscordHideMentionInstructions,
		"discord_timestamp_style":          c.DiscordTimestampStyle,
		"enable_custom_discord_footer":      c.EnableCustomDiscordFooter,
		"use_slack":                         c.UseSlack,
		"slack_webhook_url":                 c.SlackWebhookURL,
		"enable_custom_slack_footer":        c.EnableCustomSlackFooter,
		"calendar_urls":                     calURLs,
		"passed_event_handling":            c.PassedEventHandling,
		"deduplicate_events":                c.DeduplicateEvents,
		"show_date_range":                   c.ShowDateRange,
		"show_timezone_in_subheader":        c.ShowTimezoneInSubheader,
		"timezone":                          c.Timezone,
		"http_timeout":                      c.HTTPTimeout,
		"time_settings":                     c.TimeSettings,
		"schedule_settings":                 c.ScheduleSettings,
		"logging_settings":                  c.LoggingSettings,
	}

	return json.Marshal(m)
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type Alias Config
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil
	}

	c.applyFlatStringFields(raw)
	c.applyFlatBoolFields(raw)
	c.applyFlatScheduleAndLoggingFields(raw)
	c.applyFlatCalendarURLs(raw)

	return nil
}

func (c *Config) applyFlatStringFields(raw map[string]interface{}) {
	if val, ok := raw["APP_LANGUAGE"].(string); ok && val != "" {
		c.Language = val
	}
	if val, ok := raw["DISCORD_WEBHOOK_URL"].(string); ok {
		c.DiscordWebhookURL = val
	}
	if val, ok := raw["DISCORD_MENTION_ROLE_ID"].(string); ok {
		c.DiscordMentionRoleID = val
	}
	if val, ok := raw["DISCORD_TIMESTAMP_STYLE"].(string); ok {
		c.DiscordTimestampStyle = val
	}
	if val, ok := raw["SLACK_WEBHOOK_URL"].(string); ok {
		c.SlackWebhookURL = val
	}
	if val, ok := raw["PASSED_EVENT_HANDLING"].(string); ok {
		c.PassedEventHandling = val
	}
	if val, ok := raw["TZ"].(string); ok && val != "" {
		c.Timezone = val
	}
}

func (c *Config) applyFlatBoolFields(raw map[string]interface{}) {
	if val, ok := raw["USE_DISCORD"].(bool); ok {
		c.UseDiscord = val
	}
	if val, ok := raw["DISCORD_HIDE_MENTION_INSTRUCTIONS"].(bool); ok {
		c.DiscordHideMentionInstructions = val
	}
	if val, ok := raw["ENABLE_CUSTOM_DISCORD_FOOTER"].(bool); ok {
		c.EnableCustomDiscordFooter = val
	}
	if val, ok := raw["USE_SLACK"].(bool); ok {
		c.UseSlack = val
	}
	if val, ok := raw["ENABLE_CUSTOM_SLACK_FOOTER"].(bool); ok {
		c.EnableCustomSlackFooter = val
	}
	if val, ok := raw["DEDUPLICATE_EVENTS"].(bool); ok {
		c.DeduplicateEvents = val
	}
	if val, ok := raw["SHOW_DATE_RANGE"].(bool); ok {
		c.ShowDateRange = val
	}
	if val, ok := raw["SHOW_TIMEZONE_IN_SUBHEADER"].(bool); ok {
		c.ShowTimezoneInSubheader = val
	}
	if val, ok := raw["USE_24_HOUR"].(bool); ok {
		c.TimeSettings.Use24Hour = val
	}
	if val, ok := raw["ADD_LEADING_ZERO"].(bool); ok {
		c.TimeSettings.AddLeadingZero = val
	}
	if val, ok := raw["DISPLAY_TIME"].(bool); ok {
		c.TimeSettings.DisplayTime = val
	}
}

func (c *Config) applyFlatScheduleAndLoggingFields(raw map[string]interface{}) {
	if val, ok := raw["SCHEDULE_TYPE"].(string); ok && val != "" {
		c.ScheduleSettings.ScheduleType = val
	}
	if val, ok := raw["SCHEDULE_DAY"].(string); ok && val != "" {
		c.ScheduleSettings.ScheduleDay = val
	}
	if val, ok := raw["RUN_TIME"].(string); ok && val != "" {
		c.ScheduleSettings.RunTime = val
	}
	if val, ok := raw["CRON_SCHEDULE"].(string); ok {
		c.ScheduleSettings.CronSchedule = val
	}
	if val, ok := raw["RUN_ON_STARTUP"].(bool); ok {
		c.ScheduleSettings.RunOnStartup = val
	}
	if val, ok := raw["DEBUG"].(bool); ok {
		c.LoggingSettings.DebugMode = val
	}
	if val, ok := raw["HTTP_TIMEOUT"].(float64); ok {
		c.HTTPTimeout = int(val)
	}
	if val, ok := raw["LOG_MAX_SIZE_MB"].(float64); ok {
		c.LoggingSettings.MaxSizeMB = int(val)
	}
	if val, ok := raw["LOG_BACKUP_COUNT"].(float64); ok {
		c.LoggingSettings.BackupCount = int(val)
	}
}

func (c *Config) applyFlatCalendarURLs(raw map[string]interface{}) {
	if rawURLs, ok := raw["CALENDAR_URLS"].([]interface{}); ok {
		var urls []CalendarUrl
		for _, item := range rawURLs {
			if itemMap, ok := item.(map[string]interface{}); ok {
				u, _ := itemMap["url"].(string)
				t, _ := itemMap["type"].(string)
				if u != "" {
					urls = append(urls, CalendarUrl{URL: u, Type: t})
				}
			}
		}
		c.CalendarURLs = urls
	}
}
