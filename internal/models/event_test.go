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

	if !futureEv.IsMovie() {
		t.Errorf("Expected futureEv to be movie")
	}

	key := pastEv.DeduplicationKey(time.UTC)
	if key == "" {
		t.Errorf("DeduplicationKey should not be empty")
	}
}
