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
