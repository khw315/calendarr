package models

import (
	"testing"
	"time"

	"github.com/khw315/calendarr/internal/constants"
)

func TestEventMethods(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-2 * time.Hour)
	futureTime := now.Add(2 * time.Hour)

	pastEv := Event{
		Summary:    "Breaking Bad S01E01",
		StartTime:  pastTime.Add(-1 * time.Hour),
		EndTime:    pastTime,
		SourceType: constants.EventTypeTV,
	}

	futureEv := Event{
		Summary:    "Inception 2010",
		StartTime:  futureTime,
		EndTime:    futureTime.Add(2 * time.Hour),
		SourceType: constants.EventTypeMovie,
	}

	if !pastEv.IsPast(now) {
		t.Errorf("Expected pastEv to be past")
	}

	if futureEv.IsPast(now) {
		t.Errorf("Expected futureEv to not be past")
	}

	if !pastEv.IsPremiere() {
		t.Errorf("Expected Breaking Bad S01E01 to be recognized as premiere")
	}

	if !pastEv.IsTV() {
		t.Errorf("Expected pastEv to be TV")
	}

	if !futureEv.IsMovie() {
		t.Errorf("Expected futureEv to be movie")
	}

	dayKey := pastEv.DayKey(time.UTC)
	if dayKey == "" {
		t.Errorf("DayKey should not be empty")
	}

	key := pastEv.DeduplicationKey(time.UTC)
	if key == "" {
		t.Errorf("DeduplicationKey should not be empty")
	}
}

func TestConfigToJSON(t *testing.T) {
	cfg := DefaultConfig()
	data, err := cfg.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("Expected non-empty JSON byte slice")
	}
}
