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

	timeoutSec := cfg.HTTPTimeout
	if timeoutSec <= 0 {
		timeoutSec = 30
	}

	client := &http.Client{
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	parentCtx := ctx
	if parentCtx == nil {
		parentCtx = context.Background()
	}
	fetchCtx, cancel := context.WithTimeout(parentCtx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	allEvents := s.fetchConcurrent(fetchCtx, client, cfg, startDate, endDate, loc)

	// Sort events by StartTime ascending
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].StartTime.Before(allEvents[j].StartTime)
	})

	log.Printf("🗓️ Total events retrieved across feeds: %d", len(allEvents))
	return allEvents, nil
}

func (s *Service) fetchConcurrent(ctx context.Context, client *http.Client, cfg *models.Config, startDate, endDate time.Time, loc *time.Location) []*models.Event {
	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		allEvents []*models.Event
		seenKeys  = make(map[string]bool)
	)

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
			s.appendUniqueEvents(&allEvents, events, cfg, loc, seenKeys)
			mu.Unlock()
		}(calURL)
	}

	wg.Wait()
	return allEvents
}

func (s *Service) appendUniqueEvents(allEvents *[]*models.Event, events []*models.Event, cfg *models.Config, loc *time.Location, seenKeys map[string]bool) {
	for _, ev := range events {
		if cfg.DeduplicateEvents {
			key := ev.DeduplicationKey(loc)
			if seenKeys[key] {
				continue
			}
			seenKeys[key] = true
		}
		*allEvents = append(*allEvents, ev)
	}
}

func (s *Service) fetchFromURL(ctx context.Context, client *http.Client, calURL models.CalendarUrl, startDate, endDate time.Time, loc *time.Location) ([]*models.Event, error) {
	targetURL := s.buildTargetURL(calURL.URL)

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
		if vevent, ok := component.(*ics.VEvent); ok {
			if ev := s.parseVEvent(vevent, calURL.Type, startDate, endDate, loc); ev != nil {
				results = append(results, ev)
			}
		}
	}

	return results, nil
}

func (s *Service) buildTargetURL(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := parsed.Query()
	if !q.Has("pastDays") {
		q.Set("pastDays", "30")
	}
	if !q.Has("futureDays") {
		q.Set("futureDays", "30")
	}
	parsed.RawQuery = q.Encode()
	return parsed.String()
}

func (s *Service) parseVEvent(vevent *ics.VEvent, sourceType string, startDate, endDate time.Time, loc *time.Location) *models.Event {
	summaryProp := vevent.GetProperty(ics.ComponentPropertySummary)
	if summaryProp == nil || summaryProp.Value == "" {
		return nil
	}

	startTime, err := vevent.GetStartAt()
	if err != nil {
		startTime, err = s.extractTime(vevent, ics.ComponentPropertyDtStart, loc)
		if err != nil {
			return nil
		}
	} else if loc != nil {
		startTime = startTime.In(loc)
	}

	endTime, err := vevent.GetEndAt()
	if err != nil || endTime.Before(startTime) || endTime.Equal(startTime) {
		endTime = startTime.Add(1 * time.Hour)
	} else if loc != nil {
		endTime = endTime.In(loc)
	}

	if startTime.Before(startDate) || startTime.After(endDate) {
		return nil
	}

	description := ""
	if descProp := vevent.GetProperty(ics.ComponentPropertyDescription); descProp != nil {
		description = descProp.Value
	}

	return &models.Event{
		Summary:     summaryProp.Value,
		StartTime:   startTime,
		EndTime:     endTime,
		SourceType:  sourceType,
		UID:         vevent.Id(),
		Description: description,
	}
}

func (s *Service) extractTime(vevent *ics.VEvent, propName ics.ComponentProperty, loc *time.Location) (time.Time, error) {
	prop := vevent.GetProperty(propName)
	if prop == nil || prop.Value == "" {
		return time.Time{}, fmt.Errorf("missing or empty property %s", propName)
	}

	formats := []string{
		"20060102T150405Z",
		"20060102T150405",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"20060102",
		"2006-01-02",
	}

	var parsed time.Time
	var parseErr error

	val := strings.TrimSpace(prop.Value)
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
