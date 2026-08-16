package calendar

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/khw315/calendarr/internal/constants"
	"github.com/khw315/calendarr/internal/models"
)

func TestCalendarServiceFetchEvents(t *testing.T) {
	sampleICS := `BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Calendarr Test//EN
BEGIN:VEVENT
UID:test-event-1
SUMMARY:Stranger Things S04E01
DTSTART:20260816T120000Z
DTEND:20260816T130000Z
DESCRIPTION:Season 4 Premiere
END:VEVENT
END:VCALENDAR`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/calendar")
		w.Write([]byte(sampleICS))
	}))
	defer ts.Close()

	svc := NewService()
	cfg := models.DefaultConfig()
	cfg.CalendarURLs = []models.CalendarUrl{
		{URL: ts.URL, Type: constants.EventTypeTV},
	}

	startDate, _ := time.Parse("2006-01-02", "2026-08-01")
	endDate, _ := time.Parse("2006-01-02", "2026-08-31")

	events, err := svc.FetchEvents(context.Background(), cfg, startDate, endDate)
	if err != nil {
		t.Fatalf("FetchEvents failed: %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(events))
	}

	ev := events[0]
	if ev.Summary != "Stranger Things S04E01" {
		t.Errorf("Unexpected summary: %s", ev.Summary)
	}
	if ev.UID != "test-event-1" {
		t.Errorf("Unexpected UID: %s", ev.UID)
	}
	if ev.Description != "Season 4 Premiere" {
		t.Errorf("Unexpected description: %s", ev.Description)
	}
}

func TestCalendarServiceEmpty(t *testing.T) {
	svc := NewService()
	cfg := models.DefaultConfig()
	events, err := svc.FetchEvents(context.Background(), cfg, time.Now(), time.Now())
	if err != nil {
		t.Fatalf("FetchEvents on empty config failed: %v", err)
	}
	if len(events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(events))
	}
}
