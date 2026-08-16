package calendar

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/khw315/calendarr/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) FetchEvents(ctx context.Context, cfg *models.Config, startDate, endDate time.Time) ([]*models.Event, error) {
	if len(cfg.CalendarURLs) == 0 {
		return nil, nil
	}

	client := &http.Client{
		Timeout: time.Duration(cfg.HTTPTimeout) * time.Second,
	}

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		allEvents []*models.Event
		seenKeys  = make(map[string]bool)
		loc       = cfg.TimezoneLocation
	)

	if loc == nil {
		loc = time.UTC
	}

	for _, calURL := range cfg.CalendarURLs {
		if strings.TrimSpace(calURL.URL) == "" {
			continue
		}

		wg.Add(1)
		go func(cURL models.CalendarUrl) {
			defer wg.Done()

			events, err := s.fetchFromURL(ctx, client, cURL, startDate, endDate, loc)
			if err != nil {
				log.Printf("⚠️ Error fetching calendar feed [%s]: %v", cURL.URL, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			for _, ev := range events {
				if cfg.DeduplicateEvents {
					key := ev.DeduplicationKey(loc)
					if seenKeys[key] {
						continue
					}
					seenKeys[key] = true
				}
				allEvents = append(allEvents, ev)
			}
		}(calURL)
	}

	wg.Wait()

	// Sort events by StartTime ascending
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].StartTime.Before(allEvents[j].StartTime)
	})

	log.Printf("🗓️ Total events retrieved across feeds: %d", len(allEvents))
	return allEvents, nil
}

func (s *Service) fetchFromURL(ctx context.Context, client *http.Client, calURL models.CalendarUrl, startDate, endDate time.Time, loc *time.Location) ([]*models.Event, error) {
	targetURL := calURL.URL

	// Ensure pastDays and futureDays parameters exist
	parsed, err := url.Parse(targetURL)
	if err == nil {
		q := parsed.Query()
		if !q.Has("pastDays") {
			q.Set("pastDays", "30")
		}
		if !q.Has("futureDays") {
			q.Set("futureDays", "30")
		}
		parsed.RawQuery = q.Encode()
		targetURL = parsed.String()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("calendar returned non-200 HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read calendar body: %w", err)
	}

	cal, err := ics.ParseCalendar(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse iCal feed: %w", err)
	}

	var results []*models.Event

	for _, component := range cal.Components {
		vevent, ok := component.(*ics.VEvent)
		if !ok {
			continue
		}

		summaryProp := vevent.GetProperty(ics.ComponentPropertySummary)
		if summaryProp == nil || summaryProp.Value == "" {
			continue
		}

		startTime, err := s.extractTime(vevent, ics.ComponentPropertyDtStart, loc)
		if err != nil {
			continue
		}

		endTime, err := s.extractTime(vevent, ics.ComponentPropertyDtEnd, loc)
		if err != nil || endTime.Before(startTime) || endTime.Equal(startTime) {
			endTime = startTime.Add(1 * time.Hour)
		}

		// Filter events outside requested window
		if startTime.Before(startDate) || startTime.After(endDate) {
			continue
		}

		uid := vevent.Id()

		description := ""
		if descProp := vevent.GetProperty(ics.ComponentPropertyDescription); descProp != nil {
			description = descProp.Value
		}

		ev := &models.Event{
			Summary:     summaryProp.Value,
			StartTime:   startTime,
			EndTime:     endTime,
			SourceType:  calURL.Type,
			UID:         uid,
			Description: description,
		}

		results = append(results, ev)
	}

	return results, nil
}

func (s *Service) extractTime(vevent *ics.VEvent, propName ics.ComponentProperty, loc *time.Location) (time.Time, error) {
	prop := vevent.GetProperty(propName)
	if prop == nil {
		return time.Time{}, fmt.Errorf("missing property %s", propName)
	}

	val := prop.Value
	if val == "" {
		return time.Time{}, fmt.Errorf("empty property %s", propName)
	}

	// Try standard iCal time formats
	formats := []string{
		"20060102T150405Z",
		"20060102T150405",
		"20060102",
	}

	var parsed time.Time
	var parseErr error

	for _, f := range formats {
		parsed, parseErr = time.Parse(f, val)
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		return time.Time{}, parseErr
	}

	if loc != nil {
		parsed = parsed.In(loc)
	}

	return parsed, nil
}
