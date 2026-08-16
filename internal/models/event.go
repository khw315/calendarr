package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/khw315/calendarr/internal/constants"
)

var (
	premiereRegex = regexp.MustCompile(`(?i)[-\s](?:s\d+e0*1|(?:\d+x0*1))\b`)
)

type Event struct {
	Summary     string    `json:"summary"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	SourceType  string    `json:"source_type"` // "tv", "movie", "SONARR", "RADARR"
	UID         string    `json:"uid,omitempty"`
	Description string    `json:"description,omitempty"`
}

func (e *Event) IsPremiere() bool {
	return premiereRegex.MatchString(e.Summary)
}

func (e *Event) IsPast(now time.Time) bool {
	return e.EndTime.Before(now)
}

func (e *Event) DayKey(loc *time.Location) string {
	t := e.StartTime
	if loc != nil {
		t = t.In(loc)
	}
	return t.Format("2006-01-02")
}

func (e *Event) DeduplicationKey(loc *time.Location) string {
	t := e.StartTime
	if loc != nil {
		t = t.In(loc)
	}
	return strings.TrimSpace(e.Summary) + "|" + t.Format("2006-01-02")
}

func (e *Event) IsMovie() bool {
	st := strings.ToUpper(e.SourceType)
	return st == "RADARR" || st == "MOVIE" || st == strings.ToUpper(constants.EventTypeMovie)
}

func (e *Event) IsTV() bool {
	st := strings.ToUpper(e.SourceType)
	return st == "SONARR" || st == "TV" || st == strings.ToUpper(constants.EventTypeTV)
}
