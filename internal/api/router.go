package api

import (
	"context"
	"embed"
	"encoding/json"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/khw315/calendarr/internal/config"
	"github.com/khw315/calendarr/internal/localization"
	"github.com/khw315/calendarr/internal/models"
	"github.com/khw315/calendarr/internal/services/calendar"
	"github.com/khw315/calendarr/internal/services/scheduler"
	"github.com/khw315/calendarr/internal/tzdata"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

type Router struct {
	cfgMgr     *config.Manager
	schedSvc   *scheduler.Service
	calSvc     *calendar.Service
	embeddedFS embed.FS
}

func NewRouter(cfgMgr *config.Manager, schedSvc *scheduler.Service, calSvc *calendar.Service, embeddedFS embed.FS) *Router {
	return &Router{
		cfgMgr:     cfgMgr,
		schedSvc:   schedSvc,
		calSvc:     calSvc,
		embeddedFS: embeddedFS,
	}
}

func (r *Router) Setup() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Global CORS middleware for API endpoints & FE cross-origin access
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	// API Routes mounted cleanly under /api
	apiRouter := chi.NewRouter()
	apiRouter.Get("/status", r.handleGetStatus)
	apiRouter.Get("/events", r.handleGetEvents)
	apiRouter.Get("/config", r.handleGetConfig)
	apiRouter.Get("/languages", r.handleGetLanguages)
	apiRouter.Get("/schedule", r.handleGetSchedule)
	apiRouter.Get("/releases", r.handleGetReleases)
	apiRouter.Get("/past-releases", r.handleGetPastReleases)
	apiRouter.Get("/timezones", r.handleGetTimezones)
	apiRouter.Post("/config", r.handlePostConfig)
	apiRouter.Post("/trigger", r.handlePostTrigger)

	router.Mount("/api", apiRouter)

	// Static Web UI serving
	publicSub, err := fs.Sub(r.embeddedFS, "public")
	if err != nil {
		log.Printf("⚠️ Could not sub public embed directory: %v", err)
	} else {
		fileServer := http.FileServer(http.FS(publicSub))
		router.Get("/*", func(w http.ResponseWriter, req *http.Request) {
			relPath := strings.TrimPrefix(req.URL.Path, "/")
			cleanedPath := filepath.Clean(relPath)
			if cleanedPath != "." && cleanedPath != "" && !strings.HasPrefix(cleanedPath, "..") {
				if _, err := publicSub.Open(cleanedPath); err == nil {
					fileServer.ServeHTTP(w, req)
					return
				}
			}
			// Fallback to SPA index.html
			req.URL.Path = "/"
			fileServer.ServeHTTP(w, req)
		})
	}

	return router
}

func (r *Router) handleGetStatus(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	status := r.schedSvc.GetStatus()
	_ = json.NewEncoder(w).Encode(status)
}

func (r *Router) handleGetEvents(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	events := r.schedSvc.GetCachedEvents()
	if events == nil {
		events = []*models.Event{}
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
		"count":  len(events),
	})
}

func (r *Router) handleGetConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	cfg := r.cfgMgr.Get()
	_ = json.NewEncoder(w).Encode(cfg)
}

func (r *Router) handleGetLanguages(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	langs := localization.GetLanguageList()
	_ = json.NewEncoder(w).Encode(langs)
}

func (r *Router) handleGetSchedule(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	cfg := r.cfgMgr.Get()
	status := r.schedSvc.GetStatus()
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"schedule_type": cfg.ScheduleSettings.ScheduleType,
		"next_run":      status["next_run"],
		"timezone":      cfg.Timezone,
	})
}

func (r *Router) handleGetReleases(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	cfg := r.cfgMgr.Get()

	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	days := 7
	if dStr := req.URL.Query().Get("days"); dStr != "" {
		if d, err := strconv.Atoi(dStr); err == nil && d > 0 {
			days = d
		}
	}

	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	startDate := startOfDay
	endDate := startOfDay.AddDate(0, 0, days).Add(-1 * time.Nanosecond)

	var events []*models.Event
	if r.calSvc != nil && len(cfg.CalendarURLs) > 0 {
		fetched, err := r.calSvc.FetchEvents(req.Context(), cfg, startDate, endDate)
		if err == nil {
			events = fetched
		}
	}
	if events == nil {
		events = r.schedSvc.GetCachedEvents()
	}

	r.respondWithEventsDTO(w, cfg, loc, events)
}

