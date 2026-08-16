package constants

// Event Types
const (
	EventTypeTV    = "tv"
	EventTypeMovie = "movie"
)

// Platforms
const (
	PlatformDiscord = "discord"
	PlatformSlack   = "slack"
)

// Handling passed events
const (
	PassedEventDisplay = "DISPLAY"
	PassedEventHide    = "HIDE"
	PassedEventStrike  = "STRIKE"
)

// Bulk thresholds
const (
	BulkThresholdUpcoming = 2
	BulkThresholdPast     = 4
)

// Defaults
const (
	DefaultConfigPath                      = "/app/config/calendarr.json"
	DefaultLocalConfigPath                 = "config/calendarr.json"
	DefaultLogDir                          = "/app/logs"
	DefaultLogFile                         = "calendarr.log"
	DefaultRunTime                         = "09:00"
	DefaultScheduleType                    = "WEEKLY"
	DefaultScheduleDay                     = "1" // Monday
	DefaultPassedEventHandling            = PassedEventDisplay
	MaxDiscordEmbedsPerRequest           = 10
	DiscordEmbedPayloadThreshold         = 5800
	DefaultUseDiscord                     = false
	DefaultUseSlack                       = false
	DefaultShowDateRange                  = false
	DefaultShowTimezoneInSubheader        = false
	DefaultDeduplicateEvents              = false
	DefaultDiscordHideMentionInstructions = true
	DefaultUse24Hour                      = true
	DefaultAddLeadingZero                 = true
	DefaultDisplayTime                    = true
	DefaultRunOnStartup                   = false
	DefaultLogBackupCount                 = 15
	DefaultLogMaxSizeMB                   = 1
	DefaultDebugMode                      = false
	DefaultEnableCustomDiscordFooter      = false
	DefaultEnableCustomSlackFooter        = false
	DefaultHTTPTimeout                    = 30
)

// Colors for Discord & Slack
var DiscordColors = map[string]int{
	"red":    15158332,
	"orange": 15844367,
	"yellow": 16776960,
	"green":  5763719,
	"blue":   3447003,
	"indigo": 10181046,
	"violet": 9846527,
}

var SlackColors = map[string]string{
	"red":    "#E53935",
	"orange": "#FB8C00",
	"yellow": "#FFD600",
	"green":  "#43A047",
	"blue":   "#1E88E5",
	"indigo": "#5E35B1",
	"violet": "#8E24AA",
}