func (r *Router) handleGetPastReleases(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	cfg := r.cfgMgr.Get()

	loc := cfg.TimezoneLocation
	if loc == nil {
		loc = time.UTC
	}

	days := 7
	if dStr := req.URL.Query().Get("days"); dStr != "" {
		if d, err := strconv.Atoi(dStr); err == nil && d > 0 {
			days = d
		}
	}

	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	startDate := startOfDay.AddDate(0, 0, -days)
	endDate := startOfDay.Add(-1 * time.Nanosecond)

	var events []*models.Event
	if r.calSvc != nil && len(cfg.CalendarURLs) > 0 {
		fetched, err := r.calSvc.FetchEvents(req.Context(), cfg, startDate, endDate)
		if err == nil {
			events = fetched
		}
	}
	if events == nil {
		events = r.schedSvc.GetCachedEvents()
	}

	var pastEvents []*models.Event
	for _, ev := range events {
		if ev.StartTime.Before(now) {
			pastEvents = append(pastEvents, ev)
		}
	}

	r.respondWithEventsDTO(w, cfg, loc, pastEvents)
}

func (r *Router) respondWithEventsDTO(w http.ResponseWriter, cfg *models.Config, loc *time.Location, events []*models.Event) {
	grouped := make(map[string][]*models.Event)
	var dates []string

	for _, ev := range events {
		dayKey := ev.DayKey(loc)
		if _, exists := grouped[dayKey]; !exists {
			dates = append(dates, dayKey)
		}
		grouped[dayKey] = append(grouped[dayKey], ev)
	}

	type EventDTO struct {
		Title        string `json:"title"`
		Type         string `json:"type"`
		StartTime    string `json:"start_time"`
		EndTime      string `json:"end_time"`
		Date         string `json:"date"`
		Timestamp    int64  `json:"timestamp"`
		EndTimestamp int64  `json:"end_timestamp"`
		Description  string `json:"description,omitempty"`
	}

	type DayDTO struct {
		DayName string     `json:"day_name"`
		Date    string     `json:"date"`
		Events  []EventDTO `json:"events"`
	}

	var days []DayDTO
	var allEventDTOs []EventDTO
	for _, dayKey := range dates {
		dayEvents := grouped[dayKey]
		var evDTOs []EventDTO
		for _, ev := range dayEvents {
			t := ev.StartTime.In(loc)
			evType := "tv"
			if ev.IsMovie() {
				evType = "movie"
			}
			dto := EventDTO{
				Title:        ev.Summary,
				Type:         evType,
				StartTime:    t.Format("15:04"),
				EndTime:      ev.EndTime.In(loc).Format("15:04"),
				Date:         t.Format("2006-01-02"),
				Timestamp:    t.Unix(),
				EndTimestamp: ev.EndTime.In(loc).Unix(),
				Description:  ev.Description,
			}
			evDTOs = append(evDTOs, dto)
			allEventDTOs = append(allEventDTOs, dto)
		}
		dayStartTime := dayEvents[0].StartTime.In(loc)
		dateHeader := localization.FormatDateHeader(dayStartTime, "en")
		days = append(days, DayDTO{
			DayName: dateHeader,
			Date:    dayKey,
			Events:  evDTOs,
		})
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"days":   days,
		"events": allEventDTOs,
	})
}

func (r *Router) handleGetTimezones(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	allTZs := tzdata.GetTimezones()
	tzs := make(map[string]string, len(allTZs))
	for _, tz := range allTZs {
		tzs[tz] = tz
	}
	_ = json.NewEncoder(w).Encode(tzs)
}

func (r *Router) handlePostConfig(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read request body"})
		return
	}

	// Copy current active configuration so partial updates preserve all non-modified fields
	updatedCfg := *r.cfgMgr.Get()

	if err := json.Unmarshal(bodyBytes, &updatedCfg); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request JSON payload"})
		return
	}

	if err := r.cfgMgr.Save(&updatedCfg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	r.schedSvc.UpdateSchedule()

	// Trigger immediate background iCal reload on config save
	go func() {
		_, _ = r.schedSvc.TriggerRun(context.Background())
	}()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"success": true,
		"message": "Configuration saved, schedule updated, and feeds reloading in background",
		"config":  updatedCfg,
	})
}

func (r *Router) handlePostTrigger(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentTypeHeader, contentTypeJSON)

	go func() {
		_, _ = r.schedSvc.TriggerRun(context.Background())
	}()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":  "triggered",
		"message": "Manual run triggered successfully in background",
	})
}
